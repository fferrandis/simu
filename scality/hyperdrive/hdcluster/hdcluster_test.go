package hdcluster

import (
	"fmt"
	"testing"

	cfg "github.com/fferrandis/simu/scality/hyperdrive/config"
)

func TestCluster(t *testing.T) {
	/* insert 2 new servers with 6 disks eachs, whose size is 2x an extent size*/
	HDClusterSrvAdd(6, 2*cfg.HDCFG.Extent_size)
	HDClusterSrvAdd(6, 2*cfg.HDCFG.Extent_size)

	for i := 0; i < 18; i++ {
		r, load := HDClusterSrvPut(cfg.HDCFG.Extent_size / 2)
		fmt.Println("iter =", i, ";ret =", r, ";load =", load)
	}

}
