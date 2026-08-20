package order_service

import (
	"errors"
	"fmt"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/enums/enum_product"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/order_models"
	"goskeleton/app/service/agent/agent_service/agent_callback_service"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/product_service"
	"goskeleton/app/utils/cur_app"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/float_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/redis_factory"
)

func CreateCurdFactory() *OrderCurd {
	return &OrderCurd{
		masterModel:       order_models.CreateOrderFactory(),
		redisClient:       redis_factory.GetOneRedisClient(),
		redisNumberPrefix: "order_number_",
		redisNumberExpire: 86400,
	}
}

type OrderCurd struct {
	masterModel       *order_models.OrderModel
	redisClient       *redis_factory.RedisClient
	redisNumberPrefix string
	redisNumberExpire int

	lang string
}

func (crd *OrderCurd) makerOrderNo(grantType int) string {
	prefix := enum_order.OrderNoPrefixCloud
	if cur_app.GetAppUsageVersion() == variable.AppUsageVersionOffline {
		prefix = enum_order.OrderNoPrefixOffline
	}

	today := date_helper.SpecialDatetime("20060102")
	key := crd.redisNumberPrefix + today
	if val, err := crd.redisClient.Incr(key); err == nil {
		_, _ = crd.redisClient.Expire(key, crd.redisNumberExpire)
		num := int_helper.IntToStr(val)
		return prefix + int_helper.IntToStr(grantType) + today[2:] + string_helper.Pad(num, 4, "0000000", "LEFT") + int_helper.IntRandom2Str(1, 9)
	} else {
		return ""
	}
}

func (crd *OrderCurd) calcDeclineAmount(TotalAmount, RealAmount float64) float64 {
	return float_helper.Sub(TotalAmount, RealAmount, 2)
}

func (crd *OrderCurd) getExpireAt(grantDuration int) string {
	//bug[10110] 改成已订单过期为准
	return date_helper.TimestampToDatetimeByFormat(date_helper.AddTimestampByDays(0, grantDuration, 0), variable.DateFormatToday)

	/*days := 365
	if val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.OrderExpireDays); !string_helper.IsEmpty(val) {
		days = string_helper.ToInt(val)
	}

	day, _, _, _ := date_helper.CalcByDay(days)
	return day*/
}

func (crd *OrderCurd) getLicenseDefaultName() string {
	str := "import license"
	if val := dictionary_service.CreateCacheFactory().GetValById(enum_dictionary.OrderLicenseName); !string_helper.IsEmpty(val) {
		str = val
	}
	return str
}

