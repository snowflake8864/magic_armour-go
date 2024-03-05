package helpers

import "encoding/json"

// JsonEncode 转换成json
func JsonEncode(param interface{}) string {
	jsonByte, err := json.Marshal(param)
	if err != nil {
		return ""
	}
	return string(jsonByte)
}
