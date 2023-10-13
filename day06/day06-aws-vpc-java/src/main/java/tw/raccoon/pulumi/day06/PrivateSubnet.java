package tw.raccoon.pulumi.day06;

import com.pulumi.aws.ec2.*;
import com.pulumi.core.Output;

public class PrivateSubnet {
    private final String name;
    public Subnet subnet;
    public RouteTable routeTable;

    public PrivateSubnet(String name, PrivateSubnetArgs args) {
        this.name = name;
        this.subnet = new Subnet(name,
                SubnetArgs.builder()
                        .vpcId(args.VpcId)
                        .cidrBlock(args.CidrBlock)
                        .availabilityZone(args.Az)
                        .tags(args.Tags)
                        .build()
        );

        this.routeTable = new RouteTable(name + "-rt",
                RouteTableArgs.builder()
                        .vpcId(args.VpcId)
                        .build()
        );

        new RouteTableAssociation(name + "-rt-association",
                RouteTableAssociationArgs.builder()
                        .subnetId(this.subnet.id())
                        .routeTableId(this.routeTable.id())
                        .build()
        );
    }

    public void addNatGateway(Output<String> natGatewayId) {
        new Route(this.name + "-nat-gateway-route",
                RouteArgs.builder()
                        .routeTableId(this.routeTable.id())
                        .natGatewayId(natGatewayId)
                        .destinationCidrBlock("0.0.0.0/0")
                        .build()
        );
    }
}
