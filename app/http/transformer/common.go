package transformer

import (
	"fmt"
	"goskeleton/app/global/enums"
)

type Common struct {
}

// FormatMapIntOfStr2List 把map[int]string格式成标准前端数据格式
func (t *Common) FormatMapIntOfStr2List(s map[int]string, appendAll bool) (result []map[string]interface{}) {
	if s == nil {
		return result
	}
	if appendAll {
		result = append(result, map[string]interface{}{
			"val":   enums.AllDefaultInt,
			"label": "全部",
		})
	}
	for key, val := range s {
		result = append(result, map[string]interface{}{
			"val":   key,
			"label": val,
		})
	}
	fmt.Println()
	return result
}

// FormatIntList 格式成标准前端数据格式
func (t *Common) FormatIntList(keys []int, labels map[int]string, appendAll bool) (result []map[string]interface{}) {
	if keys == nil {
		return result
	}
	if appendAll {
		result = append(result, map[string]interface{}{
			"val":   enums.AllDefaultInt,
			"label": "全部",
		})
	}
	for _, val := range keys {
		if label, isOk := labels[val]; isOk {
			result = append(result, map[string]interface{}{
				"val":   val,
				"label": label,
			})
		}

	}
	return result
}

// FormatStringList 格式成标准前端数据格式
func (t *Common) FormatStringList(keys []string, labels map[string]string, appendAll bool) (result []map[string]interface{}) {
	if keys == nil {
		return result
	}
	if appendAll {
		result = append(result, map[string]interface{}{
			"val":   enums.AllDefaultStr,
			"label": "全部",
		})
	}
	for _, val := range keys {
		if label, isOk := labels[val]; isOk {
			result = append(result, map[string]interface{}{
				"val":   val,
				"label": label,
			})
		}

	}
	return result
}
