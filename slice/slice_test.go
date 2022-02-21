package slice

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUnique(t *testing.T) {
	testcases := []struct {
		name  string
		input []int
		want  []int
	}{
		{"one elem", []int{1}, []int{1}},
		{"one elem multiple times", []int{1, 1, 1}, []int{1}},
		{"1-5", []int{1, 2, 2, 3, 3, 3, 4, 4, 4, 4, 5, 5, 5, 5, 5}, []int{1, 2, 3, 4, 5}},
		{"empty slice", []int{}, nil},
	}

	for _, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			got := UniqueInt(tc.input)

			fmt.Println("got len", len(got), "want len", len(tc.want))

			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

func ExampleUniqueInt() {
	nums := []int{0, 1, 2, 2, 3, 3, 3}
	uniqueNums := UniqueInt(nums)
	fmt.Println(uniqueNums)
	// Output:
	// [0 1 2 3]
}
