// DO NOT EDIT THIS FILE. IT IS AUTOGENERATED FILE.
// ALL YOUR CHANGES WILL BE LOST.
//
// НЕ РЕДАКТИРУЙТЕ ЭТОТ ФАЙЛ. ЭТОТ ФАЙЛ СОЗДАН АВТОМАТИЧЕСКИ.
// ВСЕ ВАШИ ИЗМЕНЕНИЯ БУДУТ ПОТЕРЯНЫ

package main

import (
	"fmt"
	"time"
	"uniset"
)

// ----------------------------------------------------------------------------------
type Pump_SK struct {

	// ID
	level_s      uniset.ObjectID
	onControl_s  uniset.ObjectID
	isComplete_s uniset.ObjectID
	switchOn_c   uniset.ObjectID
	complete_c   uniset.ObjectID

	// i/o
	in_level_s      int64
	in_onControl_s  int64
	in_isComplete_s int64
	out_switchOn_c  int64
	out_complete_c  int64

	// variables
	levelLimit int64 /*!<  */

	ins  []*uniset.Int64Value // список входов
	outs []*uniset.Int64Value // список выходов

	myname     string
	id         uniset.ObjectID
	sleep_msec time.Duration
}

// ----------------------------------------------------------------------------------
func Init_Pump(sk *Pump, name string, section string) {

	sk.myname = name

	cfg, err := uniset.GetConfigParamsByName(name, section)
	if err != nil {
		panic(fmt.Sprintf("(Init_Pump): error: %s", err))
	}

	sk.id = uniset.InitObjectID(cfg, "", name)

	sk.sleep_msec = time.Duration(200 * time.Millisecond)
	sk.levelLimit = uniset.InitInt64(cfg, "levelLimit", "0")

	if sk.levelLimit < 0 {
		panic(fmt.Sprintf("%s(Init_Pump): levelLimit must be > 0\n", sk.myname))
	}

	if sk.levelLimit > 100 {
		panic(fmt.Sprintf("%s(Init_Pump): levelLimit must be < 100\n", sk.myname))
	}

	sk.level_s = uniset.InitSensorID(cfg, "level_s", "")
	sk.onControl_s = uniset.InitSensorID(cfg, "onControl_s", "")
	sk.isComplete_s = uniset.InitSensorID(cfg, "isComplete_s", "")
	sk.switchOn_c = uniset.InitSensorID(cfg, "switchOn_c", "")
	sk.complete_c = uniset.InitSensorID(cfg, "complete_c", "")

	if sk.level_s == uniset.DefaultObjectID {
		panic(fmt.Sprintf("%s(Init_Pump): Unknown ID for level_s\n", sk.myname))
	}
	if sk.onControl_s == uniset.DefaultObjectID {
		panic(fmt.Sprintf("%s(Init_Pump): Unknown ID for onControl_s\n", sk.myname))
	}
	if sk.isComplete_s == uniset.DefaultObjectID {
		panic(fmt.Sprintf("%s(Init_Pump): Unknown ID for isComplete_s\n", sk.myname))
	}
	if sk.switchOn_c == uniset.DefaultObjectID {
		panic(fmt.Sprintf("%s(Init_Pump): Unknown ID for switchOn_c\n", sk.myname))
	}
	if sk.complete_c == uniset.DefaultObjectID {
		panic(fmt.Sprintf("%s(Init_Pump): Unknown ID for complete_c\n", sk.myname))
	}

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
