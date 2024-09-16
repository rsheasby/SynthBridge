package jt4000

import (
	"errors"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
)

func NewSynth(inPort drivers.In, outPort drivers.Out) (s *Synth) {
	s = &Synth{inPort: inPort, outPort: outPort}
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

func (s *Synth) openOutPort() (err error) {
	if s.outPort == nil {
		return errors.New("no output port")
	}
	if s.outPort.IsOpen() {
		return nil
	}
	return s.outPort.Open()
}

func (s *Synth) OpenPorts() (err error) {
	if err = s.openInPort(); err != nil {
		return
	}
	return s.openOutPort()
}

func (s *Synth) inMsgListen() {
	var err error
	s.inStop, err = midi.ListenTo(s.inPort, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, vel uint8
		switch {
		case msg.GetSysEx(&bt):
			log.Println("received sysex")
			parseSinglePatch(bt)
		case msg.GetNoteStart(&ch, &key, &vel):
			fmt.Printf("starting note %s on channel %v with velocity %v\n", midi.Note(key), ch, vel)
		case msg.GetNoteEnd(&ch, &key):
			fmt.Printf("ending note %s on channel %v\n", midi.Note(key), ch)
		default:
			// ignore
		}
	}, midi.UseSysEx())

	if err != nil {
		fmt.Printf("ERROR: %s\n", err)
		return
	}

}

func parseSinglePatch(data []byte) (patch Patch, err error) {
	if len(data) != 73 {
		return patch, errors.New("patch data length is not 73")
	}
	fmt.Println("Length of data:", len(data))
	patch.Name = string(data[62:71])
	patch.Osc1Wave = OscWave(data[7])
	patch.Osc1Adj = data[9]
	patch.Osc1Coarse = data[11]
	patch.Osc1Fine = data[12]
	patch.Osc2Wave = OscWave(data[8])
	if patch.Osc2Wave == OscWaveSuperSaw {
		patch.Osc2Wave = OscWaveNoise
	}
	patch.Osc2Adj = data[10]
	patch.Osc2Coarse = data[13]
	patch.Osc2Fine = data[14]
	patch.OscBalance = data[15] - 0x40
	patch.PortamentoAmount = data[53]
	patch.PortamentoMode = PortamentoMode(data[52])
	patch.RingModEnabled = data[50] == 0x01
	patch.RingModAmount = data[51]
	patch.LFO1Wave = LFOWave(data[54])
	patch.LFO1Destination = LFODestination(data[60])
	patch.LFO1Speed = data[56]
	patch.LFO1Amount = data[57]
	patch.LFO2Wave = LFOWave(data[55])
	patch.LFO2Speed = data[58]
	patch.LFO2Amount = data[59]

	spew.Dump(patch)

	// fmt.Println(hex.Dump(data))
	return
}

func (s *Synth) GetAllPatches() (err error) {
	return
}

func (s *Synth) GetCurrentPatch() (err error) {
	log.Println("Getting current patch")
	getCurrentPatchCmd := []byte{0xF0, 0x00, 0x20, 0x32, 0x00, 0x01, 0x38, 0x00, 0x00, 0xF7}

	s.outPort.Send(getCurrentPatchCmd)
	return
}
