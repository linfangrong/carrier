package connection

import (
	"container/list"
	"fmt"
	"sync"
	"time"
)

var (
	errPoolClosed error = fmt.Errorf("pool: connection pool closed")
)

type Pool interface {
	Get() (conn Conn, err error)
	Put(conn Conn) error
	Remove(conn Conn) error
	Close() error
}

type pool struct {
	MaxIdle      int
	NewFn        func() (Conn, error)
	TestOnBorrow func(Conn, time.Time) error

	cond    *sync.Cond
	idleNum int
	idle    *list.List
	closed  bool
}

type idleItem struct {
	idleTime time.Time
	idleConn Conn
}

func NewPool(
	maxIdle int,
	newFn func() (Conn, error),
	testOnBorrow func(Conn, time.Time) error,
) Pool {
	return &pool{
		MaxIdle:      maxIdle,
		NewFn:        newFn,
		TestOnBorrow: testOnBorrow,
		cond:         sync.NewCond(&sync.Mutex{}),
		idle:         list.New(),
	}
}

func (p *pool) Get() (conn Conn, err error) {
	p.cond.L.Lock()
	for {
		// 获取idle连接
		for i, n := 0, p.idle.Len(); i < n; i++ {
			var e *list.Element
			if e = p.idle.Front(); e == nil {
				break
			}
			p.idle.Remove(e)
			p.cond.L.Unlock()
			var (
				item         idleItem                    = e.Value.(idleItem)
				TestOnBorrow func(Conn, time.Time) error = p.TestOnBorrow
			)
			if TestOnBorrow == nil || TestOnBorrow(item.idleConn, item.idleTime) == nil {
				return item.idleConn, nil
			}
			item.idleConn.Close()
			p.cond.L.Lock()
			p.idleNum--
			p.cond.Signal()
		}
		// 检测是否已被关闭
		if p.closed {
			p.cond.L.Unlock()
			return nil, errPoolClosed
		}
		// 限制
		if p.MaxIdle == 0 || p.idleNum < p.MaxIdle {
			p.idleNum++
			var NewFn func() (Conn, error) = p.NewFn
			p.cond.L.Unlock()
			if conn, err = NewFn(); err != nil {
				p.cond.L.Lock()
				p.idleNum--
				p.cond.Signal()
				p.cond.L.Unlock()
				return nil, err
			}
			return conn, err
		}
		p.cond.Wait()
	}
}

func (p *pool) Put(conn Conn) error {
	p.cond.L.Lock()
	if p.closed {
		p.cond.L.Unlock()
		return errPoolClosed
	}
	p.idle.PushFront(idleItem{idleTime: time.Now(), idleConn: conn})
	if p.idle.Len() > p.MaxIdle {
		if conn = p.idle.Remove(p.idle.Back()).(idleItem).idleConn; conn != nil {
			p.idleNum--
		}
	} else {
		conn = nil
	}
	p.cond.Signal()
	p.cond.L.Unlock()
	if conn != nil {
		conn.Close()
	}
	return nil
}

func (p *pool) Remove(conn Conn) error {
	p.cond.L.Lock()
	if p.closed {
		p.cond.L.Unlock()
		return errPoolClosed
	}
	p.idleNum--
	p.cond.Signal()
	p.cond.L.Unlock()
	if conn != nil {
		conn.Close()
	}
	return nil
}

func (p *pool) Close() error {
	p.cond.L.Lock()
	if p.closed {
		p.cond.L.Unlock()
		return errPoolClosed
	}

	var idle *list.List = p.idle
	p.closed = true
	p.idle.Init()
	p.cond.Broadcast()
	p.cond.L.Unlock()
	for e := idle.Front(); e != nil; e = e.Next() {
		e.Value.(idleItem).idleConn.Close()
	}
	return nil
}

func (p *pool) IdleNum() int {
	p.cond.L.Lock()
	var idleNum int = p.idleNum
	p.cond.L.Unlock()
	return idleNum
}

func (p *pool) Len() int {
	p.cond.L.Lock()
	var l int = p.idle.Len()
	p.cond.L.Unlock()
	return l
}
