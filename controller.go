package main

import (
	"fmt"
	"sync"
	"time"
	"uniset"
)

// ----------------------------------------------------------------------------------
type BoolValue struct {
	sid  *uniset.SensorID
	val  *bool
	prev bool
}

type Int64Value struct {
	sid  *uniset.SensorID
	val  *int64
	prev int64
}

// ----------------------------------------------------------------------------------
// Алгоритм управления
type Pump struct {
	name       string
	id         uniset.ObjectID
	evnchannel chan uniset.UMessage
	cmdchannel chan uniset.UMessage

	din  []BoolValue  // список bool-вых входов
	ain  []Int64Value // список аналоговых входов
	dout []BoolValue  // список bool-вых выходов
	aout []Int64Value // список аналоговых выходов

	// датчики
	level_s      uniset.SensorID
	onControl_s  uniset.SensorID
	switchOn_c   uniset.SensorID
	complete_c   uniset.SensorID
	isComplete_s uniset.SensorID

	// текущие значения
	in_level_s      int64
	in_onControl_s  bool
	in_isComplete_s bool
	out_switchOn_c  bool
	out_complete_c  bool // флаг завершения выполнения работы

	fill       bool  // признак того, что насос наполняющий
	levelLimit int64 // порог до которого работаем
	isWorking  bool
}

// ----------------------------------------------------------------------------------
func NewPump(name string, id uniset.ObjectID,
	onControl_s uniset.SensorID,
	level_s uniset.SensorID,
	switchOn_c uniset.SensorID,
	isComplete_s uniset.SensorID,
	complete_c uniset.SensorID,
	fill bool, levelLimit int64) *Pump {
	p := Pump{}
	p.name = name
	p.id = id
	p.evnchannel = make(chan uniset.UMessage, 10)
	p.cmdchannel = make(chan uniset.UMessage, 10)

	p.level_s = level_s
	p.switchOn_c = switchOn_c

	p.isComplete_s = isComplete_s
	p.complete_c = complete_c
	p.onControl_s = onControl_s

	p.levelLimit = levelLimit
	p.fill = fill

	p.din = []BoolValue{
		{&p.onControl_s, &p.in_onControl_s, p.in_onControl_s},
		{&p.isComplete_s, &p.in_isComplete_s, p.in_isComplete_s}	}

	p.ain = []Int64Value{{&p.level_s, &p.in_level_s, p.in_level_s}}

	p.dout = []BoolValue{
		{&p.switchOn_c, &p.out_switchOn_c, p.out_switchOn_c},
		{&p.complete_c, &p.out_complete_c, p.out_complete_c}}

	p.aout = []Int64Value{}

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

	for {

		// шаг алгоритма
		step := time.After(250 * time.Millisecond)

		// как часто обновлять выходы
		outs := time.After(250 * time.Millisecond)

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
				break
			}
		case <-step:
			p.doStep()

		case <-outs:
			p.doUpdateOutputs()
		}
	}
}

// ----------------------------------------------------------------------------------
func (p *Pump) doActivate() {

	fmt.Printf("%s: activate Ok\n", p.name)
	// заказываем датчики
	p.doAskSensors()
}

// ----------------------------------------------------------------------------------
func (p *Pump) doAskSensors() {

	for _, s := range p.din {
		uniset.AskSensor(p.cmdchannel, *s.sid)
	}

	for _, s := range p.ain {
		uniset.AskSensor(p.cmdchannel, *s.sid)
	}
}

// ----------------------------------------------------------------------------------
func (p *Pump) doFinish() {

	fmt.Printf("%s: finish..\n", p.name)
	// сбрасываем выходы в 0
	p.out_switchOn_c = false
	p.doUpdateOutputs()
}

// ----------------------------------------------------------------------------------
func (p *Pump) doUpdateInputs(sm *uniset.SensorEvent) {

	for _, s := range p.din {
		if *s.sid == sm.Id {
			if sm.Value == 0 {
				*s.val = false
			} else {
				*s.val = true
			}

			s.prev = *s.val
		}
	}

	for _, s := range p.ain {
		if *s.sid == sm.Id {
			*s.val = sm.Value
			s.prev = sm.Value
		}
	}
}

// ----------------------------------------------------------------------------------
func (p *Pump) needWork() bool {

	if p.fill {
		return p.in_level_s >= p.levelLimit
	}

	return p.in_level_s <= p.levelLimit
}

// ----------------------------------------------------------------------------------
func (p *Pump) doSensorEvent(sm *uniset.SensorEvent) {

	//fmt.Printf("%s: sensor %d = %d\n", p.name, sm.Id, sm.Value)
	if sm.Id == p.onControl_s {
		if sm.Value == 0 {
			fmt.Printf("%s: Управление отключено\n", p.name)
			p.out_switchOn_c = false
		} else {
			fmt.Printf("%s: Включено управление\n", p.name)
			if p.fill {
				fmt.Printf("%s: Начинаю наполнение\n", p.name)
				p.isWorking = true
			}
		}

	} else if sm.Id == p.level_s {

	} else if sm.Id == p.isComplete_s {
		if p.needWork() {
			p.isWorking = true
		}
	}
}

// ----------------------------------------------------------------------------------
func (p *Pump) doStep() {

	//fmt.Printf("%s (step): onControl=%d isWorking=%d fill=%d\n",p.name, p.in_onControl_s,p.isWorking,p.fill)

	// если управление отключено
	// ничего не делаем
	if !p.in_onControl_s {
		return
	}

	if !p.isWorking {
		return
	}

	if p.fill {

		// Логика наполняющего насоса

		if p.needWork() {
			fmt.Printf("%s: наполнять закончил..\n", p.name)
			p.out_switchOn_c = false
			p.out_complete_c = true
			p.isWorking = false
		} else {
			p.isWorking = true
			p.out_switchOn_c = true
			p.out_complete_c = false
		}

	} else {

		// Логика опустощающего насоса

		if p.needWork() {
			fmt.Printf("%s: опустошать закончил..\n", p.name)
			p.out_switchOn_c = false
			p.out_switchOn_c = true
			p.isWorking = false
		} else {
			p.isWorking = true
			p.out_switchOn_c = true
			p.out_switchOn_c = false
		}
	}
}

// ----------------------------------------------------------------------------------
// Проходим по выходам и если значение поменялось, относительно предыдущего
// обновляем (делаем setValue)
func (p *Pump) doUpdateOutputs() {

	for _, s := range p.dout {
		if s.prev != *s.val {
			var val int64
			if *s.val {
				val = 1
			}
			uniset.SetValue(p.cmdchannel, *s.sid, val)
			// возможно обновлять prev, стоит после подтверждения от UProxy
			// но пока для простосты обновляем сразу
			s.prev = *s.val
		}
	}

	for _, s := range p.aout {
		if s.prev != *s.val {
			uniset.SetValue(p.cmdchannel, *s.sid, *s.val)
			// возможно обновлять prev, стоит после подтверждения от UProxy
			// но пока для простосты обновляем сразу
			s.prev = *s.val
		}
	}
}
