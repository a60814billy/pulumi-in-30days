using System.Threading.Tasks;
using Pulumi;

namespace day07_input_output_01_cs;

using Aws = Pulumi.Aws;

public class Example2
{
    public static void Example2_()
    {
        // 建立 Zone
        var zone = new Aws.Route53.Zone("example", new Aws.Route53.ZoneArgs
        {
            Name = "example.com."
        });

        var www = new Aws.Route53.Record("example", new Aws.Route53.RecordArgs
        {
            ZoneId = zone.ZoneId,
            // 傳遞 async Task<string> 至 Input<string> 中
            // 因為 Output 是其中一種 Input 可以接受的型別，因此須先將 async task 轉成 Output
            Name = Output.Create(GetNameAsync()),
            Type = "TXT",
            Ttl = 300,
            Records = new[]
            {
                "Hello, World!"
            }
        });
    }

    private static async Task<string> GetNameAsync()
    {
        await Task.Delay(1);
        return "example.com.";
    }
}