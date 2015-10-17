package process

import (
	"carrier/command/command"
	"carrier/protocol"
)

func EchoCommand(cmd command.Command, msg protocol.Message) protocol.Message {
	return msg.GetArraysValue()[1]
}

func PingCommand(cmd command.Command, msg protocol.Message) protocol.Message {
	if msg.GetIntegersValue() > 2 {
		return protocol.NewMessageString(protocol.ERR_wrong_number_command, cmd.Name())
	}
	if msg.GetIntegersValue() == 2 {
		return msg.GetArraysValue()[1]
	}
	return protocol.PONG
}

func SelectCommand(cmd command.Command, msg protocol.Message) protocol.Message {
	return protocol.OK
}
