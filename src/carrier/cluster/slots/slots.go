package slots

import (
	"carrier/cluster/nodes"
	"sync"
)

type Slots interface {
	AddSlot(slot int64, nodes nodes.Nodes)
	AddSlots(begin int64, end int64, nodes nodes.Nodes)
	GetSlot(slot uint16) (nodes nodes.Nodes, ok bool)
	GetSlotsCount() int
}

type slots struct {
	sync.Map
}

func NewSlots() Slots {
	return &slots{}
}

func (s *slots) AddSlot(slot int64, nodes nodes.Nodes) {
	s.Store(uint16(slot), nodes)
}

func (s *slots) AddSlots(begin int64, end int64, nodes nodes.Nodes) {
	for ; begin <= end; begin++ {
		s.AddSlot(begin, nodes)
	}
}

func (s *slots) GetSlot(slot uint16) (_nodes nodes.Nodes, ok bool) {
	var value interface{}
	if value, ok = s.Load(slot); ok {
		_nodes = value.(nodes.Nodes)
	}
	return
}

func (s *slots) GetSlotsCount() (l int) {
	// TODO
	return 0
}
