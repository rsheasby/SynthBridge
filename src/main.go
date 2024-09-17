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

	var val int
	for {
		// fmt.Scanln()
		// val = (val + 1) % 100

		// // Set the oscillator 1 fine adjustment to the provided value
		// err := synth.SetOsc1Adj(uint8(val))
		// if err != nil {
		// 	fmt.Printf("Error setting oscillator 1 fine adjustment: %s\n", err)
		// } else {
		// 	fmt.Printf("Set oscillator 1 fine adjustment to %d\n", val)
		// }

		// fmt.Scanln()
		// val = (val + 1) % 25

		// // Set the oscillator 1 fine adjustment to the provided value
		// err := synth.SetOsc1Coarse(uint8(val))
		// if err != nil {
		// 	fmt.Printf("Error setting oscillator 1 coarse adjustment: %s\n", err)
		// } else {
		// 	fmt.Printf("Set oscillator 1 coarse adjustment to %d\n", val)
		// }

		fmt.Scanln()
		val = (val + 1) % 64

		// Set the oscillator balance to the provided value
		err := synth.SetOscBalance(uint8(val))
		if err != nil {
			fmt.Printf("Error setting oscillator balance: %s\n", err)
		} else {
			fmt.Printf("Set oscillator balance to %d\n", val)
		}
	}
}
