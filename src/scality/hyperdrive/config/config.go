package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
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

var HDCFG HdCfg

func init() {
	HDCFG.Write_speed = WRITESPEED
	HDCFG.Read_speed = READSPEED
	HDCFG.Data_scheme = DATASCHEME
	HDCFG.Coding_scheme = CODINGSCHEME
	HDCFG.Extent_size = EXTENTSIZE
	HDCFG.Network_bdwidth = NETWORKTHR

	HDCFG.Hdservers = make([]HdSrvCfg, 0)
}

func HDCfgDefault() {
	fmt.Println("no readable configuration, use default settings")
	HDCFG.Hdservers = append(HDCFG.Hdservers, HdSrvCfg{NRDISKDEFAULT, CAPACITYDISKDEFAULT})
	HDCFG.Hdservers = append(HDCFG.Hdservers, HdSrvCfg{NRDISKDEFAULT, CAPACITYDISKDEFAULT})
	HDCFG.Hdservers = append(HDCFG.Hdservers, HdSrvCfg{NRDISKDEFAULT, CAPACITYDISKDEFAULT})

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
