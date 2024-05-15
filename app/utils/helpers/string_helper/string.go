package string_helper

import (
	"fmt"
	"github.com/axgle/mahonia"
	"strconv"
	"strings"
)

// IsEmpty 字符串去掉空格后是否为空, 空返回true
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// ConvertCoder	 字符串转码
func ConvertCoder(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	fmt.Println(result)
	return result

}

// ToInt 字符串转成int
func ToInt(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	} else {
		return num
	}
}

// ToInt64 字符串转成int64
func ToInt64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	} else {
		return num
	}
}

// Pad 字符串补齐
// input string 原字符串
// padLength int 规定补齐后的字符串位数
// padString string 自定义填充字符串
// padType string 填充类型:LEFT(向左填充,自动补齐位数), 默认右侧
func Pad(input string, padLength int, padString string, padType string) string {
	output := ""
	inputLen := len(input)

	if inputLen >= padLength {
		return input
	}

	padStringLen := len(padString)
	needFillLen := padLength - inputLen

	if diffLen := padStringLen - needFillLen; diffLen > 0 {
		padString = padString[diffLen:]
	}

	for i := 1; i <= needFillLen; i += padStringLen {
		output += padString
	}
	switch padType {
	case "LEFT":
		return output + input
	default:
		return input + output
	}
}
