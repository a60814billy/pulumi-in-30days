package tw.raccoon.pulumi.hello.world;

import com.pulumi.*;
import com.pulumi.command.local.*;

public class App {

    public static void main(String[] args) {
        Pulumi.run(App::stack);
    }

    private static void stack(Context ctx) {
        var commandArgs = CommandArgs.builder()
                .create("echo 'Hello World'")
                .build();

        var command = new Command("hello-world", commandArgs);

        ctx.export("helloWorldOutput", command.stdout());
    }
}
