package helpers

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

// Int64ToStr Int64转成字符串
func Int64ToStr(num int64) string {
	return strconv.FormatInt(num, 10)
}
