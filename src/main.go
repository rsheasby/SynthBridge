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
		fmt.Print("Enter value: ")
		fmt.Scanf("%d", &val)
		value := uint8(val)

		err := synth.SetPortamentoAmount(value)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		} else {
			fmt.Println("Portamento amount set successfully.")
		}
	}
}
