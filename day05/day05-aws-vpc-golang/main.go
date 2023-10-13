package main

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type PrivateSubnetArgs struct {
	CidrBlock string
	Az        string
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		myVpc, err := ec2.NewVpc(ctx, "my-vpc", &ec2.VpcArgs{
			CidrBlock: pulumi.StringPtr("10.120.0.0/16"),
		})
		if err != nil {
			return err
		}
		igw, err := ec2.NewInternetGateway(ctx, "my-igw", &ec2.InternetGatewayArgs{
			VpcId: myVpc.ID(),
		})
		if err != nil {
			return err
		}

		defaultRt, err := ec2.NewDefaultRouteTable(ctx, "my-default-rt", &ec2.DefaultRouteTableArgs{
			DefaultRouteTableId: myVpc.DefaultRouteTableId,
		})
		if err != nil {
			return err
		}

		_, err = ec2.NewRoute(ctx, "my-route", &ec2.RouteArgs{
			DestinationCidrBlock: pulumi.StringPtr("0.0.0.0/0"),
			RouteTableId:         defaultRt.ID(),
			GatewayId:            igw.ID(),
		})
		if err != nil {
			return err
		}

		publicSubnets := make(map[string]*ec2.Subnet)

		publicSubnet1, err := ec2.NewSubnet(ctx, "my-public-subnet-1", &ec2.SubnetArgs{
			VpcId:            myVpc.ID(),
			CidrBlock:        pulumi.StringPtr("10.120.0.0/24"),
			AvailabilityZone: pulumi.StringPtr("ap-east-1a"),
		})
		if err != nil {
			return err
		}
		publicSubnets["my-public-subnet-1"] = publicSubnet1

		publicSubnet2, err := ec2.NewSubnet(ctx, "my-public-subnet-2", &ec2.SubnetArgs{
			VpcId:            myVpc.ID(),
			CidrBlock:        pulumi.StringPtr("10.120.1.0/24"),
			AvailabilityZone: pulumi.StringPtr("ap-east-1b"),
		})
		if err != nil {
			return err
		}
		publicSubnets["my-public-subnet-2"] = publicSubnet2

		for name, subnet := range publicSubnets {
			_, err := ec2.NewRouteTableAssociation(ctx, name+"-rt-association", &ec2.RouteTableAssociationArgs{
				RouteTableId: defaultRt.ID(),
				SubnetId:     subnet.ID(),
			})
			if err != nil {
				return err
			}
		}

		eip, err := ec2.NewEip(ctx, "my-eip", &ec2.EipArgs{})
		if err != nil {
			return err
		}

		natGateway, err := ec2.NewNatGateway(ctx, "my-nat-gateway", &ec2.NatGatewayArgs{
			SubnetId:     publicSubnet1.ID(),
			AllocationId: eip.ID(),
		})
		if err != nil {
			return err
		}

		privateSubnetArgs := []PrivateSubnetArgs{
			{CidrBlock: "10.120.128.0/24", Az: "ap-east-1a"},
			{CidrBlock: "10.120.129.0/24", Az: "ap-east-1b"},
		}

		privateSubnets := make(map[string]*ec2.Subnet)

		for _, args := range privateSubnetArgs {
			subnetName := "my-private-subnet-" + args.Az
			subnet, err := ec2.NewSubnet(ctx, subnetName, &ec2.SubnetArgs{
				VpcId:            myVpc.ID(),
				CidrBlock:        pulumi.StringPtr(args.CidrBlock),
				AvailabilityZone: pulumi.StringPtr(args.Az),
			})
			if err != nil {
				return err
			}
			privateSubnets[subnetName] = subnet

			rt, err := ec2.NewRouteTable(ctx, subnetName+"-rt", &ec2.RouteTableArgs{
				VpcId: myVpc.ID(),
				Routes: ec2.RouteTableRouteArray{
					&ec2.RouteTableRouteArgs{
						CidrBlock:    pulumi.StringPtr(args.CidrBlock),
						NatGatewayId: natGateway.ID(),
					},
				},
			})
			if err != nil {
				return err
			}

			_, err = ec2.NewRouteTableAssociation(ctx, subnetName+"-rt-association", &ec2.RouteTableAssociationArgs{
				RouteTableId: rt.ID(),
				SubnetId:     subnet.ID(),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
}
