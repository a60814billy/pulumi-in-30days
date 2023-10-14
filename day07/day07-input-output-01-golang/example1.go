package example1

import (
	"github.com/pulumi/pulumi-aws/sdk/v6/go/aws/route53"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func example1() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		// 建立 Zone
		zone, err := route53.NewZone(ctx, "example", &route53.ZoneArgs{
			Name: pulumi.String("example.com."),
		})
		if err != nil {
			return err
		}

		www, err := route53.NewRecord(ctx, "example", &route53.RecordArgs{
			// 將 zone 的 output 屬性當作 record 的 input 屬性
			ZoneId: zone.ZoneId,
			Name:   pulumi.String("www.example.com."),
			Type:   pulumi.String("TXT"),
			Records: pulumi.StringArray{
				pulumi.String("Hello, World"),
			},
		})
		return nil
	})
}
