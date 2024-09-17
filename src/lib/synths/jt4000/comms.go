package jt4000

import (
	"bytes"
	"encoding/hex"
	"errors"
	"fmt"
	"log"

	"github.com/davecgh/go-spew/spew"
	"github.com/fatih/color"
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
			parseIncomingSysex(bt)
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

var lastData []byte

func DumpData(data []byte) {
	red := color.New(color.FgRed).SprintFunc()

	dump := hex.Dump(data)
	lastDump := hex.Dump(lastData)
	outDump := ""

	for i := 0; i < len(data) && i < len(lastData); i++ {
		if data[i] != lastData[i] {
			fmt.Printf("Difference at position %d: %02X != %02X\n", i, data[i], lastData[i])
		}
	}

	for i := 0; i < len(dump) && i < len(lastDump); i++ {
		if dump[i] != lastDump[i] {
			outDump += red(string(dump[i]))
		} else {
			outDump += string(dump[i])
		}
	}
	fmt.Println(outDump)

	lastData = data
}

func parseIncomingSysex(data []byte) (err error) {
	if len(data) < 64 {
		return errors.New("sysex message is too short")
	}
	if bytes.HasPrefix(data, []byte{0x00, 0x20, 0x32, 0x00, 0x01, 0x38, 0x15}) {
		fmt.Println("Single Patch")
		patch, err := parsePatch(data[7:71])
		spew.Dump(patch, err)
	} else if bytes.HasPrefix(data, []byte{0x00, 0x20, 0x32, 0x00, 0x01, 0x38, 0x10}) {
		fmt.Println("All Patches")
		for i := 7; i+64 <= len(data); i += 64 {
			patch, err := parsePatch(data[i : i+64])
			if err != nil {
				return err
			}
			spew.Dump(patch)
		}
	}
	// patch := parseSinglePatch(data[:64])

	return
}

func parsePatch(data []byte) (patch Patch, err error) {
	if len(data) != 64 {
		return patch, errors.New("patch data length is not 64")
	}
	patch.Name = string(data[55:64])
	patch.Osc1Wave = OscWave(data[0])
	patch.Osc1Adj = data[2]
	patch.Osc1Coarse = data[4]
	patch.Osc1Fine = data[5]

	patch.Osc2Wave = OscWave(data[1])
	if patch.Osc2Wave == OscWaveSuperSaw {
		patch.Osc2Wave = OscWaveNoise
	}
	patch.Osc2Adj = data[3]
	patch.Osc2Coarse = data[6]
	patch.Osc2Fine = data[7]

	patch.OscBalance = data[8] - 0x40

	patch.PortamentoAmount = data[46]
	patch.PortamentoMode = PortamentoMode(data[45])

	patch.RingModEnabled = data[43] == 0x01
	patch.RingModAmount = data[44]

	patch.LFO1Wave = LFOWave(data[47])
	patch.LFO1Destination = LFODestination(data[53])
	patch.LFO1Speed = data[49]
	patch.LFO1Amount = data[50]

	patch.LFO2Wave = LFOWave(data[48])
	patch.LFO2Speed = data[51]
	patch.LFO2Amount = data[52]

	patch.VCFAttack = data[14]
	patch.VCFDecay = data[15]
	patch.VCFSustain = data[16]
	patch.VCFRelease = data[17]
	patch.VCFAmount = data[22]

	patch.VCFCutoff = data[12]
	patch.VCFResonance = data[13]
	patch.VCAAttack = data[18]
	patch.VCADecay = data[19]
	patch.VCASustain = data[20]
	patch.VCARelease = data[21]

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
