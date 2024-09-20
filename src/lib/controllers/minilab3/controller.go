package minilab3

import (
	"fmt"
	"log"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

const midiInPortName = "Minilab3 MIDI"
const midiOutPortName = "Minilab3 MIDI"

type Controller struct {
	midiInPort  drivers.In
	midiInStop  func()
	midiOutPort drivers.Out
}

func NewController() (controller *Controller, err error) {
	controller = &Controller{}
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

func (c *Controller) midiMsgListen() (err error) {
	c.midiInStop, err = midi.ListenTo(c.midiInPort, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, vel uint8
		switch {
		case msg.GetSysEx(&bt):
			log.Println("received sysex")
		case msg.GetNoteStart(&ch, &key, &vel):
			fmt.Printf("starting note %s on channel %v with velocity %v\n", midi.Note(key), ch, vel)
		case msg.GetNoteEnd(&ch, &key):
			fmt.Printf("ending note %s on channel %v\n", midi.Note(key), ch)
		default:
			// ignore
		}
	}, midi.UseSysEx(), midi.SysExBufferSize(4096))

	if err != nil {
		return fmt.Errorf("failed to listen to MIDI input port: %v", err)
	}
	return
}
