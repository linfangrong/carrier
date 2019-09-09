package protocol

import (
	"bufio"
	"bytes"
	"fmt"
)

var (
	OK             Message = NewMessageString("+OK\r\n")
	PING           Message = NewMessageString("*1\r\n$4\r\nPING\r\n")
	PONG           Message = NewMessageString("+PONG\r\n")
	ClusterSlots   Message = NewMessageString("*2\r\n$7\r\nCLUSTER\r\n$5\r\nSLOTS\r\n")
	NullBulkString Message = NewMessageString("$-1\r\n")
	READONLY       Message = NewMessageString("*1\r\n$8\r\nREADONLY\r\n")
	READWRITE      Message = NewMessageString("*1\r\n$9\r\nREADWRITE\r\n")
	ASKING         Message = NewMessageString("*1\r\n$6\r\nASKING\r\n")

	ERR_ClusterSlots        Message = NewMessageString("-ERR cluster slots\r\n")
	ERR_ClusterSlotsConn    Message = NewMessageString("-ERR cluster slots connection\r\n")
	ERR_ClusterSlotsForward Message = NewMessageString("-ERR cluster slots forward\r\n")
)

var (
	ERR_unknown_command      string = "-ERR unknown command '%s'\r\n"
	ERR_wrong_number_command string = "-ERR wrong number of arguments for '%s' command\r\n"
	ERR_forbidden_command    string = "-ERR forbidden command '%s'\r\n"
)

func NewMessageString(format string, a ...interface{}) Message {
	var (
		buf *bytes.Buffer = bytes.NewBufferString(fmt.Sprintf(format, a...))
		br  *bufio.Reader = bufio.NewReader(buf)
		msg Message       = NewMessage()
		err error
	)
	if err = msg.ReadOne(br); err != nil {
		return nil
	}
	return msg
}
