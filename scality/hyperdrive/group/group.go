package group

import (
	"sync"

	. "github.com/fferrandis/simu/scality/hyperdrive/disks"
	. "github.com/fferrandis/simu/scality/hyperdrive/extentgroup"
)

type DiskRefSet []*Disk

type Group struct {
	extentd  ExtentDataGroup
	extentc  ExtentCodingGroup
	nrdata   int
	nrcoding int
	sync.Mutex
}

func (grp *Group) PutData(datalen uint64, ts uint64) (bool, uint64) {
	r := false
	load := uint64(0)
	grp.Lock()
	defer grp.Unlock()

	r, load = grp.extentd.ExtentDataGroupPutData(datalen, ts)
	if r == false {
		/* close extents */
		grp.extentd.ExtentDataGroupClose(ts)
		grp.extentc.ExtentCodingGroupClose(ts)
	}

	return r, load
}

func (grp *Group) GroupInit(diskdata []*Disk, diskc []*Disk, esize uint64, nrdata int, nrcoding int, ts uint64) {
	grp.Lock()
	defer grp.Unlock()

	grp.nrdata = nrdata
	grp.nrcoding = nrcoding
	grp.extentd.ExtentDataGroupInit(diskdata, nrdata, esize, ts)
	grp.extentc.ExtentCodingGroupInit(diskc, nrcoding, esize, ts)
}
