package ip

import "net"

var _, reserveSubnet, _ = net.ParseCIDR("100.64.0.0/10")

func IpIsProxy(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	return reserveSubnet.Contains(ip)
}
