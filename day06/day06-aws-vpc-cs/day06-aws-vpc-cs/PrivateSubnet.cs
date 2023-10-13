using Pulumi;
using Pulumi.Aws.Ec2;
using Aws = Pulumi.Aws;

namespace day06_aws_vpc_cs
{
    public class PrivateSubnetArgs
    {
        public Input<string> VpcId { get; set; }
        public string Cidr { get; set; }
        public Input<string> Az { get; set; }
        public InputMap<string> Tags { get; set; }
    }

    public class PrivateSubnet
    {
        private string Name { get; set; }
        public Subnet Subnet { get; set; }

        public RouteTable RouteTable { get; set; }

        public PrivateSubnet(string name, PrivateSubnetArgs args)
        {
            Name = name;

            var subnetName = args.Az.Apply(az => "my-private-subnet-" + az);

            Subnet = new Subnet(Name, new SubnetArgs()
            {
                VpcId = args.VpcId,
                CidrBlock = args.Cidr,
                AvailabilityZone = args.Az,
                Tags = InputMap<string>.Merge(args.Tags, new InputMap<string>
                {
                    { "Name", subnetName }
                })
            });

            RouteTable = new RouteTable(this.Name, new RouteTableArgs
            {
                VpcId = args.VpcId,
                Tags = InputMap<string>.Merge(args.Tags, new InputMap<string>
                {
                    { "Name", subnetName.Apply(n => n + "-rt") }
                })
            });

            new Aws.Ec2.RouteTableAssociation(this.Name, new RouteTableAssociationArgs
            {
                RouteTableId = RouteTable.Id,
                SubnetId = Subnet.Id
            });
        }

        public void AddNatGateway(Input<string> natGatewayId)
        {
            new Aws.Ec2.Route(Name + "-nat-gateway-route", new RouteArgs
            {
                DestinationCidrBlock = "0.0.0.0/0",
                RouteTableId = RouteTable.Id,
                NatGatewayId = natGatewayId
            });
        }
    }
}