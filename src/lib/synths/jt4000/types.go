package jt4000

import (
	"sync"

	"gitlab.com/gomidi/midi/v2/drivers"
)

type Synth struct {
	inPort       drivers.In
	inStop       func()
	outPort      drivers.Out
	Patches      []Patch
	CurrentPatch Patch
	sync.WaitGroup
}

type Patch struct {
	Name string

	Osc1Wave   OscWave
	Osc1Adj    uint8
	Osc1Coarse uint8
	Osc1Fine   uint8

	Osc2Wave   OscWave
	Osc2Adj    uint8
	Osc2Coarse uint8
	Osc2Fine   uint8

	OscBalance uint8

	PortamentoAmount uint8
	PortamentoMode   PortamentoMode

	RingModAmount  uint8
	RingModEnabled bool

	LFO1Wave        LFOWave
	LFO1Destination LFODestination
	LFO1Speed       uint8
	LFO1Amount      uint8

	LFO2Wave   LFOWave
	LFO2Speed  uint8
	LFO2Amount uint8

	VCFAttack  uint8
	VCFDecay   uint8
	VCFSustain uint8
	VCFRelease uint8
	VCFAmount  uint8 // MIDI CC 47

	VCFCutoff    uint8
	VCFResonance uint8

	VCAAttack  uint8
	VCADecay   uint8
	VCASustain uint8
	VCARelease uint8
}

type OscWave byte

const (
	OscWaveOff      OscWave = 0x00
	OscWaveTriangle OscWave = 0x01
	OscWaveSquare   OscWave = 0x02
	OscWavePWM      OscWave = 0x03
	OscWaveSaw      OscWave = 0x04
	OscWaveSuperSaw OscWave = 0x05
	OscWaveFM       OscWave = 0x06
	OscWaveNoise    OscWave = 0x07
)

type PortamentoMode byte

const (
	PortamentoModePortamento PortamentoMode = 0x00
	PortamentoModeGlissando  PortamentoMode = 0x01
)

type LFOWave byte

const (
	LfoWaveTriangle LFOWave = 0x00
	LfoWaveSquare   LFOWave = 0x01
	LfoWaveSaw      LFOWave = 0x02
)

type LFODestination byte

const (
	LfoDestinationVCF LFODestination = 0x00
	LfoDestinationOsc LFODestination = 0x01
)