// Add 新增订单
func (crd *OrderCurd) Add(orderParam *order_models.OrderModel, operateId int64, operateName string) (errCode int, err error) {
	//企业信息
	orgModel := organize_service.CreateOrgCurdFactory().GetItemByOrg(orderParam.OrgId)
	if orgModel == nil || orgModel.Id == 0 {
		return -1, errors.New(my_errors.CurdOrgNotExists)
	}
	//产品信息
	productModel := product_service.CreateCurdFactory().GetItemById(orderParam.ProductId)
	if productModel == nil || productModel.Id == 0 {
		return -1, errors.New(my_errors.CurdProductNotExist)
	}
	if productModel.GrantType == enum_product.TypeTest {
		if orderParam.GrantDuration > 1 {
			return -1, errors.New(my_errors.CurdProductTestTypeNotGtMonth)
		}

		search := order_models.SearchConditions{
			OrgId:     orderParam.OrgId,
			GrantType: enum_product.TypeTest,
		}
		if total := crd.masterModel.Count(search); total > 0 {
			return -1, errors.New(my_errors.CurdProductTestTypeCannotRepeated)
		}
	}

	// 有未完成的订单
	search := order_models.SearchConditions{
		OrgId:      orderParam.OrgId,
		StatusList: []int{enum_order.OrderStatusUnpaid, enum_order.OrderStatusPaid},
	}
	if total := crd.masterModel.Count(search); total > 0 {
		return -1, errors.New(my_errors.CurdOrderHasUnfinished)
	}

	tx := crd.masterModel.Begin()
	if tx.Error != nil {
		return -10, errors.New(my_errors.CurdAddFailed)
	}

	//新增订单
	orderParam.GrantType = productModel.GrantType
	orderParam.OrderNo = crd.makerOrderNo(orderParam.GrantType)
	//企业相关
	orderParam.OrgName = orgModel.Name
	orderParam.UserName = orgModel.UserName
	orderParam.UserEmail = orgModel.UserEmail
	orderParam.UserPhone = orgModel.UserPhone
	//产品相关
	orderParam.ProductTitle = productModel.Title
	orderParam.GrantModule = productModel.GrantModule
	orderParam.SalePrice = productModel.Price
	orderParam.TotalAmount = float_helper.Mul(orderParam.SalePrice, orderParam.BuyNumber, 2, orderParam.GrantDuration)
	orderParam.CreatedAt = date_helper.CurDatetime()
	if orderParam.Source == enums.OperateSourcePlatform {
		orderParam.PayMethod = enum_order.PayMethodDefault
		orderParam.DeclineAmount = crd.calcDeclineAmount(orderParam.TotalAmount, orderParam.RealAmount)
		orderParam.OrderStatus = enum_order.OrderStatusCompleted
		orderParam.DeclineAt = orderParam.CreatedAt
		orderParam.PayAt = orderParam.CreatedAt
		orderParam.CompleteAt = orderParam.CreatedAt
	} else {
		orderParam.DeclineAmount = 0.00
		orderParam.OrderStatus = enum_order.OrderStatusUnpaid
	}
	orderParam.IsExpire = enums.IsNot
	orderParam.ExpireAt = crd.getExpireAt(orderParam.GrantDuration)
	if err = tx.Create(orderParam).Error; err != nil {
		tx.Rollback()
		return -1, errors.New(my_errors.CurdAddFailed)
	}

	//订单日志
	logParam := order_models.OrderLogModel{
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
		return -1, errors.New(my_errors.CurdSaveFailed)
	}
	tx.Commit()
	return errCode, err
}

