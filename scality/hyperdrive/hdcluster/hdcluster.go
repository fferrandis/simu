package hdcluster

import (
	"fmt"
	"sync"

	"github.com/fferrandis/simu/scality/hyperdrive/config"
	"github.com/fferrandis/simu/scality/hyperdrive/hdserver"
)

type HDCluster struct {
	srvs    []*hdserver.HDSrv
	srvcurr int
	sync.Mutex
	totallen uint64
}

var cluster HDCluster

func bytes2ts(totallen uint64) uint64 {
	p := float64(totallen) / float64(config.HDCFG.Network_bdwidth)

	return uint64(p * 1000000000)
}

func Init(cfg config.HdCfg) {
	cluster.srvs = make([]*hdserver.HDSrv, 0)
	cluster.srvcurr = 0
	cluster.totallen = 0

	for _, srv := range cfg.Hdservers {
		cluster.srvs = append(cluster.srvs, hdserver.New(cfg.Extent_size,
			cfg.Data_scheme,
			cfg.Coding_scheme,
			srv,
			0))
	}
}

func selectHD() *hdserver.HDSrv {
	var p *hdserver.HDSrv
	l := len(cluster.srvs)
	if l > 0 {
		if cluster.srvcurr >= l {
			cluster.srvcurr = 0
		}
		p = cluster.srvs[cluster.srvcurr]
	}
	return p
}

func HDClusterSrvPut(datalen uint64) (bool, uint64) {
	var p *hdserver.HDSrv
	var ts uint64

	cluster.Lock()
	defer cluster.Unlock()

	ts = bytes2ts(cluster.totallen)
	p = selectHD()
	cluster.totallen += datalen

	if p == nil {
		fmt.Println("cannot put data on cluster since no servers are running")
		return false, 0
	}
	r, load := p.HDSrvPutData(datalen, ts)
	return r, load
}
