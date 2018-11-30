package diskset

import (
	"fmt"
	"testing"

	. "github.com/fferrandis/simu/scality/hyperdrive/disks"
)

func TestSelectDisk(t *testing.T) {
	var model = DiskNew(2000000000, 104857600, 104857600, 0)

	set := NewDiskSet(8, model)
	data := make([]*Disk, 4)
	coding := make([]*Disk, 2)

	r := set.DiskSetSelect(data, coding, 4, 2)
	if r != true {
		t.Error("expected success")
	}
	fmt.Println(data)
	fmt.Println(coding)
}
