// Реализация управления "насосами"
// Т.к. наполняющий и опустощающий насосы отличаются только порогом срабатывания
// то для простоты введён флаг fill - означающий какую логику реализует данный насос.
// И везде, где логика различается проверяется этот флаг, в остальном веськод одинаковый
// Основная идея:
// - у каждого "объекта" есть входы и выходы (оформленные в виде список и отдельный полей)
// - имеется main loop (см. функцию Run), в котором идёт
//   - обработка сообщений об изменении датчиков и входов (doUpdateInputs)
//   - шаг алгоритма (doStep)
//   - обновление состояния выходов (doUpdateOutputs)
//

package main

import (
	"fmt"
	"sync"
	"time"
	"uniset"
)

// ----------------------------------------------------------------------------------
// Алгоритм управления
type Pump struct {
	evnchannel chan uniset.UMessage
	cmdchannel chan uniset.UMessage
	Pump_SK

	fill       bool  // признак того, что насос наполняющий
	levelLimit int64 // порог до которого работаем
	isWorking  bool
}

// ----------------------------------------------------------------------------------
func NewPump(name string, section string, fill bool, levelLimit int64) *Pump {
	p := Pump{}
	p.evnchannel = make(chan uniset.UMessage, 10)
	p.cmdchannel = make(chan uniset.UMessage, 10)

	Init_Pump(&p, name, section)

	p.levelLimit = levelLimit
	p.fill = fill
	return &p
}

// ----------------------------------------------------------------------------------
func (p *Pump) ID() uniset.ObjectID {
	return p.id
}

// ----------------------------------------------------------------------------------
func (p *Pump) UEvent() chan<- uniset.UMessage {
	return p.evnchannel
}

// ----------------------------------------------------------------------------------
func (p *Pump) UCommand() <-chan uniset.UMessage {
	return p.cmdchannel
}

// ----------------------------------------------------------------------------------
func (p *Pump) Run(wg *sync.WaitGroup) {

	defer wg.Done()

	// шаг алгоритма
	step := time.After(250 * time.Millisecond)

	// как часто обновлять выходы
	outs := time.After(250 * time.Millisecond)

	for {
		select {
		case umsg, ok := <-p.evnchannel:

			if !ok {
				break
			}

			msg, ok := umsg.PopAsSensorEvent()
			if ok {
				p.doUpdateInputs(msg)
				p.doSensorEvent(msg)
				break
			}

			_, ok = umsg.PopAsActivateEvent()
			if ok {
				p.doActivate()
				break
			}

			_, ok = umsg.PopAsFinishEvent()
			if ok {
				p.doFinish()
				fmt.Printf("%s: finish\n", p.myname)
				return
			}

		case <-step:
			p.doStep()
			step = time.After(250 * time.Millisecond)

		case <-outs:
			p.doUpdateOutputs()
			outs = time.After(250 * time.Millisecond)
		}
	}
}

// ----------------------------------------------------------------------------------
func (p *Pump) doActivate() {

	fmt.Printf("%s: activate Ok\n", p.myname)

	// инициализируем выходы
	uniset.DoReadInputs(&p.outs)
	p.doUpdateOutputs()
	// заказываем датчики
	p.doAskSensors()
}

// ----------------------------------------------------------------------------------
func (p *Pump) doAskSensors() {

	uniset.DoAskSensors(&p.ins, p.cmdchannel)
}

// ----------------------------------------------------------------------------------
func (p *Pump) doFinish() {

	fmt.Printf("%s: finish..\n", p.myname)
	// сбрасываем выходы в 0
	p.out_switchOn_c = 0
	p.out_complete_c = 0
	p.doUpdateOutputs()
	close(p.cmdchannel)
}

// ----------------------------------------------------------------------------------
func (p *Pump) doUpdateInputs(sm *uniset.SensorEvent) {

	uniset.DoUpdateInputs(&p.ins, sm)
}

// ----------------------------------------------------------------------------------
func (p *Pump) isLevelOk() bool {

	if p.fill {
		return p.in_level_s >= p.levelLimit
	}

	return p.in_level_s <= p.levelLimit
}

// ----------------------------------------------------------------------------------
func (p *Pump) doSensorEvent(sm *uniset.SensorEvent) {

	//fmt.Printf("%s: sensor %d = %d\n", p.myname, sm.Id, sm.Value)
	if sm.Id == p.onControl_s {
		if sm.Value == 0 {
			fmt.Printf("%s: Управление отключено\n", p.myname)
			p.out_switchOn_c = 0
		} else {
			fmt.Printf("%s: Включено управление\n", p.myname)
			if p.fill {
				fmt.Printf("%s: Начинаю наполнение\n", p.myname)
				p.isWorking = true
				p.out_complete_c = 0
			} else {
				p.isWorking = false
				p.out_complete_c = 0
			}
		}

	} else if sm.Id == p.isComplete_s {
		if p.in_onControl_s != 0 && p.in_isComplete_s == 0 && !p.isLevelOk() {

			if p.fill {
				fmt.Printf("%s: начинаю наполнять..\n", p.myname)
			} else {
				fmt.Printf("%s: начинаю опустошать..\n", p.myname)
			}

			p.isWorking = true
		}
	}
}

// ----------------------------------------------------------------------------------
func (p *Pump) doStep() {

	// если управление отключено ничего не делаем
	if p.in_onControl_s == 0 {
		return
	}

	if !p.isWorking {
		p.out_switchOn_c = 0
		p.out_complete_c = 0
		return
	}

	if p.isLevelOk() {
		if p.fill {
			fmt.Printf("%s: наполнять закончил\n", p.myname)
		} else {
			fmt.Printf("%s: опустошать закончил\n", p.myname)
		}
		p.out_switchOn_c = 0
		p.out_complete_c = 1
		p.isWorking = false
	} else {
		p.out_switchOn_c = 1
		p.out_complete_c = 0
	}

}

// ----------------------------------------------------------------------------------
// Проходим по выходам и если значение поменялось, относительно предыдущего
// обновляем (делаем setValue)
func (p *Pump) doUpdateOutputs() {
	uniset.DoUpdateOutputs(&p.outs, p.cmdchannel)
}
