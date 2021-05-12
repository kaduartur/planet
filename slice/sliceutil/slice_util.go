package sliceutil

import "sort"

// Merge merges two arrays and remove repeated and empty value.
func Merge(aa []string, bb ...string) []string {
	check := make(map[string]int)
	res := make([]string, 0)
	dd := append(aa, bb...)
	for _, val := range dd {
		check[val] = 1
	}

	for letter, _ := range check {
		if letter == "" {
			continue
		}
		res = append(res, letter)
	}

	sort.Strings(res)

	return res
}
