package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/model/agent/agent_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/service/organization/department_service"
	"goskeleton/app/utils/cur_userinfo"
)

// ControllerBase 基础
type ControllerBase struct {
}

// getLogin 获取登录数据
func (u *ControllerBase) getLogin(ctx *gin.Context) agent_service.AgentCacheLogin {
	return cur_userinfo.GetCurrentAgent(ctx)
}

func (u *ControllerBase) getDeptIdsByAgent(orgId int64, agentId string) (deptIds []int64) {
	agentCurdFact := agent_models.CreateAgentFactory()
	agentModel := agentCurdFact.GetItemByOrgAndAgent(orgId, agentId, agentCurdFact.ColumnSimplify())
	if agentModel == nil || agentModel.OrgId == 0 {
		return nil
	}
	//默认看全局
	deptIds = append(deptIds, 0)

	//如果有设置组织架构
	if agentModel.DeptId > 0 {
		//获取本身以及父级
		tmpIds := department_service.CreateCurdFactory().GetIdsSelfWithParentById(orgId, agentModel.DeptId)
		if len(tmpIds) > 0 {
			deptIds = append(deptIds, tmpIds...)
		}
	}
	return deptIds
}
