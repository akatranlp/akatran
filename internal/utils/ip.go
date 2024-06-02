package utils

import (
	"io"
	"net"
	"net/http"
	"strings"
)

func fetchIPv4() string {
	res, err := http.Get("https://checkip.amazonaws.com")
	if err != nil {
		return ""
	}
	defer res.Body.Close()

	var ip strings.Builder
	if _, err := io.Copy(&ip, res.Body); err != nil {
		return ""
	}

	return strings.TrimSpace(ip.String())
}

func fetchIPv6() string {
	return "::1"
}

func GetIPv4Address(ip string) net.IP {
	if ip == "" {
		ip = fetchIPv4()
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
		ip = fetchIPv6()
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
