package sysinfo

import (
	"net"
	"os"
)

// SysInfo - Stores all the information that needs to be returned
type SysInfo struct {
	IPAddress, ContainerName string
}

// GetSystemInfo - Returns information about the system
func GetSystemInfo() *SysInfo {
	info := &SysInfo{}
	info.IPAddress = getIPAddress()
	info.ContainerName = os.Getenv("CONTAINER_NAME")

	return info
}

func getIPAddress() string {
	ifaces, err := net.Interfaces()

	if err != nil {
		return "unknown"
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err != nil {
			return "unknown"
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || ip.IsLoopback() {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}
			return ip.String()
		}
	}

	return "unknown"
}
