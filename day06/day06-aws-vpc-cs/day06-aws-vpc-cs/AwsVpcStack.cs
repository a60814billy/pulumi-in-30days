using System.Collections;
using System.Collections.Generic;
using System.Linq;
using Pulumi;
using Pulumi.Aws.Ec2;
using Pulumi.Aws.Inputs;
using Aws = Pulumi.Aws;

namespace day06_aws_vpc_cs
{
    class AwsVpcStack : Pulumi.Stack
    {
        public AwsVpcStack()
        {
            var stack = Deployment.Instance;
            var defaultTags = new InputMap<string>
            {
                { "pulumi:project", stack.ProjectName },
                { "pulumi:stack", stack.StackName },
                { "ManagedBy", "Pulumi" }
            };

            // 宣告一個變數儲存 vpcCidr
            var vpcCidr = "10.120.0.0/16";
            // 使用 splitSubnets 函式，將 10.120.0.0/16 分割成 256 個 /24 的子網路
            var subnetting = SubnetUtils.SplitSubnets(vpcCidr, 24);

            // 前 128 個子網路為 public 子網路
            // 10.120.0.0/24, 10.120.1.0/24 .... 10.120.127.0/24

            var publicSubnetCidrs = subnetting.Take(128).ToList();

            // 後 128 個子網路為 private 子網路
            // 10.120.128.0/24, 10.120.129.0/24, ... 10.120.255.0/24
            var privateSubnetCidrs = subnetting.Skip(128).ToList();

            var nonLocalAvailabilityZones = Aws.GetAvailabilityZones.Invoke(new Aws.GetAvailabilityZonesInvokeArgs()
            {
                Filters = new InputList<GetAvailabilityZonesFilterInputArgs>
                {
                    new GetAvailabilityZonesFilterInputArgs
                    {
                        Name = "opt-in-status",
                        Values = new List<string>() { "opt-in-not-required" }
                    }
                }
            });

            var azNames = nonLocalAvailabilityZones.Apply(result => result.Names);

            var myVpc = new Aws.Ec2.Vpc("my-vpc", new Aws.Ec2.VpcArgs
            {
                CidrBlock = "10.120.0.0/16",
                Tags = InputMap<string>.Merge(defaultTags, new InputMap<string>
                {
                    { "Name", "my-vpc" }
                })
            });

            var igw = new Aws.Ec2.InternetGateway("my-igw", new Aws.Ec2.InternetGatewayArgs
            {
                VpcId = myVpc.Id,
                Tags = InputMap<string>.Merge(defaultTags, new InputMap<string>
                {
                    { "Name", "my-igw" }
                })
            });

            var defaultRt = new Aws.Ec2.DefaultRouteTable("my-default-rt", new Aws.Ec2.DefaultRouteTableArgs
            {
                DefaultRouteTableId = myVpc.DefaultRouteTableId,
                Tags = InputMap<string>.Merge(defaultTags, new InputMap<string>
                {
                    { "Name", "my-default-rt" }
                })
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
                        CidrBlock = publicSubnetCidrs[0],
                        AvailabilityZone = azNames.GetAt(0),
                        Tags = InputMap<string>.Merge(defaultTags, new InputMap<string>()
                        {
                            { "Name", "my-public-subnet-1" }
                        })
                    })
                },
                {
                    "my-public-subnet-2", new Aws.Ec2.Subnet("my-public-subnet-2", new Aws.Ec2.SubnetArgs
                    {
                        VpcId = myVpc.Id,
                        CidrBlock = publicSubnetCidrs[1],
                        AvailabilityZone = azNames.GetAt(1),
                        Tags = InputMap<string>.Merge(defaultTags, new InputMap<string>
                        {
                            { "Name", "my-public-subnet-2" }
                        })
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

            var eip = new Aws.Ec2.Eip("my-nat-gateway-eip", new EipArgs
            {
                Tags = InputMap<string>.Merge(defaultTags, new InputMap<string>
                {
                    { "Name", "my-nat-gateway-eip" }
                })
            });

            var natGateway = new Aws.Ec2.NatGateway("my-nat-gateway", new Aws.Ec2.NatGatewayArgs
            {
                SubnetId = publicSubnet.First().Value.Id,
                AllocationId = eip.Id,
                Tags = InputMap<string>.Merge(defaultTags, new InputMap<string>
                {
                    { "Name", "my-nat-gateway" }
                })
            });

            var privateSubnetArgs = new[]
            {
                new { Cidr = privateSubnetCidrs[0], Az = azNames.GetAt(0) },
                new { Cidr = privateSubnetCidrs[1], Az = azNames.GetAt(1) }
            };

            var privateSubnets = new Dictionary<string, Aws.Ec2.Subnet>();

            for (var i = 0; i < privateSubnetArgs.Length; i++)
            {
                var myPrivateSubnetName = "my-private-subnet-" + (i + 1);
                var args = privateSubnetArgs[i];

                var createdSubnet = new PrivateSubnet(myPrivateSubnetName, new PrivateSubnetArgs()
                {
                    VpcId = myVpc.Id,
                    Cidr = args.Cidr,
                    Az = args.Az,
                    Tags = defaultTags
                });
                createdSubnet.AddNatGateway(natGateway.Id);
                privateSubnets[myPrivateSubnetName] = createdSubnet.Subnet;
            }
        }
    }
}