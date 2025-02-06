package huanshi_platform

import (
	"errors"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/openssl_encrypt"
	"goskeleton/app/utils/redis_factory"
	"time"
)

func CreateCallbackFactory() *CallbackArmour {
	return &CallbackArmour{
		privateRsa:  variable.HuanshiPrivateRsa,
		cacheExpire: 120,
		redisClient: redis_factory.GetOneRedisClient(),
		cacheKey:    "huanshi_access_token",
		timeDiffMax: 10 * 60 * 60,
	}
}

type CallbackArmour struct {
	privateRsa  string //约定的私钥
	cacheExpire int    //缓存有效期
	redisClient *redis_factory.RedisClient
	cacheKey    string
	timeDiffMax int //时间差最大值
}

// GenerateSecret 生成秘钥
func (aa *CallbackArmour) GenerateSecret(deviceId string, timestamp int64) string {
	return md5_encrypt.MD5(aa.privateRsa + deviceId + int_helper.Int64ToStr(timestamp-1))
}

// CheckSecret 校验秘钥
func (aa *CallbackArmour) CheckSecret(deviceId string, timestamp int64, secret string) (bool, string) {
	currentTime := time.Now()

	if diff := int(currentTime.Unix() - timestamp); diff > aa.timeDiffMax {
		return false, "秘钥已无效"
	}
	if aa.GenerateSecret(deviceId, timestamp) != secret {
		return false, "秘钥不正确"
	}

	return true, ""
}

// GenerateToken 生成口令
func (aa *CallbackArmour) GenerateToken(deviceId, secret string) (interface{}, error) {
	token, _ := openssl_encrypt.CreateRsaFactory().EncryptToHex(deviceId)

	isOK, err := aa.redisClient.SetEx(aa.cacheKey, token, aa.cacheExpire)
	if !isOK || err != nil {
		return nil, err
	}

	return gin.H{
		"secret": secret,
		"expire": aa.cacheExpire,
		"token":  token,
	}, nil
}

func (aa *CallbackArmour) GetDeviceIdByToken(token string) string {
	if deviceId, err := openssl_encrypt.CreateRsaFactory().DecryptByHex(token); err != nil {
		return ""
	} else {
		return deviceId
	}
}

// CheckToken 校验口令并发货设备标识
func (aa *CallbackArmour) CheckToken(token string) (string, error) {
	val := aa.redisClient.Get(aa.cacheKey)
	if string_helper.IsEmpty(val) {
		return "", errors.New("token expired")
	}
	if val != token {
		return "", errors.New("token different")
	}

	if deviceId := aa.GetDeviceIdByToken(token); string_helper.IsEmpty(deviceId) {
		return "", errors.New("token decrypt fail")
	} else {
		_, _ = aa.redisClient.Del(aa.cacheKey)

		return deviceId, nil
	}
}
