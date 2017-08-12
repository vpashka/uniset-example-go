package main

import (
	"fmt"
	"sync"
	"time"
	"uniset"
)

// Алгоритм управления
type Pump struct {
	name       string
	id         uniset.ObjectID
	evnchannel chan uniset.UMessage
	cmdchannel chan uniset.UMessage
}

func (p *Pump) ID() uniset.ObjectID {
	return p.id
}

func (p *Pump) UEvent() chan<- uniset.UMessage {
	return p.evnchannel
}

func (p *Pump) UCommand() <-chan uniset.UMessage {
	return p.cmdchannel
}

func (p *Pump) Run(wg *sync.WaitGroup) {

	defer wg.Done()

	for {

		step := time.After(300 * time.Millisecond)

		select {
		case umsg, ok := <-p.evnchannel:

			if !ok {
				break
			}

			msg, ok := umsg.PopAsSensorEvent()
			if ok {
				p.doSensorEvent(msg)
				break
			}

			act, ok := umsg.PopAsActivateEvent()
			if ok {
				if act.Id == p.id {
					fmt.Printf("%s: activate Ok\n", p.name)
				}
			}
		case <-step:
			p.doStep()
		}
	}
}

func (p *Pump) doSensorEvent(sm *uniset.SensorEvent) {

	fmt.Printf("%s: sensor %d = %d\n", p.name, sm.Id, sm.Value)

}

func (p *Pump) doStep() {

}
