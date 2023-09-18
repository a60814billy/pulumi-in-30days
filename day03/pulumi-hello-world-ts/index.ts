import {local} from '@pulumi/command';

const helloWorld = new local.Command('hello,world', {
  create: "echo 'hello,world'"
});

export const helloWorldOutput = helloWorld.stdout;
