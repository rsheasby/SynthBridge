package main

import (
	"github.com/rsheasby/SynthBridge/lib/controllers/minilab3"
	"gitlab.com/gomidi/midi/v2"
	"gitlab.com/gomidi/midi/v2/drivers"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	inPort, outPort := getPortsByName("Minilab3 MIDI", "Minilab3 MIDI")
	controller, err := minilab3.NewController(inPort, outPort)
	if err != nil {
		panic(err)
	}

	controller.DisplayPad("Title", "value", 63, true)
}

func getPortsByName(inPortName string, outPortName string) (inPort drivers.In, outPort drivers.Out) {
	inPorts := midi.GetInPorts()
	outPorts := midi.GetOutPorts()

	for _, port := range inPorts {
		if port.String() == inPortName {
			inPort = port
			break
		}
	}

	for _, port := range outPorts {
		if port.String() == outPortName {
			outPort = port
			break
		}
	}

	return
}
