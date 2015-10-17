package connection

import (
	"bufio"
	"carrier/protocol"
	"datatypes"
	"net"
	"socket"
	"time"
)

type Conn interface {
	Close()
	IsClosed() <-chan bool
	HandelMessage(msg protocol.Message) (msgAck protocol.Message, err error)
	Asking() (err error)
	Readonly() (err error)
	Readwrite() (err error)
}

const (
	askingFlag   int = 1 << 0
	readonlyFlag int = 1 << 1
)

type conn struct {
	// Network connection
	c            net.Conn
	br           *bufio.Reader
	bw           *bufio.Writer
	readTimeout  time.Duration
	writeTimeout time.Duration

	// flag
	flag int

	// Closed
	closed *datatypes.SyncClose
}

func NewConn(address string, connectTimeout, readTimeout, writeTimeout time.Duration) (c Conn, err error) {
	var (
		clientSocket *socket.TClientSocket
		clientConn   net.Conn
	)
	if clientSocket, err = socket.NewTClientSocket(address); err != nil {
		return
	}
	if clientConn, err = clientSocket.DialTimeout(connectTimeout); err != nil {
		return
	}
	c = &conn{
		c:            clientConn,
		br:           bufio.NewReader(clientConn),
		bw:           bufio.NewWriter(clientConn),
		readTimeout:  readTimeout,
		writeTimeout: writeTimeout,
		flag:         0,
		closed:       datatypes.NewSyncClose(),
	}
	return
}

func (conn *conn) Close() {
	// 发quit
	conn.closed.Close()
}

func (conn *conn) IsClosed() <-chan bool {
	return conn.closed.IsClosed()
}

func (conn *conn) HandelMessage(msg protocol.Message) (msgAck protocol.Message, err error) {
	conn.c.SetWriteDeadline(time.Now().Add(conn.writeTimeout))
	if err = msg.WriteOne(conn.bw); err != nil {
		return
	}
	conn.c.SetReadDeadline(time.Now().Add(conn.readTimeout))
	msgAck = protocol.NewMessage()
	if err = msgAck.ReadOne(conn.br); err != nil {
		return
	}
	return
}

func (conn *conn) Asking() (err error) {
	if conn.flag&askingFlag != 0 {
		return
	}
	if _, err = conn.HandelMessage(protocol.ASKING); err != nil {
		return
	}
	conn.flag |= askingFlag
	return
}

func (conn *conn) Readonly() (err error) {
	if conn.flag&readonlyFlag != 0 {
		return
	}
	if _, err = conn.HandelMessage(protocol.READONLY); err != nil {
		return
	}
	conn.flag |= readonlyFlag
	return
}

func (conn *conn) Readwrite() (err error) {
	if conn.flag&readonlyFlag == 0 {
		return
	}
	if _, err = conn.HandelMessage(protocol.READWRITE); err != nil {
		return
	}
	conn.flag &= ^readonlyFlag
	return
}
