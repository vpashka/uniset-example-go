// ------------------------------------
// Денстрационный пример работы с uniset
// см. README.md
// ------------------------------------
package main

import (
	"sync"
	"uniset"
	"fmt"
)

func main() {

	uniset.Init("configure.xml")

	act := uniset.NewUProxy("UProxy1")

	defer act.Terminate()

	err := act.Run()
	if err != nil {
		panic(fmt.Sprintf("UProxy run error: %s",err))
	}

	pumpFill := NewPump("PumpFill", "Pump", true)
	pumpDrain := NewPump("PumpDrain", "Pump", false)
	im := NewImitator("Imitator1", "Imitator")

	act.Add(pumpFill)
	act.Add(pumpDrain)
	act.Add(im)

	var wg sync.WaitGroup

	wg.Add(3)

	go pumpFill.Run(&wg)
	go pumpDrain.Run(&wg)
	go im.Run(&wg)

	act.WaitFinish()
	wg.Wait()
}
