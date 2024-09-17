package main

import (
	"fmt"

	"github.com/rsheasby/SynthBridge/lib/synths/jt4000"
	"gitlab.com/gomidi/midi/v2"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	inPorts := midi.GetInPorts()
	outPorts := midi.GetOutPorts()

	synth := jt4000.NewSynth(inPorts[0], outPorts[0])
	waves := []jt4000.OscWave{
		jt4000.OscWaveOff,
		jt4000.OscWaveTriangle,
		jt4000.OscWaveSquare,
		jt4000.OscWavePWM,
		jt4000.OscWaveSaw,
		jt4000.OscWaveNoise,
	}

	currentWaveIndex := 0

	for {
		var input string
		fmt.Scanln(&input)
		currentWaveIndex = (currentWaveIndex + 1) % len(waves)
		err := synth.SetOsc2Wave(waves[currentWaveIndex])
		if err != nil {
			fmt.Println("Error setting Osc2 wave:", err)
		} else {
			fmt.Println("Osc2 wave set to:", waves[currentWaveIndex])
		}
	}
}
