package util

import (
	"encoding/binary"
	"net"
)

func IP2Long(ipstr string) uint32 {
	var ip net.IP = net.ParseIP(ipstr)
	if ip == nil {
		return 0
	}
	ip = ip.To4()
	return binary.BigEndian.Uint32(ip)
}

func Long2IP(ipLong uint32) string {
	var (
		ipByte []byte = make([]byte, 4)
		ip     net.IP
	)
	binary.BigEndian.PutUint32(ipByte, ipLong)
	ip = net.IP(ipByte)
	return ip.String()
}
