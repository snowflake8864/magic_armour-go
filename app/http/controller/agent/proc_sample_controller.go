package agent

import (
	"github.com/gin-gonic/gin"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/enums/enum_software"
	"goskeleton/app/global/enums/enum_task"
	"goskeleton/app/global/enums/enum_upload"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/proc_service"
	"goskeleton/app/service/agent/task_service"
	"goskeleton/app/service/platform/upload_service"
	"goskeleton/app/service/upload_file"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/response"
)

// ProcSampleController 进程或模块样本相关
type ProcSampleController struct {
	ControllerBase
}

// Get 获取信息
func (u *ProcSampleController) Get(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	result := []map[string]interface{}{} //需初始化

	//终端策略
	policyModels := proc_models.CreateProcSampleFactory().GetListByConditions(proc_models.SearchConditions{
		OrgId:   loginInfo.OrgId,
		AgentId: loginInfo.AgentId,
		Status:  enum_process.SampleStatusDefault,
	})
	if policyModels != nil && len(policyModels) > 0 {
		for _, model := range policyModels {
			dir := model.Dir
			if loginInfo.PlatformType == enum_software.PlatformTypeLinux {
				dir = path_helper.ToStandard(dir)
			} else {
				dir = path_helper.ToWindows(dir)
			}
			result = append(result, map[string]interface{}{
				"aid":    0,
				"p_dir":  dir,
				"p_hash": model.Id,
			})

		}
	}

	response.SuccessReturnData(ctx, result)
}

// Upload 上传样本
func (u *ProcSampleController) Upload(ctx *gin.Context) {
	loginInfo := u.getLogin(ctx)
	//关闭任务
	tasks := []int{enum_task.TypeProcessSampleUpload}
	task_service.CreateCurdFactory().CloseTasks(loginInfo.OrgId, loginInfo.AgentId, tasks, ctx.ClientIP())

	id, _ := ctx.GetPostForm("p_hash")
	if string_helper.IsEmpty(id) {
		response.SuccessSimplify(ctx)
	}

	if r, uploadParam := upload_file.Upload(ctx, upload_file.GetDefaultPath(), []string{".zip"}); r == true {
		uploadParam.Type = enum_upload.TypeProcSample
		uploadParam.OperateSource = enums.OperateSourceAgent
		uploadParam.OperateId = loginInfo.OrgId
		uploadParam.OperateName = loginInfo.AgentId
		if !upload_service.CreateCurdFactory().Add(&uploadParam) {
			response.FailRoutine(ctx, my_errors.AgentParamErrorCode, "保存失败")
			return
		}

		if err := proc_service.CreateSampleCurdFactory().UploadSuccess(loginInfo.OrgId, loginInfo.AgentId, id, &uploadParam); err != nil {
			response.FailRoutine(ctx, my_errors.AgentParamErrorCode, err.Error())
			return
		}
		response.SuccessSimplify(ctx)
	} else {
		//修改状态
		_ = proc_service.CreateSampleCurdFactory().UploadFail(loginInfo.OrgId, loginInfo.AgentId, id)

		response.Success(ctx, my_errors.FilesUploadFailMsg, my_errors.FilesUploadFailCode)
	}
}
