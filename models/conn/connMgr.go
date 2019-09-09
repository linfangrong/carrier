package conn

import (
	"sort"
	"time"

	"github.com/linfangrong/carrier/models/logger"
	"github.com/linfangrong/carrier/utils/datatypes"
)

var connMgr *datatypes.SafeMap = datatypes.NewSafeMap()

func init() {
	go func() {
		for range time.Tick(time.Minute) {
			checkIdleTimeout()
		}
	}()
}

func addConn(addr string, carrierConn Conn) {
	connMgr.Set(addr, carrierConn)
}

func removeConn(addr string) {
	connMgr.Erase(addr)
}

func getConn(addr string) (carrierConn Conn, ok bool) {
	var value interface{}
	if value, ok = connMgr.Get(addr); !ok {
		return nil, ok
	}
	if carrierConn, ok = value.(Conn); !ok {
		return nil, ok
	}
	return
}

func getBroadcastConn() (carrierConnList []Conn) {
	var (
		dict        map[string]interface{} = connMgr.Clone()
		value       interface{}
		carrierConn Conn
		ok          bool
	)
	carrierConnList = make([]Conn, 0, len(dict))
	for _, value = range dict {
		if carrierConn, ok = value.(Conn); !ok {
			continue
		}
		carrierConnList = append(carrierConnList, carrierConn)
	}
	return
}

// sort
type ByActiveTime []Conn

func (c ByActiveTime) Len() int      { return len(c) }
func (c ByActiveTime) Swap(i, j int) { c[i], c[j] = c[j], c[i] }
func (c ByActiveTime) Less(i, j int) bool {
	return c[i].GetLatestActiveTime().Before(c[j].GetLatestActiveTime())
}

func checkIdleTimeout() {
	var (
		carrierConnList  []Conn = getBroadcastConn()
		carrierConn      Conn
		timeoutTimestamp time.Time = time.Now().Add(-DEFAULT_IDLE_TIME_OUT)
	)
	sort.Sort(ByActiveTime(carrierConnList))
	for _, carrierConn = range carrierConnList {
		if carrierConn.GetLatestActiveTime().After(timeoutTimestamp) {
			return
		}
		logger.Infof("客户端空闲连接超时: %s", carrierConn.GetConnAddr())
		carrierConn.Close()
	}
}
