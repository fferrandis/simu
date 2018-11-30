package diskset

import (
	"sync"

	. "github.com/fferrandis/simu/scality/hyperdrive/disks"
)

type DiskSet struct {
	current_disk int
	disk         []*Disk
	sync.Mutex
}

func NewDiskSet(number_of_disk int, model_disk *Disk) *DiskSet {
	ds := &DiskSet{}
	ds.disk = make([]*Disk, number_of_disk)
	for i := 0; i < number_of_disk; i++ {
		ds.disk[i] = model_disk
	}
	return ds
}

func (ds *DiskSet) DiskSetAdd(disk *Disk) {
	ds.Lock()
	defer ds.Unlock()
	ds.disk = append(ds.disk, disk)
}

func (ds *DiskSet) DiskSetSelect(datas []*Disk,
	codings []*Disk,
	nrdata int,
	nrcoding int) bool {
	r := true

	ds.Lock()
	defer ds.Unlock()

	if nrdata+nrcoding > len(ds.disk) {
		r = false
	} else {
		for i := 0; i < nrdata; i++ {
			if ds.current_disk >= len(ds.disk) {
				ds.current_disk = 0
			}
			datas[i] = ds.disk[ds.current_disk]
			ds.current_disk++
		}
		for i := 0; i < nrcoding; i++ {
			if ds.current_disk >= len(ds.disk) {
				ds.current_disk = 0
			}
			codings[i] = ds.disk[ds.current_disk]
			ds.current_disk++

		}
	}

	return r
}
