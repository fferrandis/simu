package extentgroup

import (
	"fmt"
	"sort"
	"sync"

	. "github.com/fferrandis/simu/scality/hyperdrive/disks"
)

type Extent struct {
	usage, extentsize uint64
	diskref           *Disk
	sync.Mutex
	id int
}

func (e *Extent) ExtentPutData(datalen uint64, ts uint64) (bool, uint64) {
	r := true
	load := uint64(0)
	e.Lock()
	defer e.Unlock()

	e.usage += datalen
	r, load = e.diskref.PutData(datalen, ts)

	return r, load
}

func NewExtent(size uint64, dref *Disk, id int) *Extent {
	return &Extent{extentsize: size, diskref: dref, id: id}
}

func (e *Extent) ExtentUsageGet() (uint64, uint64) {
	u, s := uint64(0), uint64(0)
	e.Lock()
	defer e.Unlock()
	u = e.usage
	s = e.extentsize
	return u, s
}

type ExtentDataGroup struct {
	list   []*Extent
	nrdata int
	sync.Mutex
	extentsize uint64
}

type ExtentCodingGroup struct {
	list       []*Extent
	nrcoding   int
	extentsize uint64
}

func (a *ExtentDataGroup) Len() int {
	return a.nrdata
}

func (a *ExtentDataGroup) Swap(i, j int) {
	a.list[i], a.list[j] = a.list[j], a.list[i]
}

func (a *ExtentDataGroup) Less(i, j int) bool {
	u1, _ := a.list[i].ExtentUsageGet()
	u2, _ := a.list[j].ExtentUsageGet()

	return u1 < u2
}

func (e *ExtentDataGroup) ExtentDataGroupPutData(datalen uint64, ts uint64) (bool, uint64) {
	r := false
	load := uint64(0)
	e.Lock()
	defer e.Unlock()

	u1, s1 := e.list[0].ExtentUsageGet()
	if u1+datalen <= s1 {
		r, load = e.list[0].ExtentPutData(datalen, ts)
	}
	sort.Sort(e)

	return r, load
}

func (e *ExtentDataGroup) ExtentDataGroupClose(ts uint64) {
	e.Lock()
	defer e.Unlock()

	for i := 0; i < e.nrdata; i++ {
		refdisk := e.list[i].diskref
		/* add read on data disk for ece */
		refdisk.GetData(e.extentsize, ts)
	}
}

func (e *ExtentCodingGroup) ExtentCodingGroupClose(ts uint64) {
	for i := 0; i < e.nrcoding; i++ {
		refdisk := e.list[i].diskref
		ok, _ := refdisk.PutData(e.extentsize, ts)
		if !ok {
			fmt.Println("failed on ", refdisk)
		}
	}
}

func (e *ExtentDataGroup) ExtentDataGroupInit(disk []*Disk,
	nrdata int, extentsize uint64, ts uint64) {
	if cap(e.list) < nrdata {
		e.list = make([]*Extent, nrdata)
	}

	e.nrdata = nrdata
	e.extentsize = extentsize
	for i := 0; i < nrdata; i++ {
		e.list[i] = NewExtent(extentsize, disk[i], i)
		/* create file if we can */
		ret, _ := disk[i].NewFile(extentsize, ts)
		if ret != true {
			panic("disk overflow")
		}
	}
}

func (e *ExtentCodingGroup) ExtentCodingGroupInit(disk []*Disk,
	nrcoding int, extentsize uint64, ts uint64) {
	if cap(e.list) < nrcoding {
		e.list = make([]*Extent, nrcoding)
	}

	e.nrcoding = nrcoding
	e.extentsize = extentsize
	for i := 0; i < nrcoding; i++ {
		e.list[i] = NewExtent(extentsize, disk[i], i)
		ret, _ := disk[i].NewFile(extentsize, ts)
		if ret != true {
			panic("disk overflow")
		}
	}
}
