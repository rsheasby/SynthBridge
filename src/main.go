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

	for {
		var input string
		fmt.Scanln(&input)
		// Parse input as int
		var val int
		_, err := fmt.Sscanf(input, "%d", &val)
		if err != nil {
			fmt.Println("Invalid input, please enter an integer value.")
			continue
		}

		// Adjust osc1 to provided value
		err = synth.SetOsc1Adj(uint8(val))
		if err != nil {
			fmt.Printf("Error setting Osc1 adjustment: %s\n", err)
		} else {
			fmt.Printf("Osc1 adjustment set to %d\n", val)
		}
	}
}
