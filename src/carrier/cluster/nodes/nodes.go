package nodes

import (
	"carrier/cluster/connection"
	"math/rand"
)

type Nodes interface {
	Set(data []connection.Pool) Nodes
	GetMaster() connection.Pool
	GetRandom() (connection.Pool, bool)
}

type nodes struct {
	data []connection.Pool
}

func NewNodes() Nodes {
	return &nodes{}
}

func (n *nodes) Set(data []connection.Pool) Nodes {
	n.data = data
	return n
}

func (n *nodes) GetMaster() connection.Pool {
	return n.data[0]
}

func (n *nodes) GetRandom() (connection.Pool, bool) {
	var random int = rand.Intn(len(n.data))
	return n.data[random], random == 0
}
