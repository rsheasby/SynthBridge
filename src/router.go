package main

import (
	"log"
	"strconv"
	"time"

	"github.com/rsheasby/SynthBridge/lib/controllers/minilab3"
	"github.com/rsheasby/SynthBridge/lib/synths/jt4000"
)

const idleTime = 3 * time.Second

type Router struct {
	controller *minilab3.Controller
	synth      *jt4000.Synth

	NoteEventHandler      NoteEventHandler
	SelectionEventHandler SelectionEventHandler
	KnobEventHandler      KnobEventHandler
	FaderEventHandler     FaderEventHandler
	idleTimer             *time.Timer
}

type NoteEventHandler func(note minilab3.NoteEvent)
type SelectionEventHandler func(selection minilab3.SelectionEvent)
type KnobEventHandler func(knob minilab3.KnobEvent)
type FaderEventHandler func(fader minilab3.FaderEvent)

func NewRouter() (r *Router) {
	controller, err := minilab3.NewController()
	if err != nil {
		log.Fatalf("Failed to create controller: %v", err)
	}

	synth, err := jt4000.NewSynth()
	if err != nil {
		log.Fatalf("Failed to create synth: %v", err)
	}

	r = &Router{
		controller: controller,
		synth:      synth,
	}
	r.idleTimer = time.AfterFunc(0, r.DisplayCurrentPatch)

	return
}

func (r *Router) DisplayCurrentPatch() {
	r.controller.DisplayText(r.synth.CurrentPatchName, minilab3.PictogramNote, "JT-4000", minilab3.PictogramNone)
}

func (r *Router) ResetIdleTimer() {
	if r.idleTimer != nil {
		r.idleTimer.Reset(idleTime)
	}
}

func (r *Router) Run() {
	for {
		select {
		case noteEvent := <-r.controller.NoteEvents:
			r.NoteEventHandler(noteEvent)
		case selectionEvent := <-r.controller.SelectionEvents:
			r.SelectionEventHandler(selectionEvent)
		case knobEvent := <-r.controller.KnobEvents:
			r.KnobEventHandler(knobEvent)
		case faderEvent := <-r.controller.FaderEvents:
			r.FaderEventHandler(faderEvent)
		}
	}
}

func (r *Router) KnobEventDispatcher(handlers []KnobEventHandler) KnobEventHandler {
	return func(knobEvent minilab3.KnobEvent) {
		if knobEvent.KnobNumber > uint8(len(handlers)) {
			log.Printf("No handler for knob %d", knobEvent.KnobNumber)
			return
		}
		handlers[knobEvent.KnobNumber-1](knobEvent)
	}
}

func scaleValue(value int, min int, max int) uint8 {
	return uint8((float64(value-min) / float64(max-min)) * 127.0)
}

func (r *Router) KnobIntSpeedScaler(handler KnobEventHandler, speed float64) KnobEventHandler {
	threshold := 128 / speed
	currentValue := 0.0
	idleTimer := time.AfterFunc(0, func() {
		currentValue = 0.0
	})

	return func(knobEvent minilab3.KnobEvent) {
		idleTimer.Reset(idleTime)
		currentValue += float64(knobEvent.RelativeValue)
		scaledValue := int(currentValue / threshold)
		if scaledValue != 0 {
			handler(minilab3.KnobEvent{KnobNumber: knobEvent.KnobNumber, RelativeValue: scaledValue})
			currentValue -= float64(scaledValue) * threshold // This resets the currentValue but retains the fractional part
		}
	}
}

func (r *Router) KnobIntParamController(paramId string) KnobEventHandler {
	param := r.synth.IntParams[paramId]
	if param == nil {
		log.Fatalf("No such param %s", paramId)
	}

	return func(knobEvent minilab3.KnobEvent) {
		newValue := param.CurrentValue + int(knobEvent.RelativeValue)
		if newValue < param.MinValue {
			newValue = param.MinValue
		} else if newValue > param.MaxValue {
			newValue = param.MaxValue
		}
		proportionalValue := scaleValue(newValue, param.MinValue, param.MaxValue)
		err := param.SetValue(newValue)
		if err != nil {
			log.Printf("Failed to set param %s to %d: %v", paramId, newValue, err)
		}
		r.controller.DisplayKnob(param.Name(), strconv.Itoa(param.CurrentValue), proportionalValue, false)
		r.ResetIdleTimer()
	}
}

func (r *Router) KnobSelectionParamController(paramId string) KnobEventHandler {
	param := r.synth.SelectionParams[paramId]
	if param == nil {
		log.Fatalf("No such param %s", paramId)
	}

	return func(knobEvent minilab3.KnobEvent) {
		newValueIndex := param.CurrentValueIndex + int(knobEvent.RelativeValue)
		if newValueIndex < 0 {
			newValueIndex = 0
		} else if newValueIndex >= len(param.PossibleValues) {
			newValueIndex = len(param.PossibleValues) - 1
		}
		err := param.SetValueByIndex(newValueIndex)
		if err != nil {
			log.Printf("Failed to set param %s to index %d: %v", paramId, newValueIndex, err)
		}
		r.controller.DisplaySelector(param.Name(), param.CurrentValue, newValueIndex, len(param.PossibleValues)-1)
		r.ResetIdleTimer()
	}
}
