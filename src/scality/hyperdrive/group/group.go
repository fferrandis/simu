package group

import (
	. "scality/hyperdrive/disks"
	. "scality/hyperdrive/extentgroup"
	"sync"
)

type DiskRefSet []*Disk

type Group struct {
	extentd  ExtentDataGroup
	extentc  ExtentCodingGroup
	nrdata   int
	nrcoding int
	lock     sync.Mutex
}

func (this *Group) PutData(datalen uint64, ts uint64) (bool, uint64) {
	r := false
	load := uint64(0)
	this.lock.Lock()
	{
		r, load = this.extentd.ExtentDataGroupPutData(datalen, ts)
		if r == false {
			/* close extents */
			this.extentd.ExtentDataGroupClose(ts)
			this.extentc.ExtentCodingGroupClose(ts)
		}
	}
	this.lock.Unlock()

	return r, load
}

func (this *Group) GroupInit(diskdata []*Disk, diskc []*Disk, esize uint64, nrdata int, nrcoding int, ts uint64) {
	this.lock.Lock()
	{
		this.nrdata = nrdata
		this.nrcoding = nrcoding
		this.extentd.ExtentDataGroupInit(diskdata, nrdata, esize, ts)
		this.extentc.ExtentCodingGroupInit(diskc, nrcoding, esize, ts)
	}
	this.lock.Unlock()

}
