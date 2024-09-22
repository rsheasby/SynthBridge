package main

import (
	_ "gitlab.com/gomidi/midi/v2/drivers/rtmididrv"
)

func main() {
	router := Router{}
	router.Run()
}
