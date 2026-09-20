package third_user_auth_service

import (
	"fmt"
	"gorm.io/gorm"
	"goskeleton/app/global/enums"
	"goskeleton/app/model/platform/third_user_auth_models"
	"goskeleton/app/utils/helpers/date_helper"
)

func CreateConfigFactory() *Config {
	return &Config{masterModel: third_user_auth_models.CreateAuthFactory()}
}

type Config struct {
	masterModel *third_user_auth_models.AuthModel
}

// GetNewestItem 获取最新认证配置
func (crd *Config) GetNewestItem() (model *third_user_auth_models.AuthModel) {
	return crd.masterModel.GetNewestItem()
}

// ModifyAuthById 修改认证部分
func (crd *Config) ModifyAuthById(param *third_user_auth_models.AuthModel) error {
	if model := crd.masterModel.GetItemById(param.Id); model == nil || model.Id == 0 {
		return fmt.Errorf("data not exist")
	}

	updates := map[string]interface{}{
		"addr":       param.Addr,
		"app_id":     param.AppId,
		"app_key":    param.AppKey,
		"app_secret": param.AppSecret,
	}
	if isOk, _ := crd.masterModel.UpdateById(param.Id, updates); !isOk {
		return fmt.Errorf("save fail")
	}
	return nil
}

// ModifyNoticeById 修改通知部分
func (crd *Config) ModifyNoticeById(noticeSync int) error {
	model := crd.GetNewestItem()
	if model == nil || model.Id == 0 {
		return fmt.Errorf("data not exist")
	}

	updates := map[string]interface{}{
		"notice_sync": gorm.Expr(fmt.Sprintf("notice_sync+%d", noticeSync)),
		"notice_at":   date_helper.CurDatetime(),
	}
	if isOk, _ := crd.masterModel.UpdateById(model.Id, updates); !isOk {
		return fmt.Errorf("save fail")
	}
	return nil
}

// ModifySyncById 修改同步部分
func (crd *Config) ModifySyncById(param *third_user_auth_models.AuthModel) error {
	if model := crd.masterModel.GetItemById(param.Id); model == nil || model.Id == 0 {
		return fmt.Errorf("data not exist")
	}

	now := date_helper.CurDatetime()
	updates := map[string]interface{}{
		"is_init":     enums.IsYes,
		"notice_sync": 0,
		"org_seq":     param.OrgSeq,
		"org_at":      now,
		"user_seq":    param.OrgSeq,
		"user_at":     now,
	}
	if isOk, _ := crd.masterModel.UpdateById(param.Id, updates); !isOk {
		return fmt.Errorf("save fail")
	}
	return nil
}
