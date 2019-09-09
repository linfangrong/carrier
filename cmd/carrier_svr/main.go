package main

import (
	"flag"
	"log"
	"net/http"
	_ "net/http/pprof"
	"runtime"
	"time"

	"github.com/go-ini/ini"

	"github.com/linfangrong/carrier/models/cluster/cluster"
	_ "github.com/linfangrong/carrier/models/command/set"
	"github.com/linfangrong/carrier/models/logger"
	"github.com/linfangrong/carrier/utils/util"
)

var (
	confpath      = flag.String("conf", util.ExecDir()+"/../conf/carrier.ini", "配置文件路径")
	cfg           *ini.File
	carrierLogger *util.Logger
)

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	flag.Parse()

	var (
		err error
	)
	if cfg, err = ini.Load(*confpath); err != nil {
		panic(err)
	}

	if cfg.Section("PPROF").Key("UsePprof").MustBool(false) {
		go func() {
			log.Println(http.ListenAndServe(cfg.Section("PPROF").Key("PprofAddr").MustString(":6063"), nil))
		}()
	}
}

func main() {
	carrierLogger = util.NewLogger(cfg.Section("LOG").Key("Path").MustString(util.HomeDir()+"/logs/carrier.log"), util.FilenameSuffixInDay)
	logger.SetLogger(carrierLogger)

	var (
		clusterServerList []string
		err               error
		addr              *ini.Key
	)
	for _, addr = range cfg.Section("CLUSTER.HOST").Keys() {
		clusterServerList = append(clusterServerList, addr.Value())
	}

	// 初始化 cluster
	if err = cluster.InitClusterParameter(
		clusterServerList,
		time.Duration(cfg.Section("CLUSTER").Key("RefreshInterval").MustInt(300))*time.Second,
		cfg.Section("CLUSTER").Key("MaxIdle").MustInt(50),
		time.Duration(cfg.Section("CLUSTER").Key("TestOnBorrowTimeout").MustInt(150))*time.Second,
		time.Duration(cfg.Section("CLUSTER").Key("ConnectTimeout").MustInt(3))*time.Second,
		time.Duration(cfg.Section("CLUSTER").Key("ReadTimeout").MustInt(3))*time.Second,
		time.Duration(cfg.Section("CLUSTER").Key("WriteTimeout").MustInt(3))*time.Second,
	); err != nil {
		panic(err)
	}

	// 服务端口 client
	ServeForClient(cfg.Section("").Key("ClientAddr").MustString(":6679"))
}
