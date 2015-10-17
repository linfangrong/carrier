package process

import (
	"carrier/cluster/cluster"
	"carrier/cluster/crc16"
	"carrier/cluster/nodes"
	"carrier/command/command"
	"carrier/logger"
	"carrier/protocol"
)

func MGetCommand(cmd command.Command, msg protocol.Message) (msgAck protocol.Message) {
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
		msgAckMap   map[protocol.Message]protocol.Message = make(map[protocol.Message]protocol.Message)

		msgKey   int
		msgValue protocol.Message
	)
	for hashSlot, msgSpice = range classification {
		if node, ok = cluster.GetClusterParameter().GetSlot(hashSlot); !ok {
			logger.Warningf("获取不到slot节点: %d", hashSlot)
			continue
		}
		msgAckSplit = forwardMsg(msgSpice, node, cmd.CheckReadonly())
		if msgAckSplit.GetProtocolType() != protocol.ArraysType {
			logger.Warningf("命令结果预期不符: %v", msgSpice)
			continue
		}
		if msgSpice.GetIntegersValue()-1 != msgAckSplit.GetIntegersValue() {
			logger.Warningf("命令结果预期不符: %v", msgSpice)
			continue
		}
		for msgKey, msgValue = range msgSpice.GetArraysValue()[1:] {
			msgAckMap[msgValue] = msgAckSplit.GetArraysValue()[msgKey]
		}
	}
	// 整合结果
	msgAck = protocol.NewMessage()
	for _, element = range elements[1:] {
		if msgValue, ok = msgAckMap[element]; ok {
			msgAck.AppendArraysValue(msgAckMap[element])
		} else {
			msgAck.AppendArraysValue(protocol.NullBulkString)
		}
	}
	return
}

func MSetCommand(cmd command.Command, msg protocol.Message) (msgAck protocol.Message) {
	var elements []protocol.Message = msg.GetArraysValue()
	if msg.GetIntegersValue()%2 != 1 {
		return protocol.NewMessageString(protocol.ERR_wrong_number_command, cmd.Name())
	}

	var (
		keyIndex       int64
		hashSlot       uint16
		classification map[uint16]protocol.Message = make(map[uint16]protocol.Message)
		msgSpice       protocol.Message
		ok             bool
	)
	for keyIndex = 1; keyIndex < msg.GetIntegersValue(); keyIndex += 2 {
		hashSlot = crc16.HashSlot(elements[keyIndex].GetBytesValue())
		if msgSpice, ok = classification[hashSlot]; !ok {
			msgSpice = protocol.NewMessage().AppendArraysValue(elements[0])
		}
		classification[hashSlot] = msgSpice.AppendArraysValue(elements[keyIndex]).AppendArraysValue(elements[keyIndex+1])
	}
	// 按slot分命令
	var node nodes.Nodes
	for hashSlot, msgSpice = range classification {
		if node, ok = cluster.GetClusterParameter().GetSlot(hashSlot); !ok {
			logger.Warningf("获取不到slot节点: %d", hashSlot)
			continue
		}
		forwardMsg(msgSpice, node, cmd.CheckReadonly())
	}
	return protocol.OK
}
