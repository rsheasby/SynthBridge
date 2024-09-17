package utils

import (
	"encoding/hex"
	"fmt"

	"github.com/fatih/color"
)

var lastData []byte

func DumpData(data []byte) {
	red := color.New(color.FgRed).SprintFunc()

	dump := hex.Dump(data)
	lastDump := hex.Dump(lastData)
	outDump := ""

	for i := 0; i < len(data) && i < len(lastData); i++ {
		if data[i] != lastData[i] {
			fmt.Printf("Difference at position %d: %02X != %02X\n", i, data[i], lastData[i])
		}
	}

	for i := 0; i < len(dump) && i < len(lastDump); i++ {
		if dump[i] != lastDump[i] {
			outDump += red(string(dump[i]))
		} else {
			outDump += string(dump[i])
		}
	}
	fmt.Println(outDump)

	lastData = data
}
