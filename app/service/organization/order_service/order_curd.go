package order_service

import (
	"errors"
	"fmt"
	"goskeleton/app/model/model_organization"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/service/platform/admin_service"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/redis_factory"
)

func CreateCurdFactory() *OrderCurd {
	return &OrderCurd{
		masterModel: model_organization.CreateOrderFactory(),
		redisClient: redis_factory.GetOneRedisClient(),
		redisPrefix: "order_number_",
		redisExpire: 86400,
	}
}

type OrderCurd struct {
	OrderBase
	masterModel *model_organization.OrderModel
	redisClient *redis_factory.RedisClient
	redisPrefix string
	redisExpire int
}

func (crd *OrderCurd) makerOrderNo() string {
	today := helpers.DatetimeSpecial("20060102")
	key := crd.redisPrefix + today
	if val, err := crd.redisClient.Incr(key); err == nil {
		_, _ = crd.redisClient.Expire(key, crd.redisExpire)
		num := helpers.IntToStr(val)
		return "Yhj" + today[2:] + helpers.StringPad(num, 5, "0000000", "LEFT") + helpers.IntRandom2Str(10, 99)
	} else {
		return ""
	}
}

// PlatformAdd 平台新增
func (crd *OrderCurd) PlatformAdd(orderParam *model_organization.OrderModel, loginInfo admin_service.AdminCacheLogin) (errCode int, err error) {
	fmt.Println("=====================", crd.makerOrderNo())
	//订单信息
	orgModel := organize_service.CreateOrgCurdFactory().GetItemByOrg(orderParam.OrgId)
	if orgModel != nil || orgModel.Id == 0 {
		return -1, errors.New("企业版不存在")
	}
	//订单信息
	orgModel := organize_service.CreateOrgCurdFactory().GetItemByOrg(orderParam.OrgId)
	if orgModel != nil || orgModel.Id == 0 {
		return -1, errors.New("企业版不存在")
	}

	return errCode, err
}
