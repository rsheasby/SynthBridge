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
	// spew.Dump(synth)

	var val int
	fmt.Println("Updating Patch number")
	for {
		fmt.Print("Enter value: ")
		fmt.Scanf("%d", &val)

		err := synth.SetCurrentPatch(val)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
		fmt.Printf("Current patch number: %d\nCurrent patch index: %d\nCurrent patch name: %s\n", synth.CurrentPatchNumber(), synth.CurrentPatchIndex, synth.CurrentPatchName)
	}
}
