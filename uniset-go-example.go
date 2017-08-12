package main

// Пример работы с пакетом uniset
// Демонстрационная задача:
// Есть два насоса, один наполняет бак, другой опустошает.
// Как только бак наполнился, наполнявший насос должен отключиться,
// а второй (опустощающий) должен включиться. Как только бак опустошился
// опять включается наполнящий насос.
// Задача написать логику реализующую это взаимодействие, через uniset-датчики.
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

	pump1 := &Pump{"Pump1", 100, make(chan uniset.UMessage, 10), make(chan uniset.UMessage, 10)}
	pump2 := &Pump{"Pump2", 101, make(chan uniset.UMessage, 10), make(chan uniset.UMessage, 10)}

	act.Add(pump1)
	act.Add(pump2)

	var wg sync.WaitGroup

	go pump1.Run(&wg)
	go pump2.Run(&wg)

	act.WaitFinish()

	wg.Wait()
}