// ImportLicense 导入证书
func (crd *OrderCurd) ImportLicense(orderParam *order_models.OrderModel, licensePath string, operateId int64, operateName string) (errCode int, err error) {
	//解析证书
	err1, licenseResp := CreateLicenseFactory().Parse(licensePath)
	if err1 != nil {
		return -2, err1
	}

	/*licenseResp := LicenseResp{
		Serial:    "8d3080e2eff8fa3026a7cdef0c30d8d3",
		Type:      enum_product.TypeGold,
		TypeLabel: "",
		StartAt:   "2023-11-01",
		ExpireAt:  "2026-11-01",
		Duration:  36,
		Number:    12,
		Modules: enum_product.ProductGrantConfig{
			Process:   100,
			Tamper:    100,
			Blackmail: 100,
			Port:      100,
			Plug:      100,
			Usb:       100,
			Syslog:    100,
			Outreach:  100,
			Login:     100,
			Threaten:  100,
		},
	}*/

	if licenseResp.Type == enum_product.TypeTest {
		IsExpire := enums.IsNot
		search := order_models.SearchConditions{
			OrgId:       orderParam.OrgId,
			GrantType:   enum_product.TypeTest,
			OrderStatus: enum_order.OrderStatusCompleted,
			IsExpire:    &IsExpire,
		}
		if total := crd.masterModel.Count(search); total > 0 {
			return -1, errors.New(my_errors.CurdGrantLicTestTypeCannotRepeated)
		}
	}

	//企业信息
	orgModel := organize_service.CreateOrgCurdFactory().GetItemByOrg(orderParam.OrgId)
	if orgModel == nil || orgModel.Id == 0 {
		return -1, errors.New(my_errors.CurdOrgNotExists)
	}
	if licenseResp.Serial != orgModel.Code {
		return -1, errors.New(my_errors.CurdGrantLicUnbelongOrg)
	}

	// 已经导入过
	search := order_models.SearchConditions{
		OrgId:      orderParam.OrgId,
		PayOrderNo: orderParam.PayOrderNo,
	}
	if total := crd.masterModel.Count(search); total > 0 {
		return -1, errors.New(my_errors.CurdGrantLicRepeatedImport)
	}

	//补充订单信息
	orderParam.OrgName = orgModel.Name
	orderParam.UserName = orgModel.UserName
	orderParam.UserEmail = orgModel.UserEmail
	orderParam.UserPhone = orgModel.UserPhone
	orderParam.ProductTitle = crd.getLicenseDefaultName()
	orderParam.GrantType = licenseResp.Type
	//orderParam.GrantModule = productModel.GrantModule
	orderParam.GrantDuration = licenseResp.Duration
	//orderParam.BuyNumber = 0
	orderParam.SalePrice = 0.00
	orderParam.DeclineAmount = 0.00
	orderParam.TotalAmount = 0.00
	orderParam.RealAmount = 0.00
	orderParam.PayMethod = enum_order.PayMethodLicense
	orderParam.OrderStatus = enum_order.OrderStatusCompleted
	orderParam.CreatedAt = date_helper.CurDatetime()
	orderParam.DeclineAt = orderParam.CreatedAt
	orderParam.PayAt = orderParam.CreatedAt
	orderParam.CompleteAt = orderParam.CreatedAt
	orderParam.IsExpire = enums.IsNot
	orderParam.ExpireAt = licenseResp.ExpireAt

	//要插入的订单
	var insertOrders []order_models.OrderModel
	var insertLogs []order_models.OrderLogModel
	//因为每个授权模块可能不一致, 生成多个订单
	for true {
		//最小值
		minimum := licenseResp.Modules.MinGreaterThanZero()
		if minimum == 0 {
			break
		}
		orderParam.OrderNo = crd.makerOrderNo(orderParam.GrantType)
		//最小值为购买数目
		orderParam.BuyNumber = minimum
		//授权配置
		orderParam.GrantModule = json_helper.Encode(licenseResp.Modules.Num2Config())
		//每次减去最小值
		licenseResp.Modules.Sub(minimum)
		insertOrders = append(insertOrders, *orderParam)

		insertLogs = append(insertLogs, order_models.OrderLogModel{
			OrgId:         orderParam.OrgId,
			OrderNo:       orderParam.OrderNo,
			OrderStatus:   orderParam.OrderStatus,
			OperateAction: enum_order.ActionImportLicense,
			OperateSource: orderParam.Source,
			OperateId:     operateId,
			OperateName:   operateName,
			Remark:        orderParam.DeclineRemark,
			CreatedAt:     orderParam.CreatedAt,
		})
	}
	if len(insertOrders) == 0 {
		return -10, errors.New(my_errors.CurdGrantModuleNumberFail)
	}

	tx := crd.masterModel.Begin()
	if tx.Error != nil {
		return -10, errors.New(my_errors.CurdAddFailed)
	}
	if err = tx.Create(&insertOrders).Error; err != nil {
		tx.Rollback()
		return -1, errors.New(my_errors.CurdAddFailed)
	}
	if err = tx.Create(&insertLogs).Error; err != nil {
		tx.Rollback()
		return -1, errors.New(my_errors.CurdSaveFailed)
	}
	tx.Commit()

	return errCode, err
}

// GetItemByOrderNo 通过编号称获取
func (crd *OrderCurd) GetItemByOrderNo(id string) (model *order_models.OrderModel) {
	model = crd.masterModel.GetItemById(id)
	return model
}

// GetItemByOrderNos 通过编号集合称获取
func (crd *OrderCurd) GetItemByOrderNos(ids []string) (items []*order_models.OrderModel) {
	items = crd.masterModel.GetItemsByIds(ids)
	return items
}

// GetValidItemsByOrg 通过企业编号获取有效订单
func (crd *OrderCurd) GetValidItemsByOrg(orgId int64) (items []*order_models.OrderModel) {
	IsExpire := enums.IsNot
	search := order_models.SearchConditions{
		OrgId:       orgId,
		IsExpire:    &IsExpire,
		OrderStatus: enum_order.OrderStatusCompleted,
	}

	items = crd.masterModel.GetItemsByConditions(search, 0, 0, "")
	if len(items) == 0 {
		return nil
	}
	return items
}

