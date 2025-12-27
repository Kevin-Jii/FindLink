package utils

import (
	"app/utils/logger"
	"fmt"
	"net"

	"go.uber.org/zap"
)

func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}
	for _, addr := range addrs {
		if ipNet, ok := addr.(*net.IPNet); ok && !ipNet.IP.IsLoopback() {
			if ipNet.IP.To4() != nil {
				return ipNet.IP.String()
			}
		}
	}
	return ""
}

func PrintLocalIPWithEndpoint() {
	ip := GetLocalIP()
	basePath := "/api/app"
	endpoint := fmt.Sprintf("http://%s:8900%s", ip, basePath)
	logger.Debug("获取本机IP地址",
		zap.String("ip", ip),
		zap.String("endpoint", endpoint),
		zap.String("basePath", basePath),
	)
}
