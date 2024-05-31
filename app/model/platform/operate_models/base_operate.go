package operate_models

import (
	"goskeleton/app/global/enums"
	"goskeleton/app/model/model_comm"
	"goskeleton/app/utils/helpers/string_helper"
)

type SearchConditions struct {
	Module     string
	NetIp      string
	CreateName string
}

type baseOperate struct {
}

// conditions2Where 查询条件转化
func (crd *baseOperate) conditions2Where(condition SearchConditions) (wheres []model_comm.Conditions) {

	if !string_helper.IsEmpty(condition.Module) && condition.Module != enums.AllDefaultStr {
		wheres = append(wheres, model_comm.Conditions{
			Key: "module",
			Val: condition.Module,
		})
	}
	if !string_helper.IsEmpty(condition.NetIp) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "net_ip",
			Val: condition.NetIp,
		})
	}
	if !string_helper.IsEmpty(condition.CreateName) {
		wheres = append(wheres, model_comm.Conditions{
			Key: "create_name",
			Val: condition.CreateName,
		})
	}
	return wheres
}
