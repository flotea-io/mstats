package util

import "math"

func Round(x, unit float64) float64 {
	return math.Round(x/unit) * unit
}
