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

	din  []*uniset.BoolValue  // список bool-вых входов
	ain  []*uniset.Int64Value // список аналоговых входов
	dout []*uniset.BoolValue  // список bool-вых выходов
	aout []*uniset.Int64Value // список аналоговых выходов

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

	p.din = []*uniset.BoolValue{
		uniset.NewBoolValue(&p.onControl_s, &p.in_onControl_s),
		uniset.NewBoolValue(&p.isComplete_s, &p.in_isComplete_s)}

	p.ain = []*uniset.Int64Value{uniset.NewInt64Value(&p.level_s, &p.in_level_s)}

	p.dout = []*uniset.BoolValue{
		uniset.NewBoolValue(&p.switchOn_c, &p.out_switchOn_c),
		uniset.NewBoolValue(&p.complete_c, &p.out_complete_c)}

	p.aout = []*uniset.Int64Value{}

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
				fmt.Printf("%s: finish\n",p.name)
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

	fmt.Printf("%s: activate Ok\n", p.name)

	// инициализируем выходы
	uniset.DoReadBoolInputs(&p.dout)
	uniset.DoReadAnalogInputs(&p.aout)
	p.doUpdateOutputs()
	// заказываем датчики
	p.doAskSensors()
}

// ----------------------------------------------------------------------------------
func (p *Pump) doAskSensors() {

	uniset.DoAskSensorsBool(&p.din, p.cmdchannel)
	uniset.DoAskSensorsAnalog(&p.ain, p.cmdchannel)
}

// ----------------------------------------------------------------------------------
func (p *Pump) doFinish() {

	fmt.Printf("%s: finish..\n", p.name)
	// сбрасываем выходы в 0
	p.out_switchOn_c = false
	p.out_complete_c = false
	p.doUpdateOutputs()
	close(p.cmdchannel)
}

// ----------------------------------------------------------------------------------
func (p *Pump) doUpdateInputs(sm *uniset.SensorEvent) {

	uniset.DoUpdateBoolInputs(&p.din, sm)
	uniset.DoUpdateAnalogInputs(&p.ain, sm)
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
				p.out_complete_c = false
			} else {
				p.isWorking = false
				p.out_complete_c = false
			}
		}

	} else if sm.Id == p.isComplete_s {
		if p.in_isComplete_s && !p.isLevelOk() {

			if p.fill {
				fmt.Printf("%s: начинаю наполнять..\n", p.name)
			} else {
				fmt.Printf("%s: начинаю опустошать..\n", p.name)
			}

			p.isWorking = true
		}
	}
}

// ----------------------------------------------------------------------------------
func (p *Pump) doStep() {

	// если управление отключено ничего не делаем
	if !p.in_onControl_s {
		return
	}

	if !p.isWorking {
		p.out_switchOn_c = false
		p.out_complete_c = false
		return
	}

	if p.isLevelOk() {
		if p.fill {
			fmt.Printf("%s: наполнять закончил\n", p.name)
		} else {
			fmt.Printf("%s: опустошать закончил\n", p.name)
		}
		p.out_switchOn_c = false
		p.out_complete_c = true
		p.isWorking = false
	} else {
		p.out_switchOn_c = true
		p.out_complete_c = false
	}

}

// ----------------------------------------------------------------------------------
// Проходим по выходам и если значение поменялось, относительно предыдущего
// обновляем (делаем setValue)
func (p *Pump) doUpdateOutputs() {

	uniset.DoUpdateAnalogOutputs(&p.aout, p.cmdchannel)
	uniset.DoUpdateBoolOutputs(&p.dout, p.cmdchannel)
}
