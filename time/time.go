package time

import "errors"

type SimpleTime struct {
	H int
	M int
	S int
}

func SecondsToHrsMinSec(seconds int) (SimpleTime, error) {
	if seconds < 0 {
		return SimpleTime{}, errors.New("seconds must be >= 0")
	}

	m, s := divmod(seconds, 60)
	h, m := divmod(m, 60)

	time := SimpleTime{
		H: h,
		M: m,
		S: s,
	}
	return time, nil
}

// Takes two numbers as arguments and returns their quotient and remainder when using integer division.
func divmod(x, y int) (quotient, remainder int) {
	quotient = x / y
	remainder = x % y
	return
}
