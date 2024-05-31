package agent_post_service

import (
	"context"
	"errors"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/model/agent/proc_models"
	"goskeleton/app/service/agent/proc_service/job_callback_proc"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/int_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/path_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/kafka_go/producer_kafka"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
	"strings"
)

func CreateProcFactory() *ProcPost {
	return &ProcPost{
		redisClient: redis_factory.GetOneRedisClient(),
		cachePrefix: "lock_proc_upload_",
		cacheExpire: 120,
		queryClient: producer_kafka.CreateProducerFactory(),
	}
}

type ProcPost struct {
	redisClient *redis_factory.RedisClient
	cachePrefix string
	cacheExpire int
	queryClient *producer_kafka.ProducerClient
}

// 给hash加锁
func (crd *ProcPost) lockHash(orgId int64, agentId, hash string) bool {
	key := int_helper.Int64ToStr(orgId) + "_" + agentId + "_" + hash
	if isOk, _ := crd.redisClient.Lock(key, "1", crd.cacheExpire); isOk {
		return true
	} else {
		return false
	}
}

// 转换属性
// @return copyright, introduce
func (crd *ProcPost) parseAttribute(attribute, copyright, introduce string) (string, string) {
	if string_helper.IsEmpty(attribute) {
		return copyright, introduce
	}

	//CompanyName: Microsoft Corporation;
	//FileDescription: Network Store Information RPC interface;
	//FileVersion: 10.0.22621.1 (WinBuild.160101.0800);InternalName: winnsi.dll;
	//LegalCopyright: \\u00A9 Microsoft Corporation. All rights reserved.;
	//OriginalFilename: winnsi.dll;
	//ProductName: Microsoft\\u00AE Windows\\u00AE Operating System;
	//ProductVersion: 10.0.22621.1;

	for _, item := range strings.Split(attribute, ";") {
		keyAndVal := strings.Split(item, ":")
		if len(keyAndVal) != 2 {
			continue
		}
		if keyAndVal[0] == "LegalCopyright" {
			copyright = strings.Trim(keyAndVal[1], " ")
		} else if keyAndVal[0] == "FileDescription" {
			introduce = strings.Trim(keyAndVal[1], " ")
		}
	}
	return copyright, introduce
}

// 生成主表id
func (crd *ProcPost) makeMasterId(orgId int64, agentId, hash string) string {
	return md5_encrypt.MD5(int_helper.Int64ToStr(orgId) + agentId + hash)
}

// 生成运行中标id
func (crd *ProcPost) makeRunningId(orgId int64, agentId, hash string, procId int) string {
	return md5_encrypt.MD5(int_helper.Int64ToStr(orgId) + agentId + hash + int_helper.IntToStr(procId))
}

// ToQuery 格式化参数后推送给队列
func (crd *ProcPost) ToQuery(source int, task string, orgId int64, agentId, inputJson string) error {
	var params []procData
	if err := json_helper.ToAny(inputJson, &params); err != nil {
		return err
	}
	now := date_helper.CurDatetime()
	for _, itemProcess := range params {
		//进程hash全部小写
		procHash := strings.ToLower(itemProcess.Hash)

		//进程表数据
		processModel := proc_models.ProcessModel{
			OrgId:     orgId,
			AgentId:   agentId,
			ProcHash:  procHash,
			Dir:       itemProcess.Dir,
			CreatedAt: now,
		}
		processModel.Id = crd.makeMasterId(orgId, agentId, procHash)
		processModel.Copyright, processModel.Introduce = crd.parseAttribute(itemProcess.Attribute, itemProcess.Copyright, itemProcess.Introduce)
		//推送给队列
		if err := crd.sendProcess(source, task, &processModel); err != nil {
			variable.ZapLog.Sugar().Debug("sendProcess: " + err.Error())
		}

		//运行中进程表
		processRunningModel := proc_models.ProcessRunningModel{
			BaseModel: model.BaseModel{},
			OrgId:     orgId,
			AgentId:   agentId,
			ProcHash:  procHash,
			ProcId:    itemProcess.ProcId,
			ProcDir:   processModel.Dir,
			ProcUser:  itemProcess.User,
			DllNumber: itemProcess.ModuleNumber,
			DllStatus: enum_process.CollectDllStatusDefault,
			CreatedAt: now,
		}
		processRunningModel.Id = crd.makeRunningId(orgId, agentId, procHash, itemProcess.ProcId)
		//推送给队列
		if err := crd.sendProcessRunning(source, task, &processRunningModel); err != nil {
			variable.ZapLog.Sugar().Debug("sendProcessRunning: " + err.Error())
		}

		//有模块数据
		if len(itemProcess.ModuleList) > 0 {
			for _, itemDll := range itemProcess.ModuleList {
				//小写
				dllHash := strings.ToLower(itemDll.Hash)
				//模块表数据
				dllModel := proc_models.DllModel{
					OrgId:     orgId,
					AgentId:   agentId,
					DllHash:   dllHash,
					Dir:       itemDll.Dir,
					Attribute: itemDll.Attribute,
					ProcHash:  procHash,
					ProcDir:   processModel.Dir,
					CreatedAt: now,
				}
				dllModel.Id = crd.makeMasterId(orgId, agentId, dllHash)
				dllModel.Copyright, dllModel.Introduce = crd.parseAttribute(itemDll.Attribute, "", "")
				//推送给队列
				if err := crd.sendDll(source, task, &dllModel); err != nil {
					variable.ZapLog.Sugar().Debug("sendDll: " + err.Error())
				}

				//运行中模块表
				dllRunningModel := proc_models.DllRunningModel{
					BaseModel: model.BaseModel{},
					OrgId:     orgId,
					AgentId:   agentId,
					DllHash:   dllHash,
					DllDir:    itemDll.Dir,
					ProcHash:  procHash,
					ProcId:    itemProcess.ProcId,
					ProcDir:   processModel.Dir,
					CreatedAt: now,
				}
				dllRunningModel.Id = crd.makeRunningId(orgId, agentId, dllHash, dllRunningModel.ProcId)
				//推送给队列
				if err := crd.sendDllRunning(source, task, &dllRunningModel); err != nil {
					variable.ZapLog.Sugar().Debug("sendDllRunning: " + err.Error())
				}
			}

		}
	}

	return nil
}

