package jt4000

import (
	"errors"
	"fmt"
	"log"
	"sync"
	"time"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type Synth struct {
	sync.WaitGroup
	inPort          drivers.In
	inStop          func()
	outPort         drivers.Out
	MidiChannel     uint8
	Patches         []Patch
	CurrentPatch    Patch
	LivePatch       Patch
	SelectionParams map[string]SelectionParam
	IntParams       map[string]IntParam
}

func NewSynth(inPort drivers.In, outPort drivers.Out) (s *Synth) {
	s = &Synth{inPort: inPort, outPort: outPort}
	s.initParams()
	s.openPorts()
	s.GetCurrentPatch()
	s.LivePatch = s.CurrentPatch
	time.Sleep(10 * time.Millisecond)
	s.GetAllPatches()
	return
}

func (s *Synth) openInPort() (err error) {
	if s.inPort == nil {
		return errors.New("no input port")
	}
	if s.inPort.IsOpen() {
		return nil
	}
	err = s.inPort.Open()
	if err != nil {
		return
	}
	s.inMsgListen()
	return
}

func (s *Synth) inMsgListen() {
	var err error
	s.inStop, err = midi.ListenTo(s.inPort, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, vel uint8
		switch {
		case msg.GetSysEx(&bt):
			log.Println("received sysex")
			s.parseIncomingSysex(bt)
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

func (s *Synth) openOutPort() (err error) {
	if s.outPort == nil {
		return errors.New("no output port")
	}
	if s.outPort.IsOpen() {
		return nil
	}
	return s.outPort.Open()
}

func (s *Synth) openPorts() (err error) {
	if err = s.openInPort(); err != nil {
		return
	}
	return s.openOutPort()
}
