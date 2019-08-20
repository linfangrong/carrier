package cluster

import (
	"time"
)

var (
	clusterParameter Cluster
)

func InitClusterParameter(
	serverList []string, refreshInterval time.Duration,
	maxIdle int, testOnBorrowTimeout time.Duration,
	connectTimeout, readTimeout, writeTimeout time.Duration,
) (err error) {
	clusterParameter, err = NewCluster(
		serverList, refreshInterval,
		maxIdle, testOnBorrowTimeout,
		connectTimeout, readTimeout, writeTimeout,
	)
	return
}

func GetClusterParameter() Cluster {
	return clusterParameter
}
