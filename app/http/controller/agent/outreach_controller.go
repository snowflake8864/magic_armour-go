package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_outreach"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/outreach_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/response"
)

// OutreachController 外联管控相关
type OutreachController struct {
	ControllerBase
}

// GetDetect 主动探测策略
func (u *OutreachController) GetDetect(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeOutreachDetectDownload}, ctx.ClientIP())

	openSwitch := enums.SwitchOn
	result := OutreachResp{} //需初始化

	result.List = append(result.List, u.formatDetect(enum_outreach.DetectTypeDomain, "https://www.baidu.com"))
	result.List = append(result.List, u.formatDetect(enum_outreach.DetectTypeDomain, "https://hj.veda.com:10110"))
	result.List = append(result.List, u.formatDetect(enum_outreach.DetectTypeDomain, "https://192.168.15.6:443"))
	result.List = append(result.List, u.formatDetect(enum_outreach.DetectTypeDomain, "http://192.168.16.80:8030"))

	//全局策略
	globalModels := outreach_models.CreateDetectFactory().GetListByConditions(outreach_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		OpenSwitch: &openSwitch,
	})
	if globalModels != nil && len(globalModels) > 0 {
		for _, model := range globalModels {
			result.List = append(result.List, u.formatDetect(model.Type, model.Addr))
		}
	}

	response.SuccessReturnData(ctx, result)
}

type OutreachResp struct {
	List []map[string]interface{} `json:"list"`
}

func (u *OutreachController) formatDetect(Type int, addr string) map[string]interface{} {
	data := make(map[string]interface{})
	data["type"] = Type
	data["method"] = "get"
	data["addr"] = addr
	return data
}
