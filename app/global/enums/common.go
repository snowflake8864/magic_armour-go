package enums

// 数据字典
const (
	AllDefault int = -100
	IsNot      int = 0 //是否状态: 否
	IsYes      int = 1 //是否状态: 是
	SwitchOn   int = 1 //开关状态: 开
	SwitchOff  int = 0 //开关状态: 否

	AppUsageVersionOffline int = 1 //程序使用版本:离线
	AppUsageVersionCloud   int = 2 //程序使用版本:云

	OperateSourceOrganize int = 1 //操作来源：企业
	OperateSourcePlatform int = 2 //操作来源：平台
	OperateSourceAi       int = 3 //操作来源：AI
)

var EnumMapInt = map[string]map[int]string{
	"operate_source": {
		OperateSourceOrganize: "企业",
		OperateSourcePlatform: "平台",
		OperateSourceAi:       "系统AI",
	},
}

// Get1stIsStrReturnSliceOfStr 如是第一级是string,返回相应[]string
func Get1stIsStrReturnSliceOfStr(first string, s map[string][]string) []string {
	val, ok := s[first]
	if ok == false {
		return nil
	}
	return val
}

// Get1stIsStrReturnSliceOfInt 如是第一级是string,返回相应[]int
func Get1stIsStrReturnSliceOfInt(first string, s map[string][]int) []int {
	val, ok := s[first]
	if ok == false {
		return nil
	}
	return val
}

// Get1stIsStrReturnMapIntOfStr 如是第一级是string,返回相应map[int]string
func Get1stIsStrReturnMapIntOfStr(first string, s map[string]map[int]string) map[int]string {
	val, ok := s[first]
	if ok == false {
		return nil
	}
	return val
}

// Get1stIsStrAnd2ndIsIntReturnStr 如是第一级是string串且第二级是int,返回相应string
func Get1stIsStrAnd2ndIsIntReturnStr(first string, second int, s map[string]map[int]string) string {
	val, ok := s[first][second]
	if ok == false {
		return ""
	}
	return val
}

// Get1stIsStrAnd2ndIsIntReturnInt 如是第一级是string串且第二级是int,返回相应string
func Get1stIsStrAnd2ndIsIntReturnInt(first string, second int, s map[string]map[int]int) int {
	val, ok := s[first][second]
	if ok == false {
		return 0
	}
	return val
}

// Get1stIsStrAnd2ndIsIntReturnSliceInt 如是第一级是string串且第二级是int,返回相应[]int
func Get1stIsStrAnd2ndIsIntReturnSliceInt(first string, second int, s map[string]map[int][]int) []int {
	val, ok := s[first][second]
	if ok == false {
		return nil
	}
	return val
}
