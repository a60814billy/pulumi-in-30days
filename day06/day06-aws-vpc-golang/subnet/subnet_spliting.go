package subnet

import (
	"net"
	"strconv"
)

// Split splits a subnet into smaller subnets
// cidr: subnet to split
// newPrefix: new prefix for the new subnets
// returns: array of new subnets
func Split(cidr string, newPrefix int) []string {
	_, ipnet, _ := net.ParseCIDR(cidr)

	firstIP := ip2int(ipnet.IP)
	lastIP := ip2int(lastIP(ipnet))

	newSubnets := make([]string, 0)

	for i := firstIP; i <= lastIP; i += 1 << (32 - newPrefix) {
		newSubnet := int2ip(i).String() + "/" + strconv.Itoa(newPrefix)
		newSubnets = append(newSubnets, newSubnet)
	}

	return newSubnets
}

func lastIP(ipnet *net.IPNet) net.IP {
	// return last ip in subnet
	ones, _ := ipnet.Mask.Size()
	lastIp := ip2int(ipnet.IP) | (1<<(32-ones) - 1)
	return int2ip(lastIp)
}

// int2ip converts an uint32 to IP address
func int2ip(i int) net.IP {
	// convert uint32 to IP address
	return net.IPv4(byte(i>>24), byte(i>>16), byte(i>>8), byte(i))
}

// ip2int converts an IP address to uint32
func ip2int(ip net.IP) int {
	ip = ip.To4()
	return int(ip[0])<<24 + int(ip[1])<<16 + int(ip[2])<<8 + int(ip[3])
}
