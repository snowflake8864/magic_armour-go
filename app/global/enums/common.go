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

// GETFirstIsStringReturnSliceOfString 如是第一级是字符串,返回字符串的切片
func GETFirstIsStringReturnSliceOfString(first string, s map[string][]string) []string {
	val, ok := s[first]
	if ok == false {
		return nil
	}
	return val
}
