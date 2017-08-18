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
type Imitator_SK struct {

	// ID
	level_s     uniset.ObjectID
	cmdLoad_c   uniset.ObjectID
	cmdUnLoad_c uniset.ObjectID

	// i/o
	out_level_s    int64
	in_cmdLoad_c   int64
	in_cmdUnLoad_c int64

	// variables
	step int64 /*!<  */
	min  int64 /*!<  */
	max  int64 /*!<  */

	ins  []*uniset.Int64Value // список входов
	outs []*uniset.Int64Value // список выходов

	myname     string
	id         uniset.ObjectID
	sleep_msec time.Duration
}

// ----------------------------------------------------------------------------------
func Init_Imitator(sk *Imitator, name string, section string) {

	sk.myname = name

	cfg, err := uniset.GetConfigParamsByName(name, section)
	if err != nil {
		panic(fmt.Sprintf("(Init_Imitator): error: %s", err))
	}

	sk.id = uniset.InitObjectID(cfg, "", name)

	sk.sleep_msec = time.Duration(200 * time.Millisecond)
	sk.step = uniset.InitInt64(cfg, "step", "5")
	sk.min = uniset.InitInt64(cfg, "min", "0")
	sk.max = uniset.InitInt64(cfg, "max", "100")

	sk.level_s = uniset.InitSensorID(cfg, "level_s", "")
	sk.cmdLoad_c = uniset.InitSensorID(cfg, "cmdLoad_c", "")
	sk.cmdUnLoad_c = uniset.InitSensorID(cfg, "cmdUnLoad_c", "")

	if sk.level_s == uniset.DefaultObjectID {
		panic(fmt.Sprintf("%s(Init_Imitator): Unknown ID for level_s\n", sk.myname))
	}
	if sk.cmdLoad_c == uniset.DefaultObjectID {
		panic(fmt.Sprintf("%s(Init_Imitator): Unknown ID for cmdLoad_c\n", sk.myname))
	}
	if sk.cmdUnLoad_c == uniset.DefaultObjectID {
		panic(fmt.Sprintf("%s(Init_Imitator): Unknown ID for cmdUnLoad_c\n", sk.myname))
	}

	sk.ins = []*uniset.Int64Value{
		uniset.NewInt64Value(&sk.cmdLoad_c, &sk.in_cmdLoad_c),
		uniset.NewInt64Value(&sk.cmdUnLoad_c, &sk.in_cmdUnLoad_c),
	}

	sk.outs = []*uniset.Int64Value{
		uniset.NewInt64Value(&sk.level_s, &sk.out_level_s),
	}

}

// ----------------------------------------------------------------------------------
