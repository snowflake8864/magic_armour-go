package test

import (
	"fmt"
	goEncrypt "github.com/wumansgy/goEncrypt/aes"
	"goskeleton/app/global/variable"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"testing"
	"time"
)

func TestTime(t *testing.T) {
	currentTime := time.Now()
	timeStamp := currentTime.Unix()

	diff := timeStamp - 1587930069

	fmt.Println("Current Time:", currentTime)
	fmt.Println("Timestamp:", timeStamp)
	fmt.Println("diff:", diff)
}

func TestEncrypt(t *testing.T) {
	key := []byte(variable.AES256CbcKey)
	iv := []byte(variable.AES256CbcIv)

	data := "{\"o\":\"1\",\"u\":\"1\",\"t\":1695717588}"
	str, _ := goEncrypt.AesCbcEncryptBase64([]byte(data), key, iv)
	fmt.Println("Encrypted:", str)
}

func TestDecrypt(t *testing.T) {

	key := []byte(variable.AES256CbcKey)
	iv := []byte(variable.AES256CbcIv)

	data := "sa+7XiHnKKDWXWMizcCwEtg6wlsNbuU3xfNiIbTBtCqlftXUjZSUEs4+Jwil0tFK"

	str, err := goEncrypt.AesCbcDecryptByBase64(data, key, iv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(str))

}
