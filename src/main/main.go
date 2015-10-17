package main

import (
	"bufio"
	"carrier/commands"
	"carrier/protocol"
	"fmt"
	"socket"
	"time"
	"util"
)

func main() {

	c := commands.NewCommandTree()
	c.AddCommand([]byte("aBcD"), commands.NewCommand())

	fmt.Println(c.SearchCommand([]byte("ABCD")))
	fmt.Println(c.SearchCommand([]byte("AbCD")))
	fmt.Println(c.SearchCommand([]byte("ABCd")))
	fmt.Println(c.SearchCommand([]byte("ABC")))

	for i := 'a'; i <= 'z'; i++ {
		j := util.ToUpper(byte(i))
		k := util.ToLower(j)
		fmt.Println(i, j, k)
	}
	return
	socketClent, err := socket.NewTClientSocket("10.10.10.227:6679")
	//	socketClent, err := socket.NewTClientSocket("127.0.0.1:6679")
	fmt.Println(socketClent, err)
	conn, err := socketClent.DialTimeout(time.Duration(3 * time.Second))
	fmt.Println(conn, err)

	var (
		br     *bufio.Reader    = bufio.NewReader(conn)
		bw     *bufio.Writer    = bufio.NewWriter(conn)
		msgAck protocol.Message = protocol.NewMessage()
	)

	// 发数据
	var rawBytes []byte = []byte("*1\r\n$4\r\nINFO\r\n")
	//var rawBytes []byte = []byte{42, 50, 13, 10, 36, 55, 13, 10, 99, 108, 117, 115, 116, 101, 114, 13, 10, 36, 53, 13, 10, 115, 108, 111, 116, 115, 13, 10}
	//var rawBytes []byte = []byte("*2\r\n$7\r\nCLUSTER\r\n$5\r\nSLOTS\r\n")
	fmt.Println(string(rawBytes))
	if _, err = bw.Write(rawBytes); err != nil {
		fmt.Println(err)
		return
	}
	bw.Flush()

	// 收数据
	if err = msgAck.ReadOne(br); err != nil {
		fmt.Println(err)
	}
	fmt.Println(msgAck)
}
