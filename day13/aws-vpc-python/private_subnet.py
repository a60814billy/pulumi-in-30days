import pulumi
import pulumi_aws as aws


class PrivateSubnet(pulumi.ComponentResource):
    def __init__(self, name, opts: pulumi.ResourceOptions = None, **kwargs):
        super().__init__('aws-vpc-python:private_subnet:PrivateSubnet', name, kwargs, opts)
        az = kwargs['az']
        cidr = kwargs['cidr_block']
        vpc_id = kwargs['vpc_id']

        self.name = name

        self.subnet = aws.ec2.Subnet(self.name,
                                     vpc_id=vpc_id,
                                     cidr_block=cidr,
                                     availability_zone=az,
                                     opts=pulumi.ResourceOptions(parent=self))

        self.rt = aws.ec2.RouteTable(f'{self.name}-rt',
                                     vpc_id=vpc_id,
                                     opts=pulumi.ResourceOptions(parent=self))
        self.default_route = None

        aws.ec2.RouteTableAssociation(f'{self.name}-rt-association',
                                      subnet_id=self.subnet.id,
                                      route_table_id=self.rt.id,
                                      opts=pulumi.ResourceOptions(parent=self))

    def add_nat_gateway(self, nat_gateway_id):
        if self.default_route is not None:
            raise Exception('Default route already exists')

        self.default_route = aws.ec2.Route(f'{self.name}-default-route',
                                           route_table_id=self.rt.id,
                                           destination_cidr_block='0.0.0.0/0',
                                           nat_gateway_id=nat_gateway_id,
                                           opts=pulumi.ResourceOptions(parent=self))
