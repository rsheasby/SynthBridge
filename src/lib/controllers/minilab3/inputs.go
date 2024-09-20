package minilab3

import (
	"fmt"
	"log"

	"gitlab.com/gomidi/midi/v2"
)

func (c *Controller) SetKnobValue(knob uint8, value uint8) error {

	if knob < 1 || knob > 8 {
		return fmt.Errorf("knob must be between 1 and 8")
	}
	if value > 127 {
		return fmt.Errorf("value must be between 0 and 127")
	}

	knobByte := knob + 0x06

	sysexMessage := []byte{0xF0, 0x00, 0x20, 0x6B, 0x7F, 0x42, 0x02, 0x10, 0x00, knobByte, value, 0xF7}
	err := c.midiOutPort.Send(sysexMessage)
	if err != nil {
		return fmt.Errorf("failed to send sysex message: %v", err)
	}
	return nil
}

func (c *Controller) midiMsgListen() (err error) {
	c.midiInStop, err = midi.ListenTo(c.midiInPort, func(msg midi.Message, timestampms int32) {
		var bt []byte
		var ch, key, vel, cc, val uint8
		switch {
		case msg.GetSysEx(&bt):
			log.Println("received sysex")
		case msg.GetControlChange(&ch, &cc, &val):
			c.handleControlChange(ch, cc, val)
		case msg.GetNoteStart(&ch, &key, &vel):
			c.handleNoteStart(ch, key, vel)
		case msg.GetNoteEnd(&ch, &key):
			c.handleNoteEnd(ch, key)
		default:
			// ignore
		}
	}, midi.UseSysEx(), midi.SysExBufferSize(4096))

	if err != nil {
		return fmt.Errorf("failed to listen to MIDI input port: %v", err)
	}
	return
}

type NoteEventType uint8

const (
	NoteStart NoteEventType = iota
	NoteEnd
)

type NoteEvent struct {
	Type     NoteEventType
	Channel  uint8
	Key      uint8
	Velocity uint8
}

func (c *Controller) handleNoteStart(ch, key, vel uint8) {
	c.NoteEvents <- NoteEvent{Type: NoteStart, Channel: ch, Key: key, Velocity: vel}
}

func (c *Controller) handleNoteEnd(ch, key uint8) {
	c.NoteEvents <- NoteEvent{Type: NoteEnd, Channel: ch, Key: key}
}

type SelectionEventType uint8

const (
	SelectionLeft SelectionEventType = iota
	SelectionRight
	SelectionClickDown
	SelectionClickUp
)

type SelectionEvent struct {
	Type    SelectionEventType
	Channel uint8
}

const selectorKnobCC = 28
const selectorButtonCC = 118

func (c *Controller) handleControlChange(ch, cc, val uint8) {
	if cc == selectorKnobCC {
		if val < 63 {
			c.SelectionEvents <- SelectionEvent{Type: SelectionLeft, Channel: ch}
		} else if val > 63 {
			c.SelectionEvents <- SelectionEvent{Type: SelectionRight, Channel: ch}
		}
	} else if cc == selectorButtonCC {
		if val > 0 {
			c.SelectionEvents <- SelectionEvent{Type: SelectionClickDown, Channel: ch}
		} else {
			c.SelectionEvents <- SelectionEvent{Type: SelectionClickUp, Channel: ch}
		}
	}
}
