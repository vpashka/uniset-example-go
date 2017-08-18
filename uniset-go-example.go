// ------------------------------------
// Денстрационный пример работы с uniset
// см. README.md
// ------------------------------------
package main

import (
	"fmt"
	"sync"
	"uniset"
)

func main() {

	uniset.Init("configure.xml")

	act := uniset.NewUProxy("UProxy1")

	defer act.Terminate()
	act.Run()

	if !act.IsActive() {
		fmt.Print("UProxy: Not ACTIVE after run")
	}

	pumpFill := NewPump("PumpFill", "Pump", true, 100)
	pumpDrain := NewPump("PumpDrain", "Pump", false, 10)
	im := NewImitator("Imitator1", "Imitator", 5, 0, 100)

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
