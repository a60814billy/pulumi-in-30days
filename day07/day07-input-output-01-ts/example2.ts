import * as pulumi from "@pulumi/pulumi";
import * as aws from "@pulumi/aws";

// 建立 Zone
const zone = new aws.route53.Zone("example", {
    name: "example.com.",
});

const www = new aws.route53.Record("example", {
    zoneId: zone.zoneId,
    // 傳遞 Promise<string> 至 Input<string> 中
    name: Promise.resolve("www.example.com."),
    type: "TXT",
    records: [
        "Hello, World"
    ]
});
