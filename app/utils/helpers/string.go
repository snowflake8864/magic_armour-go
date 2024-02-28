package helpers

import (
	"fmt"
	"github.com/axgle/mahonia"
	"strings"
)

// IsEmptyString 字符串去掉空格后是否为空, 空返回true
func IsEmptyString(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

// ConvertToString 字符串转码
func ConvertToString(src string, srcCode string, tagCode string) string {
	srcCoder := mahonia.NewDecoder(srcCode)
	srcResult := srcCoder.ConvertString(src)
	tagCoder := mahonia.NewDecoder(tagCode)
	_, cdata, _ := tagCoder.Translate([]byte(srcResult), true)
	result := string(cdata)
	fmt.Println(result)
	return result

}
