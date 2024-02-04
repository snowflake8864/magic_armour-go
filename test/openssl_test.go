package test

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	goEncrypt "github.com/wumansgy/goEncrypt/aes"
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

func encrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	paddedData := PKCS7Padding(data, aes.BlockSize)

	cbc := cipher.NewCBCEncrypter(block, iv)
	cbc.CryptBlocks(paddedData, paddedData)

	return paddedData, nil
}

func decrypt(data []byte, key []byte, iv []byte) ([]byte, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	cbc := cipher.NewCBCDecrypter(block, iv)
	cbc.CryptBlocks(data, data)

	return PKCS7UnPadding(data), nil
}

func adjustKey(key []byte, length int) []byte {
	newKey := make([]byte, length)
	copy(newKey, key)
	return newKey
}

func PKCS5UnPadding(src []byte) []byte {
	length := len(src)
	unpadding := int(src[length-1])
	if length-unpadding < 0 {
		return []byte("")
	}
	return src[:(length - unpadding)]
}

func PKCS7Padding(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

func PKCS7UnPadding(data []byte) []byte {
	length := len(data)
	fmt.Println("length:", length)
	unpadding := int(data[length-1])
	fmt.Println("unpadding:", unpadding)
	return data[:(length - unpadding)]
}

func aesCbcKeyIvPadding(src string, aesLen uint8) []byte {
	srcByte := make([]byte, aesLen)
	srcByteTemp := []byte(src)
	for i := 0; i < len(srcByteTemp); i++ {
		srcByte[i] = srcByteTemp[i]
	}
	return srcByte
}

func TestEncrypt(t *testing.T) {

	//key := []byte("01234567890123456789012345678901") // 32字节密钥
	//iv := []byte("0123456789012345")                  // 16字节偏移量
	key := []byte("S1@HNmhV#q3224Sro6j9eyNO")
	iv := []byte("1234560123456796")
	key = adjustKey(key, 32)
	strKey := "S1@HNmhV#q3224Sro6j9eyNO"
	secretKey := aesCbcKeyIvPadding(strKey, 32)

	data := []byte("{\"o\":\"1\",\"u\":\"1\",\"t\":1695717588}") // 待加密的数据

	encrypted, err := encrypt(data, key, iv)
	fmt.Println(encrypted)
	if err != nil {
		fmt.Println("Error encrypting data:", err)
		return
	}

	fmt.Println("Encrypted:", base64.StdEncoding.EncodeToString(encrypted))

	str, _ := goEncrypt.AesCbcEncryptBase64(data, secretKey, iv)
	fmt.Println("Encrypted:", str)
}

func TestDecrypt(t *testing.T) {

	key := []byte("15538976891f8de623b12d55af7a0f4a")
	iv := []byte("1234560123456796")

	key = adjustKey(key, 32)
	strKey := "15538976891f8de623b12d55af7a0f4a"
	secretKey := aesCbcKeyIvPadding(strKey, 32)

	data := "sa+7XiHnKKDWXWMizcCwEtg6wlsNbuU3xfNiIbTBtCqlftXUjZSUEs4+Jwil0tFK"

	encrypted, _ := base64.StdEncoding.DecodeString(data)
	fmt.Println("E---------------------------")
	decrypted, err := decrypt(encrypted, key, iv)
	if err != nil {
		fmt.Println("Error decrypting data:", err)
		return
	}
	fmt.Println("Decrypted:", string(decrypted))

	str, err := goEncrypt.AesCbcDecryptByBase64(data, secretKey, iv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(str))

}

func TestAesEncrypt(t *testing.T) {

	/*originData, _ := base64.StdEncoding.DecodeString("Xpbzyi5CZLjjE6ViZs7p9Q==")
	iv := []byte("1234567890abcdef")
	key := "abc"*/

	originData, _ := base64.StdEncoding.DecodeString("sa+7XiHnKKDWXWMizcCwEtg6wlsNbuU3xfNiIbTBtCqlftXUjZSUEs4+Jwil0tFK")
	iv := []byte("1234560123456796")
	key := "15538976891f8de623b12d55af7a0f4a"

	fmt.Println(len(key), len(string(iv)))
	//初始化一个空的字节数组,这里的16为AES-128-CBC的要求的key的长度
	keyByte := [32]byte{}
	//转换key为字节数组
	keyByteTemp := []byte(key)
	//依次赋值,这里的3为key的len
	for i := 0; i < len(key); i++ {
		keyByte[i] = keyByteTemp[i]
	}
	fmt.Println(keyByteTemp, keyByte)
	//转换依次得到的keyByteString
	keyByteString := string(keyByte[:])
	fmt.Println(keyByteString)

	//解密
	cipherBlock, err := aes.NewCipher([]byte(keyByteString))
	if err != nil {
		fmt.Println(err)
	}
	cipher.NewCBCDecrypter(cipherBlock, iv).CryptBlocks(originData, originData)
	fmt.Println(string(PKCS5UnPadding(originData)))

}
