package main

// Пример работы с пакетом uniset
// Демонстрационная задача:
// Есть два насоса, один наполняет бак, другой опустошает.
// Как только бак наполнился, наполнявший насос должен отключиться,
// а второй (опустощающий) должен включиться. Как только бак опустошился
// опять включается наполнящий насос.
// Задача написать логику реализующую это взаимодействие, через uniset-датчики.
// В целом это решение на go, задачи отсюда: https://habrahabr.ru/post/278535/
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
		103)

	pumpDrain := NewPump("PumpDrain", 20005,
		100,
		101,
		102,
		103)

	act.Add(pumpFill)
	act.Add(pumpDrain)

	var wg sync.WaitGroup

	go pumpFill.Run(&wg)
	go pumpDrain.Run(&wg)

	act.WaitFinish()

	wg.Wait()
}
