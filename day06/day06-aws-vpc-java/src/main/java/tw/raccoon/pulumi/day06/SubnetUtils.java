package tw.raccoon.pulumi.day06;

import java.net.InetAddress;
import java.net.UnknownHostException;
import java.util.ArrayList;
import java.util.List;

public class SubnetUtils {

    private static InetAddress intToIpAddress(int ip) {
        byte[] ipBytes = new byte[4];
        for (int i = 0; i < 4; i++) {
            ipBytes[i] = (byte) ((ip >> (24 - (8 * i))) & 0xFF);
        }

        try {
            return InetAddress.getByAddress(ipBytes);
        } catch (UnknownHostException e) {
            e.printStackTrace();
            return null;
        }

    }

    public static List<String> splitSubnets(String baseCidr, int newPrefixLength) {
        List<String> subnets = new ArrayList<>();

        String[] parts = baseCidr.split("/");
        InetAddress baseCidrIp = null;
        try {
            baseCidrIp = InetAddress.getByName(parts[0]);
        } catch (UnknownHostException e) {
            throw new RuntimeException(e);
        }
        int baseCidrPrefixLength = Integer.parseInt(parts[1]);

        // Calculate total IPs to int type
        int totalIps = (int) Math.pow(2, 32 - baseCidrPrefixLength) - 1;

        byte[] addressBytes = baseCidrIp.getAddress();
        int ip = 0;
        for (int i = 0; i < addressBytes.length; i++) {
            ip += (addressBytes[i] & 0xFF) << (24 - (8 * i));
        }

        int lastIpAddressInt = ip + totalIps;
        int incrementIp = (int) Math.pow(2, 32 - newPrefixLength);

        while (ip < lastIpAddressInt) {
            subnets.add(intToIpAddress(ip).getHostAddress() + "/" + newPrefixLength);
            ip += incrementIp;
        }

        return subnets;
    }

}
