package software_models

import (
	"goskeleton/app/global/enums"
)

type SearchConditions struct {
	Id           int64
	PlatformType int
	OsType       int
	IsPublish    int
}
type baseSoftware struct {
}

// search2Where 查询条件转化
func (crd *baseSoftware) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {
	wheres := make(map[string]interface{})
	if search.Id > 0 {
		wheres["id"] = search.Id
	}
	if search.PlatformType > 0 && search.PlatformType != enums.AllDefaultInt {
		wheres["platform_type"] = search.PlatformType
	}
	if search.OsType > 0 && search.OsType != enums.AllDefaultInt {
		wheres["os_type"] = search.OsType
	}
	if search.IsPublish > 0 && search.IsPublish != enums.AllDefaultInt {
		wheres["is_publish"] = search.IsPublish
	}
	return wheres
}
