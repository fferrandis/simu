package disk

import (
	"fmt"
	"testing"
)

func TestPutData(t *testing.T) {
	v := Disk{capacity: 2000000000, // 200 Gb
		used:        0,
		totalw:      0,
		totalr:      0,
		load:        0,
		lastts:      0,
		write_speed: 104857600,
		read_speed:  104857600,
	}

	/* at time 0, we inject 2Mb*/
	r, load := v.PutData(2097152, 0)
	if r != true {
		t.Error("expected success")
	} else {
		fmt.Println("request should be done in ", load, " nsec")
	}
	load = v.SetTime(load / 2)
	fmt.Println("we advanced so load left is ", load, " nsec")

	load = v.SetTime(load)
	fmt.Println("we advanced so load left is ", load, " nsec")
}
