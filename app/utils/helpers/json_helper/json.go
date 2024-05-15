package json_helper

import "encoding/json"

// Encode 转成json
func Encode(param interface{}) string {
	jsonByte, err := json.Marshal(param)
	if err != nil {
		return ""
	}
	return string(jsonByte)
}

// ToMap json还原成map
func ToMap(jsonStr string) map[string]interface{} {
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

// Json2Any json还原成默认传入数据类型
// 用法：err := Json2Any([]byte(jsonStr), &person);
func Json2Any(jsonStr string, val interface{}) error {
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
