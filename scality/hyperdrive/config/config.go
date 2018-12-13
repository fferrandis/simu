package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
)

const (
	RR   = "RR"
	LOAD = "LOAD"
)

const (
	D_SELECT_ALG        = RR
	READSPEED           = 200000000 // 200Mo/s
	WRITESPEED          = 100000000 // 100Mo/s
	CODINGSCHEME        = 2
	DATASCHEME          = 4
	EXTENTSIZE          = 134200000 // 128Mio
	NETWORKTHR          = 125000000 // 1Gb/s
	NRDISKDEFAULT       = 64
	CAPACITYDISKDEFAULT = 5000000000 // 5Giga
	NRSRVDEFAULT        = 3
)

type DiskCfg struct {
	Capacity     uint64
	Write_speed  uint64
	Read_speed   uint64
	Nr_instances int
}

/* global configuration */
type HdSrvCfg struct {
	Name       string
	Diskconfig []DiskCfg
}

type HdCfg struct {
	Disk_selection_algorithm string
	Extent_size              uint64
	Data_scheme              int
	Coding_scheme            int
	Network_bdwidth          uint64
	Hdservers                []HdSrvCfg
}

var HDCFG = HdCfg{
	Disk_selection_algorithm: RR,
	Extent_size:              EXTENTSIZE,
	Data_scheme:              DATASCHEME,
	Coding_scheme:            CODINGSCHEME,
	Network_bdwidth:          NETWORKTHR,
	Hdservers:                make([]HdSrvCfg, 0),
}

func HDCfgDefault() {
	fmt.Println("no readable configuration, use default settings")

	model := DiskCfg{
		Capacity:     CAPACITYDISKDEFAULT,
		Write_speed:  WRITESPEED,
		Read_speed:   READSPEED,
		Nr_instances: NRDISKDEFAULT,
	}

	for i := uint64(0); i < NRSRVDEFAULT; i++ {
		hdsrv := HdSrvCfg{
			Name:       "hserver" + strconv.FormatUint(i, 10),
			Diskconfig: make([]DiskCfg, 0),
		}
		hdsrv.Diskconfig = append(hdsrv.Diskconfig, model)
		HDCFG.Hdservers = append(HDCFG.Hdservers, hdsrv)
	}
}

func HDCfgLoad(filename *string) {

	if filename == nil {
		HDCfgDefault()
	} else {
		file, err := os.Open(*filename)
		if err != nil {
			HDCfgDefault()
		} else {
			content, _ := ioutil.ReadAll(file)
			err := json.Unmarshal(content, &HDCFG)
			if err != nil {
				fmt.Println("cannot parse json file ", *filename, "cause : ", err)
				HDCfgDefault()
			}
			file.Close()
		}
	}
	fmt.Println("configuration used ", HDCFG)

}
