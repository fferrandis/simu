package disk

import (
	"sync"
)

type Disk struct {
	capacity, used, totalw, totalr uint64
	load                           float64
	lastts                         uint64
	write_speed                    uint64
	read_speed                     uint64
	sync.Mutex
}

func (d *Disk) GetUsage() uint64 {
	d.Lock()
	defer d.Unlock()
	return d.used
}

func (d *Disk) Dup() *Disk {
	return New(d.capacity, d.write_speed, d.read_speed, d.lastts)
}

func dataputtoload(datalen uint64, write_speed uint64) float64 {
	return float64(datalen) / float64(write_speed)
}

func datagettoload(datalen uint64, read_speed uint64) float64 {
	return float64(datalen) / float64(read_speed)
}

func (d *Disk) settime(ts uint64) {
	delta := ts - d.lastts
	d.lastts = ts

	delta_float := float64(delta) / float64(1000000000)
	if delta_float > d.load {
		d.load = 0
	} else {
		d.load -= delta_float
	}
}

// PutData and time needed to ensure that request is done
func (d *Disk) PutData(datalen uint64, ts uint64) (bool, uint64) {
	retb := false
	var retload float64

	d.Lock()
	defer d.Unlock()

	/* flush data */
	d.settime(ts)
	d.load = d.load + dataputtoload(datalen, d.write_speed)
	retb = true
	retload = d.load

	retload = retload * 1000000000
	return retb, uint64(retload)
}

func (d *Disk) GetData(datalen uint64, ts uint64) (bool, uint64) {
	retb := true
	retload := float64(0)

	d.Lock()
	{
		d.settime(ts)
		d.load = d.load + datagettoload(datalen, d.read_speed)
		retload = d.load
	}
	d.Unlock()
	retload = retload * 1000000000
	return retb, uint64(retload)
}

func (d *Disk) SetTime(ts uint64) uint64 {
	_, retload := d.GetData(0, ts)

	return retload
}

func (d *Disk) NewFile(filelen uint64, ts uint64) (bool, uint64) {
	/* for now we consider that a create operation does not bring any extra cost */
	/* XXX add config for "create" syscall cost in order to check the impact of extentsize */
	retb := true
	retload := float64(0)

	d.Lock()
	defer d.Unlock()

	if d.used+filelen > d.capacity {
		retb = false
	} else {
		d.used += filelen
	}
	d.settime(ts)
	retload = d.load

	return retb, uint64(retload * 1000000000)

}

func New(capacity uint64, write_speed uint64, read_speed uint64, ts_create uint64) *Disk {
	return &Disk{capacity: capacity,
		used:        0,
		totalw:      0,
		totalr:      0,
		lastts:      ts_create,
		write_speed: write_speed,
		read_speed:  read_speed,
	}
}
