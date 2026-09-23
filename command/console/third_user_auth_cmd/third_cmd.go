package third_user_auth_cmd

import (
	"errors"
	"fmt"
	"go.uber.org/zap"
	"goskeleton/app/global/enums"
	"goskeleton/app/global/enums/enum_organization"
	"goskeleton/app/global/enums/enum_third_user_auth"
	"goskeleton/app/global/variable"
	"goskeleton/app/model/organization/department_models"
	"goskeleton/app/model/organization/organization_models"
	"goskeleton/app/model/platform/third_user_auth_models"
	"goskeleton/app/service/organization/department_service"
	"goskeleton/app/service/organization/organize_service"
	"goskeleton/app/service/platform/third_user_auth_service"
	"goskeleton/app/utils/helpers/date_helper"
	"goskeleton/app/utils/helpers/json_helper"
	"goskeleton/app/utils/helpers/string_helper"
	"goskeleton/app/utils/redis_factory"
	"slices"
)

func CreateCmdFactory() *ThirdUserAuthCmd {
	return &ThirdUserAuthCmd{
		redisClient: redis_factory.GetOneRedisClient(),
		lockKey:     "lock_third_user_auth",
		cacheExpire: 600,
	}
}

type ThirdUserAuthCmd struct {
	redisClient    *redis_factory.RedisClient
	lockKey        string
	cacheExpire    int
	requestFactory *third_user_auth_service.RequestThird
	configFactory  *third_user_auth_service.Config
	authConfig     *third_user_auth_models.AuthModel

	orgMaxSeq  int64 //组织同步序列
	userMaxSeq int64 //用户同步序列
}

// 设置锁
func (crd *ThirdUserAuthCmd) setLock() bool {
	if isOk, _ := crd.redisClient.Lock(crd.lockKey, date_helper.CurDatetime(), crd.cacheExpire); isOk {
		return true
	} else {
		return false
	}
}

// 删除锁
func (crd *ThirdUserAuthCmd) delLock() bool {
	if _, err := crd.redisClient.Del(crd.lockKey); err == nil {
		return true
	} else {
		return false
	}
}

// 初始化工厂
func (crd *ThirdUserAuthCmd) iniFact() error {
	//获取最新配置
	crd.configFactory = third_user_auth_service.CreateConfigFactory()
	crd.authConfig = crd.configFactory.GetNewestItem()
	if crd.authConfig == nil || string_helper.IsEmpty(crd.authConfig.Addr) {
		return errors.New("auth config fail")
	}

	//请求第三方
	crd.requestFactory = third_user_auth_service.CreateRequestFactory(crd.authConfig)

	return nil
}

// Sync 同步数据
func (crd *ThirdUserAuthCmd) Sync(isNotice bool) {
	// 加锁失败：说明上一次还没执行完，直接返回
	if !crd.setLock() {
		variable.ZapLog.Info("Sync skipped: previous task not completed", zap.Bool("isNotice", isNotice))
		return
	}
	// 只有加锁成功后才 defer 解锁
	defer crd.delLock()
	variable.ZapLog.Info("Sync start", zap.Bool("isNotice", isNotice))

	//
	if err := crd.iniFact(); err != nil {
		variable.ZapLog.Info(err.Error())
		return
	}

	//初始化未完成
	if crd.authConfig.IsInit == enums.IsNot {
		crd.handleOrg()
		crd.handleUser()
	} else {
		//通知任务
		if isNotice {
			//有通知
			if crd.authConfig.NoticeSync > 0 {
				crd.handleOrg()
				crd.handleUser()
			} else {
				variable.ZapLog.Info("Not Sync", zap.Int("NoticeSync", crd.authConfig.NoticeSync))
				return
			}
		} else {
			crd.handleOrg()
			crd.handleUser()
		}
	}

	//

	variable.ZapLog.Info("Sync end", zap.Bool("isNotice", isNotice))
}

// handleOrg 处理组织相关

