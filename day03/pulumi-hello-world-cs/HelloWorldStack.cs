
using Pulumi;
using Pulumi.Command.Local;

class HelloWorldStack: Stack {

    [Output]
    public Output<string> HelloWorldOutput {get;set;}

    public HelloWorldStack () {
        var args = new CommandArgs();
        args.Create = "echo hello, world";

        var command = new Command("hello, world", args);

        this.HelloWorldOutput = command.Stdout;
    }
}