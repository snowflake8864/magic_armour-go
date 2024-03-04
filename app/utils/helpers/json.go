package helpers

import "encoding/json"

func JsonEncodeStr(param interface{}) string {
	jsonByte, err := json.Marshal(param)
	if err != nil {
		return ""
	}
	return string(jsonByte)
}
