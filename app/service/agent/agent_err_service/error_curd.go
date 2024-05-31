package agent_err_service

import (
	"errors"
	"goskeleton/app/model/agent/agent_error_models"
	"goskeleton/app/model/common/upload_models"
	"goskeleton/app/service/agent/agent_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
)

func CreateErrorCurdFactory() *ErrorCurd {
	return &ErrorCurd{
		masterModel: agent_error_models.CreateErrorFactory(),
	}
}

type ErrorCurd struct {
	masterModel *agent_error_models.ErrorModel
}

// Add 新增
func (crd *ErrorCurd) Add(loginInfo *agent_service.AgentCacheLogin, remark string, uploadModel *upload_models.UploadModel) error {
	masterModel := agent_error_models.ErrorModel{
		Id:           uploadModel.Id,
		OrgId:        loginInfo.OrgId,
		AgentId:      loginInfo.AgentId,
		OsName:       loginInfo.Os,
		OriginalName: uploadModel.OriginalName,
		StoragePath:  path_helper.ToStandard(uploadModel.StoragePath),
		FullPath:     path_helper.ToStandard(uploadModel.FullPath),
		Remark:       remark,
		CreatedAt:    date_helper.CurDatetime(),
	}

	if !crd.masterModel.InsertItem(&masterModel) {
		return errors.New("入库失败")
	}

	return nil
}

// DelByOrg 删除：企业
func (crd *ErrorCurd) DelByOrg(orgId int64) bool {
	search := agent_error_models.SearchConditions{
		OrgId: orgId,
	}
	return crd.masterModel.DeleteByConditions(search)
}

// DelByOrgAndAgent 删除：企业和终端
func (crd *ErrorCurd) DelByOrgAndAgent(orgId int64, agentId string) bool {
	search := agent_error_models.SearchConditions{
		OrgId:   orgId,
		AgentId: agentId,
	}
	return crd.masterModel.DeleteByConditions(search)
}
