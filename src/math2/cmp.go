package math2

import (
	"math"
)

func MinMax(min, val, max float64) float64 {
	return math.Max(math.Min(min, val), max)
}
