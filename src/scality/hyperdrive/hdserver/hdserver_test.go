package hdserver

import (
	"fmt"
	. "scality/hyperdrive/disks"
	"testing"
)

func TestServer(t *testing.T) {
	var hd HDSrv
	/*disk model */
	r := DiskNew(2000, 100, 100, 0)
	hd.HDSrvInit(4, 2, 1000, r, 6, 0)

	/* fill 4 extents (size = 1000) */
	fmt.Println("inject again 4x500b")
	for i := 0; i < 4; i++ {
		r, load := hd.HDSrvPutData(500, 1000000)
		fmt.Println("ret=", r, ", load=", load)
	}
	fmt.Println("reinject again 4x500b")
	for i := 0; i < 4; i++ {
		r, load := hd.HDSrvPutData(500, 2000000)
		fmt.Println("ret=", r, ", load=", load)
	}

	fmt.Println("reinject again 4x500b ==> should create a new group ")
	for i := 0; i < 4; i++ {
		r, load := hd.HDSrvPutData(500, 3000000)
		fmt.Println("ret=", r, ", load=", load)
	}

}
