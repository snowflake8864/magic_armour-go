package map_helper

import (
	"sort"
)

// MergeStringAndInterface 合并：键是string且值是interface{}
func MergeStringAndInterface(map1, map2 map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})
	for key, value := range map1 {
		result[key] = value
	}
	for key, value := range map2 {
		result[key] = value
	}
	return result
}

// SortDescIntAndInt 排序（倒叙）：键是int且值是int
func SortDescIntAndInt(map1 map[int]int) [][]int {
	var result [][]int

	var keys []int
	for k := range map1 {
		keys = append(keys, k)
	}

	sort.Slice(keys, func(i, j int) bool {
		return map1[keys[i]] > map1[keys[j]]
	})

	for _, k := range keys {
		result = append(result, []int{k, map1[k]})
	}

	return result
}
