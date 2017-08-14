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
	"sync"
	"uniset"
)

func main() {

	act := uniset.NewUProxy("UProxy1", "configure.xml", 53817)

	defer act.Terminate()
	act.Run()

	if !act.IsActive() {
		fmt.Print("UProxy: Not ACTIVE after run")
	}

	pumpFill := NewPump("PumpFill", 20004,
		100,
		101,
		102,
		104,
		105,
		true,
		100)

	pumpDrain := NewPump("PumpDrain", 20005,
		100,
		101,
		103,
		105,
		104,
		false,
		10)

	im := NewImitator("Imitator", 2006,
		101,
		102,
		103,
		5, 0, 100)

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
