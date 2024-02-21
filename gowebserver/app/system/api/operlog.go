// =================================================================================
// 系统操作日志控制器(日志查询、导出、删除...)
// =================================================================================

package api

import (
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"gowebserver/app/common/api"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/common/utils/response"
	"gowebserver/app/system/define"
	"gowebserver/app/system/service"
)

// OperlogApi 操作日志API对象
var OperlogApi = new(operlogApi)

type operlogApi struct {
	api.BaseController
}

func (a *operlogApi) Init(r *ghttp.Request) {
	a.Module = "系统操作日志管理"
}

// Get 操作日志列表分页数据
// @summary 操作日志分页获取
// @tags 	系统操作日志
// @Param   Authorization header string true "Bearer Token"
// @Param   title         query string false "系统模块(模糊匹配)"
// @Param   operName      query string false "操作人员(模糊匹配)"
// @Param   businessType  query string false "操作类型"
// @Param   status        query string false "操作结果"
// @Param   beginTime     query string false "用户创建时间范围-起始"
// @Param   endTime       query string false "用户创建时间范围-结束"
// @Param   pageNum       query int    true  "当前页码"
// @Param   pageSize      query int    true  "每页展示条数"
// @Param   orderByColumn query string false "排序字段"
// @Param   isAsc         query string false "排序方式"
// @Produce json
// @Success 200 {object} response.Response{data=define.OperlogServiceList} "操作日志列表"
// @Failure 500 {object} response.Response "请求参数错误"
// @Router 	 /monitor/operlog [GET]
func (a *operlogApi) Get(r *ghttp.Request) {
	var req *define.OperlogApiSelectPageReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}
	result := service.OperLog.GetList(r.GetCtx(), req)
	a.RespJsonExit(r, nil, result)
}

// Delete 删除操作日志
// @summary 删除操作日志
// @tags 	系统操作日志
// @Param   Authorization header string true "Bearer Token"
// @Param   ids query string true "要删除的日志记录ID, 多个ID , 分割. 例: id1,id2,id3"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	 /monitor/operlog [DELETE]
func (a *operlogApi) Delete(r *ghttp.Request) {
	var req *define.OperlogApiDeleteReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}
	rs := service.OperLog.Delete(r.GetCtx(), req.Ids)
	a.RespJsonExit(r, nil, rs)
}

// Clean 清空操作日志记录
// @summary 清空操作日志记录
// @tags 	系统操作日志
// @Param   Authorization header string true "Bearer Token"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	 /monitor/operlog/clean [DELETE]
func (a *operlogApi) Clean(r *ghttp.Request) {
	r.SetCtxVar(global.ResponseBusinessType, global.BusinessClean)
	rs := service.OperLog.Clean(r.GetCtx())
	a.RespJsonExit(r, nil, rs)
}

// Export 导出操作日志
// @summary 操作日志导出
// @tags 	系统操作日志
// @Param   Authorization header string true "Bearer Token"
// @Param   title         query string false "系统模块(模糊匹配)"
// @Param   operName      query string false "操作人员(模糊匹配)"
// @Param   businessType  query string false "操作类型"
// @Param   status        query string false "操作结果"
// @Param   beginTime     query string false "用户创建时间范围-起始"
// @Param   endTime       query string false "用户创建时间范围-结束"
// @Param   pageNum       query int    true  "当前页码"
// @Param   pageSize      query int    true  "每页展示条数"
// @Param   orderByColumn query string false "排序字段"
// @Param   isAsc         query string false "排序方式"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	 /monitor/operlog/export [GET]
func (a *operlogApi) Export(r *ghttp.Request) {
	r.SetCtxVar(global.ResponseBusinessType, global.BusinessExport)
	var req *define.OperlogApiSelectPageReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}
	url, err := service.OperLog.Export(r.GetCtx(), req)
	if err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	} else {
		a.RespJsonExit(r, nil, url)
	}
}

// OperationLog 操作日志记录
func (a *operlogApi) OperationLog(r *ghttp.Request) {
	//业务执行结束后,记录操作日志
	r.Middleware.Next()

	//如果开启调试模式,记录Get请求详细信息到日志文件
	if r.Method == "GET" && g.Cfg().GetBool("debug.enable") {
		optUser := service.Context.GetUser(r.Context())
		reqParam, err := gjson.Encode(r.GetRequestMap())
		if err != nil {
			g.Log().Error(err)
		}
		if optUser != nil {
			g.Log(global.LoggerHttpAccessRecord).Infof("\nUser： [userId:%d loginName:%s] \nURI：%s %s \nHeader：%v \nParam：%s \n",
				optUser.UserId, optUser.LoginName,
				r.Method, r.Request.RequestURI,
				r.Header,
				string(reqParam),
			)
		} else {
			g.Log(global.LoggerHttpAccessRecord).Infof("\nURI：%s %s \nHeader：%v \nParam：%s \n",
				r.Method, r.Request.RequestURI,
				r.Header,
				string(reqParam),
			)
		}
	}

	if r.Method == "GET" {
		//GET请求不记录操作日志到数据库
		return
	}

	//记录操作日志到数据库中
	module := r.GetCtxVar(global.Module).String()
	businessType := global.BusinessType(r.GetCtxVar(global.ResponseBusinessType).Int())
	resp := new(response.Response)
	if r.Response.BufferString() == "" {
		g.Log().Error("response info is null!")
		return
	}

	if err := gconv.Struct(r.Response.BufferString(), resp); err != nil {
		g.Log().Error(err, "; Resp:", r.Response.BufferString())
		return
	}

	//存操作日志到数据库
	var paramJson string
	param := r.GetMap()
	paramByte, err := gjson.Encode(param)
	if err == nil {
		paramJson = string(paramByte)
	}
	go service.OperLog.Create(r, module, paramJson, resp, businessType)
}
