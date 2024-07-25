package plug_job

import (
	"fmt"
	"goskeleton/app/global/variable"
	"goskeleton/app/service/agent/plug_service"
	"goskeleton/app/utils/helpers/date_helper"
)

func DynamicDisable(orgId int64, agentId string) {
	now := date_helper.CurDatetime()

	variable.ZapLog.Sugar().Info(fmt.Sprintf("orgId=%d, agentId=%s, now=%s", orgId, agentId, now))

	if err := plug_service.CreatePolicyDynamicCurdFactory().AutoDisable(orgId, agentId, now); err != nil {
		variable.ZapLog.Sugar().Info(err.Error())
	}
}
