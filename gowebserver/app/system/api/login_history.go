// =================================================================================
// 系统登录历史记录API控制器
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

var LoginHistoryApi = new(loginHistoryApi)

type loginHistoryApi struct {
	api.BaseController
}

func (a *loginHistoryApi) Init(r *ghttp.Request) {
	a.Module = "登录记录管理"
	r.SetCtxVar(global.Module, a.Module)
}

// Get 登录记录列表分页数据
// @summary 登录记录列表分页查询
// @tags 	登录记录管理
// @Param   Authorization header string true "Bearer Token"
// @Param   loginName      query string false "用户名称"
// @Param   ipaddr         query string false "登录IP地址"
// @Param   status         query string false "在线状态"
// @Param   beginTime      query string false "时间范围-起始"
// @Param   endTime        query string false "时间范围-结束"
// @Param   pageNum        query int    true  "当前页码"
// @Param   pageSize       query int    true  "每页展示条数"
// @Param   orderByColumn  query string false "排序字段"
// @Param   isAsc          query string false "排序方式"
// @Produce json
// @Success 200 {object} response.Response{data=define.LoginHistoryServiceList} "登录历史信息列表"
// @Failure 500 {object} response.Response "请求参数错误"
// @Router 	 /monitor/loginHistory [GET]
func (a *loginHistoryApi) Get(r *ghttp.Request) {
	var req *define.LoginHistoryApiSelectPageReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	result := service.LoginHistory.GetList(r.GetCtx(), req)
	a.RespJsonExit(r, nil, result)
}

// Delete 删除登录历史数据
// @summary 删除登录历史数据
// @tags 	登录记录管理
// @Param   Authorization header string true "Bearer Token"
// @Param   ids query string true "要删除的历史记录ID, 多个ID , 分割. 例: id1,id2,id3"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	 /monitor/loginHistory [DELETE]
func (a *loginHistoryApi) Delete(r *ghttp.Request) {
	var req *define.LoginHistoryApiDeleteReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	rs := service.LoginHistory.Delete(r.GetCtx(), req.Ids)
	if rs > 0 {
		a.RespJsonExit(r, nil, rs)
	} else {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonOperationFailed))
	}
}

// Clean 清空登录历史记录
// @summary 清空登录历史记录
// @tags 	登录记录管理
// @Param   Authorization header string true "Bearer Token"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	 /monitor/loginHistory/clean [DELETE]
func (a *loginHistoryApi) Clean(r *ghttp.Request) {
	r.SetCtxVar(global.ResponseBusinessType, global.BusinessClean)
	rs := service.LoginHistory.Clean(r.GetCtx())
	if rs > 0 {
		a.RespJsonExit(r, nil, rs)
	} else {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonOperationFailed))
	}
}

// Export 导出登录历史记录
// @summary 导出登录历史记录
// @tags 	登录记录管理
// @Param   Authorization header string true "Bearer Token"
// @Param   loginName      query string false "用户名称"
// @Param   ipaddr         query string false "登录IP地址"
// @Param   status         query string false "在线状态"
// @Param   beginTime      query string false "时间范围-起始"
// @Param   endTime        query string false "时间范围-结束"
// @Param   pageNum        query int    true  "当前页码"
// @Param   pageSize       query int    true  "每页展示条数"
// @Param   orderByColumn  query string false "排序字段"
// @Param   isAsc          query string false "排序方式"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	 /monitor/loginHistory/export [GET]
func (a *loginHistoryApi) Export(r *ghttp.Request) {
	r.SetCtxVar(global.Module, a.Module)
	r.SetCtxVar(global.ResponseBusinessType, global.BusinessExport)

	var req *define.LoginHistoryApiSelectPageReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	url, err := service.LoginHistory.Export(r.GetCtx(), req)
	if err != nil {
		a.RespJsonExit(r, err)
	} else {
		a.RespJsonExit(r, nil, url)
	}
}
