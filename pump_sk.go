package main

import (
	"fmt"
	"uniset"
)

// ----------------------------------------------------------------------------------
type Pump_SK struct {

	// ID
	level_s      uniset.SensorID
	onControl_s  uniset.SensorID
	isComplete_s uniset.SensorID
	switchOn_c   uniset.SensorID
	complete_c   uniset.SensorID

	// i/o
	in_level_s      int64
	in_onControl_s  int64
	in_isComplete_s int64
	out_switchOn_c  int64
	out_complete_c  int64

	// variables

	ins  []*uniset.Int64Value // список входов
	outs []*uniset.Int64Value // список выходов

	myname string
	id     uniset.ObjectID
}

// ----------------------------------------------------------------------------------
func Init_Pump(sk *Pump, name string, section string) {

	sk.myname = name

	cfg, err := uniset.GetConfigParamsByName(name, section)
	if err != nil {
		panic(fmt.Sprintf("(Init_Pump): error: %s", err))
	}

	sk.id = uniset.InitObjectID(cfg, name, name)

	sk.level_s = uniset.InitSensorID(cfg, "level_s", "")
	sk.onControl_s = uniset.InitSensorID(cfg, "onControl_s", "")
	sk.isComplete_s = uniset.InitSensorID(cfg, "isComplete_s", "")
	sk.switchOn_c = uniset.InitSensorID(cfg, "switchOn_c", "")
	sk.complete_c = uniset.InitSensorID(cfg, "complete_c", "")

	sk.ins = []*uniset.Int64Value{
		uniset.NewInt64Value(&sk.level_s, &sk.in_level_s),
		uniset.NewInt64Value(&sk.onControl_s, &sk.in_onControl_s),
		uniset.NewInt64Value(&sk.isComplete_s, &sk.in_isComplete_s),
	}

	sk.outs = []*uniset.Int64Value{
		uniset.NewInt64Value(&sk.switchOn_c, &sk.out_switchOn_c),
		uniset.NewInt64Value(&sk.complete_c, &sk.out_complete_c),
	}

}

// ----------------------------------------------------------------------------------