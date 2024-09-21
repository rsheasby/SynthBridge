package jt4000

import (
	"errors"
	"fmt"
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

func NewSynth() (s *Synth, err error) {
	s = &Synth{PatchNames: []string{""}, CurrentPatchIndex: 0}
	err = s.connectToMidiPorts()
	if err != nil {
		return nil, err
	}
	s.initParams()
	s.GetAllPatchNames()
	s.GetCurrentPatchDetails()
	s.guessPatchIndex()
	return
}

func (s *Synth) SendNoteOn(note uint8, velocity uint8) (err error) {
	msg := midi.NoteOn(s.MidiChannel, note, velocity)
	return s.outPort.Send(msg)
}

func (s *Synth) SendNoteOff(note uint8) (err error) {
	msg := midi.NoteOff(s.MidiChannel, note)
	return s.outPort.Send(msg)
}

const (
	midiInPortName  = "JT-4000 MICRO"
	midiOutPortName = "JT-4000 MICRO"
)

func (s *Synth) connectToMidiPorts() (err error) {
	inPorts := midi.GetInPorts()
	outPorts := midi.GetOutPorts()

	for _, port := range inPorts {
		if port.String() == midiInPortName {
			s.inPort = port
			break
		}
	}
	if s.inPort == nil {
		return fmt.Errorf(`couldn't find input port "%s"`, midiInPortName)
	}
	err = s.inPort.Open()
	if err != nil {
		return fmt.Errorf("failed to open MIDI input port: %v", err)
	}

	for _, port := range outPorts {
		if port.String() == midiOutPortName {
			s.outPort = port
			break
		}
	}
	if s.outPort == nil {
		return fmt.Errorf(`couldn't find output port "%s"`, midiOutPortName)
	}
	err = s.outPort.Open()
	if err != nil {
		return fmt.Errorf("failed to open MIDI output port: %v", err)
	}

	s.inMsgListen()

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

func (s *Synth) inMsgListen() {
	var err error
	s.inStop, err = midi.ListenTo(s.inPort, func(msg midi.Message, timestampms int32) {
		var bt []byte
		switch {
		case msg.GetSysEx(&bt):
			s.parseIncomingSysex(bt)
		default:
			// ignore
		}
	}, midi.UseSysEx(), midi.SysExBufferSize(4096))

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}
}
