import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

interface CreatePrivateSubnetArgs {
  vpcId: pulumi.Input<string>;
  cidr: string;
  az: string;
  defaultTags: Record<string, string>;
}

export class PrivateSubnet extends pulumi.ComponentResource {
  public name: string;
  public subnet: aws.ec2.Subnet;

  private routeTable: aws.ec2.RouteTable;

  private defaultRoute?: aws.ec2.Route;

  constructor(name: string, args: CreatePrivateSubnetArgs, opts: pulumi.ComponentResourceOptions = {}) {
    super("aws-vpc-ts:privateSubnet:PrivateSubnet", name, args, opts);
    this.name = name;

    this.subnet = new aws.ec2.Subnet(this.name, {
      vpcId: args.vpcId,
      cidrBlock: args.cidr,
      availabilityZone: args.az,
      tags: {
        "Name": this.name,
        ...args.defaultTags
      }
    }, {parent: this});

    this.routeTable = new aws.ec2.RouteTable(`${this.name}-rt`, {
      vpcId: args.vpcId,
      tags: {
        'Name': `${this.name}-rt`,
        ...args.defaultTags
      }
    }, {parent: this});

    new aws.ec2.RouteTableAssociation(`${this.name}-rt-association`, {
      routeTableId: this.routeTable.id,
      subnetId: this.subnet.id,
    }, {parent: this});
  }

  addNatGateway(natGatewayId: pulumi.Input<string>) {
    if (this.defaultRoute) {
      throw new Error('Cannot add a NAT gateway to a subnet that already has a default route');
    }
    this.defaultRoute = new aws.ec2.Route(`${this.name}-nat-gateway-route`, {
      routeTableId: this.routeTable.id,
      destinationCidrBlock: '0.0.0.0/0',
      natGatewayId: natGatewayId,
    }, {parent: this});
  }
}
