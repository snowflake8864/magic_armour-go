package platform_trans

import (
	"goskeleton/app/global/enums"
)

type Common struct {
}

// FormatMapIntStr2List 把map[int]string格式成标准前端数据格式
func (t *Common) FormatMapIntStr2List(s map[int]string, appendAll bool) (result []map[string]interface{}) {
	if s == nil {
		return result
	}
	if appendAll {
		result = append(result, map[string]interface{}{
			"val":   enums.AllDefault,
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
