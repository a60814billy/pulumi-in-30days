using System.Threading.Tasks;
using Pulumi;

namespace day05_aws_vpc_cs
{
    class Program
    {
        static Task<int> Main()
        {
            return Deployment.RunAsync<AwsVpcStack>();
        }
    }
}