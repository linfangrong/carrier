package set

import (
	"github.com/linfangrong/carrier/models/command/command"
	"github.com/linfangrong/carrier/models/protocol"
)

var commandSet CommandTree = NewCommandTree()

func init() {
	var cmd command.Command
	for _, cmd = range CommandTable {
		commandSet.AddCommand(cmd.Name(), cmd)
	}
}

func ProcessCommand(msg protocol.Message) protocol.Message {
	var (
		name []byte
		argc int64
		cmd  command.Command
		ok   bool
	)
	switch msg.GetProtocolType() {
	case protocol.ArraysType:
		if argc = msg.GetIntegersValue(); argc > 0 {
			name = msg.GetArraysValue()[0].GetBytesValue()
		}
	}

	if cmd, ok = commandSet.SearchCommand(name); !ok {
		return protocol.NewMessageString(protocol.ERR_unknown_command, name)
	}
	if !cmd.CheckArgc(argc) {
		return protocol.NewMessageString(protocol.ERR_wrong_number_command, cmd.Name())
	}
	if cmd.CheckForbidden() {
		return protocol.NewMessageString(protocol.ERR_forbidden_command, cmd.Name())
	}
	return cmd.Proc(cmd, msg)
}
