package minilab3

import "fmt"

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
