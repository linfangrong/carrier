package process

import (
	"github.com/linfangrong/carrier/models/cluster/cluster"
	"github.com/linfangrong/carrier/models/cluster/crc16"
	"github.com/linfangrong/carrier/models/cluster/nodes"
	"github.com/linfangrong/carrier/models/command/command"
	"github.com/linfangrong/carrier/models/logger"
	"github.com/linfangrong/carrier/models/protocol"
)

func DelCommand(cmd command.Command, msg protocol.Message) (msgAck protocol.Message) {
	var (
		elements       []protocol.Message = msg.GetArraysValue()
		element        protocol.Message
		hashSlot       uint16
		classification map[uint16]protocol.Message = make(map[uint16]protocol.Message)
		msgSpice       protocol.Message
		ok             bool
	)
	for _, element = range elements[1:] {
		hashSlot = crc16.HashSlot(element.GetBytesValue())
		if msgSpice, ok = classification[hashSlot]; !ok {
			msgSpice = protocol.NewMessage().AppendArraysValue(elements[0])
		}
		classification[hashSlot] = msgSpice.AppendArraysValue(element)
	}
	// 按slot分命令
	var (
		node        nodes.Nodes
		msgAckSplit protocol.Message
	)
	msgAck = protocol.NewMessage()
	for hashSlot, msgSpice = range classification {
		if node, ok = cluster.GetClusterParameter().GetSlot(hashSlot); !ok {
			logger.Warningf("获取不到slot节点: %d", hashSlot)
			continue
		}
		msgAckSplit = forwardMsg(msgSpice, node, cmd.CheckReadonly())
		msgAck.AppendIntegersValue(msgAckSplit)
	}
	return
}