// GetValidItemsByOrgs 通过企业编号获取有效订单
func (crd *OrderCurd) GetValidItemsByOrgs(orgIds []int64) (items []*order_models.OrderModel) {
	IsExpire := enums.IsNot
	search := order_models.SearchConditions{
		OrgIds:      orgIds,
		IsExpire:    &IsExpire,
		OrderStatus: enum_order.OrderStatusCompleted,
	}

	items = crd.masterModel.GetItemsByConditions(search, 0, 0, "")
	if len(items) == 0 {
		return nil
	}
	return items
}

// PlatformCancel 平台取消
func (crd *OrderCurd) PlatformCancel(orderNo string, operateId int64, operateName string) (errCode int, err error) {
	orderModel := crd.GetItemByOrderNo(orderNo)
	if orderModel == nil || string_helper.IsEmpty(orderModel.OrderNo) {
		return -1, errors.New(my_errors.CurdOrderNotExist)
	}
	if orderModel.OrderStatus != enum_order.OrderStatusUnpaid {
		return -1, errors.New(my_errors.CurdOrderStatusFail)
	}

	now := date_helper.CurDatetime()
	//修改
	updates := map[string]interface{}{
		"order_status": enum_order.OrderStatusCancel,
		"cancel_at":    now,
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return -1, errors.New(my_errors.CurdSaveFailed)
	}

	//订单日志
	logParam := order_models.OrderLogModel{
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
	if orderModel == nil || string_helper.IsEmpty(orderModel.OrderNo) {
		return -1, errors.New(my_errors.CurdOrderNotExist)
	}
	if orderModel.OrgId != orgId {
		return -1, errors.New(my_errors.CurdOrderOrgNotMatch)
	}
	if orderModel.OrderStatus != enum_order.OrderStatusUnpaid {
		return -1, errors.New(my_errors.CurdOrderStatusFail)
	}

	now := date_helper.CurDatetime()
	//修改
	updates := map[string]interface{}{
		"order_status": enum_order.OrderStatusCancel,
		"cancel_at":    now,
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return -1, errors.New(my_errors.CurdSaveFailed)
	}

	//订单日志
	logParam := order_models.OrderLogModel{
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
	orderModel := crd.GetItemByOrderNo(orderNo)
	if orderModel == nil || string_helper.IsEmpty(orderModel.OrderNo) {
		return -1, errors.New(my_errors.CurdOrderNotExist)
	}
	if orderModel.OrgId != orgId {
		return -1, errors.New(my_errors.CurdOrderOrgNotMatch)
	}
	if orderModel.OrderStatus != enum_order.OrderStatusUnpaid {
		return -1, errors.New(my_errors.CurdOrderStatusFail)
	}

	payInfo := enum_order.FormatPayBankcard(dictionary_service.CreateCurdFactory().GetMapByModule(enum_dictionary.ModuleOrder))

	now := date_helper.CurDatetime()
	//修改
	updates := map[string]interface{}{
		"order_status": enum_order.OrderStatusPaid,
		"pay_method":   enum_order.PayMethodBankcard,
		"pay_info":     json_helper.Encode(payInfo),
		"pay_at":       now,
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return -1, errors.New(my_errors.CurdSaveFailed)
	}

	//订单日志
	logParam := order_models.OrderLogModel{
		OrgId:         orderModel.OrgId,
		OrderNo:       orderModel.OrderNo,
		OrderStatus:   orderModel.OrderStatus,
		OperateAction: enum_order.ActionPayBankcard,
		OperateSource: enums.OperateSourceOrganize,
		OperateId:     operateId,
		OperateName:   operateName,
		CreatedAt:     now,
	}
	CreateLogFactory().Add(&logParam)

	return errCode, err
}

// Decline 直降
func (crd *OrderCurd) Decline(orderNo, remark string, realAmount float64, operateId int64, operateName string) (errCode int, err error) {

	orderModel := crd.GetItemByOrderNo(orderNo)
	if orderModel == nil || string_helper.IsEmpty(orderModel.OrderNo) {
		return -1, errors.New(my_errors.CurdOrderNotExist)
	}
	if orderModel.OrderStatus != enum_order.OrderStatusUnpaid {
		return -1, errors.New(my_errors.CurdOrderStatusFail)
	}
	if orderModel.TotalAmount < realAmount {
		return -1, errors.New(my_errors.CurdOrderPayLetOriginal)
	}

	now := date_helper.CurDatetime()
	//修改
	updates := map[string]interface{}{
		"real_amount":    crd.calcDeclineAmount(orderModel.TotalAmount, realAmount),
		"decline_amount": crd.calcDeclineAmount(orderModel.TotalAmount, realAmount),
		"decline_remark": remark,
		"decline_at":     now,
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return -1, errors.New(my_errors.CurdModifyFailed)
	}

	//订单日志
	logParam := order_models.OrderLogModel{
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
	if orderModel == nil || string_helper.IsEmpty(orderModel.OrderNo) {
		return -1, errors.New(my_errors.CurdOrderNotExist)
	}
	if orderModel.OrderStatus != enum_order.OrderStatusPaid {
		return -1, errors.New(my_errors.CurdOrderStatusFail)
	}

	now := date_helper.CurDatetime()
	//修改
	updates := map[string]interface{}{
		"order_status": enum_order.OrderStatusCompleted,
		"complete_at":  now,
		"expire_at":    crd.getExpireAt(orderModel.GrantDuration),
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return -1, errors.New(my_errors.CurdModifyFailed)
	}

	//订单日志
	logParam := order_models.OrderLogModel{
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
	if orderModel == nil || string_helper.IsEmpty(orderModel.OrderNo) {
		return -1, errors.New(my_errors.CurdOrderNotExist)
	}
	if orderModel.OrderStatus != enum_order.OrderStatusPaid {
		return -1, errors.New(my_errors.CurdOrderStatusFail)
	}

	now := date_helper.CurDatetime()
	//修改
	updates := map[string]interface{}{
		"order_status": enum_order.OrderStatusUnpaid,
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return -1, errors.New(my_errors.CurdModifyFailed)
	}

	//订单日志
	logParam := order_models.OrderLogModel{
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
	if orderModel == nil || string_helper.IsEmpty(orderModel.OrderNo) {
		return -1, errors.New(my_errors.CurdOrderNotExist)
	}
	if orderModel.OrgId != orgId {
		return -1, errors.New(my_errors.CurdOrderOrgNotMatch)
	}
	if orderModel.OrderStatus != enum_order.OrderStatusCancel {
		return -1, errors.New(my_errors.CurdOrderStatusFail)
	}

	if !crd.masterModel.DeleteById(orderNo) {
		return -1, errors.New(my_errors.CurdDeleteFailed)
	}

	//订单日志
	CreateLogFactory().DeleteByOrder(orderNo)

	return errCode, err
}

// AutoExpire 自动过期
func (crd *OrderCurd) AutoExpire(orgId int64, expireDate string) bool {
	IsExpire := enums.IsNot
	search := order_models.SearchConditions{
		OrgId:       orgId,
		OrderStatus: enum_order.OrderStatusCompleted,
		IsExpire:    &IsExpire,
		ExpireAtLte: expireDate,
	}
	items := crd.masterModel.GetItemsByConditions(search, 0, 0, "")
	if len(items) == 0 {
		return true
	}
	now := date_helper.CurDatetime()
	for _, model := range items {
		//修改订单
		updates := map[string]interface{}{
			"is_expire": enums.IsYes,
		}
		crd.masterModel.UpdateById(model.OrderNo, updates)

		//订单日志
		logParam := order_models.OrderLogModel{
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

// GrantToAgents 授权给终端集合
func (crd *OrderCurd) GrantToAgents(orgId int64, orderNo string, agentIds []string, operateId int64, operateName string) error {
	orderModel := crd.GetItemByOrderNo(orderNo)
	if orderModel == nil || string_helper.IsEmpty(orderModel.OrderNo) {
		return errors.New(my_errors.CurdOrderNotExist)
	}
	if orderModel.OrgId != orgId {
		return errors.New(my_errors.CurdOrderOrgNotMatch)
	}
	if orderModel.OrderStatus != enum_order.OrderStatusCompleted {
		return errors.New(my_errors.CurdOrderStatusFail)
	}
	if orderModel.IsExpire == enums.IsYes {
		return errors.New(my_errors.CurdOrderExpired)
	}
	if orderModel.BuyNumber <= orderModel.UsageNumber {
		return errors.New(my_errors.CurdOrderUsageNumberEqZero)
	}

	//有效终端集合
	var validIds []string
	agentFact := agent_models.CreateAgentFactory()
	agentModels := agentFact.GetItemsByOrgAndAgents(orgId, agentIds, agentFact.ColumnGrant())
	if agentModels == nil || len(agentModels) == 0 {
		return errors.New(my_errors.CurdAgentNotExist)
	}
	for _, model := range agentModels {
		// 订单编号一致 & 授权状态是已授权
		if model.GrantOrderNo == orderNo && model.GrantStatus == enum_agent.GrantStatusAuthorized {
			continue
		}
		validIds = append(validIds, model.AgentId)
	}
	//终端不需要授权
	agentCnt := len(validIds)
	if agentCnt == 0 {
		return nil
	}
	if (orderModel.BuyNumber - orderModel.UsageNumber - agentCnt) < 0 {
		return errors.New(my_errors.CurdOrderAvailableNumberLack)
	}

	//修改终端
	if !agent_callback_service.CreateCurdFactory().BatchGrant(orgId, validIds, orderNo, orderModel.GrantType,
		orderModel.GrantModule, orderModel.ProductTitle, orderModel.ExpireAt, orderModel.GrantDuration) {
		return errors.New(my_errors.CurdOrderGrantAgentFail)
	}

	//修改订单
	updates := map[string]interface{}{
		"usage_number": gorm.Expr(fmt.Sprintf("usage_number+%d", agentCnt)),
	}
	if !crd.masterModel.UpdateById(orderNo, updates) {
		return errors.New(my_errors.CurdSaveFailed)
	}

	//订单日志
	logParam := order_models.OrderLogModel{
		OrgId:         orderModel.OrgId,
		OrderNo:       orderModel.OrderNo,
		OrderStatus:   orderModel.OrderStatus,
		OperateAction: enum_order.ActionGrant2Agent,
		OperateSource: enums.OperateSourceOrganize,
		OperateId:     operateId,
		OperateName:   operateName,
		UsageNumber:   agentCnt,
		UsageAgents:   json_helper.Encode(validIds),
		CreatedAt:     date_helper.CurDatetime(),
	}
	CreateLogFactory().Add(&logParam)

	return nil
}

// ExpiredRemindForOrg 企业授权过期过期提醒
func (crd *OrderCurd) ExpiredRemindForOrg(orgId int64) (int, int, string) {
	version := cur_app.GetAppUsageVersion()
	conditions := order_models.SearchConditions{
		OrgId:       orgId,
		OrderStatus: enum_order.OrderStatusCompleted,
	}
	respRemind := crd.masterModel.CalcRemindByConditions(conditions, crd.masterModel.ColumnCalcOfRemind())
	if respRemind.TotalNumber == 0 {
		//没有完成订单: 不提示
		return 0, version, ""
	}
	if respRemind.UnexpiredNumber > 0 {
		//存在未过期订单: 不提示
		return 0, version, ""
	}
	if respRemind.GoldNumber > 0 {
		//全部过期 & 有正式授权： 提示
		if cur_app.GetAppUsageVersion() == variable.AppUsageVersionOffline {
			return 1, version, my_errors.CurdOrderExpiredToGet
		} else {
			return 1, version, my_errors.CurdOrderExpiredToBuy
		}
	} else {
		//全部过期 & 有正式授权： 提示, 根据不同版本跳转到不同页面
		if cur_app.GetAppUsageVersion() == variable.AppUsageVersionOffline {
			return 2, version, my_errors.CurdOrderExpiredToGet
		} else {
			return 2, version, my_errors.CurdOrderExpiredToBuy
		}

	}
}
