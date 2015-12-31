package cluster

import (
	"carrier/cluster/connection"
	"carrier/cluster/nodes"
	"carrier/cluster/slots"
	"carrier/logger"
	"carrier/protocol"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type Cluster interface {
	GetSlot(slot uint16) (nodes nodes.Nodes, ok bool)
	GetNodePool(server string) (connPool connection.Pool)
}

type cluster struct {
	sync.RWMutex
	refreshInterval time.Duration              // cluster 刷新间隔
	clusterNodePool map[string]connection.Pool // cluster 节点连接池
	clusterSlots    slots.Slots                // cluster slots信息

	// 连接池配置
	maxIdle             int
	testOnBorrowTimeout time.Duration
	connectTimeout      time.Duration
	readTimeout         time.Duration
	writeTimeout        time.Duration
}

func NewCluster(
	serverList []string, refreshInterval time.Duration,
	maxIdle int, testOnBorrowTimeout time.Duration,
	connectTimeout time.Duration, readTimeout time.Duration, writeTimeout time.Duration,
) (Cluster, error) {
	if len(serverList) <= 0 {
		return nil, fmt.Errorf("cluster配置节点为空")
	}

	var c *cluster = &cluster{
		refreshInterval:     refreshInterval,
		clusterNodePool:     make(map[string]connection.Pool),
		clusterSlots:        slots.NewSlots(),
		maxIdle:             maxIdle,
		testOnBorrowTimeout: testOnBorrowTimeout,
		connectTimeout:      connectTimeout,
		readTimeout:         readTimeout,
		writeTimeout:        writeTimeout,
	}
	var (
		err    error
		server string = serverList[rand.Intn(len(serverList))]
	)
	// 初始化
	if err = c.refresh(server); err != nil {
		return nil, err
	}
	// 定时刷新
	go func() {
		for range time.Tick(refreshInterval) {
			var (
				err    error
				server string = serverList[rand.Intn(len(serverList))]
			)
			if err = c.refresh(server); err != nil {
				logger.Infof("节点(%s):刷新cluster slots出错(%v)", server, err)
			}
		}
	}()
	return c, nil
}

func (c *cluster) refresh(server string) (err error) {
	var (
		connPool    connection.Pool = c.GetNodePool(server)
		conn        connection.Conn
		slotsMsgAck protocol.Message
	)
	if conn, err = connPool.Get(); err != nil {
		return
	}
	// 发送cluster slots
	if slotsMsgAck, err = conn.HandelMessage(protocol.ClusterSlots); err != nil {
		connPool.Remove(conn)
		return
	}
	connPool.Put(conn)
	// 处理返回值
	if slotsMsgAck.GetProtocolType() != protocol.ArraysType || slotsMsgAck.GetIntegersValue() <= 0 {
		return fmt.Errorf("slots数据无效")
	}
	var slotsMsgAckValue protocol.Message
	for _, slotsMsgAckValue = range slotsMsgAck.GetArraysValue() {
		if slotsMsgAckValue.GetProtocolType() != protocol.ArraysType || slotsMsgAckValue.GetIntegersValue() < 3 {
			return fmt.Errorf("slots节点数据无效")
		}
		var (
			slotsMsgAckValueArrays []protocol.Message = slotsMsgAckValue.GetArraysValue()
			nodeMsgValue           protocol.Message
			nodeMsgValueArrays     []protocol.Message
			connPoolList           []connection.Pool = make([]connection.Pool, 0, slotsMsgAckValue.GetIntegersValue()-2)
		)
		if slotsMsgAckValueArrays[0].GetProtocolType() != protocol.IntegersType || slotsMsgAckValueArrays[1].GetProtocolType() != protocol.IntegersType {
			return fmt.Errorf("slots节点起始或者结束数据无效")
		}
		for _, nodeMsgValue = range slotsMsgAckValueArrays[2:] {
			if nodeMsgValue.GetProtocolType() != protocol.ArraysType || nodeMsgValue.GetIntegersValue() < 2 {
				return fmt.Errorf("slot节点数据无效")
			}
			nodeMsgValueArrays = nodeMsgValue.GetArraysValue()
			if nodeMsgValueArrays[0].GetProtocolType() != protocol.BulkStringsType {
				return fmt.Errorf("slot节点IP无效")
			}
			if nodeMsgValueArrays[1].GetProtocolType() != protocol.IntegersType {
				return fmt.Errorf("slot节点PORT无效")
			}
			connPoolList = append(
				connPoolList,
				c.GetNodePool(fmt.Sprintf("%s:%d", nodeMsgValueArrays[0].GetBytesValue(), nodeMsgValueArrays[1].GetIntegersValue())),
			)
		}
		c.clusterSlots.AddSlots(
			slotsMsgAckValueArrays[0].GetIntegersValue(),
			slotsMsgAckValueArrays[1].GetIntegersValue(),
			nodes.NewNodes().Set(connPoolList), // TODO 同一个List用同一个Nodes
		)
	}
	return
}

func (c *cluster) GetNodePool(server string) (connPool connection.Pool) {
	var ok bool
	c.RLock()
	if connPool, ok = c.clusterNodePool[server]; ok {
		c.RUnlock()
		return
	}
	c.RUnlock()

	var (
		newFn func() (connection.Conn, error) = func() (conn connection.Conn, err error) {
			if conn, err = connection.NewConn(server, c.connectTimeout, c.readTimeout, c.writeTimeout); err != nil {
				return nil, err
			}
			return conn, err
		}

		testOnBorrow func(connection.Conn, time.Time) error = func(conn connection.Conn, borrow time.Time) (err error) {
			if borrow.Add(c.testOnBorrowTimeout).Before(time.Now()) {
				if _, err = conn.HandelMessage(protocol.PING); err != nil {
					return
				}
			}
			return
		}
	)
	connPool = connection.NewPool(c.maxIdle, newFn, testOnBorrow)

	c.Lock()
	c.clusterNodePool[server] = connPool
	c.Unlock()
	return
}

func (c *cluster) GetSlot(slot uint16) (nodes nodes.Nodes, ok bool) {
	return c.clusterSlots.GetSlot(slot)
}
