package hdcluster

import (
	"fmt"
	"testing"

	cfg "github.com/fferrandis/simu/scality/hyperdrive/config"
)

func TestCluster(t *testing.T) {
	/* insert 2 new servers with 6 disks eachs, whose size is 2x an extent size*/

	dmodel := cfg.DiskCfg{cfg.CAPACITYDISKDEFAULT,
		cfg.WRITESPEED,
		cfg.READSPEED,
		6}

	srv1 := cfg.HdSrvCfg{Name: "srv1", Diskconfig: []cfg.DiskCfg{dmodel}}
	srv2 := cfg.HdSrvCfg{Name: "srv1", Diskconfig: []cfg.DiskCfg{dmodel}}

	conf := cfg.HdCfg{Extent_size: cfg.EXTENTSIZE,
		Data_scheme:     cfg.DATASCHEME,
		Coding_scheme:   cfg.CODINGSCHEME,
		Network_bdwidth: cfg.NETWORKTHR,
		Hdservers:       []cfg.HdSrvCfg{srv1, srv2}}

	Init(conf)

	for i := 0; i < 18; i++ {
		r, load := HDClusterSrvPut(cfg.EXTENTSIZE / 2)
		fmt.Println("iter =", i, ";ret =", r, ";load =", load)
	}
}
