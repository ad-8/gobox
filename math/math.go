package math

import "math"

// RoundTo rounds a number n to the specified amount of decimal places d.
func RoundTo(n float64, d int) float64 {
	return math.Round(n*math.Pow(10, float64(d))) / math.Pow(10, float64(d))
}
