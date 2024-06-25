package job_callback_syslog

import (
	"errors"
	"goskeleton/app/model/syslog/syslog_models"
)

func CreateLoginCurdFactory() *LoginCurd {
	return &LoginCurd{
		masterModel: syslog_models.CreateLoginFactory(),
	}
}

type LoginCurd struct {
	masterModel *syslog_models.LoginModel
}

// DeleteByIds 通过id集合删除日志
func (crd *LoginCurd) DeleteByIds(orgId int64, ids []int64) error {
	conditions := syslog_models.SearchConditions{
		OrgId: orgId,
		Ids:   ids,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("删除失败")
	}
	return nil
}

// ClearByAgents 清空终端下日志
func (crd *LoginCurd) ClearByAgents(orgId int64, agentIds []string) error {
	conditions := syslog_models.SearchConditions{
		OrgId:    orgId,
		AgentIds: agentIds,
	}
	if !crd.masterModel.DeleteByConditions(conditions) {
		return errors.New("清空失败")
	}
	return nil
}
