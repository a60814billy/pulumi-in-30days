import pulumi
from pulumi_command import local

helloWorld = local.Command("hello-world",
                           create="echo 'Hello, World'")

pulumi.export('helloWorldOutput', helloWorld.stdout)
