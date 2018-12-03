package diskset

import (
	"fmt"
	"testing"

	"github.com/fferrandis/simu/scality/hyperdrive/disk"
)

func TestSelectDisk(t *testing.T) {
	var model = disk.New(2000000000, 104857600, 104857600, 0)

	set := New(8, model)
	data := make([]*disk.Disk, 4)
	coding := make([]*disk.Disk, 2)

	r := set.DiskSetSelect(data, coding, 4, 2)
	if r != true {
		t.Error("expected success")
	}
	fmt.Println(data)
	fmt.Println(coding)
}
