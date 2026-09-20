package third_user_auth_cmd

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/platform/third_user_auth_models"
	"goskeleton/app/service/platform/third_user_auth_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/redis_factory"
)

func CreateCmdFactory() *ThirdUserAuthCmd {
	return &ThirdUserAuthCmd{
		redisClient: redis_factory.GetOneRedisClient(),
		lockKey:     "lock_third_user_auth",
		cacheExpire: 600,
	}
}

type ThirdUserAuthCmd struct {
	redisClient    *redis_factory.RedisClient
	lockKey        string
	cacheExpire    int
	requestFactory *third_user_auth_service.RequestThird
	configFactory  *third_user_auth_service.Config
	authConfig     *third_user_auth_models.AuthModel
}

// 设置锁
func (crd *ThirdUserAuthCmd) setLock() bool {
	if isOk, _ := crd.redisClient.Lock(crd.lockKey, date_helper.CurDatetime(), crd.cacheExpire); isOk {
		return true
	} else {
		return false
	}
}

// 删除锁
func (crd *ThirdUserAuthCmd) delLock() bool {
	if _, err := crd.redisClient.Del(crd.lockKey); err == nil {
		return true
	} else {
		return false
	}
}

// 初始化工厂
func (crd *ThirdUserAuthCmd) iniFact() error {
	//获取最新配置
	crd.configFactory = third_user_auth_service.CreateConfigFactory()
	/*crd.authConfig := &third_user_auth_models.AuthModel{
		Addr:      "https://10.32.2.4",
		AppId:     "789546861770502144",
		AppKey:    "56de926f281e4120b932476492b4d451",
		AppSecret: "48642278b1bf4fe28b74fdd65ab799e5",
	}*/
	crd.authConfig = crd.configFactory.GetNewestItem()
	if crd.authConfig == nil || string_helper.IsEmpty(crd.authConfig.Addr) {
		return errors.New("auth config fail")
	}

	//请求第三方
	crd.requestFactory = third_user_auth_service.CreateRequestFactory(crd.authConfig)

	return nil
}

// Sync 同步数据
func (crd *ThirdUserAuthCmd) Sync(isNotice bool) {
	// 加锁失败：说明上一次还没执行完，直接返回
	if !crd.setLock() {
		variable.ZapLog.Info("Sync skipped: previous task not completed", zap.Bool("isNotice", isNotice))
		return
	}
	// 只有加锁成功后才 defer 解锁
	defer crd.delLock()
	variable.ZapLog.Info("Sync start", zap.Bool("isNotice", isNotice))

	//
	if err := crd.iniFact(); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}

	//初始化未完成
	if crd.authConfig.IsInit == enums.IsNot {
		crd.handleOrg()
		crd.handleUser()
	} else {
		//通知任务
		if isNotice {
			//有通知
			if crd.authConfig.NoticeSync > 0 {
				crd.handleOrg()
				crd.handleUser()
			} else {
				variable.ZapLog.Info("Not Sync", zap.Int("NoticeSync", crd.authConfig.NoticeSync))
				return
			}
		} else {
			crd.handleOrg()
			crd.handleUser()
		}
	}

	//

	variable.ZapLog.Info("Sync end", zap.Bool("isNotice", isNotice))
}

// 三方组织结构中的单位对应我们企业;
// 三方组织结构中的我们对应我们组织架构;
func (crd *ThirdUserAuthCmd) handleOrg() {
	variable.ZapLog.Info("handleOrg ===================================")

	//获取三方数据
	err, thirdList := crd.requestFactory.PullOrg()
	if err != nil {
		variable.ZapLog.Info(fmt.Sprintf("PullOrg Fail: %v", err))
	}
	if len(thirdList) == 0 {
		variable.ZapLog.Info("PullOrg Empty")
	}

	//三方数据
	thirdMaps := make(map[string]*third_user_auth_service.ResponseOrg)
	for _, item := range thirdList {
		thirdMaps[item.OrgId] = item
	}

}

// 处理用户
func (crd *ThirdUserAuthCmd) handleUser() {
	return
	variable.ZapLog.Info("handleUser ===================================")

	//获取三方数据
	err, thirdList := crd.requestFactory.PullUser()
	if err != nil {
		variable.ZapLog.Info(fmt.Sprintf("handleUser Fail: %v", err))
	}
	if len(thirdList) == 0 {
		variable.ZapLog.Info("handleUser Empty")
	}
	fmt.Println("userList--------------:")
	fmt.Println(json_helper.Encode(thirdList))
	fmt.Println("userList--------------:")
}
