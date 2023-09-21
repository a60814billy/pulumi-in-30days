import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

import {splitSubnets} from './cidr_subnet_utils';
import {PrivateSubnet} from "./privateSubnet";

export = async function () {
  const defaultTags = {
    "pulumi:project": pulumi.getProject(),
    "pulumi:stack": pulumi.getStack(),
    "ManagedBy": "Pulumi",
  };

  // 宣告一個變數儲存 vpcCidr
  const vpcCidr = '10.120.0.0/16';
  // 使用 splitSubnets 函式，將 10.120.0.0/16 分割成 256 個 /24 的子網路
  const subnetting = splitSubnets(vpcCidr, 24);
  // 前 128 個子網路為 public 子網路
  // 10.120.0.0/24, 10.120.1.0/24 .... 10.120.127.0/24
  const publicSubnetCidrs = subnetting.slice(0, subnetting.length / 2);
  // 後 128 個子網路為 private 子網路
  // 10.120.128.0/24, 10.120.129.0/24, ... 10.120.255.0/24
  const privateSubnetCidrs = subnetting.slice(subnetting.length / 2);

  const nonLocalAvailabilityZones = await aws.getAvailabilityZones({
    filters: [
      {
        name: "opt-in-status",
        values: ["opt-in-not-required"],
      },
    ],
  });
  const azNames = nonLocalAvailabilityZones.names;


  const myVpc = new aws.ec2.Vpc("my-vpc", {
    cidrBlock: vpcCidr,
    tags: {
      "Name": "my-vpc",
      ...defaultTags
    }
  });

  const igw = new aws.ec2.InternetGateway("my-igw", {
    vpcId: myVpc.id,
    tags: {
      "Name": "my-igw",
      ...defaultTags
    }
  });

  const defaultRT = new aws.ec2.DefaultRouteTable("my-default-rt", {
    defaultRouteTableId: myVpc.defaultRouteTableId,
    tags: {
      "Name": "my-default-rt",
      ...defaultTags
    }
  });

  new aws.ec2.Route("my-default-rt-default-route", {
    routeTableId: defaultRT.id,
    destinationCidrBlock: "0.0.0.0/0",
    gatewayId: igw.id,
  });

  const publicSubnets: Record<string, aws.ec2.Subnet> = {
    'my-public-subnet-1': new aws.ec2.Subnet("my-public-subnet-1", {
      vpcId: myVpc.id,
      cidrBlock: publicSubnetCidrs[0],
      availabilityZone: azNames[0],
      tags: {
        "Name": "my-public-subnet-1",
        ...defaultTags
      }
    }),
    'my-public-subnet-2': new aws.ec2.Subnet("my-public-subnet-2", {
      vpcId: myVpc.id,
      cidrBlock: publicSubnetCidrs[1],
      availabilityZone: azNames[1],
      tags: {
        "Name": "my-public-subnet-2",
        ...defaultTags
      }
    })
  };

  // for (const subnetName of Object.keys(publicSubnets)) {
  //   new aws.ec2.RouteTableAssociation(`${subnetName}-rt-association`, {
  //     routeTableId: defaultRT.defaultRouteTableId,
  //     subnetId: publicSubnets[subnetName].id,
  //   });
  // }

  Object.keys(publicSubnets).forEach(subnetName => {
    new aws.ec2.RouteTableAssociation(`${subnetName}-rt-association`, {
      routeTableId: defaultRT.id,
      subnetId: publicSubnets[subnetName].id,
    });
  });


  const eip = new aws.ec2.Eip("my-nat-gateway-eip", {
    tags: {
      "Name": "my-nat-gateway-eip",
      ...defaultTags
    }
  });

  const natGateway = new aws.ec2.NatGateway("my-nat-gateway", {
    subnetId: Object.values(publicSubnets)[0].id,
    allocationId: eip.id,
    tags: {
      "Name": "my-nat-gateway",
      ...defaultTags
    }
  });

  const privateSubnetArgs = [
    {cidr: privateSubnetCidrs[0], az: azNames[0]},
    {cidr: privateSubnetCidrs[1], az: azNames[1]}
  ];

  const privateSubnets: Record<string, aws.ec2.Subnet> = {};
  for (const arg of privateSubnetArgs) {
    const createdSubnet = new PrivateSubnet({
      az: arg.az,
      cidr: arg.cidr,
      vpcId: myVpc.id,
      defaultTags: defaultTags,
      natGatewayId: natGateway.id,
    });
    createdSubnet.addNatGateway(natGateway.id);
    privateSubnets[createdSubnet.name] = createdSubnet.subnet;
  }
}
