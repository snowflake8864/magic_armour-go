package enum_pw_jump

import (
	"fmt"
	"goskeleton/app/global/enums"
	"math/rand"
)

const (
	CombineNumber      int = 1 //组合: 数字
	CombineLowerLetter int = 2 //组合: 小写字母
	CombineUpperLetter int = 4 //组合: 大写字母
	CombineSpecialChar int = 8 //组合: 特殊字符

	CycleDay   int = 1 //周期: 天
	CycleWeek  int = 2 //周期: 星期
	CycleMonth int = 3 //周期: 月
)

var EnumIntMap = map[string]map[int]string{
	"combine": {
		CombineNumber:      "数字",
		CombineLowerLetter: "小写字母",
		CombineUpperLetter: "大写字母",
		CombineSpecialChar: "特殊字符",
	},
	"cycle": {
		CycleDay:   "每天",
		CycleWeek:  "每周",
		CycleMonth: "每月",
	},
}

var WebMapInt = map[string][]int{
	"combine": {
		CombineNumber,
		CombineLowerLetter,
		CombineUpperLetter,
		CombineSpecialChar,
	},
	"cycle": {
		CycleDay,
		CycleWeek,
		CycleMonth,
	},
}

// PwJumpConfig 密码跳变
type PwJumpConfig struct {
	Switch  int `json:"switch"`
	Combine int `json:"combine"` //组合, 二进制形式
	Size    int `json:"size"`    //长度
	Cycle   int `json:"cycle"`   // 周期
	Timing  int `json:"timing"`  //时机
}

func (c *PwJumpConfig) Format() {
	if c.Switch == 0 {
		c.Switch = enums.SwitchOff
	}

	if c.Combine == 0 {
		c.Combine = CombineNumber + CombineLowerLetter + CombineUpperLetter
	}
	if c.Size == 0 {
		c.Size = 10
	}
	if c.Cycle == 0 {
		c.Cycle = CycleWeek
	}
	if c.Timing == 0 {
		c.Timing = 1
	}
}

func (c *PwJumpConfig) GeneratePw() string {

	var charSets [][]rune
	var password []rune
	if (c.Combine & CombineNumber) > 0 {
		tmp := []rune("0123456789")
		charSets = append(charSets, tmp)
		//确保生成一个选中密码
		password = append(password, tmp[rand.Intn(len(tmp))])
	}
	if (c.Combine & CombineLowerLetter) > 0 {
		tmp := []rune("abcdefghijklmnopqrstuvwxyz")
		charSets = append(charSets, tmp)
		//确保生成一个选中密码
		password = append(password, tmp[rand.Intn(len(tmp))])
	}
	if (c.Combine & CombineUpperLetter) > 0 {
		tmp := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
		charSets = append(charSets, tmp)
		//确保生成一个选中密码
		password = append(password, tmp[rand.Intn(len(tmp))])
	}
	if (c.Combine & CombineSpecialChar) > 0 {
		tmp := []rune("!@#$%^&*()_+")
		charSets = append(charSets, tmp)
		//确保生成一个选中密码
		password = append(password, tmp[rand.Intn(len(tmp))])
	}

	//选中密码组合长度
	selectedTypeLen := len(charSets)
	// 填充剩余长度
	remaining := c.Size - len(password)
	for i := 0; i < remaining; i++ {
		// 从已选类型中随机选择一种
		charset := charSets[rand.Intn(selectedTypeLen)]
		password = append(password, charset[rand.Intn(len(charset))])
	}
	// 打乱字符顺序
	rand.Shuffle(len(password), func(i, j int) {
		password[i], password[j] = password[j], password[i]
	})

	return string(password)
}

// GetCycleTimingMap 获取周期和频率的对应关系
func GetCycleTimingMap(cycle int) []int {
	var list []int
	e := 1
	if cycle == CycleWeek {
		e = 7
	} else if cycle == CycleMonth {
		e = 30
	}
	for i := 1; i <= e; i++ {
		list = append(list, i)
	}
	return list
}

// GetTimingLabel 获取周频率标记
func GetTimingLabel(cycle, timing int) string {
	var data string
	if cycle == CycleWeek {
		var val string
		switch timing {
		case 1:
			val = "一"
		case 2:
			val = "二"
		case 3:
			val = "三"
		case 4:
			val = "四"
		case 5:
			val = "五"
		case 6:
			val = "六"
		case 7:
			val = "日"
		}
		data = fmt.Sprintf("周%v凌晨", val)
	} else if cycle == CycleMonth {
		data = fmt.Sprintf("%d号", timing)
	} else {
		data = "当天凌晨"
	}

	return data
}
