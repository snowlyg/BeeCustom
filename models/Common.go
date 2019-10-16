package models

import "BeeCustom/enums"

// JsonResult 用于返回ajax请求的基类
type JsonResult struct {
	Status enums.JsonResultCode `json:"status"`
	Msg    string               `json:"msg"`
	Obj    interface{}          `json:"obj"`
}

// BaseQueryParam 用于查询的类
type BaseQueryParam struct {
	Sort   string `json:"sort"`
	Order  string `json:"order"`
	Offset int64  `json:"offset"`
	Limit  int64  `json:"limit"`
}
