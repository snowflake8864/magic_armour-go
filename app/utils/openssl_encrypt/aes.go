package openssl_encrypt

import (
	goEncryptAes "github.com/wumansgy/goEncrypt/aes"
	"strings"
)

func CreateAesFactory() *AesCrypt {
	return &AesCrypt{
		key: []byte("15538976891f8de623b12d55af7a0f4a"), //md5("S1@HNmhV#q3224Sro6j9eyNO")
		iv:  []byte("1234560123456796"),
	}
}

type AesCrypt struct {
	key []byte
	iv  []byte
}

// CbcEncryptToBase64 cbc方式加密， 返回Base64
func (a *AesCrypt) CbcEncryptToBase64(data string) (result string, err error) {
	result, err = goEncryptAes.AesCbcEncryptBase64([]byte(data), a.key, a.iv)
	if err != nil {
		return "", err
	}
	return result, nil
}

// CbcDecryptByBase64 解密 CbcEncryptToBase64
func (a *AesCrypt) CbcDecryptByBase64(data string) (string, error) {
	data = strings.Replace(data, " ", "+", -1)
	plainText, err := goEncryptAes.AesCbcDecryptByBase64(data, a.key, a.iv)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}

// CbcEncryptToHex cbc方式加密， 返回Hex
func (a *AesCrypt) CbcEncryptToHex(data string) (result string, err error) {
	result, err = goEncryptAes.AesCbcEncryptHex([]byte(data), a.key, a.iv)
	if err != nil {
		return "", err
	}
	return result, nil
}

// AesCbcDecryptByHex 解密 CbcEncryptToHex
func (a *AesCrypt) AesCbcDecryptByHex(data string) (string, error) {
	plainText, err := goEncryptAes.AesCbcDecryptByHex(data, a.key, a.iv)
	if err != nil {
		return "", err
	}
	return string(plainText), nil
}