// 发送进程数据
func (crd *ProcPost) sendProcess(source int, task string, model *proc_models.ProcessModel) error {
	if source != enum_process.SourceProcess {
		return errors.New("数据主体不是进程")
	}
	if task == ProcTaskAuto && !crd.lockHash(model.OrgId, model.AgentId, model.ProcHash) {
		return errors.New("agent不能连续主动上传进程数据")
	}

	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeProcessPutMaster,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceProcess(context.Background(), msg)

	return nil
}

// 发送进程运行中数据
func (crd *ProcPost) sendProcessRunning(source int, task string, model *proc_models.ProcessRunningModel) error {
	if source != enum_process.SourceProcess {
		return errors.New("数据主体不是进程")
	}
	if task != ProcTaskRefresh {
		return errors.New("任务不是刷新")
	}

	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeProcessPutRunning,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceProcess(context.Background(), msg)

	return nil
}

// 发送模块数据
func (crd *ProcPost) sendDll(source int, task string, model *proc_models.DllModel) error {
	if source != enum_process.SourceDll {
		return errors.New("数据主体不是模块")
	}
	if task == ProcTaskAuto && !crd.lockHash(model.OrgId, model.AgentId, model.DllHash) {
		return errors.New("agent不能连续主动上传模块数据")
	}

	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeDllPutMaster,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceProcess(context.Background(), msg)

	return nil
}

// 发送模块运行中数据
func (crd *ProcPost) sendDllRunning(source int, task string, model *proc_models.DllRunningModel) error {
	if source != enum_process.SourceDll {
		return errors.New("数据主体不是模块")
	}
	if task != ProcTaskRefresh {
		return errors.New("任务不是刷新")
	}

	msg := &kafka_go.MessageStruct{
		Type: kafka_go.TypeDllPutRunning,
		Data: json_helper.Encode(&model),
	}
	crd.queryClient.ProduceProcess(context.Background(), msg)

	return nil
}

// CollectDll 收集指定进程模块
func (crd *ProcPost) CollectDll(orgId int64, agentId, inputJson string) error {
	var params []procData
	if err := json_helper.ToAny(inputJson, &params); err != nil {
		return err
	}
	for _, itemProcess := range params {
		//进程hash全部小写
		procHash := strings.ToLower(itemProcess.Hash)

		var dllList []proc_models.ProcessRunningDll
		if len(itemProcess.ModuleList) > 0 {
			for _, itemDll := range itemProcess.ModuleList {

				tmp := proc_models.ProcessRunningDll{
					Hash:      strings.ToLower(itemDll.Hash),
					Dir:       path_helper.ToStandard(itemDll.Dir),
					Attribute: itemDll.Attribute,
				}
				tmp.Copyright, tmp.Introduce = crd.parseAttribute(itemDll.Attribute, "", "")
				dllList = append(dllList, tmp)
			}
		}

		//修改
		if err := job_callback_proc.CreateProcessRunningUploadFactory().UploadDlls(orgId, agentId, procHash, itemProcess.ProcId,
			itemProcess.ModuleNumber, &dllList); err != nil {
			variable.ZapLog.Sugar().Debug("CollectDll: " + err.Error())
		}
	}

	return nil
}
