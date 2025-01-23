package xnet

import (
	"fmt"
	"net"
	"strings"
)

func isPrivateIP(ip net.IP) bool {
	// 检查是否为 IPv4 地址
	if ip.To4() == nil {
		return false // 如果不是 IPv4 地址，则返回 false
	}

	// 定义私有 IP 地址范围
	privateRanges := []net.IPNet{
		{IP: net.ParseIP("10.0.0.0"), Mask: net.CIDRMask(8, 32)},
		{IP: net.ParseIP("172.16.0.0"), Mask: net.CIDRMask(12, 32)},
		{IP: net.ParseIP("192.168.0.0"), Mask: net.CIDRMask(16, 32)},
	}

	// 检查 IP 是否在私有范围内
	for _, privateRange := range privateRanges {
		if privateRange.Contains(ip) {
			return true
		}
	}
	return false
}

func ReplaceEthIP(host string) string {
	ethIP := CheckLocalIP()
	if len(ethIP) > 0 {
		parts := strings.Split(host, ":")
		if len(parts) == 2 {
			port := parts[1] // 提取端口号
			return fmt.Sprintf("%s:%s", ethIP, port)
		} else {
			_ = fmt.Errorf("invalid host address format")
		}
	}
	return ""
}

func CheckLocalIP() string {
	fmt.Println("----------------check local ip---------------------------------")
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("Error:", err)
		return ""
	}

	ips := make([]string, 0)
	for _, i := range interfaces {
		// 只处理已激活的接口
		if i.Flags&net.FlagUp == 0 || i.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 获取接口的地址信息
		addrs, err := i.Addrs()
		if err != nil {
			fmt.Println("Error getting addresses:", err)
			continue
		}

		// 打印每个接口的 IP 地址
		for _, addr := range addrs {
			// 解析 IP 地址
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			// 排除回环地址和IPv6地址
			if ip != nil && ip.IsLoopback() == false && ip.To4() != nil {
				fmt.Printf("Interface Name: %v, Hardware(MAC) Address: %v, MTU: %v IP: %v \r\n", i.Name, i.HardwareAddr, i.MTU, ip)
				if isPrivateIP(ip) {
					ips = append(ips, ip.String())
				}
			}
		}
	}
	fmt.Println(ips)
	return ips[0]
}
