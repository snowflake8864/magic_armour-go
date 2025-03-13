package enum_pw_jump

import (
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
	Cycle   int `json:"cycle"`   // 周期
	Combine int `json:"combine"` //组合, 二进制形式
	Size    int `json:"size"`    //长度
}

func (c *PwJumpConfig) Format() {
	if c.Switch == 0 {
		c.Switch = enums.SwitchOff
	}

	if c.Cycle == 0 {
		c.Cycle = CycleWeek
	}
	if c.Combine == 0 {
		c.Combine = CombineNumber + CombineLowerLetter + CombineUpperLetter
	}
	if c.Size == 0 {
		c.Size = 10
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
		tmp := []rune("!@#$%^&*()_+-=[]{}|;:',.<>/?")
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
