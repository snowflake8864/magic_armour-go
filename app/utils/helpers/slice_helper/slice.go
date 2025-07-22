package slice_helper

import (
	"math/rand"
	"strconv"
)

// String2Map  slice 转为 map[string]struct{}
func String2Map(sl []string) map[string]struct{} {
	set := make(map[string]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// Int2Map 数字类型slice 转为 map[int]struct{}
func Int2Map(sl []int) map[int]struct{} {
	set := make(map[int]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// Int642Map 数字类型slice 转为 map[int]struct{}
func Int642Map(sl []int64) map[int64]struct{} {
	set := make(map[int64]struct{}, len(sl))
	for _, v := range sl {
		set[v] = struct{}{}
	}
	return set
}

// Any2Int 将[]interface{} 转成 []int
func Any2Int(sl []interface{}) []int {
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

// Any2Int64 将[]interface{} 转成 []int64
func Any2Int64(sl []interface{}) []int64 {
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

// Any2String 将[]interface{} 转成 []int64
func Any2String(sl []interface{}) []string {
	strSlice := make([]string, 0, len(sl))
	for _, item := range sl {
		switch v := item.(type) {
		case int64:
			strSlice = append(strSlice, strconv.FormatInt(v, 10))
		case int:
			strSlice = append(strSlice, strconv.Itoa(v))
		case string:
			strSlice = append(strSlice, v)
		}
	}
	return strSlice
}

// IntIntersect 2个数字类型slice的交集
func IntIntersect(s1, s2 []int) (intersection []int) {
	if s1 == nil || s2 == nil || len(s1) == 0 || len(s2) == 0 {
		return nil
	}
	// 创建map
	m := Int2Map(s1)

	// 查找交集
	for _, v := range s2 {
		if _, isOk := m[v]; isOk {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

// Int64Intersect 2个数字类型slice的交集
func Int64Intersect(s1, s2 []int64) (intersection []int64) {
	if s1 == nil || s2 == nil || len(s1) == 0 || len(s2) == 0 {
		return nil
	}
	// 创建map
	m := Int642Map(s1)

	// 查找交集
	for _, v := range s2 {
		if _, isOk := m[v]; isOk {
			intersection = append(intersection, v)
		}
	}
	return intersection
}

// IntRemoveRepeatByMap 通过map主键唯一的特性过滤重复元素
func IntRemoveRepeatByMap(slc []int) []int {
	var result []int
	tempMap := map[int]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// StringRemoveRepeatByMap 通过map主键唯一的特性过滤重复元素
func StringRemoveRepeatByMap(slc []string) []string {
	var result []string
	tempMap := map[string]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// Int64RemoveRepeatByMap 通过map主键唯一的特性过滤重复元素
func Int64RemoveRepeatByMap(slc []int64) []int64 {
	var result []int64
	tempMap := map[int64]byte{} // 存放不重复主键
	for _, e := range slc {
		l := len(tempMap)
		tempMap[e] = 0
		if len(tempMap) != l { // 加入map后，map长度变化，则元素不重复
			result = append(result, e)
		}
	}
	return result
}

// StringIn 是否存在
func StringIn(s string, slc []string) bool {
	toMap := String2Map(slc)

	_, isOk := toMap[s]
	return isOk
}

// ChunkSlice 泛型版本的分割函数（Go 1.18+）
func ChunkSlice[T any](slice []T, chunkSize int) [][]T {
	// 防御性检查
	if chunkSize <= 0 {
		return nil
	}

	// 预计算分块数量
	total := len(slice)
	chunks := make([][]T, 0, (total+chunkSize-1)/chunkSize)

	// 循环切割
	for i := 0; i < total; i += chunkSize {
		end := i + chunkSize
		if end > total {
			end = total
		}
		chunks = append(chunks, slice[i:end])
	}

	return chunks
}

// ShuffleStrings 使用Fisher-Yates算法打乱字符串切片顺序
func ShuffleStrings(slice []string) {
	rand.Shuffle(len(slice), func(i, j int) {
		slice[i], slice[j] = slice[j], slice[i]
	})
}
