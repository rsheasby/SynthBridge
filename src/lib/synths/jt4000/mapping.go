package jt4000

import "math"

func mapToMidi(val int, max int) uint8 {
	if val <= 0 {
		return 0
	}
	if max == 63 {
		return map63(val)
	}
	return uint8(math.Ceil(float64(val) * 127.0 / float64(max)))
}

func map63(val int) uint8 {
	// This took me way too fucking long to figure out. Why Uli? Why? You could've just rounded the values.
	result := val*2 + 2
	if result > 127 {
		result = 127
	}
	return uint8(result)
}
