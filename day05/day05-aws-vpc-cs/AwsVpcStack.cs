using System.Collections.Generic;
using System.Linq;
using Aws = Pulumi.Aws;

namespace day05_aws_vpc_cs
{
    class AwsVpcStack : Pulumi.Stack
    {
        public AwsVpcStack()
        {
            var myVpc = new Aws.Ec2.Vpc("my-vpc", new Aws.Ec2.VpcArgs
            {
                CidrBlock = "10.120.0.0/16",
            });

            var igw = new Aws.Ec2.InternetGateway("my-igw", new Aws.Ec2.InternetGatewayArgs
            {
                VpcId = myVpc.Id,
            });

            var defaultRt = new Aws.Ec2.DefaultRouteTable("my-default-rt", new Aws.Ec2.DefaultRouteTableArgs
            {
                DefaultRouteTableId = myVpc.DefaultRouteTableId,
            });

            new Aws.Ec2.Route("my-default-rt-default-route", new Aws.Ec2.RouteArgs
            {
                RouteTableId = defaultRt.Id,
                DestinationCidrBlock = "0.0.0.0/0",
                GatewayId = igw.Id,
            });

            var publicSubnet = new Dictionary<string, Aws.Ec2.Subnet>
            {
                {
                    "my-public-subnet-1", new Aws.Ec2.Subnet("my-public-subnet-1", new Aws.Ec2.SubnetArgs
                    {
                        VpcId = myVpc.Id,
                        CidrBlock = "10.120.0.0/24",
                        AvailabilityZone = "ap-east-1a",
                    })
                },
                {
                    "my-public-subnet-2", new Aws.Ec2.Subnet("my-public-subnet-2", new Aws.Ec2.SubnetArgs
                    {
                        VpcId = myVpc.Id,
                        CidrBlock = "10.120.1.0/24",
                        AvailabilityZone = "ap-east-1a",
                    })
                }
            };

            foreach (var kv in publicSubnet)
            {
                new Aws.Ec2.RouteTableAssociation(kv.Key + "-rt-association", new Aws.Ec2.RouteTableAssociationArgs
                {
                    RouteTableId = defaultRt.Id,
                    SubnetId = kv.Value.Id,
                });
            }

            var eip = new Aws.Ec2.Eip("my-nat-gateway-eip");

            var natGateway = new Aws.Ec2.NatGateway("my-nat-gateway", new Aws.Ec2.NatGatewayArgs
            {
                SubnetId = publicSubnet.First().Value.Id,
                AllocationId = eip.Id,
            });

            var privateSubnetArgs = new[]
            {
                new { Cidr = "10.120.128.0/24", Az = "ap-east-1a" },
                new { Cidr = "10.120.129.0/24", Az = "ap-east-1b" }
            };

            var privateSubnets = new Dictionary<string, Aws.Ec2.Subnet>();

            foreach (var arg in privateSubnetArgs)
            {
                var myPrivateSubnetName = "my-private-subnet-" + arg.Az;

                var subnet = new Aws.Ec2.Subnet(myPrivateSubnetName, new Aws.Ec2.SubnetArgs
                {
                    VpcId = myVpc.Id,
                    CidrBlock = arg.Cidr,
                    AvailabilityZone = arg.Az
                });
                privateSubnets.Add(myPrivateSubnetName, subnet);

                var rt = new Aws.Ec2.RouteTable(myPrivateSubnetName + "-rt", new Aws.Ec2.RouteTableArgs
                {
                    VpcId = myVpc.Id,
                    Routes = new[]
                    {
                        new Aws.Ec2.Inputs.RouteTableRouteArgs()
                            { CidrBlock = "0.0.0.0/0", NatGatewayId = natGateway.Id }
                    }
                });

                new Aws.Ec2.RouteTableAssociation(myPrivateSubnetName + "-rt-association",
                    new Aws.Ec2.RouteTableAssociationArgs
                    {
                        RouteTableId = rt.Id,
                        SubnetId = subnet.Id,
                    });
            }
        }
    }
}