package hdserver

import (
	"fmt"
	"sync"

	"github.com/fferrandis/simu/scality/hyperdrive/config"
	"github.com/fferrandis/simu/scality/hyperdrive/disks"
	"github.com/fferrandis/simu/scality/hyperdrive/diskset"
	"github.com/fferrandis/simu/scality/hyperdrive/group"
)

type HDSrv struct {
	dset *diskset.DiskSet

	group group.Group

	nrdata     int
	nrcoding   int
	extentsize uint64
	name       string
	sync.Mutex
}

func (hdsrv *HDSrv) HDSrvGroupInit(ts uint64) bool {
	datadisk := make([]*disks.Disk, hdsrv.nrdata)
	codingdisk := make([]*disks.Disk, hdsrv.nrcoding)
	hdsrv.dset.DiskSetSelect(datadisk, codingdisk, hdsrv.nrdata, hdsrv.nrcoding)
	hdsrv.group.GroupInit(datadisk, codingdisk, hdsrv.extentsize, hdsrv.nrdata, hdsrv.nrcoding, ts)

	return true
}

func New(extentsize uint64, nrd int, nrc int, srvcfg config.HdSrvCfg, ts uint64) *HDSrv {

	h := &HDSrv{
		dset:       diskset.New(srvcfg.Diskconfig, ts),
		nrdata:     nrd,
		nrcoding:   nrc,
		extentsize: extentsize,
		name:       srvcfg.Name,
	}

	h.HDSrvGroupInit(ts)

	return h
}

func (hdsrv *HDSrv) HDSrvPutData(datalen uint64, ts uint64) (bool, uint64) {
	r := true
	load := uint64(0)
	hdsrv.Lock()
	defer hdsrv.Unlock()
	r, load = hdsrv.group.PutData(datalen, ts)
	if r == false {
		fmt.Println("Group full, create new one")
		hdsrv.HDSrvGroupInit(ts)
		r, load = hdsrv.group.PutData(datalen, ts)
	}

	return r, load
}
