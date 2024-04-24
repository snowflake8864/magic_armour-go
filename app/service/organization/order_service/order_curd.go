package order_service

import (
	"errors"
	"fmt"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/model/model_organization"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/product_service"
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
		return "Yhj" + today[2:] + helpers.StringPad(num, 4, "0000000", "LEFT") + helpers.IntRandom2Str(1, 9)
	} else {
		return ""
	}
}

func (crd *OrderCurd) calcDeclineAmount(TotalAmount, RealAmount float64) float64 {
	return helpers.FloatSub(TotalAmount, RealAmount, 2)
}

func (crd *OrderCurd) getExpireAt() string {
	days := 365
	if val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.OrderExpireDays); !helpers.StringIsEmpty(val) {
		days = helpers.String2Int(val)
	}

	day, _, _, _ := helpers.DateCalcByDay(days)
	return day
}

// PlatformAdd 平台新增
func (crd *OrderCurd) PlatformAdd(orderParam *model_organization.OrderModel, operateId int64, operateName string) (errCode int, err error) {
	//订单信息
	orgModel := organize_service.CreateOrgCurdFactory().GetItemByOrg(orderParam.OrgId)
	if orgModel == nil || orgModel.Id == 0 {
		return -1, errors.New("企业不存在")
	}
	//产品信息
	productModel := product_service.CreateCurdFactory().GetItemById(orderParam.ProductId)
	if productModel == nil || productModel.Id == 0 {
		return -1, errors.New("产品不存在")
	}
	if productModel.GrantType == enum_product.TypeTest && orderParam.GrantDuration > 1 {
		return -1, errors.New("购买测试授权类型的产品不能超过一个月")
	}

	tx := crd.masterModel.Begin()
	if tx.Error != nil {
		return -10, errors.New("保存订单失败")
	}
	fmt.Println("=====================", crd.makerOrderNo())

	//新增订单
	orderParam.OrderNo = crd.makerOrderNo()
	//企业相关
	orderParam.OrgName = orgModel.Name
	orderParam.UserName = orgModel.UserName
	orderParam.UserEmail = orgModel.UserEmail
	orderParam.UserPhone = orgModel.UserPhone
	//产品相关
	orderParam.ProductTitle = productModel.Title
	orderParam.GrantType = productModel.GrantType
	orderParam.GrantModule = productModel.GrantModule
	orderParam.SalePrice = productModel.Price
	orderParam.TotalAmount = helpers.FloatMul(orderParam.SalePrice, orderParam.BuyNumber, 2, orderParam.GrantDuration)
	orderParam.CreatedAt = helpers.DatetimeCurr()
	if orderParam.Source == enums.OperateSourcePlatform {
		orderParam.DeclineAmount = crd.calcDeclineAmount(orderParam.TotalAmount, orderParam.RealAmount)
		orderParam.OrderStatus = enum_order.OrderStatusCompleted
		orderParam.DeclineAt = orderParam.CreatedAt
		orderParam.PayAt = orderParam.CreatedAt
		orderParam.CompleteAt = orderParam.CreatedAt
	} else {
		orderParam.DeclineAmount = 0.00
		orderParam.OrderStatus = enum_order.OrderStatusPaid
	}
	orderParam.IsExpire = enums.IsNot
	orderParam.ExpireAt = crd.getExpireAt()
	if err = tx.Create(orderParam).Error; err != nil {
		tx.Rollback()
		return -1, errors.New("保存订单失败")
	}

	//订单授权
	if orderParam.Source == enums.OperateSourcePlatform {
		orderGrant := CreateGrantFactory().Init(orderParam.OrderNo, orderParam.OrgId, orderParam.BuyNumber)
		if err = tx.Create(&orderGrant).Error; err != nil {
			tx.Rollback()
			return -1, errors.New("保存订单失败")
		}
	}
	tx.Commit()

	//订单日志
	logParam := model_organization.OrderLogModel{
		OrgId:         orderParam.OrgId,
		OrderNo:       orderParam.OrderNo,
		OrderStatus:   orderParam.OrderStatus,
		OperateAction: enum_order.ActionAdd,
		OperateSource: orderParam.Source,
		OperateId:     operateId,
		OperateName:   operateName,
		UsedNumber:    0,
		UsedAgents:    "",
		Remark:        orderParam.DeclineRemark,
		CreatedAt:     orderParam.CreatedAt,
	}
	CreateLogFactory().Add(&logParam)

	return errCode, err
}
