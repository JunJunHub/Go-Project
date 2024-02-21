// =================================================================================
// 中间件业务处理, 供各模块使用, 各模块自己决定功能接口需要哪些中间处理过程(路由注册)
// =================================================================================

package middleware

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/guid"
	userAgent "github.com/mssola/user_agent"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/common/service/appstate"
	"gowebserver/app/common/utils"
	"gowebserver/app/common/utils/response"
	apiSys "gowebserver/app/system/api"
	"gowebserver/app/system/define"
	"gowebserver/app/system/service"
	"io/ioutil"
	"strings"
	"time"
)

// CheckMpuApsStatus 服务状态校验
func CheckMpuApsStatus(r *ghttp.Request) {
	if !appstate.CheckHttpServerIsOk() {
		response.JsonExit(r, gerror.NewCode(errcode.ErrCommonHttpServerStartErr))
	}

	if !appstate.CheckDBConnectIsOk() {
		response.JsonExit(r, gerror.NewCode(errcode.ErrCommonDBConnectErr))
	}

	if !appstate.CheckDBInitIsOk() {
		response.JsonExit(r, gerror.NewCode(errcode.ErrCommonDBLoading))
	}

	if !appstate.CheckSyncMPUResourcesIsOk() {
		response.JsonExit(r, gerror.NewCode(errcode.ErrCommonSyncMPUResources))
	}

	if !appstate.CheckSyncReduMPUResIsOk() &&
		!strings.Contains(r.URL.Path, "notifyMysqlDataSyncOk") {
		response.JsonExit(r, gerror.NewCode(errcode.ErrCommonSyncReduMPURes))
	}
	r.Middleware.Next()
}

// CtxInit
// @summary 初始化http请求上下文
// @param1  r *ghttp.Request "http请求对象"
// @return1 nil
func CtxInit(r *ghttp.Request) {
	// 初始化, 务必最开始执行
	customCtx := &define.Context{}
	service.Context.Init(r, customCtx)

	// 将自定义的上下文对象传递到模板变量中使用
	r.Assigns(g.Map{
		"Context": customCtx,
	})

	//客户端支持的语言类型
	language := utils.ParseAcceptLanguage(r.GetHeader("Accept-Language"))
	service.Context.SetAcceptLanguage(r.Context(), language)

	// 执行下一步请求逻辑
	r.Middleware.Next()
}

// UpdateLastAccessTime
// @summary 更新数据更新时间
// @param1  r *ghttp.Request "http请求对象"
// @return1 nil
func UpdateLastAccessTime(r *ghttp.Request) {
	if r.Method != "GET" {
		//更新数据更新时间
		lastAccessTimeMs := time.Now().UnixNano() / 1e6
		if err := ioutil.WriteFile("lastAccessTime.ms", []byte(gconv.String(lastAccessTimeMs)), 0666); err != nil {
			g.Log().Error("更新最后访问时间：", err)
		}
	}

	// 执行下一步请求逻辑
	r.Middleware.Next()
}

// RequestId
// @summary 添加自定义RequestId
//          生成请求ID添加到Http请求头和响应头中
// @param1  r *ghttp.Request "http请求对象"
// @return1 nil
func RequestId(r *ghttp.Request) {
	requestId := r.GetCtxVar(global.ContextKeyRequestId).String()
	if requestId == "" {
		if v := r.Header.Get(global.ContextKeyRequestId); v != "" {
			requestId = v
		} else {
			requestId = guid.S()
		}
		r.SetCtxVar(global.ContextKeyRequestId, requestId)
		r.Response.Header().Set(global.ContextKeyRequestId, requestId)
	}
	r.Middleware.Next()
}

// CORS
// @summary 允许跨域请求
// @param1  r *ghttp.Request "http请求对象"
// @return1 nil
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()

	// 执行下一步请求逻辑
	r.Middleware.Next()
}

// LoginTokenAuth
// @summary 登录认证
//          1、校验Http请求携带的token是否有效
//          2、将token及登录用户信息设置到请求上下文
// @param1  r *ghttp.Request "http请求对象"
// @return1 nil
func LoginTokenAuth(r *ghttp.Request) {
	//校验请求令牌是否有效,并返回gtoken缓存的登录用户信息
	user := apiSys.Auth.GetGTokenCacheUser(r)
	if user == nil {
		g.Log().Noticef("request unAuthorized: %s %s %v", r.RemoteAddr, r.RequestURI, r.Header)
		response.UnAuthorizedJson(r)
	}

	//解析请求头中的token
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		g.Log().Noticef("request unAuthorized: %s %s %v", r.RemoteAddr, r.RequestURI, r.Header)
		response.UnAuthorizedJson(r)
	}
	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) == 2 && parts[0] == "Bearer" && parts[1] != "" {
		service.Context.SetToken(r.Context(), parts[1])
	} else {
		response.UnAuthorizedJson(r)
	}

	//请求上下文设置用户信息
	service.Context.SetUser(r.Context(), user)

	//执行下一步请求逻辑
	r.Middleware.Next()
}

// VerifyPermissions
// @summary 权限校验
func VerifyPermissions(r *ghttp.Request) {
	reqInterface := r.Method + " " + r.URL.Path
	if !service.InterfaceSysMenu().CheckHttpInterfacePermissions(r.GetCtx(), reqInterface) {
		response.JsonExit(r, gerror.NewCode(errcode.ErrHttpForbidden))
	}
	r.Middleware.Next()
}

// CheckClientOS
// @summary 校验客户端操作系统
func CheckClientOS(r *ghttp.Request) {
	url := r.GetUrl()
	if !strings.Contains(r.GetUrl(), "/mpuaps/v1") {
		//前端页面路由根据访问端系统类型区分返回客户端html文件
		userAgentInfo := userAgent.New(r.Header.Get("User-Agent"))
		if strings.Contains(userAgentInfo.OSInfo().FullName, "Android") {
			//Pad版本前端页面
			//r.RequestURI = "pad" + r.RequestURI
			g.Log(url)
		} else {
			//PC 版本前端页面
			//r.RequestURI = "" + r.RequestURI
			g.Log(url)
		}
	}
	r.Middleware.Next()
}
