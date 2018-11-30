package main

import (
	"flag"
	"fmt"

	hdcfg "github.com/fferrandis/simu/scality/hyperdrive/config"
	"github.com/fferrandis/simu/scality/hyperdrive/hdio"
)

func main() {
	filename := flag.String("config", "", "")
	flag.Parse()

	hdcfg.HDCfgLoad(filename)
	fmt.Println("started simulator")

	hdio.HDIoStart()
}
