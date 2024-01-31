package enums

// 数据字典
const (
	SearchAll   int = -100
	IsNot       int = 0 //是否状态: 否
	IsYes       int = 1 //是否状态: 是
	SwitchOn    int = 1 //开关状态: 开
	SwitchOff   int = 0 //开关状态: 否
	ModeMonitor int = 0 //模式：监控
	ModeProject int = 1 //模式：保护
)

// GET1stIsStrReturnSliceOfStr 如是第一级是string,返回对应的string类型切片
func GET1stIsStrReturnSliceOfStr(first string, s map[string][]string) []string {
	val, ok := s[first]
	if ok == false {
		return nil
	}
	return val
}

// GET1stIsStrReturnSliceOfInt 如是第一级是string,返回对应的int类型切片
func GET1stIsStrReturnSliceOfInt(first string, s map[string][]int) []int {
	val, ok := s[first]
	if ok == false {
		return nil
	}
	return val
}

// GET1stIsStrAnd2ndIsIntReturnStr 如是第一级是string串且第二级是int,返回对应的string
func GET1stIsStrAnd2ndIsIntReturnStr(first string, second int, s map[string]map[int]string) string {
	val, ok := s[first][second]
	if ok == false {
		return ""
	}
	return val
}

// GET1stIsStrAnd2ndIsIntReturnInt 如是第一级是string串且第二级是int,返回对应的string
func GET1stIsStrAnd2ndIsIntReturnInt(first string, second int, s map[string]map[int]int) int {
	val, ok := s[first][second]
	if ok == false {
		return 0
	}
	return val
}
