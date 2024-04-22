package transformer

import (
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
	return result
}

// FormatMapStrOStr2List 把map[string]string格式成标准前端数据格式
func (t *Common) FormatMapStrOStr2List(s map[string]string, appendAll bool) (result []map[string]interface{}) {
	if s == nil {
		return result
	}
	if appendAll {
		result = append(result, map[string]interface{}{
			"val":   enums.AllDefaultStr,
			"label": "全部",
		})
	}
	for key, val := range s {
		result = append(result, map[string]interface{}{
			"val":   key,
			"label": val,
		})
	}
	return result
}
