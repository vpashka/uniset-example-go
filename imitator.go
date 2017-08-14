// Реализация имитатора.
// Основная задача, по команда увеличивать или уменшать значение датчика уровня.

package main

import (
	"fmt"
	"sync"
	"time"
	"uniset"
)

// ----------------------------------------------------------------------------------
// Имитатор
type Imitator struct {
	name       string
	id         uniset.ObjectID
	evnchannel chan uniset.UMessage
	cmdchannel chan uniset.UMessage

	din  []*uniset.BoolValue  // список bool-вых входов
	aout []*uniset.Int64Value // список аналоговых выходов

	// датчики
	level_s     uniset.SensorID
	switchOn_c  uniset.SensorID
	switchOff_c uniset.SensorID

	// текущие значения
	out_level_s    int64
	in_switchOn_c  bool
	in_switchOff_c bool

	min  int64
	max  int64
	step int64
}

// ----------------------------------------------------------------------------------
func NewImitator(name string, id uniset.ObjectID,
	level_s uniset.SensorID,
	switchOn_c uniset.SensorID,
	switchOff_c uniset.SensorID,
	step int64, min int64, max int64) *Imitator {
	im := Imitator{}
	im.name = name
	im.id = id
	im.evnchannel = make(chan uniset.UMessage, 10)
	im.cmdchannel = make(chan uniset.UMessage, 10)

	im.level_s = level_s
	im.switchOn_c = switchOn_c
	im.switchOff_c = switchOff_c

	im.min = min
	im.max = max
	im.step = step

	im.aout = []*uniset.Int64Value{uniset.NewInt64Value(&im.level_s, &im.out_level_s)}

	im.din = []*uniset.BoolValue{
		uniset.NewBoolValue(&im.switchOn_c, &im.in_switchOn_c),
		uniset.NewBoolValue(&im.switchOff_c, &im.in_switchOff_c)}

	return &im
}

// ----------------------------------------------------------------------------------
func (im *Imitator) ID() uniset.ObjectID {
	return im.id
}

// ----------------------------------------------------------------------------------
func (im *Imitator) UEvent() chan<- uniset.UMessage {
	return im.evnchannel
}

// ----------------------------------------------------------------------------------
func (im *Imitator) UCommand() <-chan uniset.UMessage {
	return im.cmdchannel
}

// ----------------------------------------------------------------------------------
func (im *Imitator) Run(wg *sync.WaitGroup) {

	defer wg.Done()

	// шаг алгоритма
	step := time.After(250 * time.Millisecond)

	// как часто обновлять выходы
	outs := time.After(250 * time.Millisecond)

	for {
		select {
		case umsg, ok := <-im.evnchannel:

			if !ok {
				break
			}

			msg, ok := umsg.PopAsSensorEvent()
			if ok {
				im.doUpdateInputs(msg)
				im.doSensorEvent(msg)
				break
			}

			_, ok = umsg.PopAsActivateEvent()
			if ok {
				im.doActivate()
				break
			}

			_, ok = umsg.PopAsFinishEvent()
			if ok {
				im.doFinish()
				fmt.Printf("%s: finish\n",im.name)
				return
			}
		case <-step:
			im.doStep()
			step = time.After(250 * time.Millisecond)

		case <-outs:
			im.doUpdateOutputs()
			outs = time.After(250 * time.Millisecond)
		}
	}
}

// ----------------------------------------------------------------------------------
func (im *Imitator) doActivate() {

	fmt.Printf("%s: activate Ok\n", im.name)
	uniset.DoReadAnalogInputs(&im.aout)
	im.doUpdateOutputs()
	// заказываем датчики
	uniset.DoAskSensorsBool(&im.din, im.cmdchannel)
}

// ----------------------------------------------------------------------------------
func (im *Imitator) doFinish() {

	fmt.Printf("%s: finish..\n", im.name)
	close(im.cmdchannel)
}

// ----------------------------------------------------------------------------------
func (im *Imitator) doUpdateInputs(sm *uniset.SensorEvent) {

	uniset.DoUpdateBoolInputs(&im.din, sm)
}

// ----------------------------------------------------------------------------------
func (im *Imitator) doSensorEvent(sm *uniset.SensorEvent) {

}

// ----------------------------------------------------------------------------------
func (im *Imitator) doStep() {

	if im.in_switchOn_c {
		im.out_level_s += im.step
		if im.out_level_s > im.max {
			im.out_level_s = im.max
		}
	} else if im.in_switchOff_c {
		im.out_level_s -= im.step
		if im.out_level_s < im.min {
			im.out_level_s = im.min
		}
	}
}

// ----------------------------------------------------------------------------------
func (im *Imitator) doUpdateOutputs() {

	uniset.DoUpdateAnalogOutputs(&im.aout, im.cmdchannel)
}

// ----------------------------------------------------------------------------------
