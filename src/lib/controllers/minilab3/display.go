package minilab3

import "fmt"

func (c *Controller) initDisplay() error {
	sysexMessage := []byte{0xF0, 0x00, 0x20, 0x6B, 0x7F, 0x42, 0x02, 0x02, 0x40, 0x6A, 0x21, 0xF7}
	return c.outPort.Send(sysexMessage)
}

func (c *Controller) SetPadColor(pad uint8, r, g, b uint8) error {
	if pad < 1 || pad > 8 {
		return fmt.Errorf("pad must be between 1 and 8")
	}

	padId := 3 + pad

	sysexMessage := []byte{0xF0, 0x00, 0x20, 0x6B, 0x7F, 0x42, 0x02, 0x02, 0x16, padId, r, g, b, 0xF7}
	return c.outPort.Send(sysexMessage)
}

type Pictogram uint8

const (
	PictogramNone   Pictogram = 0x00
	PictogramHeart  Pictogram = 0x01
	PictogramPlay   Pictogram = 0x02
	PictogramCircle Pictogram = 0x03
	PictogramNote   Pictogram = 0x04
	PictogramTick   Pictogram = 0x05
)
