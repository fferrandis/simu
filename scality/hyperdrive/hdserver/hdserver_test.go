package hdserver

import (
	"fmt"
	"testing"

	. "github.com/fferrandis/simu/scality/hyperdrive/config"
)

func TestServer(t *testing.T) {

	cfg := HdCfg{
		Extent_size:     1000,
		Data_scheme:     4,
		Coding_scheme:   2,
		Network_bdwidth: 10000,
		Hdservers: []HdSrvCfg{HdSrvCfg{Name: "hdsrv1",
			Diskconfig: []DiskCfg{DiskCfg{Capacity: 2000,
				Write_speed:  1000,
				Read_speed:   2000,
				Nr_instances: 6}}}}}

	hd := New(cfg.Extent_size, cfg.Data_scheme, cfg.Coding_scheme, cfg.Hdservers[0], 0)

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
