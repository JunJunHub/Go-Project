// ==========================================================================
// 基础控制器, 抽象每个功能模块的API层公共代码
// ==========================================================================

package api

import (
	"github.com/gogf/gf/net/ghttp"
	"gowebserver/app/common/utils/response"
)

type BaseController struct {
	Module string // 功能模块描述
}

//RespJson
//@Summary 返回json格式
//@Param1  r *ghttp.Request  "请求信息"
//@Param2  err error  "错误信息,成功时为nil"
//@Param3  data ...interface{} "返回具体数据"
//@Return  nil
func (c *BaseController) RespJson(r *ghttp.Request, err error, data ...interface{}) {
	response.RJson(r, err, data...)
}

//RespJsonExit
//@Summary 返回json格式并结束请求处理流程
//@param1  r *ghttp.Request "请求信息"
//@param2  err error  "错误信息,成功时为nil"
//@param3  data ...interface{} "返回具体数据"
//@return  nil
func (c *BaseController) RespJsonExit(r *ghttp.Request, err error, data ...interface{}) {
	response.JsonExit(r, err, data...)
}
