package huanshi_service

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/redis_factory"
	"time"
)

func CreateAccessFactory(deviceId string) *accessArmour {
	return &accessArmour{
		privateRsa:  variable.HuanshiPrivateRsa,
		deviceId:    deviceId,
		cacheExpire: 600,
		aesFactory:  openssl_encrypt.CreateAesFactory(),
		redisClient: redis_factory.GetOneRedisClient(),
		cacheKey:    "huanshi_access_token",
	}
}

type accessArmour struct {
	privateRsa  string //约定的私钥
	deviceId    string //系统设备ID
	cacheExpire int    //缓存有效期
	aesFactory  *openssl_encrypt.AesCrypt
	redisClient *redis_factory.RedisClient
	cacheKey    string
}

// 生成秘钥
func (aa *accessArmour) GenerateSecret(timestamp int64) string {
	return md5_encrypt.MD5(aa.privateRsa + aa.deviceId + string(timestamp-1))
}

// 校验秘钥
func (aa *accessArmour) CheckSecret(timestamp int64, secret string) (bool, string) {
	currentTime := time.Now()

	if diff := int(currentTime.Unix() - timestamp); diff > aa.cacheExpire || diff < aa.cacheExpire {
		return false, "秘钥已无效"
	}
	if aa.GenerateSecret(timestamp) != secret {
		return false, "秘钥不正确"
	}

	return true, ""
}

type AccessGenerateToken struct {
	Secret   string `json:"secret"`
	Expire   int    `json:"expire"`
	Token    string `json:"token"`
	DeviceId string `json:"device_id"`
}

// 生成口令
func (aa *accessArmour) GenerateToken(secret string) (AccessGenerateToken, error) {
	token, _ := aa.aesFactory.CbcEncryptToHex(secret)

	isOK, err := aa.redisClient.SetEx(aa.cacheKey, token, aa.cacheExpire)
	if !isOK || err != nil {
		return AccessGenerateToken{}, err
	}

	return AccessGenerateToken{Secret: secret, Expire: aa.cacheExpire, Token: token, DeviceId: aa.deviceId}, nil
}
