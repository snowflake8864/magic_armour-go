package enums

// 数据字典
const (
	SearchAll int = -100
	IsNot     int = 0 //是否状态: 否
	IsYes     int = 1 //是否状态: 是
	SwitchOn  int = 1 //开关状态: 开
	SwitchOff int = 0 //开关状态: 否

)

// Get1stIsStrReturnSliceOfStr 如是第一级是string,返回对应的string类型切片
func Get1stIsStrReturnSliceOfStr(first string, s map[string][]string) []string {
	val, ok := s[first]
	if ok == false {
		return nil
	}
	return val
}

// Get1stIsStrReturnSliceOfInt 如是第一级是string,返回对应的int类型切片
func Get1stIsStrReturnSliceOfInt(first string, s map[string][]int) []int {
	val, ok := s[first]
	if ok == false {
		return nil
	}
	return val
}

// Get1stIsStrAnd2ndIsIntReturnStr 如是第一级是string串且第二级是int,返回对应的string
func Get1stIsStrAnd2ndIsIntReturnStr(first string, second int, s map[string]map[int]string) string {
	val, ok := s[first][second]
	if ok == false {
		return ""
	}
	return val
}

// Get1stIsStrAnd2ndIsIntReturnInt 如是第一级是string串且第二级是int,返回对应的string
func Get1stIsStrAnd2ndIsIntReturnInt(first string, second int, s map[string]map[int]int) int {
	val, ok := s[first][second]
	if ok == false {
		return 0
	}
	return val
}

// Get1stIsStrAnd2ndIsIntReturnSliceInt 如是第一级是string串且第二级是int,返回对应的[]int
func Get1stIsStrAnd2ndIsIntReturnSliceInt(first string, second int, s map[string]map[int][]int) []int {
	val, ok := s[first][second]
	if ok == false {
		return nil
	}
	return val
}
