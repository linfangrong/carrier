package main

import (
	carrierConn "carrier/conn"
	"net"
)

func ClientLogic(connClient net.Conn) {
	// TODO 对连接做处理(IP限制什么鬼的)
	// TODO AUTH

	var (
		connCarrier carrierConn.Conn
	)
	connCarrier = carrierConn.NewConn(connClient)
	<-connCarrier.IsClosed()
}
