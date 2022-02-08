package str

import (
	"fmt"
	"time"
)

// BuildFilename builds a filename by prefixing the current date (YYYY-MM-DD) to a given filename fn.
// Example return value: 2020-12-24_foo-bar.csv when called with foo-bar.csv.
func BuildFilename(fn string) string {
	now := time.Now()
	return fmt.Sprintf("%02d-%02d-%02d_%s", now.Year(), now.Month(), now.Day(), fn)
}
