package hdserver

import (
	"fmt"
	. "scality/hyperdrive/disks"
	. "scality/hyperdrive/diskset"
	. "scality/hyperdrive/group"
	"sync"
)

type HDSrv struct {
	dset DiskSet

	group Group

	nrdata     int
	nrcoding   int
	extentsize uint64
	lock       sync.Mutex
}

func (this *HDSrv) HDSrvGroupInit(ts uint64) bool {
	datadisk := make([]*Disk, this.nrdata)
	codingdisk := make([]*Disk, this.nrcoding)
	this.dset.DiskSetSelect(datadisk, codingdisk, this.nrdata, this.nrcoding)
	this.group.GroupInit(datadisk, codingdisk, this.extentsize, this.nrdata, this.nrcoding, ts)

	return true
}

func (this *HDSrv) HDSrvInit(nrdata int,
	nrcoding int,
	extentsize uint64,
	data Disk,
	numberof_disk int,
	ts uint64) bool {

	this.nrdata = nrdata
	this.nrcoding = nrcoding
	this.extentsize = extentsize

	this.dset.DiskSetInit(numberof_disk, data)

	this.HDSrvGroupInit(ts)

	return true
}

func (this *HDSrv) HDSrvPutData(datalen uint64, ts uint64) (bool, uint64) {
	r := true
	load := uint64(0)
	this.lock.Lock()
	{
		r, load = this.group.PutData(datalen, ts)
		if r == false {
			fmt.Println("Group full, create new one")
			this.HDSrvGroupInit(ts)
			r, load = this.group.PutData(datalen, ts)
		}
	}
	this.lock.Unlock()

	return r, load
}
