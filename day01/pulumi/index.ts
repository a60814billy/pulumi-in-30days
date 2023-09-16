import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

if (process.env['CREATE_BUCKET']) {
  const bucket = new aws.s3.Bucket("example-bucket");
}
