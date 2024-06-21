package map_helper

// MergeStringAndInterface 合并：键是字符串且值是interface{}
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
