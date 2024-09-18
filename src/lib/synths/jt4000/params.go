package jt4000

import "errors"

type SelectionParam struct {
	synth *Synth

	Name string
	// CC is the MIDI CC number. If it's -1, the parameter is not writable.
	CC             int
	PossibleValues []string
	bytePosition   int

	CurrentValue      string
	CurrentValueIndex int
	StoredValue       string
	StoredValueIndex  int
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

	Name string
	// CC is the MIDI CC number. If it's -1, the parameter is not writable.
	CC              int
	MinValue        int
	MaxValue        int
	bytePosition    int
	byteValueOffset int

	CurrentValue int
	StoredValue  int
}

func (s *Synth) initParams() {
	s.SelectionParams = []SelectionParam{
		{
			synth:          s,
			Name:           "Osc 1 Wave",
			CC:             24,
			PossibleValues: []string{"Off", "Triangle", "Square", "PWM", "Saw", "SuperSaw", "FM", "Noise"},
			bytePosition:   0,
		},
		{
			synth:          s,
			Name:           "Osc 2 Wave",
			CC:             25,
			PossibleValues: []string{"Off", "Triangle", "Square", "PWM", "Saw", "Noise"},
			bytePosition:   1,
		},
		{
			synth:          s,
			Name:           "Portamento Mode",
			CC:             -1,
			PossibleValues: []string{"Portamento", "Glissando"},
			bytePosition:   45,
		},
		{
			synth:          s,
			Name:           "Ring Mod",
			CC:             96,
			PossibleValues: []string{"Disabled", "Enabled"},
			bytePosition:   43,
		},
		{
			synth:          s,
			Name:           "LFO 1 Wave",
			CC:             54,
			PossibleValues: []string{"Triangle", "Square", "Saw"},
			bytePosition:   47,
		},
		{
			synth:          s,
			Name:           "LFO 1 Destination",
			CC:             56,
			PossibleValues: []string{"VCF", "Oscillator"},
			bytePosition:   53,
		},
		{
			synth:          s,
			Name:           "LFO 2 Wave",
			CC:             55,
			PossibleValues: []string{"Triangle", "Square", "Saw"},
			bytePosition:   48,
		},
	}
}
