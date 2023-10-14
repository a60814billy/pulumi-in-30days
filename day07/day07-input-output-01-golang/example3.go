package example1

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/ec2"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func example1() {
	pulumi.Run(func(ctx *pulumi.Context) error {

		eip, err := ec2.NewEip(ctx, 'my-eip', &ec2.EipArgs{})
		if err != nil {
			return err
		}

		// 我們從某個 Resource 拿到了 ip address，假設是 200.100.10.20 好了
		var myIp pulumi.StringOutput = eip.PublicIp

		// 我們想要將 ip 字串與其他字串做一些處理，例如取得該 ip 的網段，並產生 CIDR 格式的網段描述，類似這樣： 200.100.10.0/24
		// 產生出來後，我想要將他做為其他資源的 Input

		return nil
	})
}
