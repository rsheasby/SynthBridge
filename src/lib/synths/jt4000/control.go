package jt4000

import (
	"errors"

	"github.com/rsheasby/SynthBridge/lib/utils"
	"gitlab.com/gomidi/midi/v2"
)

func (s *Synth) SetOsc1Wave(wave OscWave) error {
	val := 0
	if wave == OscWaveOff {
		val = 0
	} else if wave == OscWaveTriangle {
		val = 30
	} else if wave == OscWaveSquare {
		val = 50
	} else if wave == OscWavePWM {
		val = 70
	} else if wave == OscWaveSaw {
		val = 80
	} else if wave == OscWaveSuperSaw {
		val = 100
	} else if wave == OscWaveFM {
		val = 120
	} else if wave == OscWaveNoise {
		val = 127
	} else {
		return errors.New("invalid wave type")
	}
	msg := midi.ControlChange(s.MidiChannel, 24, uint8(val))
	err := s.outPort.Send(msg)
	if err != nil {
		return err
	}
	s.LivePatch.Osc1Wave = wave
	return nil
}

func (s *Synth) SetOsc1Adj(val uint8) error {
	if val > 99 {
		val = 99
	}
	midiVal := utils.Uint8Map(val, 0, 99, 0, 127)
	msg := midi.ControlChange(s.MidiChannel, 113, midiVal)
	err := s.outPort.Send(msg)
	if err != nil {
		return err
	}
	s.LivePatch.Osc1Adj = val
	return nil
}

func (s *Synth) SetOsc2Wave(wave OscWave) error {
	val := 0
	if wave == OscWaveOff {
		val = 0
	} else if wave == OscWaveTriangle {
		val = 30
	} else if wave == OscWaveSquare {
		val = 50
	} else if wave == OscWavePWM {
		val = 70
	} else if wave == OscWaveSaw {
		val = 90
	} else if wave == OscWaveNoise {
		val = 127
	} else {
		return errors.New("invalid wave type")
	}
	msg := midi.ControlChange(s.MidiChannel, 25, uint8(val))
	err := s.outPort.Send(msg)
	if err != nil {
		return err
	}
	s.LivePatch.Osc2Wave = wave
	return nil
}

func (s *Synth) SetOsc2Adj(val uint8) error {
	if val > 99 {
		val = 99
	}
	midiVal := utils.Uint8Map(val, 0, 99, 0, 127)
	msg := midi.ControlChange(s.MidiChannel, 114, midiVal)
	err := s.outPort.Send(msg)
	if err != nil {
		return err
	}
	s.LivePatch.Osc1Adj = val
	return nil
}
