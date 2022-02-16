package time

import (
	"errors"
	"reflect"
	"testing"
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
