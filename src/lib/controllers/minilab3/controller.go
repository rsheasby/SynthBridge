package minilab3

import (
	"errors"
	"fmt"
	"log"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type Controller struct {
	inPort  drivers.In
	inStop  func()
	outPort drivers.Out
}

func NewController(inPort drivers.In, outPort drivers.Out) (controller *Controller, err error) {
	controller = &Controller{inPort: inPort, outPort: outPort}
	err = controller.openPorts()
	if err != nil {
		return nil, err
	}

	err = controller.initDisplay()
	if err != nil {
		return nil, err
	}
	return
}

func (c *Controller) openInPort() (err error) {
	if c.inPort == nil {
		return errors.New("no input port")
	}
	if c.inPort.IsOpen() {
		return nil
	}
	err = c.inPort.Open()
	if err != nil {
		return
	}
	c.inMsgListen()
	return
}

func (c *Controller) inMsgListen() {
	var err error
	c.inStop, err = midi.ListenTo(c.inPort, func(msg midi.Message, timestampms int32) {
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
		fmt.Printf("ERROR: %s\n", err)
		return
	}
}

func (c *Controller) openOutPort() (err error) {
	if c.outPort == nil {
		return errors.New("no output port")
	}
	if c.outPort.IsOpen() {
		return nil
	}
	return c.outPort.Open()
}

func (c *Controller) openPorts() (err error) {
	if err = c.openInPort(); err != nil {
		return
	}
	return c.openOutPort()
}
