package lego

import (
	"github.com/bytedance/go-tagexpr/validator"
)

/*
 * 表单验证相关
 * 使用例子：
 ------------------------------------------------------------------
	var vd = validator.New("vd")
	type ReqParams struct {
		Id      string // 记录id，为空则新建，否则为编辑
		Name    string `vd:"len($)>0"` // 品牌名称
		Logo    string // 品牌Logo
		Brief   string // 品牌简介
		SortIdx int    `vd:"$>0"` // 显示排序值
	}
	req := ReqParams{
		Id:      string(ctx.Peek("id")),
		Name:    string(ctx.Peek("name")),
		Logo:    string(ctx.Peek("logo")),
		Brief:   string(ctx.Peek("brief")),
		SortIdx: utils.BytesToInt(ctx.Peek("sort_idx")),
	}

	if err := vd.Validate(req); err != nil {
		code = 100
		msg = err.Error()
		break
	}
  ------------------------------------------------------------------
*/

/*
 * 提供给外部的表单验证接口
 */
 func Validate(form interface{}) error {
	 tv := validator.New("vd")
 	return tv.Validate(form)
 }