namespace day07_input_output_01_cs;

using Aws = Pulumi.Aws;

public class Example1
{
    public static void Example1_()
    {
        // 建立 Zone
        var zone = new Aws.Route53.Zone("example", new Aws.Route53.ZoneArgs
        {
            Name = "example.com."
        });

        var www = new Aws.Route53.Record("example", new Aws.Route53.RecordArgs
        {
            // 將 zone 的 output 屬性當作 record 的 input 屬性
            ZoneId = zone.ZoneId,
            Name = "www.example.com.",
            Type = "TXT",
            Ttl = 300,
            Records = new[]
            {
                "Hello, World!"
            }
        });
    }
}