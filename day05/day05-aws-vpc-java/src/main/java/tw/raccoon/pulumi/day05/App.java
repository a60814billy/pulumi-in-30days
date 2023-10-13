package tw.raccoon.pulumi.day05;

import com.pulumi.*;
import com.pulumi.aws.ec2.*;
import com.pulumi.aws.ec2.inputs.RouteTableRouteArgs;

import java.util.ArrayList;
import java.util.HashMap;

public class App {

    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    static class PrivateSubnetArgs {
        public String CidrBlock;
        public String Az;

        public PrivateSubnetArgs(String cidrBlock, String az) {
            CidrBlock = cidrBlock;
            Az = az;
        }
    }

    private static void stack(Context ctx) {
        var myVpc = new Vpc("my-vpc", VpcArgs.builder().cidrBlock("10.120.0.0/16").build());
        var igw = new InternetGateway("my-igw", InternetGatewayArgs.builder().vpcId(myVpc.id()).build());
        var defaultRt = new DefaultRouteTable("my-default-rt",
                DefaultRouteTableArgs.builder()
                        .defaultRouteTableId(myVpc.defaultRouteTableId())
                        .build());

        new Route("my-default-rt-default-route",
                RouteArgs.builder()
                        .routeTableId(defaultRt.id())
                        .destinationCidrBlock("0.0.0.0/0")
                        .gatewayId(igw.id())
                        .build());

        var publicSubnets = new HashMap<String, Subnet>();
        publicSubnets.put("my-public-subnet-1", new Subnet("my-public-subnet-1",
                SubnetArgs.builder()
                        .vpcId(myVpc.id())
                        .cidrBlock("10.120.0.0/24")
                        .availabilityZone("ap-east-1a")
                        .build()
        ));
        publicSubnets.put("my-public-subnet-2", new Subnet("my-public-subnet-2",
                SubnetArgs.builder()
                        .vpcId(myVpc.id())
                        .cidrBlock("10.120.1.0/24")
                        .availabilityZone("ap-east-1b")
                        .build()
        ));


        publicSubnets.forEach((key, value) -> new RouteTableAssociation(key + "-rt-association",
                RouteTableAssociationArgs.builder()
                        .routeTableId(defaultRt.id())
                        .subnetId(value.id())
                        .build()));


        var eip = new Eip("my-nat-gateway-eip");

        String firstSubnetName = new ArrayList<String>(publicSubnets.keySet()).get(0);
        var natGateway = new NatGateway("my-nat-gateway", NatGatewayArgs.builder()
                .subnetId(publicSubnets.get(firstSubnetName).id())
                .allocationId(eip.id())
                .build());


        var privateSubnetArgs = new ArrayList<PrivateSubnetArgs>();
        privateSubnetArgs.add(new PrivateSubnetArgs("10.120.128.0/24", "ap-east-1a"));
        privateSubnetArgs.add(new PrivateSubnetArgs("10.120.129.0/24", "ap-east-1b"));

        var privateSubnets = new HashMap<String, Subnet>();

        for (var arg : privateSubnetArgs) {
            var privateSubnetName = "my-private-subnet-" + arg.Az;
            var subnet = new Subnet(privateSubnetName,
                    SubnetArgs.builder()
                            .vpcId(myVpc.id())
                            .cidrBlock(arg.CidrBlock)
                            .availabilityZone(arg.Az)
                            .build()
            );
            privateSubnets.put(privateSubnetName, subnet);

            var rt = new RouteTable(privateSubnetName + "-rt",
                    RouteTableArgs.builder()
                            .vpcId(myVpc.id())
                            .routes(RouteTableRouteArgs.builder()
                                    .cidrBlock("0.0.0.0/0")
                                    .natGatewayId(natGateway.id())
                                    .build())
                            .build()
            );

            new RouteTableAssociation(privateSubnetName + "-rt-association",
                    RouteTableAssociationArgs.builder()
                            .subnetId(subnet.id())
                            .routeTableId(rt.id())
                            .build());
        }
    }
}
