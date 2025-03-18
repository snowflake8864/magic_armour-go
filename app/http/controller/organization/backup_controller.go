package organization

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums/enum_agent"
	"goskeleton/app/http/transformer/org_trans"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/model/agent/backup_models"
	"goskeleton/app/service/agent/backup_service"
	"goskeleton/app/utils/response"
)

// BackupController 备份回滚相关
type BackupController struct {
	ControllerBase
}

// PageForManage 分页
func (u *BackupController) PageForManage(ctx *gin.Context) {

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
		//规则统计数据
		backupFact := backup_models.CreateBackupFactory()
		calcPolicy := backupFact.CalcByConditionsGroupAgentReturnMap(backup_models.SearchConditions{
			OrgId:    u.getLogin(ctx).OrgId,
			AgentIds: agentIds,
		}, backupFact.ColumnCalc())

		for _, agentModel := range agentModels {
			gg := map[string]interface{}{
				"total_number": 0,
			}
			if item, isOk := calcPolicy[agentModel.AgentId]; isOk {
				gg["total_number"] = item.TotalNumber
			}
			extendInfo[agentModel.AgentId] = gg
		}
	}

	response.SuccessReturnData(ctx, (&org_trans.Agent{}).TransOtherPage(total, agentModels, departmentMaps, extendInfo))
}

// PageForHistory 备份历史: 分页
func (u *BackupController) PageForHistory(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	condition := backup_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: u.getAgentIdFromGin(ctx),
	}
	//分页参数和排序参数
	offset, pageSize, order := u.getPageAndSort(ctx, "")

	total, models := backup_models.CreateBackupFactory().Page(condition, offset, pageSize, order)

	response.SuccessReturnData(ctx, (&org_trans.BackupTrans{}).TransHistoryPage(total, models))
}

// AddForHistory 备份历史: 添加
func (u *BackupController) AddForHistory(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	param := backup_models.BackupModel{
		OrgId:      loginInfo.OrgId,
		AgentId:    u.getAgentIdFromGin(ctx),
		SourcePath: ctx.GetString(consts.ValidatorPrefixOrg + "sourcePath"),
		StorageDir: ctx.GetString(consts.ValidatorPrefixOrg + "storageDir"),
		Remark:     ctx.GetString(consts.ValidatorPrefixOrg + "remark"),
		BackupId:   loginInfo.UserId,
		BackupName: loginInfo.Username,
	}

	if err := backup_service.CreateCurdFactory().Add(&param); err != nil {
		response.FailSimplify(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}
