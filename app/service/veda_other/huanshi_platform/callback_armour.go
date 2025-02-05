package huanshi_platform

import (
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/redis_factory"
	"time"
)

func CreateCallbackFactory(deviceId string) *CallbackArmour {
	return &CallbackArmour{
		privateRsa:  variable.HuanshiPrivateRsa,
		deviceId:    deviceId,
		cacheExpire: 600,
		aesFactory:  openssl_encrypt.CreateAesFactory(),
		redisClient: redis_factory.GetOneRedisClient(),
		cacheKey:    "huanshi_access_token",
		timeDiffMax: 10 * 60 * 60,
	}
}

type CallbackArmour struct {
	privateRsa  string //约定的私钥
	deviceId    string //系统设备ID
	cacheExpire int    //缓存有效期
	aesFactory  *openssl_encrypt.AesCrypt
	redisClient *redis_factory.RedisClient
	cacheKey    string
	timeDiffMax int //时间差最大值
}

// GenerateSecret 生成秘钥
func (aa *CallbackArmour) GenerateSecret(timestamp int64) string {
	return md5_encrypt.MD5(aa.privateRsa + aa.deviceId + int_helper.Int64ToStr(timestamp-1))
}

// CheckSecret 校验秘钥
func (aa *CallbackArmour) CheckSecret(timestamp int64, secret string) (bool, string) {
	currentTime := time.Now()

	if diff := int(currentTime.Unix() - timestamp); diff > aa.timeDiffMax {
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

// GenerateToken 生成口令
func (aa *CallbackArmour) GenerateToken(secret string) (AccessGenerateToken, error) {
	token, _ := aa.aesFactory.CbcEncryptToHex(secret)

	isOK, err := aa.redisClient.SetEx(aa.cacheKey, token, aa.cacheExpire)
	if !isOK || err != nil {
		return AccessGenerateToken{}, err
	}

	return AccessGenerateToken{Secret: secret, Expire: aa.cacheExpire, Token: token, DeviceId: aa.deviceId}, nil
}
