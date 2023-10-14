using System.Threading.Tasks;
using Pulumi;

namespace day07_input_output_01_cs;

using Aws = Pulumi.Aws;

public class Example3
{
    public static void Example3_()
    {
        var eip = new Aws.Ec2.Eip("my-eip");

        // 我們從某個 Resource 拿到了 ip address，假設是 200.100.10.20 好了
        Output<string> myIp = eip.PublicIp;

        // 我們想要將 ip 字串與其他字串做一些處理，例如取得該 ip 的網段，並產生 CIDR 格式的網段描述，類似這樣： 200.100.10.0/24
        // 產生出來後，我想要將他做為其他資源的 Input
    }
}