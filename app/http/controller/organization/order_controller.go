package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/enums/enum_order"
	"goskeleton/app/global/enums/enum_upload"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/http/transformer"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/organization/order_models"
	"goskeleton/app/service/organization/order_service"
	"goskeleton/app/service/platform/dictionary_service"
	"goskeleton/app/service/platform/product_service"
	"goskeleton/app/service/platform/upload_service"
	"goskeleton/app/service/upload_file"
	"goskeleton/app/utils/response"
)

// OrderController  订单授权相关
type OrderController struct {
	ControllerBase
}

// OrderStatus  订单状态
func (u *OrderController) OrderStatus(ctx *gin.Context) {

	keys := enums.Get1stIsStrReturnSliceOfInt("order_status", enum_order.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("order_status", enum_order.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// ExpireStatus 过期状态
func (u *OrderController) ExpireStatus(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("expire_status", enum_order.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("expire_status", enum_order.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, true))
}

// GrantDuration 授权时长
func (u *OrderController) GrantDuration(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("grant_duration", enum_order.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("grant_duration", enum_order.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// PayMethod 支付方式
func (u *OrderController) PayMethod(ctx *gin.Context) {
	keys := enums.Get1stIsStrReturnSliceOfInt("pay_method", enum_order.WebMapInt)
	labels := enums.Get1stIsStrReturnMapIntOfStr("pay_method", enum_order.EnumMap)
	response.SuccessReturnData(ctx, (&transformer.Common{}).FormatIntList(keys, labels, false))
}

// PayBankcard 支付银行卡信息
func (u *OrderController) PayBankcard(ctx *gin.Context) {

	maps := dictionary_service.CreateCurdFactory().GetMapByModule(enum_dictionary.ModuleOrder)

	response.SuccessReturnData(ctx, enum_order.FormatPayBankcard(maps))
}

// PublishedForProduct 产品:已发布的
func (u *OrderController) PublishedForProduct(ctx *gin.Context) {
	models := product_service.CreateCurdFactory().GetPublished()

	response.SuccessReturnData(ctx, (&org_trans.Product{}).TransList(models))
}

// ExpiredRemind 到期提醒
func (u *OrderController) ExpiredRemind(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	status, version, message := order_service.CreateCurdFactory().ExpiredRemindForOrg(loginInfo.OrgId)

	response.SuccessReturnData(ctx, gin.H{
		"status":  status,
		"version": version,
		"message": message,
	})
}

// ValidList 有效订单列表
func (u *OrderController) ValidList(ctx *gin.Context) {

	loginInfo := u.getLogin(ctx)
	conditions := order_models.SearchConditions{}
	conditions.OrgId = loginInfo.OrgId
	isExpire := enums.IsNot
	conditions.IsExpire = &isExpire
	conditions.OrderStatus = enum_order.OrderStatusCompleted
	orderModels := order_models.CreateOrderFactory().GetListByConditions(conditions)
	response.SuccessReturnData(ctx, (&org_trans.Order{}).TransList(orderModels))
}

func (u *OrderController) getOrderSearchConditions(ctx *gin.Context) order_models.SearchConditions {
	loginInfo := u.getLogin(ctx)
	conditions := order_models.SearchConditions{}
	conditions.OrgId = loginInfo.OrgId

	conditions.OrderStatus = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "orderStatus"))
	conditions.ProductTitle = ctx.GetString(consts.ValidatorPrefixOrg + "productName")
	conditions.CreateAtGte = ctx.GetString(consts.ValidatorPrefixOrg + "startTime")
	conditions.CreateAtLte = ctx.GetString(consts.ValidatorPrefixOrg + "endTime")
	if val, isOk := ctx.Get(consts.ValidatorPrefixOrg + "isExpire"); isOk && val != nil {
		isExpire := int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "isExpire"))
		conditions.IsExpire = &isExpire
	}
	return conditions
}

// PageForPurchase 购买：分页
func (u *OrderController) PageForPurchase(ctx *gin.Context) {
	//查询条件
	conditions := u.getOrderSearchConditions(ctx)

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	//终端分页
	total, orderModels := order_models.CreateOrderFactory().Page(conditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.Order{}).TransPage(total, orderModels))
}

// InfoForPurchase 购买：信息
func (u *OrderController) InfoForPurchase(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	orderNo := ctx.GetString(consts.ValidatorPrefixOrg + "orderNo")

	model := order_models.CreateOrderFactory().GetItemByOrgAndId(loginInfo.OrgId, orderNo)
	response.SuccessReturnData(ctx, (&org_trans.Order{}).TransDetail(model))
}

// LogsForPurchase 购买：日志
func (u *OrderController) LogsForPurchase(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	conditions := order_models.SearchConditions{
		OrderNo: ctx.GetString(consts.ValidatorPrefixOrg + "orderNo"),
		OrgId:   loginInfo.OrgId,
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	//日志分页
	total, logModels := order_models.CreateLogFactory().Page(conditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.Order{}).TransLogsPage(total, logModels))
}

// AddForPurchase 购买：添加
func (u *OrderController) AddForPurchase(ctx *gin.Context) {

	loginInfo := u.getLogin(ctx)
	param := order_models.OrderModel{}
	param.OrgId = loginInfo.OrgId
	param.ProductId = int64(ctx.GetFloat64(consts.ValidatorPrefixOrg + "productId"))
	param.BuyNumber = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "buyNumber"))
	param.GrantDuration = int(ctx.GetFloat64(consts.ValidatorPrefixOrg + "grantDuration"))
	param.Source = enums.OperateSourceOrganize

	if errCode, err := order_service.CreateCurdFactory().Add(&param, loginInfo.UserId, loginInfo.Username); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}

