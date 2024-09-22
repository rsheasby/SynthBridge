package jt4000

import (
	"errors"
)

type SelectionParam struct {
	synth *Synth

	StaticName  string
	dynamicName func() string
	// CC is the MIDI CC number. If it's -1, the parameter is not writable.
	CC             int
	PossibleValues []string
	bytePosition   int

	CurrentValue      string
	CurrentValueIndex int
	StoredValue       string
	StoredValueIndex  int
}

func (sp *SelectionParam) Name() string {
	if sp.dynamicName != nil {
		return sp.dynamicName()
	}
	return sp.StaticName
}

func (sp *SelectionParam) SetValue(value string) error {
	for i, v := range sp.PossibleValues {
		if v == value {
			return sp.SetValueByIndex(i)
		}
	}
	return errors.New("invalid value label")
}

func (sp *SelectionParam) SetValueByIndex(valueIndex int) error {
	if valueIndex < 0 || valueIndex >= len(sp.PossibleValues) {
		return errors.New("value out of range")
	}
	sp.synth.setValue(uint8(sp.CC), mapToMidi(valueIndex, len(sp.PossibleValues)-1))
	sp.CurrentValueIndex = valueIndex
	sp.CurrentValue = sp.PossibleValues[valueIndex]
	return nil
}

type IntParam struct {
	synth *Synth

	StaticName  string
	dynamicName func() string
	// CC is the MIDI CC number. If it's -1, the parameter is not writable.
	CC              int
	MinValue        int
	MaxValue        int
	bytePosition    int
	byteValueOffset int

	CurrentValue int
	StoredValue  int
}

func (ip *IntParam) Name() string {
	if ip.dynamicName != nil {
		return ip.dynamicName()
	}
	return ip.StaticName
}

func (ip *IntParam) SetValue(value int) error {
	if value < ip.MinValue || value > ip.MaxValue {
		return errors.New("value out of range")
	}
	ip.synth.setValue(uint8(ip.CC), mapToMidi(value, ip.MaxValue))
	ip.CurrentValue = value
	return nil
}

