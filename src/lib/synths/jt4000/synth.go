package jt4000

import (
	"errors"
	"fmt"
	"log"
	"sync"

	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

type Synth struct {
	wg                sync.WaitGroup
	inPort            drivers.In
	inStop            func()
	outPort           drivers.Out
	CurrentPatchIndex int
	CurrentPatchName  string

	MidiChannel     uint8
	PatchNames      []string
	SelectionParams map[string]SelectionParam
	IntParams       map[string]IntParam
}

func NewSynth(inPort drivers.In, outPort drivers.Out) (s *Synth) {
	s = &Synth{inPort: inPort, outPort: outPort, PatchNames: []string{""}, CurrentPatchIndex: 0}
	s.initParams()
	s.openPorts()
	s.GetAllPatchNames()
	s.GetCurrentPatchDetails()
	return
}

func (s *Synth) SetCurrentPatch(patch int) error {
	if patch < 0 || patch >= len(s.PatchNames) {
		return errors.New("patch out of range")
	}
	programChangeMsg := midi.ProgramChange(s.MidiChannel, uint8(patch))
	err := s.outPort.Send(programChangeMsg)
	if err != nil {
		return err
	}
	s.GetCurrentPatchDetails()
	s.CurrentPatchIndex = patch
	s.PatchNames[patch] = s.CurrentPatchName
	return nil
}

func (s *Synth) CurrentPatchNumber() (patchNumber int) {
	return s.CurrentPatchIndex + 1
}

func (s *Synth) setValue(cc, val uint8) error {
	msg := midi.ControlChange(s.MidiChannel, cc, val)
	return s.outPort.Send(msg)
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
