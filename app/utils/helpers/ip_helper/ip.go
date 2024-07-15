package ip_helper

import (
	"net"
	"strings"
)

// IsIPInRange 检查IP地址是否在IP段中(支持iptv6)
func IsIPInRange(ip string, ipWithMask string) bool {

	if strings.Contains(ipWithMask, "/") {
		if ip == ipWithMask {
			return true
		} else {
			return false
		}

	}
	_, netNet, _ := net.ParseCIDR(ipWithMask)
	return netNet.Contains(net.ParseIP(ip))
}
