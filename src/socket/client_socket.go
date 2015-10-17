package socket

import (
	"net"
	"time"
)

type TClientSocket struct {
	tcpAddr *net.TCPAddr
}

func NewTClientSocket(listenAddr string) (socket *TClientSocket, err error) {
	var addr *net.TCPAddr
	if addr, err = net.ResolveTCPAddr("tcp", listenAddr); err != nil {
		return
	}
	socket = &TClientSocket{tcpAddr: addr}
	return
}

func (p *TClientSocket) Dial() (conn net.Conn, err error) {
	if conn, err = net.DialTCP("tcp", nil, p.tcpAddr); err != nil {
		return
	}
	return
}

func (p *TClientSocket) DialTimeout(timeout time.Duration) (conn net.Conn, err error) {
	if conn, err = net.DialTimeout("tcp", p.tcpAddr.String(), timeout); err != nil {
		return
	}
	return
}
