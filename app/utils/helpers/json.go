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

// Json2Map json还原成map
func Json2Map(jsonStr string) map[string]interface{} {
	var data map[string]interface{}

	// 进行反序列化操作
	err := json.Unmarshal([]byte(jsonStr), &data)
	if err != nil {
		return nil
	}
	return data
}

// Json2Struct json还原Struct
// 用法：err := Json2Struct([]byte(jsonStr), &person);println("Name: ", person.Name)
func Json2Struct(data []byte, val interface{}) error {
	return json.Unmarshal(data, &val)
}
