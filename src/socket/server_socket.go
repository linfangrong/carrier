package socket

import (
	"fmt"
	"net"
)

type TServerSocket struct {
	listener net.Listener
	addr     net.Addr
}

func NewTServerSocket(listenAddr string) (socket *TServerSocket, err error) {
	var addr net.Addr
	if addr, err = net.ResolveTCPAddr("tcp", listenAddr); err != nil {
		return
	}
	socket = &TServerSocket{addr: addr}
	return
}

func (p *TServerSocket) IsListening() bool {
	return p.listener != nil
}

func (p *TServerSocket) Listen() (err error) {
	if p.IsListening() {
		return
	}
	var listener net.Listener
	if listener, err = net.Listen(p.addr.Network(), p.addr.String()); err != nil {
		return
	}
	p.listener = listener
	return
}

func (p *TServerSocket) Accept() (conn net.Conn, err error) {
	if p.listener == nil {
		err = fmt.Errorf("no underlying server socket")
		return
	}
	if conn, err = p.listener.Accept(); err != nil {
		return
	}
	return
}

func (p *TServerSocket) Close() (err error) {
	defer func() {
		p.listener = nil
	}()
	if p.IsListening() {
		return p.listener.Close()
	}
	return
}

func (p *TServerSocket) Addr() net.Addr {
	if p.listener != nil {
		return p.listener.Addr()
	}
	return p.addr
}
