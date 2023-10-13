package tw.raccoon.pulumi.day06;

import com.pulumi.*;
import com.pulumi.aws.ec2.*;
import com.pulumi.aws.*;
import com.pulumi.aws.inputs.GetAvailabilityZonesArgs;
import com.pulumi.aws.inputs.GetAvailabilityZonesFilterArgs;
import com.pulumi.aws.outputs.GetAvailabilityZonesResult;
import com.pulumi.core.Output;

import java.util.ArrayList;
import java.util.HashMap;

public class App {

    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    static class PrivateSubnetData {
        public String CidrBlock;
        public Output<String> Az;

        public PrivateSubnetData(String cidrBlock, Output<String> az) {
            CidrBlock = cidrBlock;
            Az = az;
        }
    }

    private static HashMap<String, String> getTagWithName(Context ctx, String name) {
        var tags = new HashMap<String, String>();
        tags.put("pulumi:project", ctx.projectName());
        tags.put("pulumi:stack", ctx.stackName());
        tags.put("ManagedBy", "Pulumi");
        tags.put("Name", name);
        return tags;
    }

    private static void stack(Context ctx) {
        var vpcCidr = "10.120.0.0/16";
        var splittedSubnets = SubnetUtils.splitSubnets(vpcCidr, 24);
        String[] subnets = new String[splittedSubnets.size()];
        splittedSubnets.toArray(subnets);
        // 前 128 個子網路為 public 子網路
        // 10.120.0.0/24, 10.120.1.0/24 .... 10.120.127.0/24
        var halfSize = subnets.length / 2;
        String[] publicSubnetCidrs = new String[halfSize];
        System.arraycopy(subnets, 0, publicSubnetCidrs, 0, 128);
        // 後 128 個子網路為 private 子網路
        // 10.120.128.0/24, 10.120.129.0/24, ... 10.120.255.0/24
        String[] privateSubnetCidrs = new String[halfSize];
        System.arraycopy(subnets, halfSize, privateSubnetCidrs, 0, 128);

        Output<GetAvailabilityZonesResult> availabilityZones = AwsFunctions.getAvailabilityZones(
                GetAvailabilityZonesArgs.builder()
                        .state("available")
                        .filters(GetAvailabilityZonesFilterArgs.builder()
                                .name("opt-in-status")
                                .values("opt-in-not-required")
                                .build())
                        .build()
        );


        var myVpc = new Vpc("my-vpc", VpcArgs.builder()
                .cidrBlock(vpcCidr)
                .tags(getTagWithName(ctx, "my-vpc"))
                .build());

        var igw = new InternetGateway("my-igw", InternetGatewayArgs.builder()
                .vpcId(myVpc.id())
                .tags(getTagWithName(ctx, "my-igw"))
                .build());
        var defaultRt = new DefaultRouteTable("my-default-rt",
                DefaultRouteTableArgs.builder()
                        .defaultRouteTableId(myVpc.defaultRouteTableId())
                        .tags(getTagWithName(ctx, "my-default-rt"))
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
                        .cidrBlock(publicSubnetCidrs[0])
                        .availabilityZone(availabilityZones.applyValue(azs -> azs.names().get(0)))
                        .tags(getTagWithName(ctx, "my-public-subnet-1"))
                        .build()
        ));
        publicSubnets.put("my-public-subnet-2", new Subnet("my-public-subnet-2",
                SubnetArgs.builder()
                        .vpcId(myVpc.id())
                        .cidrBlock(publicSubnetCidrs[1])
                        .availabilityZone(availabilityZones.applyValue(azs -> azs.names().get(1)))
                        .tags(getTagWithName(ctx, "my-public-subnet-2"))
                        .build()
        ));


        publicSubnets.forEach((key, value) -> new RouteTableAssociation(key + "-rt-association",
                RouteTableAssociationArgs.builder()
                        .routeTableId(defaultRt.id())
                        .subnetId(value.id())
                        .build()));


        var eip = new Eip("my-nat-gateway-eip",
                EipArgs.builder()
                        .tags(getTagWithName(ctx, "my-nat-gateway-eip"))
                        .build()
        );

        String firstSubnetName = new ArrayList<>(publicSubnets.keySet()).get(0);
        var natGateway = new NatGateway("my-nat-gateway", NatGatewayArgs.builder()
                .subnetId(publicSubnets.get(firstSubnetName).id())
                .allocationId(eip.id())
                .tags(getTagWithName(ctx, "my-nat-gateway"))
                .build()
        );


        var privateSubnetArgs = new ArrayList<PrivateSubnetData>();
        privateSubnetArgs.add(new PrivateSubnetData(privateSubnetCidrs[0], availabilityZones.applyValue(azs -> azs.names().get(0))));
        privateSubnetArgs.add(new PrivateSubnetData(privateSubnetCidrs[1], availabilityZones.applyValue(azs -> azs.names().get(1))));

        var privateSubnets = new HashMap<String, Subnet>();

        for (int i = 0; i < privateSubnetArgs.size(); i++) {
            var arg = privateSubnetArgs.get(i);
            var privateSubnetName = "my-private-subnet-" + (i + 1);

            var createdSubnet = new PrivateSubnet(privateSubnetName, PrivateSubnetArgs.builder()
                    .vpcId(myVpc.id())
                    .cidrBlock(arg.CidrBlock)
                    .az(arg.Az)
                    .tags(getTagWithName(ctx, privateSubnetName))
                    .build()
            );
            createdSubnet.addNatGateway(natGateway.id());

            privateSubnets.put(privateSubnetName, createdSubnet.subnet);
        }
    }
}
