using day06_aws_vpc_cs;
using Xunit;


namespace UnitTest
{
    public class SubnetUtilsUnitTest
    {
        [Fact]
        public void ShouldSplitSubnetCorrectly()
        {
            var splitSubnets = SubnetUtils.SplitSubnets("10.120.0.0/16", 22);

            Assert.Equal(64, splitSubnets.Count);
            Assert.Equal("10.120.0.0/22", splitSubnets[0]);
            Assert.Equal("10.120.4.0/22", splitSubnets[1]);
        }
    }
}