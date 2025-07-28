package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/consts"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_tamper"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/model/agent/tamper_models"
	"goskeleton/app/service/agent/agent_post_service"
	"goskeleton/app/service/agent/dir_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
	"strings"
)

// TamperController 防篡改相关
type TamperController struct {
	ControllerBase
}

// GetDir 获取查询目录
func (u *TamperController) GetDir(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeDirStructureUpload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	result := []map[string]interface{}{} //需初始化

	searchDir := dir_service.CreateDirCurdFactory().GetBySearch(loginInfo.OrgId, loginInfo.AgentId)
	if string_helper.IsEmpty(searchDir) || searchDir == enum_tamper.RootDir {
		pType := enum_tamper.TypeFolder
		if loginInfo.PlatformType == enum_software.PlatformTypeWindows {
			pType = enum_tamper.TypeVolume
		}
		result = append(result, map[string]interface{}{
			"dir":  enum_tamper.RootDir,
			"type": pType,
			"pid":  0, //字段没有使用，但是需要保留
		})
	} else {
		if loginInfo.PlatformType == enum_software.PlatformTypeWindows {
			searchDir = path_helper.ToWindows(searchDir)
		} else {
			searchDir = path_helper.ToStandard(searchDir)
		}
		result = append(result, map[string]interface{}{
			"dir":  searchDir,
			"type": enum_tamper.TypeFolder,
			"pid":  0, //字段没有使用，但是需要保留
		})
	}
	response.SuccessReturnData(ctx, result)
}

// PutDir 上传目录信息
func (u *TamperController) PutDir(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)

	inputJson := ctx.GetString(consts.ValidatorPrefixAgent + "dir")
	err := agent_post_service.CreateAgentFactory().SaveDir(loginInfo.OrgId, loginInfo.AgentId, inputJson)
	if err != nil {
		response.SuccessReturnData(ctx, err.Error())
		return
	}
	response.SuccessSimplify(ctx)
}

// GetPolicy 下载防篡改策略
func (u *TamperController) GetPolicy(ctx *gin.Context) {
	//关闭任务
	loginInfo := u.getLogin(ctx)
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, []int{enum_task.TypeTamperPolicyDownload}, ctx.ClientIP())

	openSwitch := enums.SwitchOn
	var existsPolicy = make(map[string]struct{})
	result := []*TamperResp{} //需初始化

	//终端策略
	policyModels := tamper_models.CreatePolicyFactory().GetListByConditions(tamper_models.SearchConditions{
		OrgId:      loginInfo.OrgId,
		AgentId:    loginInfo.AgentId,
		OpenSwitch: &openSwitch,
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			policy := u.format(loginInfo.PlatformType, model.Type, model.ProtectRw, model.IsExtend, model.Dir,
				model.IncludeSuffix, model.ExcludeSuffix, model.ExcludeFolder, model.Processes)
			if _, isOK := existsPolicy[policy.Dir]; !isOK {
				existsPolicy[policy.Dir] = struct{}{}
				result = append(result, policy)
			}
		}
	}

	//全局策略
	globalModels := tamper_models.CreateGlobalFactory().GetListByConditions(tamper_models.SearchConditions{
		OrgId:        loginInfo.OrgId,
		PlatformType: loginInfo.PlatformType,
		DeptIds:      u.getDeptIdsByAgent(loginInfo.OrgId, loginInfo.AgentId),
		OpenSwitch:   &openSwitch,
	})
	if globalModels != nil && len(globalModels) > 0 {
		for _, model := range globalModels {
			policy := u.format(loginInfo.PlatformType, model.Type, model.ProtectRw, model.IsExtend, model.Dir,
				model.IncludeSuffix, model.ExcludeSuffix, model.ExcludeFolder, model.Processes)
			if _, isOK := existsPolicy[policy.Dir]; !isOK {
				existsPolicy[policy.Dir] = struct{}{}
				result = append(result, policy)
			}
		}
	}

	response.SuccessReturnData(ctx, result)
}

type TamperResp struct {
	Id            int                      `json:"id"`
	Dir           string                   `json:"dir"`
	Type          int                      `json:"type"`
	Hash          string                   `json:"hash"`
	ProtectRw     int                      `json:"protect_rw"`
	IncludeSuffix string                   `json:"include_file"`
	ExcludeSuffix string                   `json:"protect_file"`
	ExcludeFolder string                   `json:"protect_folder"`
	IsExtend      int                      `json:"is_extend"`
	Processes     []map[string]interface{} `json:"process"`
}

func (u *TamperController) format(platformType int, pType int, protectRw, isExtend int, dir string,
	includeSuffix, excludeSuffix, excludeFolder, processStr string) *TamperResp {
	data := &TamperResp{}
	data.Type = pType
	data.ProtectRw = protectRw
	data.IsExtend = isExtend
	data.IncludeSuffix = includeSuffix
	data.ExcludeSuffix = excludeSuffix

	//路径
	if platformType == enum_software.PlatformTypeLinux {
		data.Dir = path_helper.ToStandard(dir)
		data.Dir = "/" + strings.TrimLeft(data.Dir, "/")
		data.Dir = strings.TrimRight(data.Dir, "/")
	} else {
		data.Dir = path_helper.ToWindows(dir)
		data.Dir = strings.TrimRight(data.Dir, "\\")
	}

	//排除目录
	if !string_helper.IsEmpty(excludeFolder) {
		glue := "|"
		var fullArr []string
		childArr := strings.Split(excludeFolder, glue)
		for _, child := range childArr {
			tmp := data.Dir
			if platformType == enum_software.PlatformTypeLinux {
				tmp = tmp + "/" + child
				tmp = "/" + strings.TrimLeft(tmp, "/")
				tmp = strings.TrimRight(tmp, "/")
			} else {
				tmp = tmp + "\\" + child
			}
			fullArr = append(fullArr, tmp)
		}
		data.ExcludeFolder = strings.Join(fullArr, glue)
	}

	//信任进程
	var processArr []tamper_models.Process
	processMap := []map[string]interface{}{}
	if err := json_helper.ToAny(processStr, &processArr); err == nil {
		for _, item := range processArr {
			processMap = append(processMap, map[string]interface{}{
				"hash": item.Hash,
			})
		}
	}
	data.Processes = processMap

	return data

}
