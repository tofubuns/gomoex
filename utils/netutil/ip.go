package netutil

import (
	"net"
	"strings"
)

func GetLocalIPV4() (ip string) {
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
						return ipnet.IP.String()
					}
				}
			}
		}
	}

	return
}
