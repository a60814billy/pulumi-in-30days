package tw.raccoon.pulumi.day07.example;

import com.pulumi.Pulumi;
import com.pulumi.aws.route53.Record;
import com.pulumi.aws.route53.RecordArgs;
import com.pulumi.aws.route53.Zone;
import com.pulumi.aws.route53.ZoneArgs;

public class Example1 {
    public static void example1(String[] args) {
        Pulumi.run(ctx -> {
            var zone = new Zone("example-com",
                    ZoneArgs.builder()
                            .name("example.com.")
                            .build());

            var www = new Record("example",
                    RecordArgs.builder()
                            // 將 zone 的 output 屬性當作 record 的 input 屬性，
                            // 不過在 Java 中，所有的 Input 都會使用 Output class 標註。
                            // 並不像其他程式語言有特別區分 Input 與 Output
                            .zoneId(zone.id())
                            .name("www.example.com.")
                            .type("TXT")
                            .build());
        });
    }
}