func (crd *ThirdUserAuthCmd) handleOrg() {
	variable.ZapLog.Info("handleOrg ===================================")

	//获取三方数据
	err, thirdList := crd.requestFactory.PullOrg()
	if err != nil {
		variable.ZapLog.Info(fmt.Sprintf("PullOrg Fail: %v", err))
	}
	if len(thirdList) == 0 {
		variable.ZapLog.Info("PullOrg Empty")
	}
	//fmt.Println("orgList--------------:")
	//fmt.Println(json_helper.Encode(thirdList))
	//fmt.Println("orgList--------------:")

	//格式化三方返回数据列表
	orgSaveList, deptSaveList, deptDeleteCodes := crd.formatThirdOrg(thirdList)
	fmt.Println("orgSaveList--------------:")
	fmt.Println(json_helper.Encode(orgSaveList))
	fmt.Println("deptSaveList--------------:")
	fmt.Println(json_helper.Encode(deptSaveList))
	fmt.Println("deptDeleteCodes--------------:")
	fmt.Println(json_helper.Encode(deptDeleteCodes))
	//处理企业
	if len(orgSaveList) > 0 {
		if err := organize_service.CreateKoalFactory().Sync(orgSaveList); err != nil {
			//同步序列
			crd.orgMaxSeq = 0

			variable.ZapLog.Info(fmt.Sprintf("org sync fail: %s", err.Error()))
			return
		}
	}

	//处理部门
	deptFact := department_service.CreateKoalFactory()
	if len(deptSaveList) > 0 {
		for _, model := range deptSaveList {
			if err := deptFact.Sync(model); err != nil {
				//同步序列
				crd.orgMaxSeq = 0

				variable.ZapLog.Info(fmt.Sprintf("dept id=%d sync fail: %v", model.Id, err))
				return
			}
		}
	}
	if len(deptDeleteCodes) > 0 {
		if err := deptFact.DeleteByCodes(deptDeleteCodes); err != nil {
			//同步序列
			crd.orgMaxSeq = 0

			variable.ZapLog.Info(fmt.Sprintf("dept [%s] delete fail: %v", json_helper.Encode(deptDeleteCodes), err))
			return
		}
	}

}

// formatThirdOrg 格式化三方返回数据列表
// 返回：
//   - orgSaveList: 企业（单位）要保存列表
//   - deptSaveList: 组织架构（部门）要保存列表, 已按照级别从高到底排序(1 > 2 >3)
//   - deptDeleteCodes: 组织架构（部门）要删除集合
func (crd *ThirdUserAuthCmd) formatThirdOrg(thirdList []*third_user_auth_service.ResponseOrg) (orgSaveList []*organization_models.OrganizationModel,
	deptSaveList []*department_models.DepartmentModel, deptDeleteCodes []string) {

	//三方数据
	thirdMaps := make(map[string]*third_user_auth_service.ResponseOrg, len(thirdList))
	for _, item := range thirdList {
		thirdMaps[item.OrgId] = item

		//同步序列
		crd.orgMaxSeq = string_helper.ToInt64(item.SyncSequence)
	}

	//组织架构（部门）最深级别
	deptMaxLevel := 1
	//组织架构（部门） 按级别处理
	deptSaveMap := make(map[int][]*department_models.DepartmentModel)

	// 格式化成需要的数据
	for _, item := range thirdList {
		createdAt := date_helper.MilliTimestampToDatetime(string_helper.ToInt64(item.CreateTime))
		//是企业
		if item.OrgType == enum_third_user_auth.OrgTypeOrg {
			tmpModel := organization_models.OrganizationModel{
				Code:              item.OrgId,
				Name:              item.OrgName,
				Type:              enum_organization.TypeOrg,
				Status:            enum_organization.StatusVerifying,
				IsClean:           enums.IsNot,
				IsDel:             enums.IsNot,
				Remark:            "用户统一认证",
				LoginPhoneCheck:   enums.IsNot,
				HuanshiSwitch:     enums.IsYes,
				TrusteeshipSwitch: enums.IsYes,
				Source:            enums.OperateSourcekoal,
				CreatedAt:         createdAt,
			}

			if item.Status == enum_third_user_auth.StatusDeleted {
				tmpModel.IsDel = enums.IsYes
				tmpModel.Status = enum_organization.StatusLocking
			} else if item.Status == enum_third_user_auth.StatusLocking {
				tmpModel.Status = enum_organization.StatusLocking
			} else if item.SyncStatus == enum_third_user_auth.SyncStatusRevoke {
				tmpModel.Status = enum_organization.StatusLocking
			}

			orgSaveList = append(orgSaveList, &tmpModel)
		} else if item.OrgType == enum_third_user_auth.OrgTypeDept {
			if item.Status == enum_third_user_auth.StatusNormal && item.SyncStatus == enum_third_user_auth.SyncStatusNormal {
				//所有父级
				allParents := crd.findAllParentOrgs(thirdMaps, item.OrgId)
				//完整路径
				var fullName []string
				if len(allParents.ParentPath) > 0 {
					//翻转路径
					slices.Reverse(allParents.ParentPath)
					for _, v := range allParents.ParentPath {
						fullName = append(fullName, v.OrgName)
					}
				}
				//父级
				parentId := string_helper.ToInt64(item.ParentId)
				if item.ParentId == allParents.NoParentOrgId {
					parentId = 0
				}
				//级别
				level := len(allParents.ParentPath)
				if level > deptMaxLevel {
					deptMaxLevel = level
				}

				tmpModel := department_models.DepartmentModel{
					OrgId:     string_helper.ToInt64(allParents.NoParentOrgId),
					ParentId:  parentId,
					Name:      item.OrgName,
					Code:      item.OrgId,
					FullCode:  "",
					FullName:  json_helper.Encode(fullName),
					Level:     level,
					Sort:      item.ShowOrder,
					CreatedAt: createdAt,
				}
				deptSaveMap[level] = append(deptSaveMap[level], &tmpModel)
			} else {
				deptDeleteCodes = append(deptDeleteCodes, item.OrgId)
			}
		}
	} //end thirdList

	// 已按照级别从高到底排序(1 > 2 >3)
	if len(deptSaveMap) > 0 {
		for lv := 1; lv <= deptMaxLevel; lv++ {
			if models, ok := deptSaveMap[lv]; ok && len(models) > 0 {
				deptSaveList = append(deptSaveList, models...)
			}
		}
	}

	return orgSaveList, deptSaveList, deptDeleteCodes
}

