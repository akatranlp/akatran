package utils

import "net"

func GetIPv4Address(ip string) net.IP {
	if ip == "" {
		ip = "127.0.0.1"
	}
	addr := net.ParseIP(ip)

	if addr == nil {
		return nil
	}
	if addr.To4() == nil {
		return nil
	}
	return addr
}

func GetIPv6Address(ip string) net.IP {
	if ip == "" {
		ip = "::1"
	}
	addr := net.ParseIP(ip)

	if addr == nil {
		return nil
	}
	if addr.To16() == nil {
		return nil
	}
	return addr
}
