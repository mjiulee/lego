package utils

import (
	"fmt"
	"net"
)

// 获取本地ip地址
func GetLocalIpAddress() string {
	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		return ""
	}

	ip := ""
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				//fmt.Println(ipnet.IP.String())
				ip = ipnet.IP.String()
				break
			}
		}
	}
	return ip
}
