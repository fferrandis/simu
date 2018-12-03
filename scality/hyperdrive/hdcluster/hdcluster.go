package hdcluster

import (
	"fmt"
	"sync"

	cfg "github.com/fferrandis/simu/scality/hyperdrive/config"
	"github.com/fferrandis/simu/scality/hyperdrive/disk"
	"github.com/fferrandis/simu/scality/hyperdrive/hdserver"
)

type HDCluster struct {
	srvs     []*hdserver.HDSrv
	srvcurr  int
	totallen uint64
	sync.Mutex
}

var cluster = HDCluster{
	srvs:    make([]*hdserver.HDSrv, 0, 8),
	srvcurr: 0,
}

func bytes2ts(totallen uint64) uint64 {
	p := float64(totallen) / float64(cfg.HDCFG.Network_bdwidth)

	return uint64(p * 1000000000)
}

func HDClusterSrvAdd(nrdisk int, capacity uint64) {
	now := bytes2ts(cluster.totallen)
	d := disk.New(capacity, cfg.HDCFG.Write_speed, cfg.HDCFG.Read_speed, now)
	newsrv := hdserver.NewHDSrv(cfg.HDCFG.Data_scheme,
		cfg.HDCFG.Coding_scheme,
		cfg.HDCFG.Extent_size,
		d,
		nrdisk,
		now)
	cluster.Lock()
	defer cluster.Unlock()
	cluster.srvs = append(cluster.srvs, newsrv)
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
	cluster.Lock()
	defer cluster.Unlock()

	ts := bytes2ts(cluster.totallen)
	p := selectHD()
	cluster.totallen += datalen

	if p == nil {
		fmt.Println("cannot put data on cluster since no servers are running")
		return false, 0
	}
	r, load := p.PutData(datalen, ts)
	return r, load
}
