from ipaddress import ip_network, IPv4Address


def split_subnets(base_cidr, new_prefix_length):
    subnets = []
    base_subnet = ip_network(base_cidr)
    first_address = int(base_subnet.network_address)
    last_address = int(base_subnet.broadcast_address)
    mask = ip_network(f"0.0.0.0/{new_prefix_length}")
    increment = int(mask.netmask) - int(IPv4Address('255.255.255.255'))

    i = first_address
    while i <= last_address:
        subnet = f"{IPv4Address(i)}/{new_prefix_length}"
        subnets.append(subnet)
        i += abs(increment) + 1

    return subnets
