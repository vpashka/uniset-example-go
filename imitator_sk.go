package main

import (
	"fmt"
	"uniset"
)

// ----------------------------------------------------------------------------------
type Imitator_SK struct {

	// ID
	level_s     uniset.SensorID
	cmdLoad_c   uniset.SensorID
	cmdUnLoad_c uniset.SensorID

	// i/o
	out_level_s    int64
	in_cmdLoad_c   int64
	in_cmdUnLoad_c int64

	// variables

	ins  []*uniset.Int64Value // список входов
	outs []*uniset.Int64Value // список выходов

	myname string
	id     uniset.ObjectID
}

// ----------------------------------------------------------------------------------
func Init_Imitator(sk *Imitator, name string, section string) {

	sk.myname = name

	cfg, err := uniset.GetConfigParamsByName(name, section)
	if err != nil {
		panic(fmt.Sprintf("(Init_Imitator): error: %s", err))
	}

	sk.id = uniset.InitObjectID(cfg, name, name)

	sk.level_s = uniset.InitSensorID(cfg, "level_s", "")
	sk.cmdLoad_c = uniset.InitSensorID(cfg, "cmdLoad_c", "")
	sk.cmdUnLoad_c = uniset.InitSensorID(cfg, "cmdUnLoad_c", "")

	sk.ins = []*uniset.Int64Value{
		uniset.NewInt64Value(&sk.cmdLoad_c, &sk.in_cmdLoad_c),
		uniset.NewInt64Value(&sk.cmdUnLoad_c, &sk.in_cmdUnLoad_c),
	}

	sk.outs = []*uniset.Int64Value{
		uniset.NewInt64Value(&sk.level_s, &sk.out_level_s),
	}

}

// ----------------------------------------------------------------------------------
