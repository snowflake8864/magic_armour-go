package int_helper

import (
	"math/rand"
	"strconv"
)

// IntRandom 随机数
func IntRandom(minInt, maxInt int) int {
	return minInt + int(float64(maxInt-minInt)*rand.Float64())
}

// IntRandom2Str 随机数
func IntRandom2Str(minInt, maxInt int) string {
	randomInt := IntRandom(minInt, maxInt)

	return strconv.Itoa(randomInt)
}

// Unt32Random 随机数
func Unt32Random(minInt, maxInt uint32) uint32 {
	delta := maxInt - minInt + 1
	if delta == 0 { // 处理 delta 溢出为 0 的情况（即 min=0, max=MaxUint32）
		return minInt
	}
	return minInt + uint32(rand.Uint32()%delta)
}

// Int64ToStr Int64转成字符串
func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}

// IntToStr Int转成字符串
func IntToStr(num int) string {
	return strconv.Itoa(num)
}
