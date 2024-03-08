package agent_service

import (
	"goskeleton/app/model/model_platform"
)

func CreateCurdFactory() *AdminCurd {
	return &AdminCurd{masterModel: model_platform.CreateAdminFactory()}
}

type AdminCurd struct {
	masterModel *model_platform.AdminModel
}

// ColumnConfig 字段:全部配置相关
func (crd *AdminCurd) ColumnConfig(needId bool) string {
	return ""
}
