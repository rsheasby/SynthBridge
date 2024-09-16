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
	err := synth.OpenPorts()
	if err != nil {
		fmt.Println("Error opening ports:", err)
		return
	}
	synth.GetCurrentPatch()
	select {}
}
