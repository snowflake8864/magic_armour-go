package threaten_service

import (
	"goskeleton/app/model/model_agent"
)

func CreatePushFactory() *ThreatenPush {
	return &ThreatenPush{masterModel: model_agent.CreateThreatenPushFactory()}
}

type ThreatenPush struct {
	ThreatenBase
	masterModel *model_agent.ThreatenPushModel
}

// Page 分页
func (crd *ThreatenPush) Page(search SearchConditions, offset, pageSize int, orderField, orderSort string) (total int64, items []*model_agent.ThreatenPushModel) {
	wheres := crd.search2Where(search)
	order := orderField + " " + orderSort
	total, items = crd.masterModel.Page(wheres, offset, pageSize, order)
	return total, items
}

// Clear 清空
func (crd *ThreatenPush) Clear() bool {
	return crd.masterModel.TruncateTable()
}
