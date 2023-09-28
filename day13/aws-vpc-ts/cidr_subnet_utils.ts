import * as ip from 'ip';

export function splitSubnets(baseCidr: string, newPrefixLength: number): string[] {
  let subnets: string[] = [];
  let baseSubnet = ip.cidrSubnet(baseCidr);
  let firstAddress = ip.toLong(baseSubnet.firstAddress);
  let lastAddress = ip.toLong(baseSubnet.lastAddress);
  let mask = ip.fromPrefixLen(newPrefixLength);
  let increment = ip.toLong(mask) - ip.toLong('255.255.255.255');

  for (let i = firstAddress; i <= lastAddress; i += Math.abs(increment) + 1) {
    let subnet = ip.fromLong(i - 1) + '/' + newPrefixLength;
    subnets.push(subnet);
  }

  return subnets;
}
