package utils

import "math"

func Map24(val uint8) uint8 {
	return uint8(math.Ceil(float64(val) * 127.0 / 24.0))
}

func Map63(val uint8) uint8 {
	// This took me way too fucking long to figure out. Why Uli? Why? You could've just rounded the values.
	result := val*2 + 2
	if result > 127 {
		result = 127
	}
	return result
}

func Map99(val uint8) uint8 {
	return uint8(math.Ceil(float64(val) * 127.0 / 99.0))
}
