package test

import (
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"github.com/stretchr/testify/assert"
	goEncryptAes "github.com/wumansgy/goEncrypt/aes"
	goEncryptRsa "github.com/wumansgy/goEncrypt/rsa"
	"goskeleton/app/global/variable"
	_ "goskeleton/bootstrap" //  为了保证单元测试与正常启动效果一致，记得引入该包
	"log"
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

func TestAesEncrypt(t *testing.T) {
	key := []byte(variable.AES256CbcKey)
	iv := []byte(variable.AES256CbcIv)

	data := "{\"o\":\"1\",\"u\":\"1\",\"t\":1695717588}"
	str, _ := goEncryptAes.AesCbcEncryptBase64([]byte(data), key, iv)
	fmt.Println("Encrypted:", str)
}

func TestAesDecrypt(t *testing.T) {

	key := []byte(variable.AES256CbcKey)
	iv := []byte(variable.AES256CbcIv)

	data := "sa+7XiHnKKDWXWMizcCwEtg6wlsNbuU3xfNiIbTBtCqlftXUjZSUEs4+Jwil0tFK"

	str, err := goEncryptAes.AesCbcDecryptByBase64(data, key, iv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(str))

}

func TestRsaEncrypt(t *testing.T) {

	/*base64Key2048, err := goEncryptRsa.GenerateRsaKeyBase64(2048)
	assert.Nil(t, err)
	fmt.Println(base64Key2048.PublicKey)
	fmt.Println("------------------------------------------------")
	fmt.Println(base64Key2048.PrivateKey)
	fmt.Println("------------------------------------------------")*/

	base64PubKey2048 := "" +
		"MIIBCgKCAQEAwlCBJWEwP97VDdmu1XVxAm+NY9xgLM7sSH6n2CCfgoWQtbNE01tj" +
		"i/zndl/zVnd2Ci06aX8aqssTAgqe3MKV5RIhD/gx0XVpM+5/1igFh2BuxdrFPNO1" +
		"+GfBu6qnnWfPeSCikKI0yY6TKvM9GOiQXca1Sym8oouwY8oV3eWC/Fg5voNmsh78" +
		"iDdsHjB/2YbEEQX+Io++3j3sxzeekNe9K21k0eCH4NoumsZ8S2xsocaEVDNelPYX" +
		"Zv4emKgIuJCO+Qi6LrjTQ5MDXFLc0kWTHAKWwO3iF00V0dee/HEj4DAYQcWFcNIl" +
		"Nwx2FKxXL089ZZbWmrkPiEzXpjVIQwdClQIDAQAB" +
		""
	base64CipherText, err := goEncryptRsa.RsaEncryptToBase64([]byte("root"), base64PubKey2048)
	fmt.Println("用户名:")
	fmt.Println(base64CipherText)
	base64CipherText, err = goEncryptRsa.RsaEncryptToBase64([]byte("veda@2022"), base64PubKey2048)
	fmt.Println("密码:")
	fmt.Println(base64CipherText)

	base64PriKey2048 := "" +
		"MIIEpQIBAAKCAQEAwlCBJWEwP97VDdmu1XVxAm+NY9xgLM7sSH6n2CCfgoWQtbNE" +
		"01tji/zndl/zVnd2Ci06aX8aqssTAgqe3MKV5RIhD/gx0XVpM+5/1igFh2BuxdrF" +
		"PNO1+GfBu6qnnWfPeSCikKI0yY6TKvM9GOiQXca1Sym8oouwY8oV3eWC/Fg5voNm" +
		"sh78iDdsHjB/2YbEEQX+Io++3j3sxzeekNe9K21k0eCH4NoumsZ8S2xsocaEVDNe" +
		"lPYXZv4emKgIuJCO+Qi6LrjTQ5MDXFLc0kWTHAKWwO3iF00V0dee/HEj4DAYQcWF" +
		"cNIlNwx2FKxXL089ZZbWmrkPiEzXpjVIQwdClQIDAQABAoIBAQCUyl5JyCTRE4ki" +
		"5OfqTdZuHgsvR0VhqkTlsg5GAZ0CrdHEQHRYP2d3fv2h+BD1yd3aJw3OUGSUGZIB" +
		"gy5oLpvm3z3v44YKNJSNDnhW/ShWpDHOURms3RQ/FPwU/rUAY8hLoQC6LWqsV/86" +
		"6SyEkThShmC2IExwTvSPROP/mMn2/SUbCHIcKOoYQ6WS/viMjQhkoAyckmKt19MH" +
		"9muNC0eqtsF6Vz3lZnyuZ90yl0NSu88o6VY1b5cQ0xU7fYn0qloFBayRlIpPXDKy" +
		"4AiF4KRt9xDlbZH9seklqP4uexLdcZaDozpWoqVYMRIg9npakzNfRwuu7cMo3YBy" +
		"2sIXDZthAoGBAPBiKBGV+yU6hcOM0SGWiBRi7IG5SzLsrlJz8DW9yYelo+J3juBV" +
		"mo1ADEQ/ddqHSxP4yAuuZ8DK2H/+fhQqbEJl63HlR7YWq4Nq7g4BZC5q2UT7jBmO" +
		"+u9hdA3DUBpluJFQ3IAN40jyPmMY978AAbvZVb4iCSOFJXlOKy9+dz9ZAoGBAM7w" +
		"K4mRo4f2K/qdRriNdLyESrpX4JmiZ/nyGfbWhSnQqYxi+jOAL9n5f5dksszMop4e" +
		"EXAThnbjmGy2oBg+O6ZGR5w1IgYFXxrw/K6/20E+sfwXeF1iAVx4KaMsHUAE54rV" +
		"DE/pl8kzaZTY+5KfwjwkUARl2ICm0krErm+0pJGdAoGBAIdg8Gn/CqML03uvGbSN" +
		"33ODxujEHxq9+usKynSpvw769VI8vSWQcGsub/obt1bsLRdBXABgJ6SHWOi0vfLD" +
		"dZbraW8Kjpl5KPzpQsp2dCJEILrd5Cgcwyv/uWlq3vkO4qRMYvJnX+KcMxCCZjyW" +
		"V2Og4c4jwvvTJ6iqol2sHqJ5AoGAMD3RCHsJMnLzAbOUfwWLBz+X5CUEWQyBTMVh" +
		"xCdyV7MrM20HguYJQblmQzbYb+Hsrwu3aPdES4b3S7CvVbd8VGDh5MqGKf0qlK9F" +
		"WL+lL2qcN3aJOsvUIWhp3H42O2xpRfYwrph4KBFyNZ731tCjiT214bCUbjl1oLSF" +
		"Lclg7AkCgYEAtfbgbkeOyKK1n6BFKwRW7GHkfD0FIRA/x1B6VoLeQ5/bLshnb49N" +
		"DVDmYCEeMDAdYejNRp6V/Ljua2uyWdld9dLA0VqBRspysgcHinJ7z2tjJp2YTVtR" +
		"Zsf2gGUj/YHasTI54jEAlOjt6uUor4hbtSPLy4TsjXmNyWQMDImta48=" +
		""
	// good
	plainText, err := goEncryptRsa.RsaDecryptByBase64(base64CipherText, base64PriKey2048)
	assert.Nil(t, err)
	fmt.Println("解密")
	fmt.Println(string(plainText))

	// 将 DER 编码的公钥转换为 PEM 格式
	pubDER, err := base64.StdEncoding.DecodeString(base64PubKey2048)
	if err != nil {
		log.Fatalf("公钥编码失败: %v", err)
	}
	pubPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: pubDER,
	})

	// 将 DER 编码的公钥转换为 PEM 格式
	privDER, err := base64.StdEncoding.DecodeString(base64PriKey2048)
	if err != nil {
		log.Fatalf("公钥编码失败: %v", err)
	}
	privPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privDER,
	})

	fmt.Println("公钥:")
	fmt.Println(string(pubPEM))
	fmt.Println("私钥:")
	fmt.Println(string(privPEM))
}
