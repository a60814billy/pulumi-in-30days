import pulumi_aws as aws

# 建立 Zone
zone = aws.route53.Zone("example", name="example.com.")

www = aws.route53.Record("example",
                         # 將 zone 的 output 屬性當作 record 的 input 屬性
                         zone_id=zone.zone_id,
                         name="www.example.com.",
                         type="TXT",
                         records=["Hello, World"])
