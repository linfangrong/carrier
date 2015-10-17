package process

import (
	"bytes"
	"carrier/cluster/cluster"
	"carrier/cluster/connection"
	"carrier/cluster/crc16"
	"carrier/cluster/nodes"
	"carrier/command/command"
	"carrier/logger"
	"carrier/protocol"
)

var (
	ASK   []byte = []byte("ASK")
	MOVED []byte = []byte("MOVED")
)

func ProcKeyCommand(cmd command.Command, msg protocol.Message) (msgAck protocol.Message) {
	var (
		slot uint16 = crc16.HashSlot(msg.GetArraysValue()[1].GetBytesValue()) // 第一个参数是Key
		node nodes.Nodes
		ok   bool
	)
	if node, ok = cluster.GetClusterParameter().GetSlot(slot); !ok {
		logger.Warningf("获取不到slot节点: %d", slot)
		return protocol.ERR_ClusterSlots
	}
	return forwardMsg(msg, node, cmd.CheckReadonly())
}

func forwardMsg(msg protocol.Message, node nodes.Nodes, readonly bool) (msgAck protocol.Message) {
	var (
		clusterConnPool connection.Pool
		isMaster        bool
	)
	if readonly {
		clusterConnPool, isMaster = node.GetRandom()
	} else {
		clusterConnPool = node.GetMaster()
		isMaster = true
	}
	msgAck = forwardMsgToPool(msg, clusterConnPool, isMaster, false)
	if msgAck.GetProtocolType() == protocol.ErrorsStringsType {
		var msgAckBytesValueSplit [][]byte = bytes.Fields(msgAck.GetBytesValue())
		if len(msgAckBytesValueSplit) == 3 {
			switch {
			case bytes.EqualFold(msgAckBytesValueSplit[0], MOVED):
				msgAck = forwardMsgToPool(msg, cluster.GetClusterParameter().GetNodePool(string(msgAckBytesValueSplit[2])), true, false)
			case bytes.EqualFold(msgAckBytesValueSplit[0], ASK):
				msgAck = forwardMsgToPool(msg, cluster.GetClusterParameter().GetNodePool(string(msgAckBytesValueSplit[2])), true, true)
			}
		}
	}
	return
}

func forwardMsgToPool(msg protocol.Message, clusterConnPool connection.Pool, isMaster bool, isAsking bool) (msgAck protocol.Message) {
	var (
		clusterConn connection.Conn
		err         error
	)
	if clusterConn, err = clusterConnPool.Get(); err != nil {
		logger.Warningf("从连接池中获取句柄出错: %v", err)
		return protocol.ERR_ClusterSlotsConn
	}
	if isMaster {
		if err = clusterConn.Readwrite(); err != nil {
			goto failure
		}
	} else {
		if err = clusterConn.Readonly(); err != nil {
			goto failure
		}
	}
	if isAsking {
		if err = clusterConn.Asking(); err != nil {
			goto failure
		}
	}
	if msgAck, err = clusterConn.HandelMessage(msg); err != nil {
		goto failure
	}
	clusterConnPool.Put(clusterConn)
	return
failure:
	clusterConnPool.Remove(clusterConn)
	logger.Warningf("转发消息出错: %v", err)
	return protocol.ERR_ClusterSlotsForward
}
