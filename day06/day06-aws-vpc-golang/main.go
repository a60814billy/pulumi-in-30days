package main

import (
	"github.com/a60814billy/pulumi-in-30days/day06/day06-aws-vpc-golang/private_subnet"
	"github.com/a60814billy/pulumi-in-30days/day06/day06-aws-vpc-golang/subnet"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws"
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type PrivateSubnetArgs struct {
	CidrBlock string
	Az        string
}

func AppendNameTag(tags pulumi.StringMap, name string) pulumi.StringMap {
	newTags := pulumi.StringMap(make(map[string]pulumi.StringInput))
	for k, v := range tags {
		newTags[k] = v
	}
	newTags["Name"] = pulumi.String(name)
	return newTags
}

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		defaultTags := pulumi.StringMap{
			"pulumi:project": pulumi.String(ctx.Project()),
			"pulumi:stack":   pulumi.String(ctx.Stack()),
			"ManagedBy":      pulumi.String("Pulumi"),
		}

		// 宣告一個變數儲存 vpcCidr
		vpcCidr := "10.120.0.0/16"
		// 使用 splitSubnets 函式，將 10.120.0.0/16 分割成 256 個 /24 的子網路
		subnetting := subnet.Split(vpcCidr, 24)
		// 前 128 個子網路為 public 子網路
		// 10.120.0.0/24, 10.120.1.0/24 .... 10.120.127.0/24
		publicSubnetCidrs := subnetting[:len(subnetting)/2]

		// 後 128 個子網路為 private 子網路
		// 10.120.128.0/24, 10.120.129.0/24, ... 10.120.255.0/24
		privateSubnetCidrs := subnetting[len(subnetting)/2:]

		nonLocalAvailabilityZones, err := aws.GetAvailabilityZones(ctx, &aws.GetAvailabilityZonesArgs{
			Filters: []aws.GetAvailabilityZonesFilter{
				{
					Name:   "opt-in-status",
					Values: []string{"opt-in-not-required"},
				},
			},
		})
		if err != nil {
			return err
		}
		azNames := nonLocalAvailabilityZones.Names

		myVpc, err := ec2.NewVpc(ctx, "my-vpc", &ec2.VpcArgs{
			CidrBlock: pulumi.StringPtr("10.120.0.0/16"),
			Tags:      AppendNameTag(defaultTags, "my-vpc"),
		})
		if err != nil {
			return err
		}
		igw, err := ec2.NewInternetGateway(ctx, "my-igw", &ec2.InternetGatewayArgs{
			VpcId: myVpc.ID(),
			Tags:  AppendNameTag(defaultTags, "my-igw"),
		})
		if err != nil {
			return err
		}

		defaultRt, err := ec2.NewDefaultRouteTable(ctx, "my-default-rt", &ec2.DefaultRouteTableArgs{
			DefaultRouteTableId: myVpc.DefaultRouteTableId,
			Tags:                AppendNameTag(defaultTags, "my-default-rt"),
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
			CidrBlock:        pulumi.StringPtr(publicSubnetCidrs[0]),
			AvailabilityZone: pulumi.StringPtr(azNames[0]),
			Tags:             AppendNameTag(defaultTags, "my-public-subnet-1"),
		})
		if err != nil {
			return err
		}
		publicSubnets["my-public-subnet-1"] = publicSubnet1

		publicSubnet2, err := ec2.NewSubnet(ctx, "my-public-subnet-2", &ec2.SubnetArgs{
			VpcId:            myVpc.ID(),
			CidrBlock:        pulumi.StringPtr(publicSubnetCidrs[1]),
			AvailabilityZone: pulumi.StringPtr(azNames[1]),
			Tags:             AppendNameTag(defaultTags, "my-public-subnet-2"),
		})
		if err != nil {
			return err
		}
		publicSubnets["my-public-subnet-2"] = publicSubnet2

		for name, s := range publicSubnets {
			_, err := ec2.NewRouteTableAssociation(ctx, name+"-rt-association", &ec2.RouteTableAssociationArgs{
				RouteTableId: defaultRt.ID(),
				SubnetId:     s.ID(),
			})
			if err != nil {
				return err
			}
		}

		eip, err := ec2.NewEip(ctx, "my-nat-gateway-eip", &ec2.EipArgs{
			Tags: AppendNameTag(defaultTags, "my-nat-gateway-eip"),
		})
		if err != nil {
			return err
		}

		natGateway, err := ec2.NewNatGateway(ctx, "my-nat-gateway", &ec2.NatGatewayArgs{
			SubnetId:     publicSubnet1.ID(),
			AllocationId: eip.ID(),
			Tags:         AppendNameTag(defaultTags, "my-nat-gateway"),
		})
		if err != nil {
			return err
		}

		privateSubnetArgs := []PrivateSubnetArgs{
			{CidrBlock: privateSubnetCidrs[0], Az: azNames[0]},
			{CidrBlock: privateSubnetCidrs[1], Az: azNames[1]},
		}

		privateSubnets := make(map[string]*ec2.Subnet)

		for _, args := range privateSubnetArgs {
			subnetName := "my-private-subnet-" + args.Az
			privateSubnet, err := private_subnet.NewPrivateSubnet(ctx, subnetName, &private_subnet.PrivateSubnetArgs{
				VpcID:       myVpc.ID(),
				Cidr:        args.CidrBlock,
				AZ:          args.Az,
				DefaultTags: defaultTags,
			})
			if err != nil {
				return err
			}
			err = privateSubnet.AddNatGateway(natGateway.ID())
			if err != nil {
				return err
			}
			privateSubnets[subnetName] = privateSubnet.Subnet
		}
		return nil
	})
}
