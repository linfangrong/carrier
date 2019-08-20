package datatypes

import (
	"sync"
)

type SyncClose struct {
	sync.Once
	closed chan bool
}

func NewSyncClose() *SyncClose {
	return &SyncClose{
		closed: make(chan bool),
	}
}

func (sc *SyncClose) Close() {
	sc.Do(func() {
		close(sc.closed)
	})
}

func (sc *SyncClose) IsClosed() <-chan bool {
	return sc.closed
}
