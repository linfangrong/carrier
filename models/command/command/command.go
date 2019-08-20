package command

import (
	"github.com/linfangrong/carrier/models/protocol"
)

type Command interface {
	Name() []byte
	Proc(cmd Command, msg protocol.Message) protocol.Message
	CheckArgc(argc int64) bool
	CheckReadonly() bool
	CheckForbidden() bool
}

type command struct {
	name      []byte
	proc      func(Command, protocol.Message) protocol.Message //处理函数
	arity     int64                                            //参数个数
	readonly  bool                                             //只读命令
	forbidden bool                                             //禁用命令
}

func NewCommand(
	name []byte,
	proc func(Command, protocol.Message) protocol.Message,
	arity int64,
	readonly bool,
	forbidden bool,
) Command {
	return &command{
		name:      name,
		proc:      proc,
		arity:     arity,
		readonly:  readonly,
		forbidden: forbidden,
	}
}

func (c *command) Name() []byte {
	return c.name
}

func (c *command) Proc(cmd Command, msg protocol.Message) protocol.Message {
	return c.proc(cmd, msg)
}

func (c *command) CheckArgc(argc int64) bool {
	if c.arity > 0 && c.arity != argc {
		return false
	}
	if argc < -c.arity {
		return false
	}
	return true
}

func (c *command) CheckReadonly() bool {
	return c.readonly
}

func (c *command) CheckForbidden() bool {
	return c.forbidden
}
