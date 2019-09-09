package conn

import (
	"bufio"
	"fmt"
	"net"
	"time"

	CommandSet "github.com/linfangrong/carrier/models/command/set"
	"github.com/linfangrong/carrier/models/logger"
	"github.com/linfangrong/carrier/models/protocol"
	"github.com/linfangrong/carrier/utils/datatypes"
)

const (
	DEFAULT_SEND_BUFFER_SIZE    = 5
	DEFAULT_RECEIVE_BUFFER_SIZE = 5
	DEFAULT_IDLE_TIME_OUT       = time.Second * 300
)

type Conn interface {
	Close()
	IsClosed() <-chan bool
	GetConnAddr() string
	GetLatestActiveTime() time.Time
}

type conn struct {
	// send message buffer
	sendMessageQueue chan protocol.Message
	// receive message buffer
	receiveMessageQueue chan protocol.Message

	// Network connection
	c    net.Conn
	br   *bufio.Reader
	bw   *bufio.Writer
	addr string

	// latestActiveTime
	latestActiveTime time.Time

	// Closed
	closed *datatypes.SyncClose
}

func NewConn(c net.Conn) Conn {
	conn := &conn{
		sendMessageQueue:    make(chan protocol.Message, DEFAULT_SEND_BUFFER_SIZE),
		receiveMessageQueue: make(chan protocol.Message, DEFAULT_RECEIVE_BUFFER_SIZE),

		c:    c,
		br:   bufio.NewReader(c),
		bw:   bufio.NewWriter(c),
		addr: fmt.Sprintf("%s_%s", c.LocalAddr().String(), c.RemoteAddr().String()),

		latestActiveTime: time.Now(),

		closed: datatypes.NewSyncClose(),
	}
	go conn.sendLoop()
	go conn.readLoop()
	go conn.handleLoop()
	addConn(conn.addr, conn)
	return conn
}

func (conn *conn) Close() {
	removeConn(conn.addr)
	conn.closed.Close()
}

func (conn *conn) IsClosed() <-chan bool {
	return conn.closed.IsClosed()
}

func (conn *conn) GetConnAddr() string {
	return conn.addr
}

func (conn *conn) GetLatestActiveTime() time.Time {
	return conn.latestActiveTime
}

func (conn *conn) updateActiveTime() {
	conn.latestActiveTime = time.Now()
}

func (conn *conn) error(err error, desc string) {
	if err == nil {
		return
	}
	if nerr, ok := err.(net.Error); ok {
		switch {
		case nerr.Timeout(), nerr.Temporary():
			logger.Infof("客户端连接: %s %v", conn.c.RemoteAddr().String(), nerr)
			return
		}
	}
	logger.Infof("客户端连接出错: %s %s %v", conn.c.RemoteAddr().String(), desc, err)
	conn.Close()
}

func (conn *conn) sendMessage(msg protocol.Message) (err error) {
	err = msg.WriteOne(conn.bw)
	conn.error(err, "sendMessage 出错")
	return
}

func (conn *conn) readMessage() (msg protocol.Message, err error) {
	msg = protocol.NewMessage()
	err = msg.ReadOne(conn.br)
	conn.error(err, "readMessage 出错")
	return
}

// send loop
func (conn *conn) sendLoop() {
	var msg protocol.Message
	for {
		select {
		case <-conn.closed.IsClosed():
			return
		case msg = <-conn.sendMessageQueue:
			conn.sendMessage(msg)
		}
	}
}

// read loop
func (conn *conn) readLoop() {
	var (
		msg protocol.Message
		err error
	)
	for {
		select {
		case <-conn.closed.IsClosed():
			return
		default:
			if msg, err = conn.readMessage(); err == nil {
				conn.receiveMessageQueue <- msg
			}
		}
	}
}

// handle loop
func (conn *conn) handleLoop() {
	var msg protocol.Message
	for {
		select {
		case <-conn.closed.IsClosed():
			return
		case msg = <-conn.receiveMessageQueue:
			conn.updateActiveTime()
			conn.sendMessageQueue <- CommandSet.ProcessCommand(msg)
		}
	}
}
