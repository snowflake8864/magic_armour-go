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
		redisClient:   redis_factory.GetOneRedisClient(),
		cachePrefix:   "lock_proc_upload_",
		cacheExpire:   200,
		queryClient:   producer_kafka.CreateProducerFactory(),
		supportNumber: 30,
	}
}

type ProcPost struct {
	redisClient   *redis_factory.RedisClient
	cachePrefix   string
	cacheExpire   int
	queryClient   *producer_kafka.ProducerClient
	supportNumber int
}

// 给hash加锁
func (crd *ProcPost) lockHash(orgId int64, agentId, hash string) bool {
	key := crd.cachePrefix + int_helper.Int64ToStr(orgId) + "_" + agentId + "_" + hash
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
	//ProductTitle: Microsoft\\u00AE Windows\\u00AE Operating System;
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

// ToQueryByAll 格式化参数后推送给队列: 上传全部进程模块
func (crd *ProcPost) ToQueryByAll(source int, task string, orgId int64, agentId, inputJson string) error {
	var params []procData
	if err := json_helper.ToAny(inputJson, &params); err != nil {
		return err
	}
	if len(params) == 0 {
		return nil
	}

	now := date_helper.CurDatetime()
	var processModels []*proc_models.ProcessModel
	var dllModels []*proc_models.DllModel

	for _, item := range params {
		//hash全部小写
		hash := strings.ToLower(item.Hash)
		if !crd.lockHash(orgId, agentId, hash) {
			continue
		}

		if source == enum_process.DataFromProcess {
			processModel := proc_models.ProcessModel{
				OrgId:     orgId,
				AgentId:   agentId,
				ProcHash:  hash,
				Dir:       item.Dir,
				CreatedAt: now,
			}
			processModel.Id = crd.makeMasterId(orgId, agentId, hash)
			processModel.Copyright, processModel.Introduce = crd.parseAttribute(item.Attribute, item.Copyright, item.Introduce)

			processModels = append(processModels, &processModel)
		} else {
			dllModel := proc_models.DllModel{
				OrgId:     orgId,
				AgentId:   agentId,
				DllHash:   hash,
				Dir:       item.Dir,
				Attribute: item.Attribute,
				CreatedAt: now,
			}
			dllModel.Id = crd.makeMasterId(orgId, agentId, hash)
			dllModel.Copyright, dllModel.Introduce = crd.parseAttribute(item.Attribute, item.Copyright, item.Introduce)

			dllModels = append(dllModels, &dllModel)
		}
	}

	if err := crd.batchSendProcess(orgId, agentId, source, processModels); err != nil {
		variable.ZapLog.Sugar().Info("sendProcess: " + err.Error())
	}

	if err := crd.batchSendDll(orgId, agentId, source, dllModels); err != nil {
		variable.ZapLog.Sugar().Info("batchSendDll: " + err.Error())
	}

	return nil
}

// ToQueryByCommon 格式化参数后推送给队列: 主动上传进程模块 | 刷新进程模块
func (crd *ProcPost) ToQueryByCommon(source int, task string, orgId int64, agentId, inputJson string) error {
	var params []procData
	if err := json_helper.ToAny(inputJson, &params); err != nil {
		return err
	}
	if len(params) == 0 {
		return nil
	}

	var processModels []*proc_models.ProcessModel
	var processRunningModels []proc_models.ProcessRunningModel
	var dllModels []*proc_models.DllModel
	var dllRunningModels []proc_models.DllRunningModel

	now := date_helper.CurDatetime()
	for _, itemProcess := range params {
		//进程hash全部小写
		procHash := strings.ToLower(itemProcess.Hash)

		//进程短时间内没有上传过
		if source == enum_process.DataFromProcess && crd.lockHash(orgId, agentId, procHash) {
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

			processModels = append(processModels, &processModel)
		}

		//进程刷新
		if source == enum_process.DataFromProcess && task == ProcTaskRefresh {
			//运行中进程表
			processRunningModel := proc_models.ProcessRunningModel{
				BaseModel: model.BaseModel{},
				OrgId:     orgId,
				AgentId:   agentId,
				ProcHash:  procHash,
				ProcId:    itemProcess.ProcId,
				ProcDir:   itemProcess.Dir,
				ProcUser:  itemProcess.User,
				DllNumber: itemProcess.ModuleNumber,
				DllStatus: enum_process.CollectDllStatusDefault,
				CreatedAt: now,
			}
			processRunningModel.Id = crd.makeRunningId(orgId, agentId, procHash, itemProcess.ProcId)

			processRunningModels = append(processRunningModels, processRunningModel)
		}

		//有模块数据
		if len(itemProcess.ModuleList) > 0 {
			for _, itemDll := range itemProcess.ModuleList {
				//小写
				dllHash := strings.ToLower(itemDll.Hash)
				//模块短时间内没有上传过
				if source == enum_process.DataFromDll && crd.lockHash(orgId, agentId, dllHash) {
					//模块表数据
					dllModel := proc_models.DllModel{
						OrgId:     orgId,
						AgentId:   agentId,
						DllHash:   dllHash,
						Dir:       itemDll.Dir,
						Attribute: itemDll.Attribute,
						ProcHash:  procHash,
						ProcDir:   itemProcess.Dir,
						CreatedAt: now,
					}
					dllModel.Id = crd.makeMasterId(orgId, agentId, dllHash)
					dllModel.Copyright, dllModel.Introduce = crd.parseAttribute(itemDll.Attribute, "", "")

					dllModels = append(dllModels, &dllModel)
				}

				//模块刷新
				if source == enum_process.DataFromDll && task == ProcTaskRefresh {
					//运行中模块表
					dllRunningModel := proc_models.DllRunningModel{
						BaseModel: model.BaseModel{},
						OrgId:     orgId,
						AgentId:   agentId,
						DllHash:   dllHash,
						DllDir:    itemDll.Dir,
						ProcHash:  procHash,
						ProcId:    itemProcess.ProcId,
						ProcDir:   itemProcess.Dir,
						CreatedAt: now,
					}
					dllRunningModel.Id = crd.makeRunningId(orgId, agentId, dllHash, dllRunningModel.ProcId)

					dllRunningModels = append(dllRunningModels, dllRunningModel)
				}
			}

		}
	}

	if err := crd.batchSendProcess(orgId, agentId, source, processModels); err != nil {
		variable.ZapLog.Sugar().Info("sendProcess: " + err.Error())
	}

	if err := crd.batchSendDll(orgId, agentId, source, dllModels); err != nil {
		variable.ZapLog.Sugar().Info("batchSendDll: " + err.Error())
	}

	if err := crd.batchSendProcessRunning(source, task, &processRunningModels); err != nil {
		variable.ZapLog.Sugar().Info("batchSendProcessRunning: " + err.Error())
	}

	if err := crd.batchSendDllRunning(source, task, &dllRunningModels); err != nil {
		variable.ZapLog.Sugar().Info("batchSendDllRunning: " + err.Error())
	}
	return nil
}

// 发送进程数据
func (crd *ProcPost) batchSendProcess(orgId int64, agentId string, source int, models []*proc_models.ProcessModel) error {
	if source != enum_process.DataFromProcess {
		return nil
	}
	modelLength := len(models)
	if modelLength == 0 {
		return nil
	}

	cxt := context.Background()
	for i := 0; i < modelLength; i += crd.supportNumber {
		end := i + crd.supportNumber
		if end > modelLength {
			end = modelLength
		}
		//proc_job.BatchUploadProcess(orgId, agentId, json_helper.Encode(models[i:end]))

		msg := &kafka_go.MessageStruct{
			Type:    kafka_go.TypeProcessBatchPutMaster,
			OrgId:   orgId,
			AgentId: agentId,
			Data:    json_helper.Encode(models[i:end]),
		}
		crd.queryClient.ProduceProcess(cxt, msg)
	}

	return nil
}

// 发送模块数据
func (crd *ProcPost) batchSendDll(orgId int64, agentId string, source int, models []*proc_models.DllModel) error {
	if source != enum_process.DataFromDll {
		return nil
	}
	modelLength := len(models)
	if modelLength == 0 {
		return nil
	}

	cxt := context.Background()
	for i := 0; i < modelLength; i += crd.supportNumber {
		end := i + crd.supportNumber
		if end > modelLength {
			end = modelLength
		}
		//proc_job.BatchUploadDll(orgId, agentId, json_helper.Encode(models[i:end]))

		msg := &kafka_go.MessageStruct{
			Type:    kafka_go.TypeDllBatchPutMaster,
			OrgId:   orgId,
			AgentId: agentId,
			Data:    json_helper.Encode(models[i:end]),
		}
		crd.queryClient.ProduceDll(cxt, msg)
	}

	return nil
}

// 发送进程运行中数据
func (crd *ProcPost) batchSendProcessRunning(source int, task string, models *[]proc_models.ProcessRunningModel) error {
	if source != enum_process.DataFromProcess {
		return nil
	}
	if task != ProcTaskRefresh {
		return nil
	}
	if len(*models) == 0 {
		return nil
	}

	if !job_callback_proc.CreateProcessRunningUploadFactory().BatchAdd(models) {
		return errors.New("保存失败")
	}
	return nil
}

// 发送模块运行中数据
func (crd *ProcPost) batchSendDllRunning(source int, task string, models *[]proc_models.DllRunningModel) error {
	if source != enum_process.DataFromDll {
		return nil
	}
	if task != ProcTaskRefresh {
		return nil
	}

	if len(*models) == 0 {
		return nil
	}

	if !job_callback_proc.CreateDllRunningUploadFactory().BatchAdd(models) {
		return errors.New("保存失败")
	}
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
			variable.ZapLog.Sugar().Info("CollectDll: " + err.Error())
		}
	}

	return nil
}
