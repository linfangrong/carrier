package slots

import (
	"sync"

	"github.com/linfangrong/carrier/models/cluster/nodes"
)

type Slots interface {
	AddSlot(slot int64, nodes nodes.Nodes)
	AddSlots(begin int64, end int64, nodes nodes.Nodes)
	GetSlot(slot uint16) (nodes nodes.Nodes, ok bool)
	GetSlotsCount() int
}

type slots struct {
	sync.RWMutex
	data map[uint16]nodes.Nodes
}

func NewSlots() Slots {
	return &slots{
		data: make(map[uint16]nodes.Nodes),
	}
}

func (s *slots) AddSlot(slot int64, nodes nodes.Nodes) {
	s.Lock()
	s.data[uint16(slot)] = nodes
	s.Unlock()
}

func (s *slots) AddSlots(begin int64, end int64, nodes nodes.Nodes) {
	for ; begin <= end; begin++ {
		s.AddSlot(begin, nodes)
	}
}

func (s *slots) GetSlot(slot uint16) (nodes nodes.Nodes, ok bool) {
	s.RLock()
	nodes, ok = s.data[slot]
	s.RUnlock()
	return
}

func (s *slots) GetSlotsCount() (l int) {
	s.RLock()
	l = len(s.data)
	s.RUnlock()
	return
}
