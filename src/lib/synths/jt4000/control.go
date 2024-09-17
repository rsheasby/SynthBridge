package jt4000

import (
	"errors"

	"github.com/rsheasby/SynthBridge/lib/utils"
	"gitlab.com/gomidi/midi/v2"
)

func (s *Synth) setValue(cc, maxRange, percentVal uint8) error {
	if percentVal > maxRange {
		percentVal = maxRange
	}
	midiVal := utils.Uint8Map(percentVal, 0, maxRange, 0, 127)
	msg := midi.ControlChange(s.MidiChannel, cc, midiVal)
	return s.outPort.Send(msg)
}

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
	err := s.setValue(113, 99, val)
	if err != nil {
		return err
	}
	s.LivePatch.Osc1Adj = val
	return nil
}

func (s *Synth) SetOsc1Coarse(val uint8) error {
	err := s.setValue(115, 24, val)
	if err != nil {
		return err
	}
	s.LivePatch.Osc1Coarse = val
	return nil
}

func (s *Synth) SetOsc1Fine(val uint8) error {
	err := s.setValue(111, 99, val)
	if err != nil {
		return err
	}
	s.LivePatch.Osc1Fine = val
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
	err := s.setValue(114, 99, val)
	if err != nil {
		return err
	}
	s.LivePatch.Osc2Adj = val
	return nil
}

func (s *Synth) SetOsc2Coarse(val uint8) error {
	err := s.setValue(116, 24, val)
	if err != nil {
		return err
	}
	s.LivePatch.Osc2Coarse = val
	return nil
}

func (s *Synth) SetOsc2Fine(val uint8) error {
	err := s.setValue(112, 99, val)
	if err != nil {
		return err
	}
	s.LivePatch.Osc2Fine = val
	return nil
}
