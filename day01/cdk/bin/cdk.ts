#!/usr/bin/env node

import * as cdk from 'aws-cdk-lib';
import * as s3 from 'aws-cdk-lib/aws-s3';
import {Construct} from 'constructs';

export class S3BucketStack extends cdk.Stack {
  constructor(scope: Construct, id: string, props?: cdk.StackProps) {
    super(scope, id, props);
    if (process.env['CREATE_BUCKET']) {
      const s3Bucket = new s3.Bucket(this, 'ExampleBucket', {});
    }
  }
}

const app = new cdk.App();
new S3BucketStack(app, 'CdkStack', {});
