package math

import (
	"fmt"
	"testing"
)

func TestRoundTo(t *testing.T) {
	testCases := []struct{
		name string
		got float64
		want float64
	} {
		{"1 digit up", RoundTo(1.25, 1), 1.3},
		{"1 digit down", RoundTo(1.24, 1), 1.2},
		{"1 digit up neg", RoundTo(-1.25, 1), -1.3},
		{"1 digit down neg", RoundTo(-1.24, 1), -1.2},
		{"0 digits", RoundTo(1.5058, 0), 2.0},
		{"0 digits down", RoundTo(1.4958, 0), 1.0},
		{"3 digits", RoundTo(1.5058, 3), 1.506},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func( *testing.T) {
			if tc.got != tc.want {
				t.Errorf("got %v, want %v\n", tc.got, tc.want)
			}
		})
	}
}

func ExampleRoundTo() {
	const PI = 3.14159

	fmt.Println(RoundTo(PI, 0))
	fmt.Println(RoundTo(PI, 1))
	fmt.Println(RoundTo(PI, 3))
	// Output:
	// 3
	// 3.1
	// 3.142
}
