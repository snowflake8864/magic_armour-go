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

// ToStruct json还原Struct
// 用法：err := ToStruct([]byte(jsonStr), &person);println("Name: ", person.Name)
func ToStruct(jsonStr string, val interface{}) error {
	return json.Unmarshal([]byte(jsonStr), &val)
}

// ToAny json还原成默认传入数据类型
// 用法：err := ToAny([]byte(jsonStr), &person);
func ToAny(jsonStr string, val interface{}) error {
	return json.Unmarshal([]byte(jsonStr), &val)
}

// ToSliceInt json还原[]int
func ToSliceInt(jsonStr string) []int {
	var numbers []int

	// 进行反序列化操作
	err := json.Unmarshal([]byte(jsonStr), &numbers)
	if err != nil {
		return nil
	}
	return numbers
}

// ToSliceInt64 json还原[]int
func ToSliceInt64(jsonStr string) []int64 {
	var numbers []int64

	// 进行反序列化操作
	err := json.Unmarshal([]byte(jsonStr), &numbers)
	if err != nil {
		return nil
	}
	return numbers
}

// ToSliceString json还原[]string
func ToSliceString(jsonStr string) []string {
	var strList []string

	// 进行反序列化操作
	err := json.Unmarshal([]byte(jsonStr), &strList)
	if err != nil {
		return nil
	}
	return strList
}

// Number2Int64 json.Number to int64
func Number2Int64(number json.Number) int64 {
	if val, err := number.Int64(); err == nil {
		return val
	} else {
		return 0
	}
}

// Number2Int json.Number to int
func Number2Int(number json.Number) int {
	if val, err := number.Int64(); err == nil {
		return int(val)
	} else {
		return 0
	}
}

// Number2Float64 json.Number to float64
func Number2Float64(number json.Number) float64 {
	if val, err := number.Float64(); err == nil {
		return val
	} else {
		return 0
	}
}

// Number2String json.Number to string
func Number2String(number json.Number) string {
	return number.String()
}
