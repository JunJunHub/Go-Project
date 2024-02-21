package api

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"gowebserver/app/common/api"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"gowebserver/app/system/service"
)

var SysParamApi = new(sysParamApi)

type sysParamApi struct {
	api.BaseController
}

func (a *sysParamApi) Init(r *ghttp.Request) {
	a.Module = "系统参数"
	r.SetCtxVar(global.Module, a.Module)
}

// GetSysParamCfgList
// @summary 查询系统参数配置信息
// @tags 	系统参数
// @Param   Authorization header string true "Bearer Token"
// @Param   keyword       query string false "关键字查询"
// @Param   orderByColumn query string false "排序字段"
// @Param   isAsc         query string false "排序方式"
// @Param   pageNum       query int    true  "当前页码"
// @Param   pageSize      query int    true  "每页展示条数"
// @Produce json
// @Success 200 {object} response.Response{data=define.SysParamSelectRes} "系统参数配置信息"
// @Failure 500 {object} response.Response "请求参数错误"
// @Router 	/param/getList [GET]
func (a *sysParamApi) GetSysParamCfgList(r *ghttp.Request) {
	var req *define.SysParamSelectReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	//分页查询配置列表
	pageInfo, cfgList, err := service.InterfaceSysParamMgr().GetListSearch(r.Context(), req)
	if err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}
	//请求返回信息
	var resp define.SysParamSelectRes
	if err = gconv.Struct(cfgList, &resp.List); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInternalError, err))
	}
	resp.Total = pageInfo.Total
	resp.Page = pageInfo.PageNum
	resp.Size = pageInfo.PageSize
	a.RespJsonExit(r, nil, resp)
}

// EditSysParamCfg
// @summary 编辑系统参数配置
// @tags 	系统参数
// @Param   Authorization header string true "Bearer Token"
// @Produce json
// @Success 200 {object} response.Response "返回结果"
// @Router 	/system/param/edit [PUT]
func (a *sysParamApi) EditSysParamCfg(r *ghttp.Request) {
	var req *define.SysParamEditReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	if err := service.InterfaceSysParamMgr().SetSysParamCfg(r.GetCtx(), req.ParamKey, req.ParamValue); err != nil {
		a.RespJsonExit(r, err)
	}
	a.RespJsonExit(r, nil)
}

// GetLogStorageStrategy
// @summary 查询日志存储策略
// @tags 	系统参数
// @Param   Authorization header string true "Bearer Token"
// @Param   logType       query int    true  "当前页码"
// @Produce json
// @Success 200 {object} response.Response "返回结果"
// @Router 	/monitor/log/storageStrategy [PUT]
func (a *sysParamApi) GetLogStorageStrategy(r *ghttp.Request) {
	logType := r.GetInt("logType")
	logStorageCfg, err := service.InterfaceSysParamMgr().GetLogStorageStrategy(r.GetCtx(), define.ELogType(logType))
	if err != nil {
		a.RespJsonExit(r, err)
	}
	a.RespJsonExit(r, nil, logStorageCfg)
}

// SetLogStorageStrategy
// @summary 设置日志存储策略
// @tags 	系统参数
// @Param   Authorization header string true "Bearer Token"
// @Produce json
// @Success 200 {object} response.Response "返回结果"
// @Router 	/monitor/log/storageStrategy [PUT]
func (a *sysParamApi) SetLogStorageStrategy(r *ghttp.Request) {
	var req define.SysParamLogStorageStrategy
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}

	var err error
	var cfgInfo *model.SysParam
	var cfgKey string
	switch req.LogType {
	case 1: //操作日志
		cfgKey = global.SysParamKey_OperLogStorageStrategy
		cfgInfo, err = service.InterfaceSysParamMgr().GetSysParamCfg(r.GetCtx(), global.SysParamKey_OperLogStorageStrategy)
	case 2: //登录日志
		cfgKey = global.SysParamKey_LoginHistoryLogStorageStrategy
		cfgInfo, err = service.InterfaceSysParamMgr().GetSysParamCfg(r.GetCtx(), global.SysParamKey_LoginHistoryLogStorageStrategy)
	case 3: //告警日志
		cfgKey = global.SysParamKey_AlarmLogStorageStrategy
		cfgInfo, err = service.InterfaceSysParamMgr().GetSysParamCfg(r.GetCtx(), global.SysParamKey_AlarmLogStorageStrategy)
	case 4: //设备上下线日志
		cfgKey = global.SysParamKey_DeviceOnlineLogStorageStrategy
		cfgInfo, err = service.InterfaceSysParamMgr().GetSysParamCfg(r.GetCtx(), global.SysParamKey_DeviceOnlineLogStorageStrategy)
	default:
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if cfgInfo == nil {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}

	//保存配置
	cfgValue, _ := gjson.Encode(req)
	if err = service.InterfaceSysParamMgr().SetSysParamCfg(r.GetCtx(), cfgKey, string(cfgValue)); err != nil {
		a.RespJsonExit(r, err)
	}
	a.RespJsonExit(r, nil)
}
