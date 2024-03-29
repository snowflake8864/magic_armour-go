package helpers

import (
	"strconv"
)

// SliceString2Map  slice 转为 map[string]struct{}
func SliceString2Map(sl []string) map[string]struct{} {
	set := make(map[string]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// SliceInt2Map 数字类型slice 转为 map[int]struct{}
func SliceInt2Map(sl []int) map[int]struct{} {
	set := make(map[int]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// SliceInt642Map 数字类型slice 转为 map[int]struct{}
func SliceInt642Map(sl []int64) map[int64]struct{} {
	set := make(map[int64]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// SliceAny2Int 将[]interface{} 转成 []int
func SliceAny2Int(sl []interface{}) []int {
	intSlice := make([]int, 0, len(sl))
	for _, item := range sl {
		switch v := item.(type) {
		case int:
			intSlice = append(intSlice, v)
		case float64:
			intSlice = append(intSlice, int(v))
		case string:
			num, err := strconv.Atoi(v)
			if err == nil {
				intSlice = append(intSlice, num)
			}
		}
	}
	return intSlice
}

// SliceAny2Int64 将[]interface{} 转成 []int64
func SliceAny2Int64(sl []interface{}) []int64 {
	intSlice := make([]int64, 0, len(sl))
	for _, item := range sl {
		switch v := item.(type) {
		case int64:
			intSlice = append(intSlice, v)
		case int:
			intSlice = append(intSlice, int64(v))
		case float64:
			intSlice = append(intSlice, int64(v))
		case string:
			num, err := strconv.Atoi(v)
			if err == nil {
				intSlice = append(intSlice, int64(num))
			}
		}
	}
	return intSlice
}

// SliceIntIntersect 2个数字类型slice的交集
func SliceIntIntersect(s1, s2 []int) (intersection []int) {
	if s1 == nil || s2 == nil || len(s1) == 0 || len(s2) == 0 {
		return nil
	}
	// 创建map
	m := SliceInt2Map(s1)

	// 查找交集
	for _, v := range s2 {
		if _, isOk := m[v]; isOk {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

// SliceInt64Intersect 2个数字类型slice的交集
func SliceInt64Intersect(s1, s2 []int64) (intersection []int64) {
	if s1 == nil || s2 == nil || len(s1) == 0 || len(s2) == 0 {
		return nil
	}
	// 创建map
	m := SliceInt642Map(s1)

	// 查找交集
	for _, v := range s2 {
		if _, isOk := m[v]; isOk {
			intersection = append(intersection, v)
		}
	}
	return intersection
}
