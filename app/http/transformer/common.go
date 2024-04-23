package transformer

import (
	"goskeleton/app/global/enums"
)

type Common struct {
}

// FormatIntList 格式成标准前端数据格式：val是int
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

// FormatStringList 格式成标准前端数据格式:val是string
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
