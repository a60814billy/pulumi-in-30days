package private_subnet

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

type PrivateSubnetArgs struct {
	VpcID       pulumi.StringInput
	Cidr        string
	AZ          string
	DefaultTags pulumi.StringMap
}

type PrivateSubnet struct {
	ctx        *pulumi.Context
	name       string
	Subnet     *ec2.Subnet
	RouteTable *ec2.RouteTable
}

func (p *PrivateSubnet) AddNatGateway(natGatewayID pulumi.StringInput) error {
	_, err := ec2.NewRoute(p.ctx, p.name+"-nat-gateway-route", &ec2.RouteArgs{
		RouteTableId:         p.RouteTable.ID(),
		DestinationCidrBlock: pulumi.String("0.0.0.0/0"),
		NatGatewayId:         natGatewayID,
	})
	if err != nil {
		return err
	}
	return nil
}

func NewPrivateSubnet(ctx *pulumi.Context, name string, args *PrivateSubnetArgs) (*PrivateSubnet, error) {

	privateSubnet := &PrivateSubnet{
		ctx:  ctx,
		name: name,
	}

	subnet, err := ec2.NewSubnet(ctx, name, &ec2.SubnetArgs{
		VpcId:            args.VpcID,
		CidrBlock:        pulumi.StringPtr(args.Cidr),
		AvailabilityZone: pulumi.StringPtr(args.AZ),
	})
	if err != nil {
		return nil, err
	}
	privateSubnet.Subnet = subnet

	rt, err := ec2.NewRouteTable(ctx, name+"-rt", &ec2.RouteTableArgs{
		VpcId: args.VpcID,
	})
	if err != nil {
		return nil, err
	}
	privateSubnet.RouteTable = rt

	_, err = ec2.NewRouteTableAssociation(ctx, name+"-rt-association", &ec2.RouteTableAssociationArgs{
		RouteTableId: rt.ID(),
		SubnetId:     subnet.ID(),
	})
	if err != nil {
		return nil, err
	}

	return privateSubnet, nil
}
