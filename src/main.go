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

	var val string
	var param = synth.SelectionParams[3]
	fmt.Printf("Updating Param: %s\n", param.Name)
	for {
		fmt.Print("Enter value: ")
		fmt.Scanf("%s", &val)

		err := param.SetValue(val)
		if err != nil {
			fmt.Printf("Error: %s\n", err)
		}
	}
}
