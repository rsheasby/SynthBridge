package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rsheasby/SynthBridge/lib/controllers/minilab3"
	"github.com/rsheasby/SynthBridge/lib/synths/jt4000"
)

func RouteMinilab3ToJT4000() (err error) {
	synth, err := jt4000.NewSynth()
	if err != nil {
		return err
	}

	controller, err := minilab3.NewController()
	if err != nil {
		return err
	}

	go routeNotes(controller, synth)
	go routePatchSelection(controller, synth)

	select {}
}

func routeNotes(controller *minilab3.Controller, synth *jt4000.Synth) {
	log.Printf("routing notes")
	for noteMsg := range controller.NoteEvents {
		if noteMsg.Type == minilab3.NoteStart {
			err := synth.SendNoteOn(noteMsg.Key, noteMsg.Velocity)
			if err != nil {
				log.Printf("error sending note on: %v", err)
			}
		} else if noteMsg.Type == minilab3.NoteEnd {
			err := synth.SendNoteOff(noteMsg.Key)
			if err != nil {
				log.Printf("error sending note off: %v", err)
			}
		}
	}
}

func routePatchSelection(controller *minilab3.Controller, synth *jt4000.Synth) {
	const resetDisplayAfter = 3 * time.Second
	selectingIndex := synth.CurrentPatchIndex
	resetDisplayTimer := time.AfterFunc(0, func() {
		selectingIndex = synth.CurrentPatchIndex
		displayCurrentPatch(controller, synth)
	})
	for selectionEvent := range controller.SelectionEvents {
		if selectionEvent.Type == minilab3.SelectionLeft {
			selectingIndex--
			if selectingIndex < 0 {
				selectingIndex = 0
			}
			err := controller.DisplaySelector(synth.PatchNames[selectingIndex], fmt.Sprintf("%d / 32", selectingIndex+1), selectingIndex, 31)
			if err != nil {
				log.Printf("error displaying selector: %v", err)
			}
			resetDisplayTimer.Reset(resetDisplayAfter)
		} else if selectionEvent.Type == minilab3.SelectionRight {
			selectingIndex++
			if selectingIndex >= len(synth.PatchNames) {
				selectingIndex = len(synth.PatchNames) - 1
			}
			err := controller.DisplaySelector(synth.PatchNames[selectingIndex], fmt.Sprintf("%d / 32", selectingIndex+1), selectingIndex, 31)
			if err != nil {
				log.Printf("error displaying selector: %v", err)
			}
			resetDisplayTimer.Reset(resetDisplayAfter)
		} else if selectionEvent.Type == minilab3.SelectionClickDown {
			synth.SetCurrentPatch(selectingIndex)
			resetDisplayTimer.Reset(0)
		}
	}
}

func displayCurrentPatch(controller *minilab3.Controller, synth *jt4000.Synth) {
	controller.DisplayText(synth.CurrentPatchName, minilab3.PictogramNote, "JT-4000", minilab3.PictogramNone)
}
