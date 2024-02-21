package api

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"gowebserver/app/common/api"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"gowebserver/app/system/service"
	"strings"
)

var SysBlackListApi = new(sysBlackListApi)

type sysBlackListApi struct {
	api.BaseController
}

func (a *sysBlackListApi) Init(r *ghttp.Request) {
	a.Module = "系统黑白名单"
	r.SetCtxVar(global.Module, a.Module)
}

// GetList
// @summary 分页获取系统访问黑白名单配置
// @tags 	系统访问黑白名单
// @Param   authorization header string true  "Bearer Token"
// @Param   type          query  int    true  "黑白名单类型: 1黑名单 2白名单"
// @Param   keyword       query  string false "搜索查询关键字"
// @Param   pageNum       query  int    true  "当前页码"
// @Param   pageSize      query  int    true  "每页展示条数"
// @Produce json
// @Success 200 {object} response.Response{data=define.SysBlackListInfoRes} "配置信息列表"
// @Failure 500 {object} response.Response "请求参数错误"
// @Router 	/system/blacklist/getList [GET]
func (a *sysBlackListApi) GetList(r *ghttp.Request) {
	var req *define.SysBlackListSelectReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	//分页查询配置列表
	pageInfo, cfgList, err := service.InterfaceSysBlackList().GetListSearch(r.Context(), req)
	if err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}

	//请求返回信息
	var resp define.SysBlackListInfoRes
	if err = gconv.Struct(cfgList, &resp.List); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInternalError, err))
	}
	resp.Total = pageInfo.Total
	resp.Page = pageInfo.PageNum
	resp.Size = pageInfo.PageSize
	a.RespJsonExit(r, nil, resp)
}

// Add
// @summary 新增黑白名单配置
// @tags 	系统访问黑白名单
// @Router 	/system/blacklist/add [POST]
func (a *sysBlackListApi) Add(r *ghttp.Request) {
	var req *define.SysBlackListAddReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}

	//配置信息
	var cfgInfo model.SysBlackList
	_ = gconv.Struct(req, &cfgInfo)
	cfgInfo.CreateBy = service.Context.GetUserName(r.GetCtx())
	cfgInfo.CreateTime = gtime.Now()
	if err := service.InterfaceSysBlackList().Add(r.GetCtx(), &cfgInfo); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}
	a.RespJsonExit(r, nil, cfgInfo)
}

// Delete
// @summary 删除黑白名单配置
// @tags 	系统访问黑白名单
// @Router 	/system/blacklist/delete [DELETE]
func (a *sysBlackListApi) Delete(r *ghttp.Request) {
	strDleCfgIds := r.GetString("ids")
	params := strings.Split(strDleCfgIds, ",")
	delCfgIds := make([]int, 0, len(params))
	for _, cfgId := range params {
		delCfgIds = append(delCfgIds, gconv.Int(cfgId))
	}
	if err := service.InterfaceSysBlackList().Del(r.GetCtx(), delCfgIds); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}
	a.RespJsonExit(r, nil)
}

// GetFilterMode
// @summary 查询黑白名单模式
// @tags 	系统访问黑白名单
// @Router 	/system/blacklist/mode [GET]
func (a *sysBlackListApi) GetFilterMode(r *ghttp.Request) {
	a.RespJsonExit(r, nil, service.InterfaceSysBlackList().GetMode(r.GetCtx()))
}

// SetFilterMode
// @summary 修改黑白名单模式
// @tags 	系统访问黑白名单
// @Router 	/system/blacklist/mode [PUT]
func (a *sysBlackListApi) SetFilterMode(r *ghttp.Request) {
	mode := r.GetInt("type")
	if err := service.InterfaceSysBlackList().SetMode(r.GetCtx(), mode); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}
	a.RespJsonExit(r, nil)
}