// orgPathItem 路径上的一个节点
type orgPathItem struct {
	OrgId   string `json:"orgId"`
	OrgName string `json:"orgName"`
}

// findAllParentOrg 向上查找的返回结果
type findAllParentOrg struct {
	// 没有父级的 orgId（最终追到的根节点，只有一个）
	NoParentOrgId string `json:"noParentOrgId"`
	// 是否有根节点（防止 orgId 本身就不存在的情况）
	HasNoParent bool `json:"hasNoParent"`
	// 有父级的节点路径：从自己 → 最顶层
	ParentPath []orgPathItem `json:"parentPath"`
}

// findAllParentOrgs 从 orgId 开始，向上查找所有父级
// 返回：
//   - NoParentOrgId: 没有父级的 orgId（根节点，只有一个）
//   - ParentPath:    有父级的节点路径（从自己 → 最顶层），包含 orgId 和 orgName
func (crd *ThirdUserAuthCmd) findAllParentOrgs(orgMaps map[string]*third_user_auth_service.ResponseOrg, orgId string) findAllParentOrg {
	result := findAllParentOrg{
		ParentPath: make([]orgPathItem, 0),
	}

	currentId := orgId
	visited := make(map[string]bool)

	for currentId != "" && currentId != "0" {
		if visited[currentId] {
			break // 出现环
		}
		visited[currentId] = true

		org, ok := orgMaps[currentId]
		if !ok {
			// 当前节点不在 map 里，无法继续向上
			// 把它当作"没有父级"处理
			result.NoParentOrgId = currentId
			result.HasNoParent = true
			break
		}

		parentId := org.ParentId

		// 没有父级：记录唯一的根节点
		if parentId == "" || parentId == "0" {
			result.NoParentOrgId = org.OrgId
			result.HasNoParent = true
			break
		}

		// 有父级：记录 orgId + orgName
		result.ParentPath = append(result.ParentPath, orgPathItem{
			OrgId:   org.OrgId,
			OrgName: org.OrgName,
		})

		currentId = parentId
	}
	return result
}

// 处理用户
func (crd *ThirdUserAuthCmd) handleUser() {
	return
	variable.ZapLog.Info("handleUser ===================================")

	//获取三方数据
	err, thirdList := crd.requestFactory.PullUser()
	if err != nil {
		variable.ZapLog.Info(fmt.Sprintf("handleUser Fail: %v", err))
	}
	if len(thirdList) == 0 {
		variable.ZapLog.Info("handleUser Empty")
	}
	fmt.Println("userList--------------:")
	fmt.Println(json_helper.Encode(thirdList))
	fmt.Println("userList--------------:")
}
