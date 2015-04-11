package genutil

import (
	"sort"
)

// SortedKeys_String2Int returns sorted keys of that type
func SortedKeys_String2Int(_mp *map[string]int) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// SortedKeys_String2Int64 returns sorted keys of that type
func SortedKeys_String2Int64(_mp *map[string]int64) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// SortedKeys_String2Float64 returns sorted keys of that type
func SortedKeys_String2Float64(_mp *map[string]float64) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// SortedKeys_String2Bool returns sorted keys of that type
func SortedKeys_String2Bool(_mp *map[string]bool) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// SortedKeys_String2String returns sorted keys of that type
func SortedKeys_String2String(_mp *map[string]string) []string {
	keys := make([]string, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Strings(keys)
	return keys
}

// SortedKeys_Int2Int returns sorted keys of that type
func SortedKeys_Int2Int(_mp *map[int]int) []int {
	keys := make([]int, len(*_mp))
	ii := 0
	for kk := range *_mp {
		keys[ii] = kk
		ii++
	}
	sort.Ints(keys)
	return keys
}
