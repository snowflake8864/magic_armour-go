package order_service

import (
	"errors"
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

// Add 平台新增
func (crd *OrderCurd) Add(orderParam *model_organization.OrderModel, operateId int64, operateName string) (errCode int, err error) {
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

	//订单日志
	logParam := model_organization.OrderLogModel{
		OrgId:         orderParam.OrgId,
		OrderNo:       orderParam.OrderNo,
		OrderStatus:   orderParam.OrderStatus,
		OperateAction: enum_order.ActionAdd,
		OperateSource: orderParam.Source,
		OperateId:     operateId,
		OperateName:   operateName,
		Remark:        orderParam.DeclineRemark,
		CreatedAt:     orderParam.CreatedAt,
	}
	if err = tx.Create(&logParam).Error; err != nil {
		tx.Rollback()
		return -1, errors.New("保存日志失败")
	}
	tx.Commit()
	return errCode, err
}

// GetItemByOrderNo 通过编号称获取
func (crd *OrderCurd) GetItemByOrderNo(id string) (model *model_organization.OrderModel) {
	model, _ = crd.masterModel.GetItemById(id)
	return model
}

// GetItemByOrderNos 通过编号集合称获取
func (crd *OrderCurd) GetItemByOrderNos(ids []string) (items []*model_organization.OrderModel) {
	items, _ = crd.masterModel.GetItemsByIds(ids)
	return items
}

// GetValidItemsByOrg 通过企业编号获取有效订单
func (crd *OrderCurd) GetValidItemsByOrg(orgId int64) (items []*model_organization.OrderModel) {
	search := SearchConditions{
		OrgId:       orgId,
		IsExpire:    enums.IsNot,
		OrderStatus: enum_order.OrderStatusCompleted,
	}

	items, err := crd.masterModel.GetItemsByConditions(crd.search2Where(search), 0, 0, "")
	if err != nil || len(items) == 0 {
		return nil
	}
	return items
}

// Page 分页
func (crd *OrderCurd) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []*model_organization.OrderModel) {
	where := crd.search2Where(search)
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(where, offset, pageSize, order)
	return total, items
}

// PlatformCancel 平台取消
func (crd *OrderCurd) PlatformCancel(orderNo string, operateId int64, operateName string) (errCode int, err error) {
	orderModel := crd.GetItemByOrderNo(orderNo)
	if orderModel == nil || helpers.StringIsEmpty(orderModel.OrderNo) {
		return -1, errors.New("订单不存在")
	}
	if orderModel.OrderStatus != enum_order.OrderStatusUnpaid {
		return -1, errors.New("订单状态不正确")
	}

	now := helpers.DatetimeCurr()
	//修改
	updates := map[string]interface{}{
		"order_status": enum_order.OrderStatusCancel,
		"cancel_at":    now,
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return -1, errors.New("取消订单失败")
	}

	//订单日志
	logParam := model_organization.OrderLogModel{
		OrgId:         orderModel.OrgId,
		OrderNo:       orderModel.OrderNo,
		OrderStatus:   orderModel.OrderStatus,
		OperateAction: enum_order.ActionCancel,
		OperateSource: enums.OperateSourcePlatform,
		OperateId:     operateId,
		OperateName:   operateName,
		CreatedAt:     now,
	}
	CreateLogFactory().Add(&logParam)

	return errCode, err
}

// OrgCancel 企业取消
func (crd *OrderCurd) OrgCancel(orderNo string, orgId, operateId int64, operateName string) (errCode int, err error) {
	orderModel := crd.GetItemByOrderNo(orderNo)
	if orderModel == nil || helpers.StringIsEmpty(orderModel.OrderNo) {
		return -1, errors.New("订单不存在")
	}
	if orderModel.OrgId != orgId {
		return -1, errors.New("订单错误")
	}
	if orderModel.OrderStatus != enum_order.OrderStatusUnpaid {
		return -1, errors.New("订单状态不正确")
	}

	now := helpers.DatetimeCurr()
	//修改
	updates := map[string]interface{}{
		"order_status": enum_order.OrderStatusCancel,
		"cancel_at":    now,
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return -1, errors.New("取消订单失败")
	}

	//订单日志
	logParam := model_organization.OrderLogModel{
		OrgId:         orderModel.OrgId,
		OrderNo:       orderModel.OrderNo,
		OrderStatus:   orderModel.OrderStatus,
		OperateAction: enum_order.ActionCancel,
		OperateSource: enums.OperateSourceOrganize,
		OperateId:     operateId,
		OperateName:   operateName,
		CreatedAt:     now,
	}
	CreateLogFactory().Add(&logParam)

	return errCode, err
}

// PayByBankcard 通过银行卡支付
func (crd *OrderCurd) PayByBankcard(orderNo string, orgId int64, operateId int64, operateName string) (errCode int, err error) {

	return errCode, err
}

