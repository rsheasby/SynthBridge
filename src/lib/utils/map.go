package utils

import "math"

func Uint8Map(x, inMin, inMax, outMin, outMax uint8) uint8 {
	return uint8(math.Ceil(float64(x-inMin)*float64(outMax-outMin)/float64(inMax-inMin) + float64(outMin)))
}
