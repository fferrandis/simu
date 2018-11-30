package diskset

import (
	. "scality/hyperdrive/disks"
	"sync"
)

type DiskSet struct {
	current_disk int
	disk         []Disk
	lock         sync.Mutex
}

func (this *DiskSet) DiskSetInit(number_of_disk int,
	model_disk Disk) {
	this.disk = make([]Disk, number_of_disk)
	for i := 0; i < number_of_disk; i++ {
		this.disk[i] = model_disk
	}
}

func (this *DiskSet) DiskSetAdd(disk Disk) {
	this.lock.Lock()
	{
		this.disk = append(this.disk, disk)
	}
	this.lock.Unlock()
}

func (this *DiskSet) DiskSetSelect(datas []*Disk,
	codings []*Disk,
	nrdata int,
	nrcoding int) bool {
	r := true

	this.lock.Lock()
	{
		if nrdata+nrcoding > len(this.disk) {
			r = false
		} else {
			for i := 0; i < nrdata; i++ {
				if this.current_disk >= len(this.disk) {
					this.current_disk = 0
				}
				datas[i] = &this.disk[this.current_disk]
				this.current_disk += 1
			}
			for i := 0; i < nrcoding; i++ {
				if this.current_disk >= len(this.disk) {
					this.current_disk = 0
				}
				codings[i] = &this.disk[this.current_disk]
				this.current_disk += 1

			}
		}
	}
	this.lock.Unlock()
	return r
}
