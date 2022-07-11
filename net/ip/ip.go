package ip

import (
	"net"
	"net/http"
)

var _, reserveSubnet, _ = net.ParseCIDR("100.64.0.0/10")

func IpIsProxy(ipstr string) bool {
	ip := net.ParseIP(ipstr)
	return reserveSubnet.Contains(ip)
}

func RemoteAddress(r *http.Request) string {
	ipAddress := r.Header.Get("X-Real-Ip")
	if ipAddress == "" {
		ipAddress = r.Header.Get("X-Forwarded-For")
	}
	if ipAddress == "" {
		ipAddress = r.RemoteAddr
	}
	return ipAddress
}
