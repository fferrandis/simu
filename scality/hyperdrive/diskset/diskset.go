package diskset

import (
	hdcfg "github.com/fferrandis/simu/scality/hyperdrive/config"
	hddisk "github.com/fferrandis/simu/scality/hyperdrive/disks"
	. "github.com/fferrandis/simu/scality/hyperdrive/diskstat"
	"sync"
)

type DiskSet struct {
	current_disk int
	disk         []*hddisk.Disk
	sync.Mutex
}

func New(cfg []hdcfg.DiskCfg, now uint64) *DiskSet {
	ds := &DiskSet{}

	ds.disk = make([]*hddisk.Disk, 0)
	for _, d := range cfg {
		ds.disk = append(ds.disk, hddisk.New(d, now)...)
	}

	return ds
}

func (ds *DiskSet) DiskSetAdd(disk *hddisk.Disk) {
	ds.Lock()
	defer ds.Unlock()
	ds.disk = append(ds.disk, disk)
}

func (ds *DiskSet) DiskSetSelect(datas []*hddisk.Disk,
	codings []*hddisk.Disk,
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

func (ds *DiskSet) DiskSetStatsGet(ts uint64) []DiskStat {
	r := make([]DiskStat, len(ds.disk))

	for i, disk := range ds.disk {
		r[i] = disk.DiskStatsGet(ts)
	}
	return r
}
