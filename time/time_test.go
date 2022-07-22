package time

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestSecondsToHrsMinSec(t *testing.T) {
	testCases := []struct {
		name string
		got  int
		want SimpleTime
	}{
		{"00:00:00", 0, SimpleTime{0, 0, 0}},
		{"1:02:03", 3723, SimpleTime{1, 2, 3}},
		{"3:25:45", 12345, SimpleTime{3, 25, 45}},
		{"34293:33:9", 123456789, SimpleTime{34293, 33, 9}},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got, err := SecondsToHrsMinSec(tc.got)
			if err != nil {
				t.Errorf("got error %q, but wanted nil", err)
			}
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("\ngot %v\nwant %v", tc.got, tc.want)
			}
		})
	}
}

func TestSecondsToHrsMinSecWithNegativeInput(t *testing.T) {
	_, err := SecondsToHrsMinSec(-1)

	if !errors.Is(err, ErrNegativeSeconds) {
		t.Errorf("got '%v' but wanted '%v'", err, ErrNegativeSeconds)
	}
}

func TestGetCalendarWeeks(t *testing.T) {
	start := time.Date(2021, 12, 20, 0, 0, 0, 0, time.Local)
	end := time.Date(2022, 03, 07, 0, 0, 0, 0, time.Local)
	want := []CalendarWeek{{2021 ,51}, {2021, 52}, {2022, 1}, {2022, 2},
		{2022, 3}, {2022, 4}, {2022, 5}, {2022, 6},
		{2022, 7}, {2022, 8}, {2022, 9}, {2022, 10}}
	got, _ := GetCalendarWeeks(start, end)
	if !reflect.DeepEqual(got, want) {
		t.Errorf("\ngot %v\nwant %v", got, want)
	}
}

func TestGetCalendarWeeksStartAfterEnd(t *testing.T) {
	start := time.Date(2021, 12, 20, 0, 0, 0, 0, time.Local)
	end := time.Date(2021, 12, 19, 0, 0, 0, 0, time.Local)

	_, err := GetCalendarWeeks(start, end)
	if !errors.Is(err, ErrStartDateAfterEndDate) {
		t.Errorf("got '%v' but wanted '%v'", err, ErrStartDateAfterEndDate)
	}
}
