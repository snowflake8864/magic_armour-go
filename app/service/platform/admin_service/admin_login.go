package agent_service

import (
	"goskeleton/app/model/model_platform"
	"goskeleton/app/utils/redis_factory"
)

func CreateLoginFactory() *AdminLogin {
	return &AdminLogin{
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

type AdminLogin struct {
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
	AdminId  int    `json:"admin_id"`
	Username int    `json:"username"`
	LoginAt  string `json:"login_at"`
}
