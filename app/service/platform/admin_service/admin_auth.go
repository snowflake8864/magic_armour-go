package admin_service

import (
	"encoding/json"
	"errors"
	"goskeleton/app/model/model_platform"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/redis_factory"
)

func CreateAuthFactory() *AdminAuth {
	return &AdminAuth{
		masterModel:   model_platform.CreateAdminFactory(),
		redisClient:   redis_factory.GetOneRedisClient(),
		tokenPrefix:   "admin_auth_token_",
		lockingPrefix: "admin_auth_locking_",
		sliderPrefix:  "admin_auth_slider_",
		sliderExpire:  300,
		smsPrefix:     "admin_auth_sms_",
		smsExpire:     180,
	}
}

type AdminAuth struct {
	masterModel   *model_platform.AdminModel
	redisClient   *redis_factory.RedisClient
	tokenPrefix   string
	lockingPrefix string
	sliderPrefix  string
	sliderExpire  int
	smsPrefix     string
	smsExpire     int
}

type AdminCacheLogin struct {
	AdminId  int64  `json:"admin_id"`
	Username string `json:"username"`
	LoginAt  string `json:"login_at"`
}

func (crd *AdminAuth) keyByToken(val string) string {
	return crd.tokenPrefix + val
}

func (crd *AdminAuth) keyByLocking(val string) string {
	return crd.lockingPrefix + val
}

func (crd *AdminAuth) keyBySlider(val string) string {
	return crd.sliderPrefix + val
}

func (crd *AdminAuth) keyBySms(val string) string {
	return crd.smsPrefix + val
}

// LoginByPw 密码登录
func (crd *AdminAuth) LoginByPw(username, pw string) {

}

// DecryptByLogin 登录信息解密
func (crd *AdminAuth) DecryptByLogin(token string, isDelay bool) (loginInfo AdminCacheLogin, err error) {
	jsonStr := crd.redisClient.Get(crd.keyByToken(token))
	if helpers.StringIsEmpty(jsonStr) {
		return AdminCacheLogin{}, errors.New("token expired")
	}
	err = json.Unmarshal([]byte(jsonStr), &loginInfo)
	if err != nil {
		return AdminCacheLogin{}, errors.New("token json fail")
	}

	if isDelay {
		_, _ = crd.redisClient.Expire(crd.keyByToken(token), 60)
	}
	return loginInfo, nil
}
