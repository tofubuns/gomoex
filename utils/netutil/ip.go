package netutil

import (
	"net"
	"strings"
)

// privateIPV4Block 内网的ipv4范围
var privateIPV4Block = []string{
	"10.0.0.0/8",     // 10.0.0.0 到 10.255.255.255
	"172.16.0.0/12",  // 172.16.0.0 到 172.31.255.255
	"192.168.0.0/16", // 192.168.0.0 到 192.168.255.255
	"169.254.0.0/16", // 169.254.0.0 到 169.254.255.255
}

// SetPrivateIPV4Block 设置或修改默认的内网 ip 块
// 请注意, 此函数带来的副作用是用永久的
func SetPrivateIPV4Block(fn func([]string) []string) {
	privateIPV4Block = fn(privateIPV4Block)
}

// GetLocalIPV4 获取本机的 IPV4 地址
//
// isIntranet 参数表示是否需要返回内网 IP
func GetLocalIPV4(isIntranet bool) (ip string) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return
	}
	for _, inter := range interfaces {
		if inter.Flags&net.FlagUp != 0 && !strings.HasPrefix(inter.Name, "lo") {
			addrs, err := inter.Addrs()
			if err != nil {
				continue
			}
			for _, addr := range addrs {
				if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
					if ipnet.IP.To4() != nil {
						// 当前 ip 是否为内网 ip
						intranetIpFlag := ipnet.IP.String() == "127.0.0.1" && !IsIntranetIpv4(ipnet)

						// 如果需要内网 IP，且当前 IP 是内网 IP，则返回
						if isIntranet && intranetIpFlag {
							ip = ipnet.IP.String()
							return
						}

						// 如果需要外网 IP，且当前 IP 不是内网 IP，则返回
						if !isIntranet && !intranetIpFlag {
							ip = ipnet.IP.String()
							return
						}
					}
				}
			}
		}
	}

	return
}

// IsIntranetIpv4 是否是内网 Ipv4
func IsIntranetIpv4(ipnet *net.IPNet) bool {
	if ipnet.IP.To4() == nil {
		return false
	}

	// 解析这些内网ip块, 如果验证 ip 在这些代码块中
	// 就认为验证 ip 为内网
	for _, block := range privateIPV4Block {
		_, cidr, _ := net.ParseCIDR(block)
		if cidr.Contains(ipnet.IP) {
			return true
		}
	}
	return false
}

// IsIntranetIpv4ByString 判断字符类型的 ip 是否为内网
func IsIntranetIpv4ByString(ip string) bool {
	_, ipnet, err := net.ParseCIDR(ip)
	if err != nil {
		return false
	}

	return IsIntranetIpv4(ipnet)
}
