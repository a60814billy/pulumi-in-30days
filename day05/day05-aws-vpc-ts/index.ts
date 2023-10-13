import * as aws from "@pulumi/aws";

const myVpc = new aws.ec2.Vpc("my-vpc", {
  cidrBlock: "10.120.0.0/16"
});

const igw = new aws.ec2.InternetGateway("my-igw", {
  vpcId: myVpc.id
});

const defaultRT = new aws.ec2.DefaultRouteTable("my-default-rt", {
  defaultRouteTableId: myVpc.defaultRouteTableId,
});

new aws.ec2.Route("my-default-rt-default-route", {
  routeTableId: defaultRT.id,
  destinationCidrBlock: "0.0.0.0/0",
  gatewayId: igw.id,
});

const publicSubnets: Record<string, aws.ec2.Subnet> = {
  'my-public-subnet-1': new aws.ec2.Subnet("my-public-subnet-1", {
    vpcId: myVpc.id,
    cidrBlock: '10.120.0.0/24',
    availabilityZone: 'ap-east-1a',
  }),
  'my-public-subnet-2': new aws.ec2.Subnet("my-public-subnet-2", {
    vpcId: myVpc.id,
    cidrBlock: '10.120.1.0/24',
    availabilityZone: 'ap-east-1b',
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


const eip = new aws.ec2.Eip("my-nat-gateway-eip", {});

const natGateway = new aws.ec2.NatGateway("my-nat-gateway", {
  subnetId: Object.values(publicSubnets)[0].id,
  allocationId: eip.id,
});

const privateSubnetArgs = [
  {cidr: '10.120.128.0/24', az: 'ap-east-1a'},
  {cidr: '10.120.129.0/24', az: 'ap-east-1b'}
];

const privateSubnets: Record<string, aws.ec2.Subnet> = {};
for (const arg of privateSubnetArgs) {
  const subnet = new aws.ec2.Subnet(`my-private-subnet-${arg.az}`, {
    vpcId: myVpc.id,
    cidrBlock: arg.cidr,
    availabilityZone: arg.az,
  });
  privateSubnets[`my-private-subnet-${arg.az}`] = subnet;

  const rt = new aws.ec2.RouteTable(`my-private-subnet-${arg.az}-rt`, {
    vpcId: myVpc.id,
    routes: [{cidrBlock: '0.0.0.0/0', natGatewayId: natGateway.id}]
  });

  new aws.ec2.RouteTableAssociation(`my-private-subnet-${arg.az}-rt-association`, {
    routeTableId: rt.id,
    subnetId: subnet.id,
  });
}