// Decline 直降
func (crd *OrderCurd) Decline(orderNo, remark string, realAmount float64, operateId int64, operateName string) (errCode int, err error) {

	orderModel := crd.GetItemByOrderNo(orderNo)
	if orderModel == nil || helpers.StringIsEmpty(orderModel.OrderNo) {
		return -1, errors.New("订单不存在")
	}
	if orderModel.OrderStatus != enum_order.OrderStatusUnpaid {
		return -1, errors.New("订单状态不正确")
	}
	if orderModel.TotalAmount < realAmount {
		return -1, errors.New("支付金额不能大于原始金额")
	}

	now := helpers.DatetimeCurr()
	//修改
	updates := map[string]interface{}{
		"real_amount":    crd.calcDeclineAmount(orderModel.TotalAmount, realAmount),
		"decline_amount": crd.calcDeclineAmount(orderModel.TotalAmount, realAmount),
		"decline_remark": remark,
		"decline_at":     now,
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return -1, errors.New("订单修改失败")
	}

	//订单日志
	logParam := model_organization.OrderLogModel{
		OrgId:         orderModel.OrgId,
		OrderNo:       orderModel.OrderNo,
		OrderStatus:   orderModel.OrderStatus,
		OperateAction: enum_order.ActionDecline,
		OperateSource: enums.OperateSourcePlatform,
		OperateId:     operateId,
		OperateName:   operateName,
		Remark:        remark,
		CreatedAt:     now,
	}
	CreateLogFactory().Add(&logParam)

	return errCode, err
}

// ConfirmPaid 确认已支付
func (crd *OrderCurd) ConfirmPaid(orderNo string, operateId int64, operateName string) (errCode int, err error) {
	orderModel := crd.GetItemByOrderNo(orderNo)
	if orderModel == nil || helpers.StringIsEmpty(orderModel.OrderNo) {
		return -1, errors.New("订单不存在")
	}
	if orderModel.OrderStatus != enum_order.OrderStatusPaid {
		return -1, errors.New("订单状态不正确")
	}

	now := helpers.DatetimeCurr()
	//修改
	updates := map[string]interface{}{
		"order_status": enum_order.OrderStatusCompleted,
		"complete_at":  now,
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return -1, errors.New("订单修改失败")
	}

	//订单日志
	logParam := model_organization.OrderLogModel{
		OrgId:         orderModel.OrgId,
		OrderNo:       orderModel.OrderNo,
		OrderStatus:   orderModel.OrderStatus,
		OperateAction: enum_order.ActionConfirmPaid,
		OperateSource: enums.OperateSourcePlatform,
		OperateId:     operateId,
		OperateName:   operateName,
		CreatedAt:     now,
	}
	CreateLogFactory().Add(&logParam)

	return errCode, err
}

// ConfirmUnpaid 确认未付款
func (crd *OrderCurd) ConfirmUnpaid(orderNo string, operateId int64, operateName string) (errCode int, err error) {
	orderModel := crd.GetItemByOrderNo(orderNo)
	if orderModel == nil || helpers.StringIsEmpty(orderModel.OrderNo) {
		return -1, errors.New("订单不存在")
	}
	if orderModel.OrderStatus != enum_order.OrderStatusPaid {
		return -1, errors.New("订单状态不正确")
	}

	now := helpers.DatetimeCurr()
	//修改
	updates := map[string]interface{}{
		"order_status": enum_order.OrderStatusUnpaid,
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return -1, errors.New("订单修改失败")
	}

	//订单日志
	logParam := model_organization.OrderLogModel{
		OrgId:         orderModel.OrgId,
		OrderNo:       orderModel.OrderNo,
		OrderStatus:   orderModel.OrderStatus,
		OperateAction: enum_order.ActionConfirmUnpaid,
		OperateSource: enums.OperateSourcePlatform,
		OperateId:     operateId,
		OperateName:   operateName,
		CreatedAt:     now,
	}
	CreateLogFactory().Add(&logParam)

	return errCode, err
}

// OrgDelete 企业删除
func (crd *OrderCurd) OrgDelete(orderNo string, orgId int64) (errCode int, err error) {
	orderModel := crd.GetItemByOrderNo(orderNo)
	if orderModel == nil || helpers.StringIsEmpty(orderModel.OrderNo) {
		return -1, errors.New("订单不存在")
	}
	if orderModel.OrgId != orgId {
		return -1, errors.New("订单错误")
	}
	if orderModel.OrderStatus != enum_order.OrderStatusCancel {
		return -1, errors.New("订单状态不正确")
	}

	if !crd.masterModel.DeleteById(orderNo) {
		return -1, errors.New("订单删除失败")
	}

	//订单日志
	CreateLogFactory().DeleteByOrder(orderNo)

	return errCode, err
}

// AutoExpire 自动过期
func (crd *OrderCurd) AutoExpire(orgId int64, expireDate string) bool {
	search := SearchConditions{
		OrgId:       orgId,
		OrderStatus: enum_order.OrderStatusCompleted,
		IsExpire:    enums.IsNot,
		ExpireAtGte: expireDate,
	}
	items, err := crd.masterModel.GetItemsByConditions(crd.search2Where(search), 0, 0, "")
	if err != nil || len(items) == 0 {
		return true
	}
	now := helpers.DatetimeCurr()
	for _, model := range items {
		//修改订单
		updates := map[string]interface{}{
			"is_expire": enums.IsYes,
		}
		crd.masterModel.UpdateById(model.OrderNo, updates)

		//订单日志
		logParam := model_organization.OrderLogModel{
			OrgId:         model.OrgId,
			OrderNo:       model.OrderNo,
			OrderStatus:   model.OrderStatus,
			OperateAction: enum_order.ActionAutoExpire,
			OperateSource: enums.OperateSourceAi,
			CreatedAt:     now,
		}
		CreateLogFactory().Add(&logParam)
	}

	return true
}
