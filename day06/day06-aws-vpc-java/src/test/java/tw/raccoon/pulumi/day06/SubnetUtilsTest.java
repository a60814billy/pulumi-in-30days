package tw.raccoon.pulumi.day06;


import org.junit.Assert;
import org.junit.Test;

import java.util.List;


public class SubnetUtilsTest {

    @Test
    public void test1() {
        List<String> subnets = null;
        subnets = SubnetUtils.splitSubnets("10.120.0.0/16", 22);
        Assert.assertEquals(64, subnets.size());
        Assert.assertEquals("10.120.0.0/22", subnets.get(0));
        Assert.assertEquals("10.120.4.0/22", subnets.get(1));
    }

}
