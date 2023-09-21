import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

interface CreatePrivateSubnetArgs {
  az: string;
  cidr: string;
  vpcId: pulumi.Input<string>;
  natGatewayId?: pulumi.Input<string>;
  defaultTags: Record<string, string>;
}

export class PrivateSubnet {
  public name: string;
  public subnet: aws.ec2.Subnet;

  private routeTable: aws.ec2.RouteTable;

  constructor(args: CreatePrivateSubnetArgs) {
    const subnet = new aws.ec2.Subnet(`my-private-subnet-${args.az}`, {
      vpcId: args.vpcId,
      cidrBlock: args.cidr,
      availabilityZone: args.az,
      tags: {
        "Name": `my-private-subnet-${args.az}`,
        ...args.defaultTags
      }
    });

    const rt = new aws.ec2.RouteTable(`my-private-subnet-${args.az}-rt`, {
      vpcId: args.vpcId,
      tags: {
        'Name': `my-private-subnet-${args.az}-rt`,
        ...args.defaultTags
      }
    });

    new aws.ec2.RouteTableAssociation(`my-private-subnet-${args.az}-rt-association`, {
      routeTableId: rt.id,
      subnetId: subnet.id,
    });

    this.name = `my-private-subnet-${args.az}`;
    this.subnet = subnet;
    this.routeTable = rt;
  }

  addNatGateway(natGatewayId: pulumi.Input<string>) {
    new aws.ec2.Route(`${this.name}-nat-gateway-route`, {
      routeTableId: this.routeTable.id,
      destinationCidrBlock: '0.0.0.0/0',
      natGatewayId: natGatewayId,
    });
  }
}
