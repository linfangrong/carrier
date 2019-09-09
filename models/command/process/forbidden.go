package process

import (
	"github.com/linfangrong/carrier/models/command/command"
	"github.com/linfangrong/carrier/models/protocol"
)

func ForbiddenCommand(cmd command.Command, msg protocol.Message) protocol.Message {
	return protocol.NewMessageString(protocol.ERR_forbidden_command, cmd.Name())
}
