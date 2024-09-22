package main

import (
	"fmt"
	"log"
	"time"

	"github.com/rsheasby/SynthBridge/lib/controllers/minilab3"
	"github.com/rsheasby/SynthBridge/lib/synths/jt4000"
)

const resetDisplayAfter = 3 * time.Second

// first 8 are knobs, last 4 are faders
var controlMap = []string{
	"osc1Wave",
	"osc1Adjustment",
	"osc1Coarse",
	"osc1Fine",
	"lfo1Wave",
	"lfo1Destination",
	"lfo1Speed",
	"lfo1Amount",
	"vcaAttack",
	"vcaDecay",
	"vcaSustain",
	"vcaRelease",
}

type Router struct {
	controller        *minilab3.Controller
	synth             *jt4000.Synth
	selectingIndex    int
	resetDisplayTimer *time.Timer
}

func (r *Router) Run() (err error) {
	synth, err := jt4000.NewSynth()
	if err != nil {
		return err
	}
	r.synth = synth

	controller, err := minilab3.NewController()
	if err != nil {
		return err
	}
	r.controller = controller

	r.resetDisplayTimer = time.AfterFunc(0, r.resetDisplayFunc)

	go r.routeNotes()
	go r.routePatchSelection()
	go r.routeControls()
	select {}
}

func (r *Router) resetDisplayFunc() {
	r.selectingIndex = r.synth.CurrentPatchIndex
	r.displayCurrentPatch()
}

func (r *Router) routeNotes() {
	log.Printf("routing notes")
	for noteMsg := range r.controller.NoteEvents {
		if noteMsg.Type == minilab3.NoteStart {
			err := r.synth.SendNoteOn(noteMsg.Key, noteMsg.Velocity)
			if err != nil {
				log.Printf("error sending note on: %v", err)
			}
		} else if noteMsg.Type == minilab3.NoteEnd {
			err := r.synth.SendNoteOff(noteMsg.Key)
			if err != nil {
				log.Printf("error sending note off: %v", err)
			}
		}
	}
}

func (r *Router) routePatchSelection() {
	r.selectingIndex = r.synth.CurrentPatchIndex
	for selectionEvent := range r.controller.SelectionEvents {
		if selectionEvent.Type == minilab3.SelectionLeft {
			r.selectingIndex--
			if r.selectingIndex < 0 {
				r.selectingIndex = 0
			}
			err := r.controller.DisplaySelector(r.synth.PatchNames[r.selectingIndex], fmt.Sprintf("%d / 32", r.selectingIndex+1), r.selectingIndex, 31)
			if err != nil {
				log.Printf("error displaying selector: %v", err)
			}
			r.resetDisplayTimer.Reset(resetDisplayAfter)
		} else if selectionEvent.Type == minilab3.SelectionRight {
			r.selectingIndex++
			if r.selectingIndex >= len(r.synth.PatchNames) {
				r.selectingIndex = len(r.synth.PatchNames) - 1
			}
			err := r.controller.DisplaySelector(r.synth.PatchNames[r.selectingIndex], fmt.Sprintf("%d / 32", r.selectingIndex+1), r.selectingIndex, 31)
			if err != nil {
				log.Printf("error displaying selector: %v", err)
			}
			r.resetDisplayTimer.Reset(resetDisplayAfter)
		} else if selectionEvent.Type == minilab3.SelectionClickDown {
			r.synth.SetCurrentPatch(r.selectingIndex)
			err := r.synth.GetCurrentPatchDetails()
			if err != nil {
				log.Printf("error getting current patch details: %v", err)
			} else {
				for i := 1; i <= 8; i++ {
					paramName := controlMap[i-1]
					if selectionParam, exists := r.synth.SelectionParams[paramName]; exists {
						mappedValue := uint8((float64(selectionParam.CurrentValueIndex) / float64(len(selectionParam.PossibleValues)-1)) * 127.0)
						r.controller.SetKnobValue(uint8(i), mappedValue)
					} else if intParam, exists := r.synth.IntParams[paramName]; exists {
						r.controller.SetKnobValue(uint8(i), uint8((float64(intParam.CurrentValue-intParam.MinValue)/float64(intParam.MaxValue-intParam.MinValue))*127.0))
					}
				}
			}
			r.resetDisplayTimer.Reset(0)
		}
	}
}

func (r *Router) displayCurrentPatch() {
	r.controller.DisplayText(r.synth.CurrentPatchName, minilab3.PictogramNote, "JT-4000", minilab3.PictogramNone)
}

func (r *Router) routeControls() {
	for controlEvent := range r.controller.ControlEvents {
		// Locate the correct parameter based on the control type and index
		var selectionParam *jt4000.SelectionParam = nil
		var intParam *jt4000.IntParam = nil
		controlIndex := int(controlEvent.InputNumber) - 1
		if controlEvent.ControlType == minilab3.ControlFader {
			controlIndex += 8
		}
		if controlIndex >= len(controlMap) {
			log.Printf("invalid control index: %d", controlIndex)
			continue
		}
		paramName := controlMap[controlIndex]
		param, isSelectionParam := r.synth.SelectionParams[paramName]
		if isSelectionParam {
			selectionParam = param
		} else {
			param, isIntParam := r.synth.IntParams[paramName]
			if isIntParam {
				intParam = param
			}
		}

		// Set the parameter value based on the control type and index
		if selectionParam != nil {
			newValueIndex := int((float64(controlEvent.Value) / 127.0) * float64(len(selectionParam.PossibleValues)-1))
			selectionParam.SetValueByIndex(newValueIndex)
			r.controller.DisplaySelector(selectionParam.Name(), selectionParam.CurrentValue, newValueIndex, len(selectionParam.PossibleValues)-1)
			r.resetDisplayTimer.Reset(resetDisplayAfter)
		} else if intParam != nil {
			mappedValue := int((float64(controlEvent.Value)/127.0)*float64(intParam.MaxValue-intParam.MinValue) + float64(intParam.MinValue))
			intParam.SetValue(mappedValue)
			if controlEvent.ControlType == minilab3.ControlKnob {
				r.controller.DisplayKnob(intParam.Name(), fmt.Sprintf("%d", mappedValue), controlEvent.Value, false)
			} else {
				r.controller.DisplayFader(intParam.Name(), fmt.Sprintf("%d", mappedValue), controlEvent.Value, false)
			}
			r.resetDisplayTimer.Reset(resetDisplayAfter)
		}
	}
}
