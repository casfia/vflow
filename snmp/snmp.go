package snmp

import (
	"time"
	"github.com/alouca/gosnmp"
	"log"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	"../vlogger"
)

var ifNameOid = ".1.3.6.1.2.1.31.1.1.1.1"
var ifIndexOid = ".1.3.6.1.2.1.2.2.1.1"
var ifDesOid = ".1.3.6.1.2.1.31.1.1.1.18"

type WalkTask struct {
	snmpConfigs DeviceSnmpConfig
}
/***

{
    "interval":30,
    "devices":[

        {
            "DeviceAddress": "159.226.8.131",
            "Community":"cst*net"
        },
        {
            "DeviceAddress": "10.0.0.2",
            "Community":"public"
        },
        {
            "DeviceAddress": "10.0.0.3",
            "Community":"public"
        }
    ]
}



 */

type DeviceSnmpConfig struct {
	Interval     int32             `json:"interval"`
	DeviceCfg []CommunityConfig `json:"devices"`
}

type CommunityConfig struct {
	DeviceAddress string `json:"DeviceAddress"`
	Community     string `json:"Community"`
}

var snmpTaskInstance *WalkTask
var snmpCfgFile string
var cfg DeviceSnmpConfig

func Init(cfgFile string) (*WalkTask,error) {
	snmpCfgFile = cfgFile
	snmpTaskInstance = new(WalkTask)

	b, err := ioutil.ReadFile(cfgFile)
	if err != nil {
		vlogger.Logger.Printf("No SNMP config file is defined. \n")
		fmt.Printf("No SNMP config file is defined. \n")
		return nil,err
	}

	fmt.Printf("config is %s",string(b))
	err = json.Unmarshal(b, &cfg)
	if err != nil {
		vlogger.Logger.Printf("SNMP config file is worong, exit! \n")
		fmt.Printf("SNMP config file is worong,exit! \n")
		os.Exit(-1)
		return  nil,err
	}
	fmt.Printf("delay is %d. device length is %d\n",cfg.Interval, len(cfg.DeviceCfg))
	snmpTaskInstance.snmpConfigs = cfg
	return snmpTaskInstance,nil
}

func (task *WalkTask) Run() {
	go func() {
		duration := time.Duration(time.Duration(task.snmpConfigs.Interval) * time.Second)
		timer1 := time.NewTicker(duration)
		for {
			select {
			case <-timer1.C:
				task.task()
			}
		}
	}()
}

func (task *WalkTask) task() {
	for _, dev := range task.snmpConfigs.DeviceCfg {
		task.walkIndex(dev.DeviceAddress, dev.Community)
	}
}

type NameIndex struct {
	IfName  string
	IfIndex string
}

func (task *WalkTask) walkIndex(DeviceAddress string, Community string) {
	s, err := gosnmp.NewGoSNMP(DeviceAddress, Community, gosnmp.Version2c, 5)
	if err != nil {
		log.Fatal(err)
	}
	indexResp, err := s.Walk(ifIndexOid)

	if err == nil {
		for _, v := range indexResp {

			log.Printf("Response: %s : %s : %s \n",
				v.Name, v.Value.(string), v.Type.String())

		}
	} else {
		log.Printf("snmp walk err %e", err)
	}

	nameResp, err := s.Walk(ifNameOid)
	if err == nil {
		for _, v := range nameResp {
			log.Printf("Response: %s : %s : %s \n",
				v.Name, v.Value.(string), v.Type.String())
		}
	} else {
		log.Printf("snmp walk err %e", err)
	}

	desResp, err := s.Walk(ifNameOid)
	if err == nil {
		for _, v := range desResp {
			log.Printf("Response: %s : %s : %s \n",
				v.Name, v.Value.(string), v.Type.String())
		}
	} else {
		log.Printf("snmp walk err %e", err)
	}
}


func (task *WalkTask) AddConfig(DeviceCfg CommunityConfig) (int, string) {
	for _, addr := range task.snmpConfigs.DeviceCfg {
		if addr.DeviceAddress == DeviceCfg.DeviceAddress {
			return -1, "config exist!"
		}
	}
	task.snmpConfigs.DeviceCfg = append(task.snmpConfigs.DeviceCfg, DeviceCfg)
	return len(task.snmpConfigs.DeviceCfg), "add success!"
}

func (task *WalkTask) DeleteConfig(deviceAddr string) (int, string) {
	index := -1
	for i, addr := range task.snmpConfigs.DeviceCfg {
		if addr.DeviceAddress == deviceAddr {
			index = i
			break
		}
	}
	if index == -1 {
		return -1, "can not find address " + deviceAddr
	}
	task.snmpConfigs.DeviceCfg = append(task.snmpConfigs.DeviceCfg[:index],
		task.snmpConfigs.DeviceCfg[index+1:]...)
	err := saveConfigToFile()
	if err != nil {
		return -1, "save config to file error"
	}
	return len(task.snmpConfigs.DeviceCfg), "delete success!"
}

func saveConfigToFile() error {
	b, err := json.MarshalIndent(snmpTaskInstance.snmpConfigs, "", "    ")
	if err == nil {
		return ioutil.WriteFile(snmpCfgFile, b, 0x777)
	} else {
		return err
	}
}