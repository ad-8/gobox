package slice

// UniqueInt returns the unique elements in nums.
func UniqueInt(nums []int) []int { // TODO generic version when Go generics are available
	var unique []int
	var contains = make(map[int]bool)

	for _, num := range nums {
		if _, ok := contains[num]; !ok {
			unique = append(unique, num)
			contains[num] = true
		}
	}
	return unique
}
