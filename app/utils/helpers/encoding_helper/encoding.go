package encoding_helper

import (
	"github.com/axgle/mahonia"
	"github.com/saintfish/chardet"
)

// DetectWithSaintfish 使用【saintfish】库: 检测字符编码
func DetectWithSaintfish(raw []byte) string {
	detector := chardet.NewTextDetector()
	result, _ := detector.DetectBest(raw)
	return result.Charset // 返回如"GB-18030", "UTF-8"
}

// ConvertWithMahonia  使用【Mahonia】库: 把指定编码数据转成UTF-8字符串
func ConvertWithMahonia(raw []byte, srcEnc string) string {
	dec := mahonia.NewDecoder(srcEnc)
	return dec.ConvertString(string(raw))
}

// ConvertToUTF8 将已知编码的字节流转换为UTF-8
//
// 参数：
//
//	rawData : 原始字节数据
//
// 返回：
//
//	utf8Data : UTF-8编码数据
//	err : 错误信息
func ConvertToUTF8(rawData []byte) (string, error) {
	// 步骤1：获取对应编码的编码器
	enc := DetectWithSaintfish(rawData)
	if enc == "UTF-8" {
		//return string(rawData), nil
	}

	return ConvertWithMahonia(rawData, enc), nil
}
