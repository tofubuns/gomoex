package netutil

import (
	"net"
)

// GetAvailablePort 获得可用的随机端口
func GetAvailablePort() int {
	l, _ := net.Listen("tcp", ":0") // 监听一个随机的端口
	defer l.Close()
	port := l.Addr().(*net.TCPAddr).Port
	return port
}
