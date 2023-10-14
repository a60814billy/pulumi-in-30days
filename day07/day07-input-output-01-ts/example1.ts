import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

// 建立 Zone
const zone = new aws.route53.Zone("example", {
    name: "example.com.",
});

const www = new aws.route53.Record("example", {
    // 將 zone 的 output 屬性當作 record 的 input 屬性
    zoneId: zone.zoneId,
    name: "www.example.com.",
    type: "TXT",
    records: [
        "Hello, World"
    ]
});
