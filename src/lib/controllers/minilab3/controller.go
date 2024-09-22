package minilab3

import (
	"fmt"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

const midiInPortName = "Minilab3 MIDI"
const midiOutPortName = "Minilab3 MIDI"

type Controller struct {
	midiInPort      drivers.In
	midiInStop      func()
	midiOutPort     drivers.Out
	NoteEvents      chan NoteEvent
	SelectionEvents chan SelectionEvent
	ControlEvents   chan ControlEvent
}

func NewController() (controller *Controller, err error) {
	controller = &Controller{
		NoteEvents:      make(chan NoteEvent, 1),
		SelectionEvents: make(chan SelectionEvent, 1),
		ControlEvents:   make(chan ControlEvent, 1),
	}
	err = controller.connectToMidiPorts()
	if err != nil {
		return nil, err
	}

	err = controller.initDisplay()
	if err != nil {
		return nil, err
	}

	err = controller.midiMsgListen()
	if err != nil {
		return nil, err
	}

	return
}

func (c *Controller) connectToMidiPorts() (err error) {
	inPorts := midi.GetInPorts()
	outPorts := midi.GetOutPorts()

	for _, port := range inPorts {
		if port.String() == midiInPortName {
			c.midiInPort = port
			break
		}
	}
	if c.midiInPort == nil {
		return fmt.Errorf(`couldn't find input port "%s"`, midiInPortName)
	}
	err = c.midiInPort.Open()
	if err != nil {
		return fmt.Errorf("failed to open MIDI input port: %v", err)
	}

	for _, port := range outPorts {
		if port.String() == midiOutPortName {
			c.midiOutPort = port
			break
		}
	}
	if c.midiOutPort == nil {
		return fmt.Errorf(`couldn't find output port "%s"`, midiOutPortName)
	}
	err = c.midiOutPort.Open()
	if err != nil {
		return fmt.Errorf("failed to open MIDI output port: %v", err)
	}

	return
}
