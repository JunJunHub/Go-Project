// ==========================================================================
// 定制功能模块
// @summary 客户定制化功能开发
//          1、环境设备管理
// @author  LiYongJun@2022.09.09
// @reviser LiYongJun@2022.09.09
//          模块初始化
// ==========================================================================

package customization

import (
	"context"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"gowebserver/app/common/global"
	"gowebserver/app/common/middleware"
	"gowebserver/app/customization/api"
	apiSys "gowebserver/app/system/api"
	"gowebserver/app/system/service/publish"
)

func init() {
	g.Log().Infof("%s module load", global.LoggerCustomization)
	loggerInit()
	routerInit()
	publish.WsPublish.SetWsClientDisConnectHandelFunc(global.ModuleMpuSrv, wsClientDisConnectHandel)
}

func loggerInit() {
	global.LoggerInit(global.LoggerCustomization)
}

//wsClientDisConnectHandel
//@summary 登录用户Ws订阅链路断开连接,显控业务处理接口
//         注意：此接口不能阻塞
//@param1  ctx context.Context "上下文"
//@param2  disConnectClientInfo publish.WSClientInfo "断开连接的客户端信息"
//@return  nil
func wsClientDisConnectHandel(ctx context.Context, disConnectClientInfo publish.WSClientInfo) {
	if disConnectClientInfo.IsCascadeSubChn {
		//上级平台断开连接
		g.Log().Noticef("检测到上级平台已断开连接 cascadeId:%s, isCascadeSubChn: %t, remoteAddr: %v",
			disConnectClientInfo.CascadeId, disConnectClientInfo.IsCascadeSubChn, disConnectClientInfo.WSSubLinkMap)
		//nothing todo
	} else {
		//ws订阅断开连接
		g.Log().Noticef("检测到客户端ws链路已断开连接  isCascadeSubChn: %t, remoteAddr: %v",
			disConnectClientInfo.IsCascadeSubChn, disConnectClientInfo.WSSubLinkMap)
		//nothing todo
	}
}

func routerInit() {
	httpServer := g.Server()
	prefix := g.Cfg().GetString("server.Prefix")

	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CheckMpuApsStatus)
		group.Middleware(middleware.CtxInit)
		group.Middleware(middleware.RequestId)
		routerMpuDebug(group)

		//请求token认证、操作日志记录
		// to do
		group.Middleware(middleware.LoginTokenAuth /*middleware.VerifyPermissions,*/, apiSys.OperlogApi.OperationLog)
		routerMpuApplianceManager(group)
	})
}

//routerMpuDebug
//@summary 调试接口路由注册
//@param1  routerGroup *ghttp.RouterGroup "路由组"
//@return  nil
func routerMpuDebug(routerGroup *ghttp.RouterGroup) {
	routerGroup.Group("/debug", func(group *ghttp.RouterGroup) {
	})
}

//routerMpuApplianceManager
//@summary 家电设备管理接口路由注册
//@param1  routerGroup *ghttp.RouterGroup "路由组"
//@return  nil
func routerMpuApplianceManager(routerGroup *ghttp.RouterGroup) {
	routerGroup.Group("", func(group *ghttp.RouterGroup) {
		group.GET("/appliance/control/scheme", api.MpuApplianceSchemeApi.GetList)
		group.GET("/appliance/control/schemeDetailInfo", api.MpuApplianceSchemeApi.GetDetailInfo)
		group.PUT("/appliance/control/scheme/load", api.MpuApplianceSchemeApi.Load)
		group.PUT("/appliance/control/scheme/unload", api.MpuApplianceSchemeApi.UnLoad)
		group.PUT("/appliance/control", api.MpuApplianceSchemeApi.Control)
	})
}