func (s *Synth) initParams() {
	s.SelectionParams = map[string]*SelectionParam{
		"osc1Wave": {
			synth:          s,
			StaticName:     "Osc 1 Wave",
			CC:             24,
			PossibleValues: []string{"Off", "Triangle", "Square", "PWM", "Saw", "SuperSaw", "FM", "Noise"},
			bytePosition:   0,
		},
		"osc2Wave": {
			synth:          s,
			StaticName:     "Osc 2 Wave",
			CC:             25,
			PossibleValues: []string{"Off", "Triangle", "Square", "PWM", "Saw", "Noise"},
			bytePosition:   1,
		},
		"portamentoMode": {
			synth:          s,
			StaticName:     "Portamento Mode",
			CC:             -1,
			PossibleValues: []string{"Portamento", "Glissando"},
			bytePosition:   45,
		},
		"ringModEnabled": {
			synth:          s,
			StaticName:     "Ring Mod",
			CC:             96,
			PossibleValues: []string{"Disabled", "Enabled"},
			bytePosition:   43,
		},
		"lfo1Wave": {
			synth:          s,
			StaticName:     "LFO 1 Wave",
			CC:             54,
			PossibleValues: []string{"Triangle", "Square", "Saw"},
			bytePosition:   47,
		},
		"lfo1Destination": {
			synth:          s,
			StaticName:     "LFO 1 Destination",
			CC:             56,
			PossibleValues: []string{"VCF", "Oscillator"},
			bytePosition:   53,
		},
		"lfo2Wave": {
			synth:          s,
			StaticName:     "LFO 2 Wave",
			CC:             55,
			PossibleValues: []string{"Triangle", "Square", "Saw"},
			bytePosition:   48,
		},
	}
	s.IntParams = map[string]*IntParam{
		"osc1Adjustment": {
			synth:      s,
			StaticName: "Osc 1 PWM/Detune/Feedback",
			dynamicName: func() string {
				osc1Type := s.SelectionParams["osc1Wave"].CurrentValue
				if osc1Type == "PWM" {
					return "Osc 1 Pulse Width"
				} else if osc1Type == "SuperSaw" {
					return "Osc 1 Detune"
				} else if osc1Type == "FM" {
					return "Osc 1 Feedback"
				} else {
					return s.IntParams["osc1Adjustment"].StaticName
				}
			},
			CC:           113,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 2,
		},
		"osc1Coarse": {
			synth:        s,
			StaticName:   "Osc 1 Coarse",
			CC:           115,
			MinValue:     0,
			MaxValue:     24,
			bytePosition: 4,
		},
		"osc1Fine": {
			synth:        s,
			StaticName:   "Osc 1 Fine",
			CC:           111,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 5,
		},
		"osc2Adjustment": {
			synth:        s,
			StaticName:   "Osc 2 Pulse Width",
			CC:           114,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 3,
		},
		"osc2Coarse": {
			synth:        s,
			StaticName:   "Osc 2 Coarse",
			CC:           116,
			MinValue:     0,
			MaxValue:     24,
			bytePosition: 6,
		},
		"osc2Fine": {
			synth:        s,
			StaticName:   "Osc 2 Fine",
			CC:           112,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 7,
		},
		"oscBalance": {
			synth:           s,
			StaticName:      "Osc Balance",
			CC:              29,
			MinValue:        0,
			MaxValue:        63,
			bytePosition:    8,
			byteValueOffset: -0x40,
		},
		"portamentoAmount": {
			synth:        s,
			StaticName:   "Portamento Amount",
			CC:           5,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 46,
		},
		"ringModAmount": {
			synth:        s,
			StaticName:   "Ring Mod Amount",
			CC:           95,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 44,
		},
		"lfo1Speed": {
			synth:        s,
			StaticName:   "LFO 1 Speed",
			CC:           72,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 49,
		},
		"lfo1Amount": {
			synth:        s,
			StaticName:   "LFO 1 Amount",
			CC:           70,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 50,
		},
		"lfo2Speed": {
			synth:        s,
			StaticName:   "LFO 2 Speed",
			CC:           73,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 51,
		},
		"lfo2Amount": {
			synth:        s,
			StaticName:   "LFO 2 Amount",
			CC:           28,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 52,
		},
		"vcfAttack": {
			synth:        s,
			StaticName:   "VCF Attack",
			CC:           85,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 14,
		},
		"vcfDecay": {
			synth:        s,
			StaticName:   "VCF Decay",
			CC:           86,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 15,
		},
		"vcfSustain": {
			synth:        s,
			StaticName:   "VCF Sustain",
			CC:           87,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 16,
		},
		"vcfRelease": {
			synth:        s,
			StaticName:   "VCF Release",
			CC:           88,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 17,
		},
		"vcfAmount": {
			synth:        s,
			StaticName:   "VCF Amount",
			CC:           47,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 22,
		},
		"vcfCutoff": {
			synth:        s,
			StaticName:   "VCF Cutoff",
			CC:           74,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 12,
		},
		"vcfResonance": {
			synth:        s,
			StaticName:   "VCF Resonance",
			CC:           71,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 13,
		},
		"vcaAttack": {
			synth:        s,
			StaticName:   "VCA Attack",
			CC:           81,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 18,
		},
		"vcaDecay": {
			synth:        s,
			StaticName:   "VCA Decay",
			CC:           82,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 19,
		},
		"vcaSustain": {
			synth:        s,
			StaticName:   "VCA Sustain",
			CC:           83,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 20,
		},
		"vcaRelease": {
			synth:        s,
			StaticName:   "VCA Release",
			CC:           84,
			MinValue:     0,
			MaxValue:     99,
			bytePosition: 21,
		},
	}
}
