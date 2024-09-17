package jt4000

import (
	"errors"

	"github.com/rsheasby/SynthBridge/lib/utils"
	"gitlab.com/gomidi/midi/v2"
)

func (s *Synth) setValue(cc, val uint8) error {
	msg := midi.ControlChange(s.MidiChannel, cc, val)
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
	err := s.setValue(113, utils.Map99(val))
	if err != nil {
		return err
	}
	s.LivePatch.Osc1Adj = val
	return nil
}

func (s *Synth) SetOsc1Coarse(val uint8) error {
	err := s.setValue(115, utils.Map24(val))
	if err != nil {
		return err
	}
	s.LivePatch.Osc1Coarse = val
	return nil
}

func (s *Synth) SetOsc1Fine(val uint8) error {
	err := s.setValue(111, utils.Map99(val))
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
	err := s.setValue(114, utils.Map99(val))
	if err != nil {
		return err
	}
	s.LivePatch.Osc2Adj = val
	return nil
}

func (s *Synth) SetOsc2Coarse(val uint8) error {
	err := s.setValue(116, utils.Map24(val))
	if err != nil {
		return err
	}
	s.LivePatch.Osc2Coarse = val
	return nil
}

func (s *Synth) SetOsc2Fine(val uint8) error {
	err := s.setValue(112, utils.Map99(val))
	if err != nil {
		return err
	}
	s.LivePatch.Osc2Fine = val
	return nil
}

func (s *Synth) SetOscBalance(val uint8) error {
	err := s.setValue(29, utils.Map63(val))
	if err != nil {
		return err
	}
	s.LivePatch.OscBalance = val
	return nil
}

func (s *Synth) BruteforceSet(startCC, endCC, value uint8) error {
	if startCC > endCC {
		return errors.New("startCC must be less than or equal to endCC")
	}

	for cc := startCC; cc <= endCC; cc++ {
		msg := midi.ControlChange(s.MidiChannel, cc, value)
		err := s.outPort.Send(msg)
		if err != nil {
			return err
		}
	}

	return nil
}
