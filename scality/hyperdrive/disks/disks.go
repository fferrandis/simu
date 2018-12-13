package disks

import (
	hdcfg "github.com/fferrandis/simu/scality/hyperdrive/config"
	"sync"
)

type Disk struct {
	capacity, used, totalw, totalr uint64
	load                           float64
	lastts                         uint64
	write_speed                    uint64
	read_speed                     uint64
	mutex                          sync.Mutex
}

func (this *Disk) DiskUsageGet() (uint64, float64) {
	r := uint64(0)
	l := float64(0)
	this.mutex.Lock()
	{
		r = this.used
		l = this.load
	}
	this.mutex.Unlock()
	return r, l
}

func dataputtoload(datalen uint64, write_speed uint64) float64 {
	return float64(datalen) / float64(write_speed)
}

func datagettoload(datalen uint64, read_speed uint64) float64 {
	return float64(datalen) / float64(read_speed)
}

func (this *Disk) settime(ts uint64) {
	delta := ts - this.lastts
	this.lastts = ts

	delta_float := float64(delta) / float64(1000000000)
	if delta_float > this.load {
		this.load = 0
	} else {
		this.load -= delta_float
	}
}

/* Put data and time needed to ensure that request is done  */
func (this *Disk) PutData(datalen uint64, ts uint64) (bool, uint64) {
	retb := false
	retload := float64(0)

	this.mutex.Lock()
	{
		/* flush data */
		this.settime(ts)
		this.load = this.load + dataputtoload(datalen, this.write_speed)
		retb = true
		retload = this.load
	}
	this.mutex.Unlock()
	retload = retload * 1000000000
	return retb, uint64(retload)
}

func (this *Disk) GetData(datalen uint64, ts uint64) (bool, uint64) {
	retb := true
	retload := float64(0)

	this.mutex.Lock()
	{
		this.settime(ts)
		this.load = this.load + datagettoload(datalen, this.read_speed)
		retload = this.load
	}
	this.mutex.Unlock()
	retload = retload * 1000000000
	return retb, uint64(retload)
}

func (this *Disk) SetTime(ts uint64) uint64 {
	_, retload := this.GetData(0, ts)

	return retload
}

func (this *Disk) NewFile(filelen uint64, ts uint64) (bool, uint64) {
	/* for now we consider that a create operation does not bring any extra cost */
	/* XXX add config for "create" syscall cost in order to check the impact of extentsize */
	retb := true
	retload := float64(0)

	this.mutex.Lock()
	{
		if this.used+filelen > this.capacity {
			retb = false
		} else {
			this.used += filelen
		}
		this.settime(ts)
		retload = this.load
	}
	this.mutex.Unlock()

	return retb, uint64(retload * 1000000000)

}

func New(d hdcfg.DiskCfg, ts_create uint64) []*Disk {
	nr_instance := d.Nr_instances
	if nr_instance == 0 {
		nr_instance = 1
	}
	r := make([]*Disk, 0, nr_instance)

	for i := 0; i < nr_instance; i++ {
		r = append(r, &Disk{
			capacity:    d.Capacity,
			used:        0,
			totalw:      0,
			totalr:      0,
			load:        0,
			lastts:      ts_create,
			write_speed: d.Write_speed,
			read_speed:  d.Read_speed,
		})
	}
	return r
}
