package main

import (
	"flag"
	"fmt"
	hdcfg "scality/hyperdrive/config"
	"scality/hyperdrive/hdio"
)

func main() {
	filename := flag.String("config", "", "")
	flag.Parse()

	hdcfg.HDCfgLoad(filename)
	fmt.Println("started simulator")

	hdio.HDIoStart()
}
