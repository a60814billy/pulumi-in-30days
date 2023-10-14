import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

const eip = new aws.ec2.Eip('my-eip');

// 我們從某個 Resource 拿到了 ip address，假設是 200.100.10.20 好了
const myIp: pulumi.Output<string> = eip.publicIp;

// 我們想要將 ip 字串與其他字串做一些處理，例如取得該 ip 的網段，並產生 CIDR 格式的網段描述，類似這樣： 200.100.10.0/24
// 產生出來後，我想要將他做為其他資源的 Input

