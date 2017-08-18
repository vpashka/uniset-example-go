package main

// Пример работы с пакетом uniset
// Демонстрационная задача:
// Есть два насоса, один наполняет бак, другой опустошает.
// Как только бак наполнился, наполнявший насос должен отключиться,
// а второй (опустощающий) должен включиться. Как только бак опустошился
// опять включается наполнящий насос (а-ля "пинг-понг")
// Задача написать логику реализующую это взаимодействие, через uniset-датчики.
//
// ----------------
import (
	"fmt"
	"uniset"
	"sync"
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
	pumpDrain := NewPump("PumpDrain", "Pump", false,10)
	im := NewImitator("Imitator1", "Imitator",	5, 0, 100)

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
