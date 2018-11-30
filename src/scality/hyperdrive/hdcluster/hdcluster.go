package hdcluster

import (
	"fmt"
	cfg "scality/hyperdrive/config"
	. "scality/hyperdrive/disks"
	. "scality/hyperdrive/hdserver"
	"sync"
)

type HDCluster struct {
	srvs     []HDSrv
	srvcurr  int
	lock     sync.Mutex
	totallen uint64
}

var cluster HDCluster

func bytes2ts(totallen uint64) uint64 {
	p := float64(totallen) / float64(cfg.HDCFG.Network_bdwidth)

	return uint64(p * 1000000000)
}

func HDClusterSrvAdd(nrdisk int, capacity uint64) {
	now := bytes2ts(cluster.totallen)
	d := DiskNew(capacity, cfg.HDCFG.Write_speed, cfg.HDCFG.Read_speed, now)
	var newsrv HDSrv
	newsrv.HDSrvInit(cfg.HDCFG.Data_scheme,
		cfg.HDCFG.Coding_scheme,
		cfg.HDCFG.Extent_size,
		d,
		nrdisk,
		now)
	cluster.lock.Lock()
	{
		cluster.srvs = append(cluster.srvs, newsrv)
	}
	cluster.lock.Unlock()
}

func init() {
	cluster.srvs = make([]HDSrv, 0, 8)
	cluster.srvcurr = 0
}

func selectHD() *HDSrv {
	var p *HDSrv = nil
	l := len(cluster.srvs)
	if l > 0 {
		if cluster.srvcurr >= l {
			cluster.srvcurr = 0
		}
		p = &cluster.srvs[cluster.srvcurr]
	}
	return p
}

func HDClusterSrvPut(datalen uint64) (bool, uint64) {
	var p *HDSrv = nil
	ts := uint64(0)
	cluster.lock.Lock()
	{
		ts = bytes2ts(cluster.totallen)
		p = selectHD()
		cluster.totallen += datalen
	}
	cluster.lock.Unlock()

	if p == nil {
		fmt.Println("cannot put data on cluster since no servers are running")
		return false, 0
	}
	r, load := p.HDSrvPutData(datalen, ts)
	return r, load
}
