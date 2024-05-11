package agent_post_service

import (
	"fmt"
	"goskeleton/app/global/enums/enum_process"
	"goskeleton/app/model"
	"goskeleton/app/model/model_agent"
	"goskeleton/app/utils/helpers"
	"goskeleton/app/utils/md5_encrypt"
	"goskeleton/app/utils/redis_factory"
	"strings"
)

func CreateProcFactory() *ProcPost {
	return &ProcPost{
		redisClient: redis_factory.GetOneRedisClient(),
		cachePrefix: "lock_proc_upload_",
		cacheExpire: 120,
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
}

type procData struct {
	ProcId       int       `json:"id"`
	Dir          string    `json:"dir"`
	User         string    `json:"user"`
	Hash         string    `json:"hash"`
	Introduce    string    `json:"introduce"`
	Copyright    string    `json:"copyright"`
	Attribute    string    `json:"attribute"`
	Signature    string    `json:"signature"`
	ModuleNumber int       `json:"module_number"`
	ModuleList   []dllData `json:"module"`

	//Finish int `json:"finish"` //没有用
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
func (crd *ProcPost) ToQuery(source int, task string, orgId int64, agentId, inputJson string) (err error) {
	var params []procData
	err = helpers.Json2Any(inputJson, &params)
	if err != nil {
		return err
	}
	now := helpers.DatetimeCurr()
	fmt.Println("params========================================")
	fmt.Println(helpers.JsonEncode(params))
	for _, process := range params {
		//全部小写
		procHash := strings.ToLower(process.Hash)

		//来源：进程
		if source == enum_process.SourceProcess {
			//进程表数据
			processModel := model_agent.ProcessModel{
				OrgId:     orgId,
				AgentId:   agentId,
				ProcHash:  procHash,
				Dir:       process.Dir,
				CreatedAt: now,
			}
			processModel.Id = crd.makeMasterId(orgId, agentId, processModel.ProcHash)
			processModel.Copyright, processModel.Introduce = crd.parseAttribute(process.Attribute, process.Copyright, process.Introduce)
			fmt.Println("processModel========================================")
			fmt.Println(helpers.JsonEncode(processModel))
			//推送给队列

			//任务：刷新
			if task == TaskRefresh {
				//运行中进程表
				processRunningModel := model_agent.ProcessRunningModel{
					BaseModel: model.BaseModel{},
					Id:        "",
					OrgId:     orgId,
					AgentId:   agentId,
					ProcHash:  procHash,
					ProcId:    process.ProcId,
					ProcDir:   process.Dir,
					ProcUser:  process.User,
					DllNumber: process.ModuleNumber,
					DllStatus: enum_process.ProcGetDllStatusDefault,
					CreatedAt: now,
				}
				processRunningModel.Id = crd.makeRunningId(orgId, agentId, procHash, process.ProcId)
				fmt.Println("processRunningModel========================================")
				fmt.Println(helpers.JsonEncode(processRunningModel))
				//推送给队列
			}
		}

		//来源：模块
		if source == enum_process.SourceDll && len(process.ModuleList) > 0 {
			for _, process := range params {

			}
		}
	}

	return nil
}
