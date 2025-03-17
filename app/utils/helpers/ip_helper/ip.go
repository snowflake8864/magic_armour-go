package ip_helper

import (
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

// IsIPInRange 检查needle否在haystack中(支持ipv6)
func IsIPInRange(needle string, haystack string) bool {

	if !strings.Contains(haystack, "/") {
		if needle == haystack {
			return true
		} else {
			return false
		}
	}

	if _, ipNet, err := net.ParseCIDR(haystack); err != nil {
		return false
	} else {
		return ipNet.Contains(net.ParseIP(needle))
	}

}

// IPv4ToUint32 将IPv4地址转换为无符号32位整数
func IPv4ToUint32(ipStr string) (uint32, error) {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ipStr)
	}

	ipv4 := ip.To4()
	if ipv4 == nil {
		return 0, fmt.Errorf("not an IPv4 address: %s", ipStr)
	}

	return binary.BigEndian.Uint32(ipv4), nil
}

// IPv4Range IPv4地址范围
type IPv4Range struct {
}

// CalcRange 计算IPv4地址范围，包括网络地址、广播地址、第一可用地址、最后可用地址
func (s *IPv4Range) CalcRange(ipWithMask string) (err error, networkIP, broadcastIP, firstUsableIp, lastUsableIp string) {
	// 尝试解析为CIDR格式
	ip, ipNet, err := net.ParseCIDR(ipWithMask)
	if err == nil {
		if ip.To4() == nil {
			return fmt.Errorf("only IPv4 is supported"), "", "", "", ""
		}
		mask := ipNet.Mask
		networkAddr := ipNet.IP.To4()
		broadcastAddr := s.calcBroadcast(networkAddr, mask)

		ones, _ := mask.Size()
		first, last := s.calcUsableIPs(networkAddr, broadcastAddr, ones)

		return nil, networkAddr.String(), broadcastAddr.String(), first, last
	}

	// 尝试其他格式（如"192.168.1.0 255.255.255.0"）
	parts := strings.FieldsFunc(ipWithMask, func(r rune) bool { return r == '/' || r == ' ' })
	if len(parts) != 2 {
		return fmt.Errorf("invalid IP/mask format"), "", "", "", ""
	}

	// 解析IP
	ip = net.ParseIP(parts[0])
	if ip == nil || ip.To4() == nil {
		return fmt.Errorf("invalid IPv4 address"), "", "", "", ""
	}
	ip = ip.To4()

	// 解析掩码
	mask, err := s.parseMask(parts[1])
	if err != nil {
		return err, "", "", "", ""
	}

	// 计算网络地址和广播地址
	networkAddr := ip.Mask(mask).To4()
	broadcastAddr := s.calcBroadcast(networkAddr, mask)

	// 计算可用地址
	ones, _ := mask.Size()
	first, last := s.calcUsableIPs(networkAddr, broadcastAddr, ones)

	return nil, networkAddr.String(), broadcastAddr.String(), first, last
}

// 解析子网掩码（支持CIDR数字或点分十进制）
func (s *IPv4Range) parseMask(maskStr string) (net.IPMask, error) {
	// 尝试解析为CIDR数字
	if cidr, err := strconv.Atoi(maskStr); err == nil {
		if cidr < 0 || cidr > 32 {
			return nil, fmt.Errorf("invalid CIDR prefix")
		}
		mask := net.CIDRMask(cidr, 32)
		if ones, _ := mask.Size(); ones != cidr {
			return nil, fmt.Errorf("invalid mask")
		}
		return mask, nil
	}

	// 尝试解析为点分十进制
	maskIP := net.ParseIP(maskStr)
	if maskIP == nil {
		return nil, fmt.Errorf("invalid mask format")
	}
	maskIP = maskIP.To4()
	if maskIP == nil {
		return nil, fmt.Errorf("invalid IPv4 mask")
	}

	mask := net.IPv4Mask(maskIP[0], maskIP[1], maskIP[2], maskIP[3])
	if ones, bits := mask.Size(); bits != 32 || ones == 0 {
		return nil, fmt.Errorf("non-contiguous mask")
	}
	return mask, nil
}

// 计算广播地址
func (s *IPv4Range) calcBroadcast(network net.IP, mask net.IPMask) net.IP {
	broadcast := make(net.IP, len(network))
	for i := range network {
		broadcast[i] = network[i] | ^mask[i]
	}
	return broadcast
}

// 计算可用IP范围
func (s *IPv4Range) calcUsableIPs(network, broadcast net.IP, maskBits int) (string, string) {
	switch {
	case maskBits == 32: // 单个地址
		return "", ""
	case maskBits == 31: // RFC 3021特殊处理
		return network.String(), s.addIP(network, 1).String()
	default: // 常规情况
		first := s.addIP(network, 1)
		last := s.subIP(broadcast, 1)
		if s.ipToUint(first) > s.ipToUint(last) {
			return "", ""
		}
		return first.String(), last.String()
	}
}

// IP地址转换辅助函数
func (s *IPv4Range) ipToUint(ip net.IP) uint32 {
	ip = ip.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

func (s *IPv4Range) uintToIP(n uint32) net.IP {
	return net.IPv4(
		byte(n>>24),
		byte(n>>16),
		byte(n>>8),
		byte(n),
	)
}

func (s *IPv4Range) addIP(ip net.IP, inc uint32) net.IP {
	return s.uintToIP(s.ipToUint(ip) + inc)
}

func (s *IPv4Range) subIP(ip net.IP, dec uint32) net.IP {
	return s.uintToIP(s.ipToUint(ip) - dec)
}
