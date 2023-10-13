package tw.raccoon.pulumi.day06;

import com.pulumi.core.Output;

import java.util.Map;

public class PrivateSubnetArgs {
    public Output<String> VpcId;
    public String CidrBlock;
    public Output<String> Az;
    public Map<String, String> Tags;

    public PrivateSubnetArgs() {
    }

    public static Builder builder() {
        return new Builder();
    }

    public static class Builder {
        private Output<String> vpcId;
        private String cidrBlock;
        private Output<String> az;
        private Map<String, String> tags;

        public Builder vpcId(Output<String> vpcId) {
            this.vpcId = vpcId;
            return this;
        }

        public Builder cidrBlock(String cidrBlock) {
            this.cidrBlock = cidrBlock;
            return this;
        }

        public Builder az(Output<String> az) {
            this.az = az;
            return this;
        }

        public Builder tags(Map<String, String> tags) {
            this.tags = tags;
            return this;
        }

        public tw.raccoon.pulumi.day06.PrivateSubnetArgs build() {
            tw.raccoon.pulumi.day06.PrivateSubnetArgs args = new tw.raccoon.pulumi.day06.PrivateSubnetArgs();
            args.VpcId = this.vpcId;
            args.CidrBlock = this.cidrBlock;
            args.Az = this.az;
            args.Tags = this.tags;
            return args;
        }
    }

}
