package main

import (
	"github.com/rsheasby/SynthBridge/lib/controllers/minilab3"
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	controller, err := minilab3.NewController()
	if err != nil {
		panic(err)
	}

	controller.DisplayText("Hello", minilab3.PictogramTick, "World", minilab3.PictogramHeart)
}
