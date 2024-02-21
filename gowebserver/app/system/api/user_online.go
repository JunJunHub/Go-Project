// =================================================================================
// 在线用户管理API控制器
// =================================================================================

package api

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"gowebserver/app/common/api"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/system/define"
	"gowebserver/app/system/service"
)

var OnlineUserUserApi = new(onlineUserApi)

type onlineUserApi struct {
	api.BaseController
}

func (a *onlineUserApi) Init(r *ghttp.Request) {
	a.Module = "在线用户管理"
	r.SetCtxVar(global.Module, a.Module)
}

// Get 在线用户列表分页查询
// @summary 在线用户列表分页查询
// @tags 	在线用户管理
// @Param   Authorization header string true "Bearer Token"
// @Param   loginLocation  query string false "登录地点"
// @Param   loginName      query string false "用户名称"
// @Param   ipaddr         query string false "登录IP地址"
// @Param   browser        query string false "浏览器类型"
// @Param   os             query string false "操作系统"
// @Param   status         query string false "在线状态"
// @Param   beginTime      query string false "时间范围-起始"
// @Param   endTime        query string false "时间范围-结束"
// @Param   pageNum        query int    true  "当前页码"
// @Param   pageSize       query int    true  "每页展示条数"
// @Param   orderByColumn  query string false "排序字段"
// @Param   isAsc          query string false "排序方式"
// @Produce json
// @Success 200 {object} response.Response{data=define.OnlineUserServiceList} "在线用户信息列表"
// @Failure 500 {object} response.Response "请求参数错误"
// @Router 	 /monitor/online [GET]
func (a *onlineUserApi) Get(r *ghttp.Request) {
	var req *define.OnlineUserApiSelectPageReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}
	result := service.OnlineUser.GetList(req)
	a.RespJsonExit(r, nil, result)
}

// ForceLogout 用户强退
// @summary 强制用户登出
// @tags 	在线用户管理
// @Param   Authorization header string true "Bearer Token"
// @Param   tokenStr query string true "强退用户token"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/system/user [DELETE]
func (a *onlineUserApi) ForceLogout(r *ghttp.Request) {
	r.SetCtxVar(global.ResponseBusinessType, global.BusinessForceLogout)
	tokenStr := r.GetString("tokenStr")
	if tokenStr == "" {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}
	Auth.DelGTokenCacheUser(tokenStr)
	err := Auth.UserLogOutHandle(tokenStr)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	a.RespJsonExit(r, nil)
}
