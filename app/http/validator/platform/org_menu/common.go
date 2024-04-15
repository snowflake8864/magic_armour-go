package org_menu

type Common struct {
	Type   int    `form:"type" json:"type" binding:"required,max=10000"`
	Name   string `form:"name" json:"name" binding:"required,max=100"`
	Code   string `form:"code" json:"code" binding:"required,max=100"`
	Sort   *int   `form:"sort" json:"sort" binding:"required,min=0,max=9999"`
	IsShow *int   `form:"isShow" json:"isShow" binding:"required,min=0,max=1"`
	IsRoot *int   `form:"isRoot" json:"isRoot" binding:"required,min=0,max=1"`
	Remark string `form:"remark" json:"remark" binding:"omitempty,max=100"`
}
