package minilab3

import (
	"bytes"
	"fmt"
)

func (c *Controller) initDisplay() error {
	sysexMessage := []byte{0xF0, 0x00, 0x20, 0x6B, 0x7F, 0x42, 0x02, 0x02, 0x40, 0x6A, 0x21, 0xF7}
	return c.midiOutPort.Send(sysexMessage)
}

type Pictogram byte

const (
	PictogramNone   Pictogram = 0x00
	PictogramHeart  Pictogram = 0x01
	PictogramPlay   Pictogram = 0x02
	PictogramCircle Pictogram = 0x03
	PictogramNote   Pictogram = 0x04
	PictogramTick   Pictogram = 0x05
)

func (c *Controller) DisplayText(topText string, topPictogram Pictogram, bottomText string, bottomPictogram Pictogram) error {
	var buffer bytes.Buffer

	buffer.Write([]byte{0xF0, 0x00, 0x20, 0x6B, 0x7F, 0x42, 0x04, 0x02, 0x60, 0x1F, 0x07, 0x01})
	buffer.WriteByte(byte(topPictogram))
	buffer.WriteByte(byte(bottomPictogram))
	buffer.Write([]byte{0x01, 0x00, 0x01})

	buffer.WriteString(topText)
	buffer.Write([]byte{0x00, 0x02})

	buffer.WriteString(bottomText)
	buffer.Write([]byte{0x00, 0x02, 0xF7})

	return c.midiOutPort.Send(buffer.Bytes())
}

func (c *Controller) DisplayKnob(topText, bottomText string, knobPosition uint8, autoHide bool) error {
	if knobPosition > 127 {
		return fmt.Errorf("knob position must be between 0 and 127")
	}

	var buffer bytes.Buffer

	autoHideByte := byte(0x00)
	if autoHide {
		autoHideByte = byte(0x02)
	}

	buffer.Write([]byte{0xF0, 0x00, 0x20, 0x6B, 0x7F, 0x42, 0x04, 0x02, 0x60, 0x1F, 0x03, autoHideByte, knobPosition, 0x00, 0x00, 0x01})
	buffer.WriteString(topText)
	buffer.Write([]byte{0x00, 0x02})
	buffer.WriteString(bottomText)
	buffer.WriteByte(0xF7)

	return c.midiOutPort.Send(buffer.Bytes())
}

func (c *Controller) DisplayFader(topText, bottomText string, faderPosition uint8, autoHide bool) error {
	if faderPosition > 127 {
		return fmt.Errorf("fader position must be between 0 and 127")
	}

	var buffer bytes.Buffer

	autoHideByte := byte(0x00)
	if autoHide {
		autoHideByte = byte(0x02)
	}

	buffer.Write([]byte{0xF0, 0x00, 0x20, 0x6B, 0x7F, 0x42, 0x04, 0x02, 0x60, 0x1F, 0x04, autoHideByte, faderPosition, 0x00, 0x00, 0x01})
	buffer.WriteString(topText)
	buffer.Write([]byte{0x00, 0x02})
	buffer.WriteString(bottomText)
	buffer.WriteByte(0xF7)

	return c.midiOutPort.Send(buffer.Bytes())
}

func (c *Controller) DisplayPad(topText, bottomText string, padPosition uint8, autoHide bool) error {
	if padPosition > 127 {
		return fmt.Errorf("pad position must be between 0 and 127")
	}

	var buffer bytes.Buffer

	autoHideByte := byte(0x00)
	if autoHide {
		autoHideByte = byte(0x02)
	}

	buffer.Write([]byte{0xF0, 0x00, 0x20, 0x6B, 0x7F, 0x42, 0x04, 0x02, 0x60, 0x1F, 0x05, autoHideByte, padPosition, 0x00, 0x00, 0x01})
	buffer.WriteString(topText)
	buffer.Write([]byte{0x00, 0x02})
	buffer.WriteString(bottomText)
	buffer.WriteByte(0xF7)

	return c.midiOutPort.Send(buffer.Bytes())
}

type Button int

const (
	ButtonShift Button = iota
	ButtonHold
	ButtonOctaveDown
	ButtonOctaveUp
)

func (c *Controller) SetButtonColor(button Button, r, g, b uint8) error {
	if button < 0 || button > 3 {
		return fmt.Errorf("invalid button")
	}
	// The way Arturia mapped the buttons in the sysex is confusing, so I'm reordering them as left to right, then top to bottom.
	buttonMap := []byte{0x00, 0x02, 0x01, 0x03}

	sysexMessage := []byte{0xF0, 0x00, 0x20, 0x6B, 0x7F, 0x42, 0x02, 0x02, 0x16, buttonMap[button], r, g, b, 0xF7}
	return c.midiOutPort.Send(sysexMessage)
}

func (c *Controller) SetPadColor(pad uint8, r, g, b uint8) error {
	if pad < 1 || pad > 8 {
		return fmt.Errorf("pad must be between 1 and 8")
	}

	padId := 3 + pad

	sysexMessage := []byte{0xF0, 0x00, 0x20, 0x6B, 0x7F, 0x42, 0x02, 0x02, 0x16, padId, r, g, b, 0xF7}
	return c.midiOutPort.Send(sysexMessage)
}
