package group

import (
	. "scality/hyperdrive/disks"
	. "scality/hyperdrive/diskset"
	"testing"
)

func TestGroupOVerflow(t *testing.T) {
	var set DiskSet
	var model = DiskNew(2000, 10, 5, 0)

	set.DiskSetInit(12, model)
	data := make([]*Disk, 4)
	coding := make([]*Disk, 2)

	var group Group

	for count := 0; count < 2; count++ {

		set.DiskSetSelect(data, coding, 4, 2)
		/* check that we select new disk ONLY*/
		for _, disk := range data {
			load := disk.SetTime(0)
			if load != 0 {
				t.Error("this is not a new disk that has been selected")
			}
		}

		for _, disk := range coding {
			load := disk.SetTime(0)
			if load != 0 {
				t.Error("this is not a new disk that has been selected")
			}
		}
		group.GroupInit(data, coding, 1000, 4, 2, 0)
		/* fill the group */
		for i := 0; i < 4; i++ {
			ok, _ := group.PutData(1000, 0)
			if ok != true {
				t.Error("group put data failed, whereas it should not")
			}
		}
		load := coding[0].SetTime(0)
		if load != 0 {
			t.Error("coding disk has activity whereas it should not")
		}
		ok, _ := group.PutData(1000, 0)
		if ok != false {
			t.Error("group put data succeess whereas it should not")
		}

		load = coding[0].SetTime(0)
		if load == 0 {
			t.Error("coding disk has NO activity whereas it should ")
		}

	}

}
