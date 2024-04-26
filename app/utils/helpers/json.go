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
func Json2Struct(jsonStr string, val interface{}) error {
	return json.Unmarshal([]byte(jsonStr), &val)
}

// Json2SliceInt json还原[]int
func Json2SliceInt(jsonStr string) []int {
	var numbers []int

	// 进行反序列化操作
	err := json.Unmarshal([]byte(jsonStr), &numbers)
	if err != nil {
		return nil
	}
	return numbers
}

// Json2SliceInt64 json还原[]int
func Json2SliceInt64(jsonStr string) []int64 {
	var numbers []int64

	// 进行反序列化操作
	err := json.Unmarshal([]byte(jsonStr), &numbers)
	if err != nil {
		return nil
	}
	return numbers
}

// Json2SliceString json还原[]string
func Json2SliceString(jsonStr string) []string {
	var strList []string

	// 进行反序列化操作
	err := json.Unmarshal([]byte(jsonStr), &strList)
	if err != nil {
		return nil
	}
	return strList
}
