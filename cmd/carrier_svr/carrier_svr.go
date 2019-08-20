package main

import (
	"net"

	"github.com/linfangrong/carrier/models/logger"
	"github.com/linfangrong/carrier/utils/socket"
)

func ServeForClient(listenAddr string) {
	var (
		serverSocket *socket.TServerSocket
		serverConn   net.Conn
		err          error
	)
	if serverSocket, err = socket.NewTServerSocket(listenAddr); err != nil {
		panic(err)
	}
	if err = serverSocket.Listen(); err != nil {
		panic(err)
	}
	defer serverSocket.Close()

	// AcceptLoop
	for {
		if serverConn, err = serverSocket.Accept(); err != nil {
			logger.Notice("Socket To Client Accept: ", err)
			continue
		}
		go func(connClient net.Conn) {
			defer func() {
				logger.Info("客户端断开连接: ", connClient.RemoteAddr().String())
				connClient.Close()
			}()
			logger.Info("客户端新连接: ", connClient.RemoteAddr().String())
			ClientLogic(connClient)
		}(serverConn)
	}
}
