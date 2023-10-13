import pulumi_aws as aws

from cidr_subnet_utils import split_subnets
from private_subnet import PrivateSubnet

vpc_cidr = '10.120.0.0/16'
subnetting = split_subnets(vpc_cidr, 24)
public_subnets = subnetting[:(len(subnetting) // 2)]
private_subnets = subnetting[(len(subnetting) // 2):]

non_local_availability_zones = aws.get_availability_zones(filters=[
    aws.GetAvailabilityZonesFilterArgs(name='opt-in-status', values=['opt-in-not-required'])
])
az_names = non_local_availability_zones.names

# 建立 VPC
vpc = aws.ec2.Vpc('my-vpc',
                  cidr_block='10.120.0.0/16',
                  enable_dns_hostnames=True)

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
                                         cidr_block=public_subnets[0],
                                         availability_zone=az_names[0]),
    'my-public-subnet-2': aws.ec2.Subnet('my-public-subnet-2',
                                         vpc_id=vpc.id,
                                         cidr_block=public_subnets[1],
                                         availability_zone=az_names[1])
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
    {'cidr': private_subnets[0], 'az': az_names[0]},
    {'cidr': private_subnets[1], 'az': az_names[1]}
]

# 建立 Private Subnet
privateSubnets = {}
for arg in privateSubnetArgs:
    subnet = PrivateSubnet(f'my-private-subnet-{arg["az"]}', vpc_id=vpc.id, cidr_block=arg['cidr'], az=arg['az'])
    subnet.add_nat_gateway(nat_gateway.id)
    privateSubnets[subnet.name] = subnet
