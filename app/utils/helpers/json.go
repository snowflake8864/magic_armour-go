package helpers

import "encoding/json"

// JsonEncode 转成json
func JsonEncode(param interface{}) string {
	jsonByte, err := json.Marshal(param)
	if err != nil {
		return ""
	}
	return string(jsonByte)
}

// JsonToMap json还原成map
func JsonToMap(jsonStr string) map[string]interface{} {
	var data map[string]interface{}

	// 进行反序列化操作
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil
	}
	return data
}

// JsonToStruct json还原Struct
// 用法：err := JsonToStruct([]byte(jsonStr), &person);println("Name: ", person.Name)
func JsonToStruct(data []byte, val interface{}) error {
	return json.Unmarshal(data, &val)
}
