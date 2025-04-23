package job_callback_syslog

import (
	"errors"
	"goskeleton/app/global/enums/enum_dictionary"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/syslog/syslog_models"
	"goskeleton/app/service/veda_other/huanshi_platform"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

func CreateProcessUploadFactory() *SyslogProcessUpload {
	return &SyslogProcessUpload{
		masterModel: syslog_models.CreateProcessFactory(),
	}
}

type SyslogProcessUpload struct {
	syslogBase
	masterModel *syslog_models.ProcessModel
}

// Add 添加
func (crd *SyslogProcessUpload) Add(model *syslog_models.ProcessModel) error {
	if string_helper.IsEmpty(model.CreatedAt) {
		model.CreatedAt = date_helper.CurDatetime()
	}
	model.Id = variable.SnowFlake.GetId()
	model.ProcHash = strings.ToLower(model.ProcHash)
	model.ProcDir = path_helper.ToStandard(model.ProcDir)
	model.ProcParam = path_helper.ToStandard(model.ProcParam)
	model.ParentHash = strings.ToLower(model.ParentHash)
	model.ParentDir = path_helper.ToStandard(model.ParentDir)
	model.ParentParam = path_helper.ToStandard(model.ParentParam)

	//保存
	if !crd.masterModel.InsertItem(model) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiLogSwitchProcess); hsCli != nil {
		hsData := huanshi_platform.LogProcess{
			LogBase: huanshi_platform.LogBase{
				OrgId:     model.OrgId,
				AgentId:   model.AgentId,
				CreatedAt: model.CreatedAt,
			},
			ProcHash:    model.ProcHash,
			ProcId:      model.ProcId,
			ProcDir:     model.ProcDir,
			ProcParam:   model.ProcParam,
			ParentHash:  model.ParentHash,
			ParentId:    model.ParentId,
			ParentDir:   model.ParentDir,
			ParentParam: model.ParentParam,
		}
		if err := hsCli.SyncLogProcess(hsData); err != nil {
			variable.ZapLog.Warn(err.Error())
		}
	}

	return nil
}

// BatchAdd 新增
func (crd *SyslogProcessUpload) BatchAdd(list *[]syslog_models.ProcessModel) error {
	if len(*list) == 0 {
		return nil
	}

	createdAt := date_helper.CurDatetime()
	var inserts []syslog_models.ProcessModel
	for _, model := range *list {
		model.Id = variable.SnowFlake.GetId()
		model.ProcHash = strings.ToLower(model.ProcHash)
		model.ProcDir = path_helper.ToStandard(model.ProcDir)
		model.ProcParam = path_helper.ToStandard(model.ProcParam)
		model.ParentHash = strings.ToLower(model.ParentHash)
		model.ParentDir = path_helper.ToStandard(model.ParentDir)
		model.ParentParam = path_helper.ToStandard(model.ParentParam)

		if string_helper.IsEmpty(model.CreatedAt) {
			model.CreatedAt = createdAt
		}
		inserts = append(inserts, model)
	}

	if !crd.masterModel.InsertBatch(&inserts) {
		return errors.New("保存失败")
	}

	//推送到幻势,
	if hsCli := crd.getHuanshiCli(enum_dictionary.HuanshiLogSwitchProcess); hsCli != nil {
		for _, model := range inserts {
			hsData := huanshi_platform.LogProcess{
				LogBase: huanshi_platform.LogBase{
					OrgId:     model.OrgId,
					AgentId:   model.AgentId,
					CreatedAt: model.CreatedAt,
				},
				ProcHash:    model.ProcHash,
				ProcId:      model.ProcId,
				ProcDir:     model.ProcDir,
				ProcParam:   model.ProcParam,
				ParentHash:  model.ParentHash,
				ParentId:    model.ParentId,
				ParentDir:   model.ParentDir,
				ParentParam: model.ParentParam,
			}
			if err := hsCli.SyncLogProcess(hsData); err != nil {
				variable.ZapLog.Warn(err.Error())
			}
		}
	}

	return nil
}
