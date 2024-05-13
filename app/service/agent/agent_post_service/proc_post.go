package agent_post_service

import (
	"context"
	"errors"
	"fmt"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/global/variable"
	"goskeleton/app/model"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/kafka_go"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
	"strings"
)

func CreateProcFactory() *ProcPost {
	return &ProcPost{
		redisClient: redis_factory.GetOneRedisClient(),
		cachePrefix: "lock_proc_upload_",
		cacheExpire: 120,
		queryClient: kafka_go.CreateProducerFactory(),
	}
}

const (
	TaskAll     string = "all"
	TaskAuto    string = "auto"
	TaskRefresh string = "refresh"
)

type ProcPost struct {
	redisClient *redis_factory.RedisClient
	cachePrefix string
	cacheExpire int
	queryClient *kafka_go.ProducerClient
}

type procData struct {
	ProcId       int       `json:"id"`
	Dir          string    `json:"dir"`
	User         string    `json:"user"`
	Hash         string    `json:"hash"`
	Introduce    string    `json:"introduce"`
	Copyright    string    `json:"copyright"`
	Attribute    string    `json:"attribute"`
	ModuleNumber int       `json:"module_number"`
	ModuleList   []dllData `json:"module"`
	Signature    string    `json:"signature"` //不在使用 @todo bug6830
	Finish       int       `json:"finish"`    //没有使用
}

type dllData struct {
	Dir       string `json:"name"`
	Hash      string `json:"hash"`
	Attribute string `json:"attribute"`
}

// 给hash加锁
func (crd *ProcPost) lockHash(orgId int64, agentId, hash string) bool {
	key := helpers.Int64ToStr(orgId) + "_" + agentId + "_" + hash
	if isOk, _ := crd.redisClient.Lock(key, "1", crd.cacheExpire); isOk {
		return true
	} else {
		return false
	}
}

// 转换属性
// @return copyright, introduce
func (crd *ProcPost) parseAttribute(attribute, copyright, introduce string) (string, string) {
	if helpers.StringIsEmpty(attribute) {
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
	return md5_encrypt.MD5(helpers.Int64ToStr(orgId) + agentId + hash)
}

// 生成运行中标id
func (crd *ProcPost) makeRunningId(orgId int64, agentId, hash string, procId int) string {
	return md5_encrypt.MD5(helpers.Int64ToStr(orgId) + agentId + hash + helpers.IntToStr(procId))
}

// ToQuery 格式化参数后推送给队列
func (crd *ProcPost) ToQuery(source int, task string, orgId int64, agentId, inputJson string) error {
	var params []procData
	if err := helpers.Json2Any(inputJson, &params); err != nil {
		return err
	}
	now := helpers.DatetimeCurr()
	fmt.Println("params========================================")
	fmt.Println(helpers.JsonEncode(params))
	for _, itemProcess := range params {
		//进程hash全部小写
		procHash := strings.ToLower(itemProcess.Hash)

		//进程表数据
		processModel := model_agent.ProcessModel{
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
		processRunningModel := model_agent.ProcessRunningModel{
			BaseModel: model.BaseModel{},
			OrgId:     orgId,
			AgentId:   agentId,
			ProcHash:  procHash,
			ProcId:    itemProcess.ProcId,
			ProcDir:   itemProcess.Dir,
			ProcUser:  itemProcess.User,
			DllNumber: itemProcess.ModuleNumber,
			DllStatus: enum_process.ProcGetDllStatusDefault,
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
				dllHash := strings.ToLower(itemDll.Hash)
				//模块表数据
				dllModel := model_agent.DllModel{
					OrgId:     orgId,
					AgentId:   agentId,
					DllHash:   dllHash,
					Dir:       itemDll.Dir,
					CreatedAt: now,
				}
				dllModel.Id = crd.makeMasterId(orgId, agentId, dllHash)
				dllModel.Copyright, dllModel.Introduce = crd.parseAttribute(itemDll.Attribute, "", "")
				//推送给队列

				//运行中模块表
				dllRunningModel := model_agent.DllRunningModel{
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
				//推送给队列
			}

		}
	}

	return nil
}

// 发送进程数据
func (crd *ProcPost) sendProcess(source int, task string, model *model_agent.ProcessModel) error {
	if source != enum_process.SourceProcess {
		return errors.New("数据主体不是进程")
	}
	if task == TaskAuto && !crd.lockHash(model.OrgId, model.AgentId, model.ProcHash) {
		return errors.New("agent主动上传不能连续主动上传进程数据")
	}
	fmt.Println("sendProcess =============================")
	fmt.Println(helpers.JsonEncode(&model))

	msg := &kafka_go.MessageStruct{
		Sign: kafka_go.SignProc,
		Type: kafka_go.TypeProcMaster,
		Data: helpers.JsonEncode(&model),
	}
	crd.queryClient.ProduceProcess(context.Background(), msg)

	return nil
}

// 发送进程运行中数据
func (crd *ProcPost) sendProcessRunning(source int, task string, model *model_agent.ProcessRunningModel) error {
	if source != enum_process.SourceProcess {
		return errors.New("数据主体不是进程")
	}
	if task != TaskRefresh {
		return errors.New("任务不是刷新")
	}
	fmt.Println("sendProcessRunning =============================")
	fmt.Println(helpers.JsonEncode(&model))

	msg := &kafka_go.MessageStruct{
		Sign: kafka_go.SignProc,
		Type: kafka_go.TypeProcRunning,
		Data: helpers.JsonEncode(&model),
	}
	crd.queryClient.ProduceProcess(context.Background(), msg)

	return nil
}
