package diskset

import (
	"sync"

	"github.com/fferrandis/simu/scality/hyperdrive/disk"
)

type DiskSet struct {
	current_disk int
	disk         []*disk.Disk
	sync.Mutex
}

func New(number_of_disk int, model_disk *disk.Disk) *DiskSet {
	ds := &DiskSet{}
	ds.disk = make([]*disk.Disk, number_of_disk)
	for i := 0; i < number_of_disk; i++ {
		ds.disk[i] = model_disk.Dup()
	}
	return ds
}

func (ds *DiskSet) Add(disk *disk.Disk) {
	ds.Lock()
	defer ds.Unlock()
	ds.disk = append(ds.disk, disk)
}

func (ds *DiskSet) Select(datas []*disk.Disk,
	codings []*disk.Disk,
	nrdata int,
	nrcoding int) bool {
	r := true

	ds.Lock()
	defer ds.Unlock()

	if nrdata+nrcoding > len(ds.disk) {
		return false
	}

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
	return r
}
