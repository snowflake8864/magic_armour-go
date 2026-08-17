package organization

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/global/enums/enum_antivirus_engine_task"
	"goskeleton/app/global/variable"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/antivirus_engine_task_models"
	"goskeleton/app/service/agent/antivirus_engine_task_service"
	"goskeleton/app/service/platform/antivirus_engine_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/response"
)

// AntivirusEngineController 杀毒引擎升级任务相关
type AntivirusEngineController struct {
	ControllerBase
}

// PageForManage 管理: 分页
func (u *AntivirusEngineController) PageForManage(ctx *gin.Context) {
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
	agentIds, departmentMaps := u.parseAgentModels(agentModels, true)

	var extendInfo = make(map[string]interface{})
	if len(agentIds) > 0 {
		//最新引擎
		newestEngineModels := antivirus_engine_service.CreateCurdFactory().GetMapsByNewest()
		for _, agentModel := range agentModels {
			gg := map[string]interface{}{
				"can_upgrade_engine": 0,
			}

			if item, isOk := newestEngineModels[agentModel.PlatformType]; isOk && item.Ver != agentModel.AntivirusEngineVer {
				gg["can_upgrade_engine"] = 1

			}
			extendInfo[agentModel.AgentId] = gg
		}
	}
	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, extendInfo))
}

// UpgradeForManage 管理: 单终端升级
func (u *AntivirusEngineController) UpgradeForManage(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	agentId := u.getAgentIdFromGin(ctx)

	param := &antivirus_engine_task_models.UpgradeTaskModel{
		OrgId:       loginInfo.OrgId,
		Title:       fmt.Sprintf("终端[%s]的升级任务", agentId),
		TargetType:  enum_antivirus_engine_task.TargetTypeSingleAgent,
		TargetIds:   agentId,
		UpgradeMode: enum_antivirus_engine_task.UpgradeModeManual,
		TaskType:    enum_antivirus_engine_task.TaskTypeManualInstantly,
		ExecDate:    date_helper.SpecialDatetime(variable.DateFormatToday),
		ExecHis:     date_helper.SpecialDatetime("15:04:00"),
		CreateId:    loginInfo.UserId,
		CreateName:  loginInfo.Username,
		UpdateId:    loginInfo.UserId,
		UpdateName:  loginInfo.Username,
	}

	if err := antivirus_engine_task_service.CreatCurdFactory().Add(param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
