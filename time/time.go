package time

import (
	"errors"
	"fmt"
	"time"
)

var ErrNegativeSeconds = errors.New("seconds must be >= 0")
var ErrStartDateAfterEndDate = errors.New("start date can't be after end date")

type SimpleTime struct {
	H int
	M int
	S int
}

type CalendarWeek struct {
	Year int
	Week int
}

func SecondsToHrsMinSec(seconds int) (SimpleTime, error) {
	if seconds < 0 {
		return SimpleTime{}, fmt.Errorf("SecondsToHrsMinSec failed: %w", ErrNegativeSeconds)
	}

	m, s := divmod(seconds, 60)
	h, m := divmod(m, 60)

	return SimpleTime{
		H: h,
		M: m,
		S: s,
	}, nil
}

// Takes two numbers as arguments and returns their quotient and remainder when using integer division.
func divmod(x, y int) (quotient, remainder int) {
	quotient = x / y
	remainder = x % y
	return
}

func GetCalendarWeeks(start, end time.Time) ([]CalendarWeek, error) {
	if start.After(end) {
		return nil, ErrStartDateAfterEndDate
	}
	var weeks []CalendarWeek
	for start.Before(end.Add(time.Second)) { // add a second to include week of end date
		year, week := start.ISOWeek()
		weeks = append(weeks, CalendarWeek{year, week})
		oneWeek := 1 * 7 * 24 * time.Hour
		start = start.Add(oneWeek)
	}
	return weeks, nil
}
