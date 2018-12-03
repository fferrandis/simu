package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

const (
	READSPEED           = 200000000 // 200Mo/s
	WRITESPEED          = 100000000 // 100Mo/s
	CODINGSCHEME        = 2
	DATASCHEME          = 4
	EXTENTSIZE          = 134200000 // 128Mio
	NETWORKTHR          = 125000000 // 1Gb/s
	NRDISKDEFAULT       = 64
	CAPACITYDISKDEFAULT = 5000000000 // 5Giga

)

type HdSrvCfg struct {
	Nr_disk  int
	Capacity uint64
}

type HdCfg struct {
	Write_speed     uint64
	Read_speed      uint64
	Extent_size     uint64
	Data_scheme     int
	Coding_scheme   int
	Network_bdwidth uint64
	Hdservers       []HdSrvCfg
}

var HDCFG = HdCfg{
	Write_speed:     WRITESPEED,
	Read_speed:      READSPEED,
	Data_scheme:     DATASCHEME,
	Coding_scheme:   CODINGSCHEME,
	Extent_size:     EXTENTSIZE,
	Network_bdwidth: NETWORKTHR,
	Hdservers:       make([]HdSrvCfg, 0),
}

func HDCfgDefault() {
	fmt.Println("no readable configuration, use default settings")
	HDCFG.Hdservers = append(HDCFG.Hdservers, HdSrvCfg{NRDISKDEFAULT, CAPACITYDISKDEFAULT},
		HdSrvCfg{NRDISKDEFAULT, CAPACITYDISKDEFAULT},
		HdSrvCfg{NRDISKDEFAULT, CAPACITYDISKDEFAULT})
}

func HDCfgLoad(filename *string) {
	if filename == nil {
		HDCfgDefault()
		return
	} else {
		content, err := ioutil.ReadFile(*filename)
		if err != nil {
			HDCfgDefault()
		} else {
			if err := json.Unmarshal(content, &HDCFG); err != nil {
				fmt.Println("cannot parse json file ", *filename, "cause : ", err)
				HDCfgDefault()
			}
		}
	}
	fmt.Println("configuration used ", HDCFG)
}
