package cmn

import "net"

var localIpAddres string

// 取本机IP地址（IPv4）
func GetLocalIp() string {
	if localIpAddres != "" {
		return localIpAddres
	}

	addrs, err := net.InterfaceAddrs()
	if err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
				localIpAddres = ipnet.IP.String()
			}
		}
	}
	return localIpAddres
}
