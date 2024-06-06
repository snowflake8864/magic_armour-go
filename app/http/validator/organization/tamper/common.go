package tamper

import (
	"errors"
	"goskeleton/app/utils/helpers/string_helper"
)

type Dir struct {
	Dir string `form:"dir" json:"dir" binding:"required,max=512"`
}

type PlatformType struct {
	PlatformType int `form:"platformType" json:"platformType" binding:"required,oneof=1 2"`
}

type Id struct {
	Id int `form:"id" json:"id" binding:"required"`
}

type IsExtend struct {
	IsExtend *int `form:"isExtend" json:"isExtend" binding:"required,oneof=0 1"`
}

type Remark struct {
	Remark string `form:"remark" json:"remark" binding:"omitempty,max=128"`
}
type IncludeSuffix struct {
	IncludeSuffix string `form:"includeSuffix" json:"includeSuffix" binding:"omitempty,max=256"`
}

type ExcludeSuffix struct {
	ExcludeSuffix string `form:"excludeSuffix" json:"excludeSuffix" binding:"omitempty,max=256"`
}

type ExcludeFolder struct {
	ExcludeFolder string `form:"excludeFolder" json:"excludeFolder" binding:"omitempty,max=512"`
}

type Processes struct {
	Processes string `form:"processes" json:"processes" binding:"omitempty,min=32"`
}

type ProtectRw struct {
	ProtectRw string `form:"protectRw" json:"protectRw" binding:"omitempty,min=3,max=32"`
}

type DeptId struct {
	DeptId int `form:"deptId" json:"deptId"`
}

func checkSuffix(include, exclude string) error {
	if !string_helper.IsEmpty(include) && !string_helper.IsEmpty(exclude) {
		return errors.New("包含文件类型和排除文件类型不能同时存在")
	}

	return nil
}
