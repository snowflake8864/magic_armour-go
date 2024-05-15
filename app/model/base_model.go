package model

import (
	"fmt"
	"gorm.io/gorm"
	"goskeleton/app/global/my_errors"
	"goskeleton/app/global/variable"
	"goskeleton/app/utils/helpers/string_helper"
	"strings"
)

type BaseModel struct {
	*gorm.DB `gorm:"-" json:"-"`
	//Id       int64 `gorm:"primaryKey" json:"id"`
	//CreatedAt string `json:"CreatedAt"` //日期时间字段统一设置为字符串即可
	//UpdatedAt string `json:"updated_at"`
	//DeletedAt gorm.DeletedAt `json:"deleted_at"`   // 如果开发者需要使用软删除功能，打开本行注释掉的代码即可，同时需要在数据库的所有表增加字段deleted_at 类型为 datetime
}

func UseDbConn(selectDb string) *gorm.DB {
	var db *gorm.DB
	selectDb = strings.Trim(selectDb, " ")
	if selectDb == "" {
		selectDb = variable.ConfigGormv2Yml.GetString("Gormv2.DefaultSelectDb")
	}
	switch selectDb {
	case variable.DbNamePlatform:
		if variable.GormDbPlatform == nil {
			variable.ZapLog.Fatal(fmt.Sprintf(my_errors.ErrorsGormNotInitGlobalPointer, selectDb, selectDb))
		}
		db = variable.GormDbPlatform
	case variable.DbNameAlarm:
		if variable.GormDbAlarm == nil {
			variable.ZapLog.Fatal(fmt.Sprintf(my_errors.ErrorsGormNotInitGlobalPointer, selectDb, selectDb))
		}
		db = variable.GormDbAlarm
	case variable.DbNameAgent:
		if variable.GormDbAgent == nil {
			variable.ZapLog.Fatal(fmt.Sprintf(my_errors.ErrorsGormNotInitGlobalPointer, selectDb, selectDb))
		}
		db = variable.GormDbAgent
	default:
		variable.ZapLog.Error(my_errors.ErrorsDbDriverNotExists + selectDb)
	}
	return db
}

func (m *BaseModel) TrimSql(sql string) string {
	sql = strings.ReplaceAll(sql, "	", " ")
	sql = strings.ReplaceAll(sql, "\t", "")
	sql = strings.ReplaceAll(sql, "\n", "")
	//sql = strings.ReplaceAll(sql, "((", "(")
	//sql = strings.ReplaceAll(sql, "))", ")")
	//fmt.Println(sql)
	return sql
}

func (m *BaseModel) GetColumn(column string) string {
	if string_helper.IsEmpty(column) {
		column = "*"
	}
	return column
}

// WhereBuild 封装map 查询条件
func (m *BaseModel) WhereBuild(conditions map[string]interface{}) (whereSQL string, vals []interface{}, err error) {
	for conditionKey, conditionVal := range conditions {
		fieldAndWhere := strings.Split(conditionKey, " ")
		if len(fieldAndWhere) > 2 {
			return "", nil, fmt.Errorf("Error in query condition: %s. ", conditionKey)
		}
		if whereSQL != "" {
			whereSQL += " AND "
		}
		switch len(fieldAndWhere) {
		case 1:
			whereSQL += fmt.Sprint(fieldAndWhere[0], "=?")
			vals = append(vals, conditionVal)
			break
		case 2:
			field := fieldAndWhere[0]
			switch fieldAndWhere[1] {
			case "=":
				whereSQL += fmt.Sprint(field, "=?")
				vals = append(vals, conditionVal)
				break
			case ">":
				whereSQL += fmt.Sprint(field, ">?")
				vals = append(vals, conditionVal)
				break
			case ">=":
				whereSQL += fmt.Sprint(field, ">=?")
				vals = append(vals, conditionVal)
				break
			case "<":
				whereSQL += fmt.Sprint(field, "<?")
				vals = append(vals, conditionVal)
				break
			case "<=":
				whereSQL += fmt.Sprint(field, "<=?")
				vals = append(vals, conditionVal)
				break
			case "!=":
				whereSQL += fmt.Sprint(field, "!=?")
				vals = append(vals, conditionVal)
				break
			case "<>":
				whereSQL += fmt.Sprint(field, "!=?")
				vals = append(vals, conditionVal)
				break
			case "in":
				whereSQL += fmt.Sprint(field, " IN ?")
				vals = append(vals, conditionVal)
				break
			case "like":
				whereSQL += fmt.Sprint(field, " LIKE ?")
				vals = append(vals, conditionVal)
				break
			case "IsNull":
				whereSQL += fmt.Sprint(field, " IS NULL")
				break
			case "IsNotNull":
				whereSQL += fmt.Sprint(field, " IS NOT NULL")
				break
			}
			break
		}
	}
	return whereSQL, vals, err
}

// 在 ginskeleton项目中如果在业务 model 设置了回调函数，请看以下说明
// 注意：gorm 的自动回调函数（BeforeCreate、BeforeUpdate 等），不是由本项目的 Create ... 函数先初始化然后调用的，而是gorm自动直接调用的，
// 所以 接收器 b 的所有参数都是没有赋值的,因此这里需要给 b.DB 赋予回调的 gormDb
// baseModel 的代码执行顺序晚于其他业务 model 的回调函数，如果回调函数名称相同，会被普通业务model的同名回调函数覆盖
// gorm 支持的自动回调函数清单：https://github.com/go-gorm/gorm/blob/master/callbacks/interfaces.go

//func (b *BaseModel) BeforeCreate(gormDB *gorm.DB) error {
//	 第一步必须反向将 gormDB 赋值给 b.DB
//	b.DB = gormDB
//	 后续的代码就可以像普通业务 model 一样操作，
//	  b.Exec(sql,参数1，参数2，...)
//	 b.Raw(sql,参数1，参数2，...)
//	return nil
//}

//  BeforeUpdate、BeforeSave 函数都会因为 更新类的操作而被触发
//  如果baseModel 和 普通业务 model 都想使用回调函数，那么请设置不同的回调函数名，例如：这里设置 BeforeUpdate、普通业务model 设置 BeforeSave 即可
//func (b *BaseModel) BeforeUpdate(gormDB *gorm.DB) error {
//	 第一步必须反向将 gormDB 赋值给 b.DB
//	b.DB = gormDB
//	 后续的代码就可以像普通业务 model 一样操作，
//	  b.Exec(sql,参数1，参数2，...)
//	 b.Raw(sql,参数1，参数2，...)
//	return nil
//}
