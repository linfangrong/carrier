package process

import (
	"carrier/command/command"
	"carrier/protocol"
)

func ForbiddenCommand(cmd command.Command, msg protocol.Message) protocol.Message {
	return protocol.NewMessageString(protocol.ERR_forbidden_command, cmd.Name())
}
