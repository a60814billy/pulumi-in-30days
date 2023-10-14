package tw.raccoon.pulumi.day07.example;

import com.pulumi.Pulumi;
import com.pulumi.aws.route53.Record;
import com.pulumi.aws.route53.RecordArgs;
import com.pulumi.aws.route53.Zone;
import com.pulumi.aws.route53.ZoneArgs;
import com.pulumi.core.Output;

import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutorService;
import java.util.concurrent.Executors;

public class Example2 {
    private static final ExecutorService executor = Executors.newCachedThreadPool();

    public static void example2(String[] args) {
        Pulumi.run(ctx -> {
            var zone = new Zone("example-com",
                    ZoneArgs.builder()
                            .name("example.com.")
                            .build());

            var www = new Record("example",
                    RecordArgs.builder()
                            .zoneId(zone.id())
                            // 傳遞 CompletableFuture<String> 至 Output<String> 中
                            .name(Output.of(getRecordName()))
                            .type("TXT")
                            .build());
        });
    }

    static CompletableFuture<String> getRecordName() {
        var future = new CompletableFuture<String>();

        executor.submit(() -> {
            try {
                future.complete("www.example.com.");
            } catch (Exception e) {
                future.completeExceptionally(e);
            }
        });

        return future;
    }
}
