package main

import (
	"github.com/pulumi/pulumi-command/sdk/go/command/local"
	"github.com/pulumi/pulumi/sdk/v3/go/pulumi"
)

func main() {
	pulumi.Run(func(ctx *pulumi.Context) error {
		command, err := local.NewCommand(ctx, "hello", &local.CommandArgs{
			Create: pulumi.StringPtr("echo 'Hello, World'"),
		})
		if err != nil {
			return err
		}

		ctx.Export("helloWorldOutput", command.Stdout)
		return nil
	})
}
