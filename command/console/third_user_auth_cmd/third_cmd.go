package third_user_auth_cmd

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/platform/third_user_auth_models"
	"goskeleton/app/service/platform/third_user_auth_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/redis_factory"
)

func CreateCmdFactory() *ThirdUserAuthCmd {
	return &ThirdUserAuthCmd{
		redisClient: redis_factory.GetOneRedisClient(),
		lockKey:     "lock_third_user_auth",
		cacheExpire: 300,
	}
}

type ThirdUserAuthCmd struct {
	redisClient *redis_factory.RedisClient
	lockKey     string
	cacheExpire int
}

func (crd *ThirdUserAuthCmd) SetLock() bool {
	if isOk, _ := crd.redisClient.Lock(crd.lockKey, date_helper.CurDatetime(), crd.cacheExpire); isOk {
		return true
	} else {
		return false
	}
}

// Sync 同步数据
func (crd *ThirdUserAuthCmd) Sync(isNotice bool) {
	variable.ZapLog.Info(fmt.Sprintf("AddFromProcForTask start isNotice: %v", isNotice))

	param := &third_user_auth_models.AuthModel{
		Addr:      "https://10.32.2.4",
		AppId:     "789546861770502144",
		AppKey:    "56de926f281e4120b932476492b4d451",
		AppSecret: "48642278b1bf4fe28b74fdd65ab799e5",
	}
	third_user_auth_service.CreateRequestFactory(param).PullOrg()

	variable.ZapLog.Info(fmt.Sprintf("AddFromProcForTask end isNotice: %v", isNotice))
}
