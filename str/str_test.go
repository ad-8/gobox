package str

import (
	"fmt"
	"testing"
	"time"
)

func TestBuildFilename(t *testing.T) {
	now := time.Now()

	testcases := []struct {
		name, got, want string
	}{
		{"basic",
			BuildFilename("foo.go"),
			fmt.Sprintf("%02d-%02d-%02d_%s", now.Year(), now.Month(), now.Day(), "foo.go"),
		},
		{
			"empty",
			BuildFilename(""),
			fmt.Sprintf("%02d-%02d-%02d_%s", now.Year(), now.Month(), now.Day(), ""),
		},
		{
			"whitespace",
			BuildFilename("foo bar baz.txt"),
			fmt.Sprintf("%02d-%02d-%02d_%s", now.Year(), now.Month(), now.Day(), "foo bar baz.txt"),

		},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.got != tc.want {
				t.Errorf("got %q but wanted %q", tc.got, tc.want)
			}
		})
	}
}