// PayByBankcardForPurchase 购买：银卡支付
func (u *OrderController) PayByBankcardForPurchase(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	orderNo := ctx.GetString(consts.ValidatorPrefixOrg + "orderNo")

	if errCode, err := order_service.CreateCurdFactory().PayByBankcard(orderNo, loginInfo.OrgId, loginInfo.UserId, loginInfo.Username); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// CancelForPurchase 购买：取消
func (u *OrderController) CancelForPurchase(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	orderNo := ctx.GetString(consts.ValidatorPrefixOrg + "orderNo")

	if errCode, err := order_service.CreateCurdFactory().OrgCancel(orderNo, loginInfo.OrgId, loginInfo.UserId, loginInfo.Username); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// DeleteForPurchase 购买：删除
func (u *OrderController) DeleteForPurchase(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	orderNo := ctx.GetString(consts.ValidatorPrefixOrg + "orderNo")

	if errCode, err := order_service.CreateCurdFactory().OrgDelete(orderNo, loginInfo.OrgId); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForGrant 终端授权：分页
func (u *OrderController) PageForGrant(ctx *gin.Context) {
	//查询条件
	conditions, err := u.getAgentSearchConditions(ctx)
	if err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}

	//分页参数和排序参数
	offset, pageSize, _ := u.getPageAndSort(ctx, "")

	//终端分页
	total, agentModels := agent_models.CreateAgentFactory().Page(conditions, offset, pageSize, enum_agent.AgentFieldDefault)
	//转化数据
	_, departmentMaps := u.parseAgentModels(agentModels, true)

	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, nil))
}

// SaveForGrant 终端授权：保存
func (u *OrderController) SaveForGrant(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentIds := u.getAgentIdsFromGin(ctx)
	orderNo := ctx.GetString(consts.ValidatorPrefixOrg + "orderNo")

	if err := order_service.CreateCurdFactory().GrantToAgents(loginInfo.OrgId, orderNo, agentIds, loginInfo.UserId, loginInfo.Username); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// PageForUsed 已使用：分页
func (u *OrderController) PageForUsed(ctx *gin.Context) {
	//查询条件
	conditions := u.getOrderSearchConditions(ctx)
	conditions.OrderStatus = enum_order.OrderStatusCompleted

	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	//终端分页
	total, orderModels := order_models.CreateOrderFactory().Page(conditions, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.Order{}).TransPage(total, orderModels))
}

// ImportForLicense  授权许可:导入
func (u *OrderController) ImportForLicense(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	//上传文件
	isOk, uploadParam := upload_file.Upload(ctx, upload_file.GetDefaultPath(), []string{".lic"})
	if !isOk {
		response.FailSimplify(ctx, my_errors.FilesUploadFailMsg)
		return
	}
	uploadParam.Type = enum_upload.TypeLicense
	uploadParam.OperateSource = enums.OperateSourceOrganize
	uploadParam.OrgId = loginInfo.OrgId
	uploadParam.OperateId = loginInfo.UserId
	uploadParam.OperateName = loginInfo.Username
	if !upload_service.CreateCurdFactory().Add(&uploadParam) {
		response.FailSimplify(ctx, my_errors.CurdSaveFailed)
		return
	}

	//根据授权证书生成订单
	orderParam := order_models.OrderModel{}
	orderParam.OrgId = loginInfo.OrgId
	orderParam.Source = enums.OperateSourceOrganize
	orderParam.PayOrderNo = uploadParam.FileHash
	if errCode, err := order_service.CreateCurdFactory().ImportLicense(&orderParam, uploadParam.FullPath, loginInfo.UserId, loginInfo.Username); err != nil {
		response.FailRoutine(ctx, errCode, err.Error())
		return
	}

	response.SuccessSimplify(ctx)
}
