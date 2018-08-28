package sysinfo

import (
	"errors"
	"net"
)

// SysInfo - Stores all the information that needs to be returned
type SysInfo struct {
	IPAddress string
}

// GetSystemInfo - Returns information about the system
func GetSystemInfo() (*SysInfo, error) {
	info := &SysInfo{}
	ifaces, err := net.Interfaces()

	if err != nil {
		return nil, err
	}

	for _, i := range ifaces {
		addrs, err := i.Addrs()

		if err != nil {
			return nil, err
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
			info.IPAddress = ip.String()
			return info, nil
		}
	}

	return nil, errors.New("couldn't get the ip address")
}
