// Реализация имитатора.
// Основная задача, по командам увеличивать или уменьшать значение датчика уровня.
// основные поля генерируются при помощи uniset-codegen-go

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
	evnchannel chan uniset.UMessage
	cmdchannel chan uniset.UMessage

	// датчики
	Imitator_SK
}

// ----------------------------------------------------------------------------------
func NewImitator(name string, section string) *Imitator {
	im := Imitator{}
	im.evnchannel = make(chan uniset.UMessage, 10)
	im.cmdchannel = make(chan uniset.UMessage, 10)

	Init_Imitator(&im, name, section)

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
	step := time.After(im.sleep_msec)

	// как часто обновлять выходы
	outs := time.After(im.sleep_msec)

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
				fmt.Printf("%s: finish\n", im.myname)
				return
			}
		case <-step:
			im.doStep()
			step = time.After(im.sleep_msec)

		case <-outs:
			im.doUpdateOutputs()
			outs = time.After(im.sleep_msec)
		}
	}
}

// ----------------------------------------------------------------------------------
func (im *Imitator) doActivate() {

	fmt.Printf("%s: activate Ok\n", im.myname)
	uniset.DoReadInputs(&im.outs)
	im.doUpdateOutputs()
	// заказываем датчики
	uniset.DoAskSensors(&im.ins, im.cmdchannel)
}

// ----------------------------------------------------------------------------------
func (im *Imitator) doFinish() {

	fmt.Printf("%s: finish..\n", im.myname)
	close(im.cmdchannel)
}

// ----------------------------------------------------------------------------------
func (im *Imitator) doUpdateInputs(sm *uniset.SensorEvent) {

	uniset.DoUpdateInputs(&im.ins, sm)
}

// ----------------------------------------------------------------------------------
func (im *Imitator) doSensorEvent(sm *uniset.SensorEvent) {

}

// ----------------------------------------------------------------------------------
func (im *Imitator) doStep() {

	if im.in_cmdLoad_c != 0 {
		im.out_level_s += im.step
		if im.out_level_s > im.max {
			im.out_level_s = im.max
		}
	} else if im.in_cmdUnLoad_c != 0 {
		im.out_level_s -= im.step
		if im.out_level_s < im.min {
			im.out_level_s = im.min
		}
	}
}

// ----------------------------------------------------------------------------------
func (im *Imitator) doUpdateOutputs() {

	uniset.DoUpdateOutputs(&im.outs, im.cmdchannel)
}

// ----------------------------------------------------------------------------------
