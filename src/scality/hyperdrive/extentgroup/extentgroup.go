package extentgroup

import (
	"fmt"
	. "scality/hyperdrive/disks"
	"sort"
	"sync"
)

type Extent struct {
	usage, extentsize uint64
	diskref           *Disk
	lock              sync.Mutex
	id                int
}

func (this *Extent) ExtentPutData(datalen uint64, ts uint64) (bool, uint64) {
	r := true
	load := uint64(0)
	this.lock.Lock()
	{
		this.usage += datalen
		r, load = this.diskref.PutData(datalen, ts)
	}
	this.lock.Unlock()
	return r, load
}

func (this *Extent) ExtentInit(extentsize uint64, dref *Disk) {
	this.usage = 0
	this.extentsize = extentsize
	this.diskref = dref
}

func (this *Extent) ExtentUsageGet() (uint64, uint64) {
	u, s := uint64(0), uint64(0)
	this.lock.Lock()
	{
		u = this.usage
		s = this.extentsize
	}
	this.lock.Unlock()
	return u, s
}

type ExtentDataGroup struct {
	list       []Extent
	nrdata     int
	lock       sync.Mutex
	extentsize uint64
}

type ExtentCodingGroup struct {
	list       []Extent
	nrcoding   int
	extentsize uint64
}

func (a ExtentDataGroup) Len() int {
	return a.nrdata
}

func (a ExtentDataGroup) Swap(i, j int) {
	a.list[i], a.list[j] = a.list[j], a.list[i]
}

func (a ExtentDataGroup) Less(i, j int) bool {
	u1, _ := a.list[i].ExtentUsageGet()
	u2, _ := a.list[j].ExtentUsageGet()

	return u1 < u2
}

func (this *ExtentDataGroup) ExtentDataGroupPutData(datalen uint64, ts uint64) (bool, uint64) {
	r := false
	load := uint64(0)
	this.lock.Lock()
	{
		u1, s1 := this.list[0].ExtentUsageGet()
		if u1+datalen <= s1 {
			r, load = this.list[0].ExtentPutData(datalen, ts)
		}
		sort.Sort(this)
	}
	this.lock.Unlock()
	return r, load
}

func (this *ExtentDataGroup) ExtentDataGroupClose(ts uint64) {
	this.lock.Lock()
	{
		for i := 0; i < this.nrdata; i++ {
			refdisk := this.list[i].diskref
			/* add read on data disk for ece */
			refdisk.GetData(this.extentsize, ts)
		}
	}
	this.lock.Unlock()
}

func (this *ExtentCodingGroup) ExtentCodingGroupClose(ts uint64) {
	for i := 0; i < this.nrcoding; i++ {
		refdisk := this.list[i].diskref
		ok, _ := refdisk.PutData(this.extentsize, ts)
		if ok != true {
			fmt.Println("failed on ", refdisk)
		}
	}
}

func (this *ExtentDataGroup) ExtentDataGroupInit(disk []*Disk,
	nrdata int, extentsize uint64, ts uint64) {
	if cap(this.list) < nrdata {
		this.list = make([]Extent, nrdata)
	}
	this.nrdata = nrdata
	this.extentsize = extentsize
	for i := 0; i < nrdata; i++ {
		this.list[i].diskref = disk[i]
		this.list[i].usage = 0
		this.list[i].extentsize = extentsize
		this.list[i].id = i

		/* create file if we can */
		ret, _ := disk[i].NewFile(extentsize, ts)
		if ret != true {
			panic("disk overflow")
		}
	}
}

func (this *ExtentCodingGroup) ExtentCodingGroupInit(disk []*Disk,
	nrcoding int, extentsize uint64, ts uint64) {
	if cap(this.list) < nrcoding {
		this.list = make([]Extent, nrcoding)
	}

	this.nrcoding = nrcoding
	this.extentsize = extentsize
	for i := 0; i < nrcoding; i++ {
		this.list[i].diskref = disk[i]
		this.list[i].usage = 0
		this.list[i].extentsize = extentsize
		this.list[i].id = i
		ret, _ := disk[i].NewFile(extentsize, ts)
		if ret != true {
			panic("disk overflow")
		}
	}
}
