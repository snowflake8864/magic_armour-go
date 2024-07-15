package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/enums/enum_usb"
	"goskeleton/app/model/agent/usb_models"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/response"
)

// UsbController 外设相关
type UsbController struct {
	ControllerBase
}

// Put 添加设备
func (u *UsbController) Put(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeUsbUpload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "data")
	err := agent_post_service.CreateUsbFactory().Save(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// GetBlack 获取禁止列表
func (u *UsbController) GetBlack(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeUsbBlackDownload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	var existsPolicy = make(map[string]struct{})
	result := []*UsbResp{}

	//终端策略
	policyModels := usb_models.CreateUsbFactory().GetListByConditions(usb_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		AgentId:    loginInfo.AgentId,
		StatusList: []int{enum_usb.StatusWhite, enum_usb.StatusBlack},
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			policy := u.format(model.Id, model.Name)
			if _, isOK := existsPolicy[policy.Hash]; !isOK {
				existsPolicy[policy.Hash] = struct{}{}
				if model.Status == enum_usb.StatusBlack {
					result = append(result, policy)
				}
			}
		}
	}

	//全局策略
	globalModels := usb_models.CreateGlobalFactory().GetListByConditions(usb_models.SearchConditions{
		OrgId:  loginInfo.OrgId,
		Status: enum_usb.StatusBlack,
	})
	if globalModels != nil && len(globalModels) > 0 {
		for _, model := range globalModels {
			policy := u.format(model.Id, model.Name)
			if _, isOK := existsPolicy[policy.Hash]; !isOK {
				existsPolicy[policy.Hash] = struct{}{}
				result = append(result, policy)
			}
		}
	}

	response.SuccessReturnData(ctx, result)
}

// GetWhite 获取信任列表
func (u *UsbController) GetWhite(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeUsbWhiteDownload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	var existsPolicy = make(map[string]struct{})
	result := []*UsbResp{} //需初始化

	//终端策略
	policyModels := usb_models.CreateUsbFactory().GetListByConditions(usb_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		AgentId:    loginInfo.AgentId,
		StatusList: []int{enum_usb.StatusWhite, enum_usb.StatusBlack},
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			policy := u.format(model.Id, model.Name)
			if _, isOK := existsPolicy[policy.Hash]; !isOK {
				existsPolicy[policy.Hash] = struct{}{}
				if model.Status == enum_usb.StatusWhite {
					result = append(result, policy)
				}
			}
		}
	}

	//全局策略
	globalModels := usb_models.CreateGlobalFactory().GetListByConditions(usb_models.SearchConditions{
		OrgId:  loginInfo.OrgId,
		Status: enum_usb.StatusWhite,
	})
	if globalModels != nil && len(globalModels) > 0 {
		for _, model := range globalModels {
			policy := u.format(model.Id, model.Name)
			if _, isOK := existsPolicy[policy.Hash]; !isOK {
				existsPolicy[policy.Hash] = struct{}{}
				result = append(result, policy)
			}
		}
	}

	response.SuccessReturnData(ctx, result)
}

type UsbResp struct {
	Hash string `json:"perpheral_eid"`
	Name string `json:"perpheral_name"`
}

func (u *UsbController) format(hash string, name string) *UsbResp {
	data := &UsbResp{}

	data.Hash = hash
	data.Name = name

	return data
}
