import pulumi
import pulumi_aws as aws

# 建立 VPC
vpc = aws.ec2.Vpc('my-vpc',
                  cidr_block='10.120.0.0/16',
                  enable_dns_hostnames=True,
                  opts=pulumi.ResourceOptions(
                      protect=True,
                  ))

# 建立 IGW
igw = aws.ec2.InternetGateway('my-igw',
                              vpc_id=vpc.id)

# 取得 Default Route Table
default_rt = aws.ec2.DefaultRouteTable('my-default-rt',
                                       default_route_table_id=vpc.default_route_table_id)

# 設定 Default Route Table 的預設路由到 IGW
aws.ec2.Route('my-default-rt-default-route',
              route_table_id=default_rt.id,
              destination_cidr_block='0.0.0.0/0',
              gateway_id=igw.id)

# 建立 Public Subnets
public_subnet = {
    'my-public-subnet-1': aws.ec2.Subnet('my-public-subnet-1',
                                         vpc_id=vpc.id,
                                         cidr_block='10.120.0.0/24',
                                         availability_zone='ap-east-1a'),
    'my-public-subnet-2': aws.ec2.Subnet('my-public-subnet-2',
                                         vpc_id=vpc.id,
                                         cidr_block='10.120.1.0/24',
                                         availability_zone='ap-east-1b')
}

# 使用迴圈將 Public Subnets 與 Default Route Table 進行關聯
for name, subnet in public_subnet.items():
    aws.ec2.RouteTableAssociation(f'my-{name}-rt-assoc',
                                  subnet_id=subnet.id,
                                  route_table_id=default_rt.id)

# 建立 NAT Gateway 用的 Elastic IP
eip = aws.ec2.Eip('my-nat-gateway-eip')

# 建立 NAT Gateway
nat_gateway = aws.ec2.NatGateway('my-nat-gateway',
                                 allocation_id=eip.id,
                                 subnet_id=public_subnet[next(iter(public_subnet))].id)

# Private Subnet 的參數
privateSubnetArgs = [
    {'cidr': '10.120.128.0/24', 'az': 'ap-east-1a'},
    {'cidr': '10.120.129.0/24', 'az': 'ap-east-1b'}
]

# 建立 Private Subnet
privateSubnets = {}
for arg in privateSubnetArgs:
    privateSubnets[f'my-private-subnet-{arg["az"]}'] = subnet = aws.ec2.Subnet(f'my-private-subnet-{arg["az"]}',
                                                                               vpc_id=vpc.id,
                                                                               cidr_block=arg['cidr'],
                                                                               availability_zone=arg['az'])
    rt = aws.ec2.RouteTable(f'my-private-subnet-{arg["az"]}-rt',
                            vpc_id=vpc.id,
                            routes=[aws.ec2.RouteTableRouteArgs(
                                cidr_block='0.0.0.0/0',
                                nat_gateway_id=nat_gateway.id)
                            ])
    aws.ec2.RouteTableAssociation(f'my-private-subnet-{arg["az"]}-rt-association',
                                  subnet_id=subnet.id,
                                  route_table_id=rt.id)
