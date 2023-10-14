import pulumi_aws as aws


async def get_name():
    return "example.com."


# 建立 Zone
zone = aws.route53.Zone("example", name="example.com.")

www = aws.route53.Record("example",
                         zone_id=zone.zone_id,
                         # 傳遞 Awaitable 至 Input[str] 中
                         name=get_name(),
                         type="TXT",
                         records=["Hello, World"])
