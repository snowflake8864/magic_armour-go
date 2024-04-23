package helpers

import (
	"fmt"
	"github.com/axgle/mahonia"
	"strconv"
	"strings"
)

// StringIsEmpty 字符串去掉空格后是否为空, 空返回true
func StringIsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// StringConvertCoder	 字符串转码
func StringConvertCoder(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	fmt.Println(result)
	return result

}

// String2Int 字符串转成int
func String2Int(str string) int {
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	} else {
		return num
	}
}

// String2Int64 字符串转成int64
func String2Int64(str string) int64 {
	num, err := strconv.ParseInt(str, 10, 64)
	if err != nil {
		return 0
	} else {
		return num
	}
}

// StringPad
// input string 原字符串
// padLength int 规定补齐后的字符串位数
// padString string 自定义填充字符串
// padType string 填充类型:LEFT(向左填充,自动补齐位数), 默认右侧
func StringPad(input string, padLength int, padString string, padType string) string {
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
