package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_blackmail"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/blackmail_models"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/response"
	"strings"
)

// BlackmailController 勒索相关
type BlackmailController struct {
	ControllerBase
}

// GetPolicy 下载勒索保护策略
func (u *BlackmailController) GetPolicy(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	go task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeBlackmailPolicyDownload}, ctx.ClientIP())

	openSwitch := enums.SwitchOn
	var existsPolicy = make(map[string]struct{})
	var result []map[string]interface{}

	//终端策略
	policyModels := blackmail_models.CreatePolicyFactory().GetListByConditions(blackmail_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		AgentId:    loginInfo.AgentId,
		OpenSwitch: &openSwitch,
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			policies := u.format(loginInfo.PlatformType, model.Type, model.TypeVal, model.Processes)
			if policies == nil || len(policies) == 0 {
				continue
			}
			for _, policy := range policies {
				suffix := policy["file_suffix"].(string)
				if _, isOK := existsPolicy[suffix]; !isOK {
					existsPolicy[suffix] = struct{}{}
					result = append(result, policy)
				}
			}
		}
	}

	//全局策略
	globalModels := blackmail_models.CreateGlobalFactory().GetListByConditions(blackmail_models.SearchConditions{
		OrgId:        loginInfo.OrgId,
		PlatformType: loginInfo.PlatformType,
		DeptIds:      u.getDeptIdsByAgent(loginInfo.OrgId, loginInfo.AgentId),
		OpenSwitch:   &openSwitch,
	})
	if globalModels != nil && len(globalModels) > 0 {
		for _, model := range globalModels {
			policies := u.format(loginInfo.PlatformType, model.Type, model.TypeVal, model.Processes)
			if policies == nil || len(policies) == 0 {
				continue
			}
			for _, policy := range policies {
				suffix := policy["file_suffix"].(string)
				if _, isOK := existsPolicy[suffix]; !isOK {
					existsPolicy[suffix] = struct{}{}
					result = append(result, policy)
				}
			}
		}
	}

	response.SuccessReturnData(ctx, result)
}

func (u *BlackmailController) format(platformType int, pType int, typeVal string, processStr string) []map[string]interface{} {

	//后缀
	var suffixes []string
	if pType == enum_blackmail.TypeDir {
		if platformType == enum_software.PlatformTypeLinux {
			suffixes = append(suffixes, path_helper.ToStandard(typeVal))
		} else {
			suffixes = append(suffixes, path_helper.ToWindows(typeVal))
		}
	} else {
		suffixes = strings.Split(typeVal, ",")
	}

	var processArr []blackmail_models.Process
	var processMap []map[string]string
	if err := json_helper.ToAny(processStr, &processArr); err != nil {
		return nil
	} else {
		for _, item := range processArr {
			processMap = append(processMap, map[string]string{
				"hash": item.Hash,
				"name": "",
			})
		}
	}

	var list []map[string]interface{}
	for _, suffix := range suffixes {

		list = append(list, map[string]interface{}{
			"type":        pType,
			"file_suffix": suffix,
			"process":     processMap,
		})
	}

	return list
}
