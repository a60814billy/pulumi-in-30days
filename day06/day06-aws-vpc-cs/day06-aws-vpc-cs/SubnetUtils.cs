using System;
using System.Collections.Generic;
using System.Net;

namespace day06_aws_vpc_cs
{
    public class SubnetUtils
    {
        private static IPAddress IntToIpAddress(int ip)
        {
            var ipBytes = new byte[4];
            for (var i = 0; i < 4; i++)
            {
                ipBytes[i] = (byte)(ip >> (24 - (8 * i)) & 0xFF);
            }

            return new IPAddress(ipBytes);
        }

        public static List<string> SplitSubnets(string baseCidr, int newPrefixLength)
        {
            var subnets = new List<string>();

            var baseCidrIp = IPAddress.Parse(baseCidr.Split('/')[0]);
            var baseCidrPrefixLength = int.Parse(baseCidr.Split('/')[1]);

            // calculate total ips to int type
            var totalIps = (int)Math.Pow(2, 32 - baseCidrPrefixLength) - 1;


            var addressBytes = baseCidrIp.GetAddressBytes();
            int ip = 0;
            for (int i = 0; i < addressBytes.Length; i++)
            {
                ip += addressBytes[i] << (24 - (8 * i));
            }

            var lastIpAddressInt = ip + totalIps;

            var incrementIp = (int)Math.Pow(2, 32 - newPrefixLength);

            while (ip < lastIpAddressInt)
            {
                subnets.Add(IntToIpAddress(ip) + "/" + newPrefixLength);
                ip += incrementIp;
            }

            return subnets;
        }
    }
}