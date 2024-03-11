package helpers

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
