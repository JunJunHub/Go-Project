// ==========================================================================
// 系统模块入口.
// ==========================================================================

package system

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"gowebserver/app/common/global"
	"gowebserver/app/common/middleware"
	"gowebserver/app/common/service/appstate"
	"gowebserver/app/common/task"
	"gowebserver/app/system/api"
	"gowebserver/app/system/service"
	"gowebserver/app/system/service/publish"
	"time"
)

// 系统功能路由注册
func routerInit() {
	httpServer := g.Server()
	prefix := g.Cfg().GetString("server.Prefix")

	// ===================调试API=============================================
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CtxInit)
		group.Middleware(middleware.RequestId)
		group.ALL("/debug/removeLoginLockUser", api.DebugApi.RemoveLoginLockUser)
		group.ALL("/debug/checkPasswordRight", api.DebugApi.CheckPasswordRight)
		group.GET("/debug/printOnlineUserInfo", api.DebugApi.PrintOnlineUserInfo)
		group.GET("/debug/printWsSubLinkInfo", api.DebugApi.PrintWsSubLinkInfo)
		group.GET("/debug/printUserPermissions", api.DebugApi.PrintUserPermissions)
	})

	// ============================路由注册=====================================
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		//请求上下文初始化
		group.Middleware(middleware.CheckMpuApsStatus)
		group.Middleware(middleware.CtxInit)
		group.Middleware(middleware.RequestId)

		//ws订阅接口
		//websocket协议在握手阶段借用了HTTP的协议
		//但是在前端JavaScript websocket API中并没有修改请求头的方法,前端不方便在Http头中把token带过来
		//让前端把token放在url中带过来,在Subscribe接口中单独做登录鉴权
		group.ALL("/ws/subscribe/*token", api.Auth.Subscribe)

		//gtoken鉴权拦截+登录登出路由注册
		if err := api.GfToken.Middleware(group); err != nil {
			g.Log().Error(err)
		}
		//刷新token有效期
		group.ALL("/auth/refreshToken", api.Auth.RefreshToken)

		//操作日志记录拦截器
		group.Middleware(middleware.VerifyPermissions, middleware.UpdateLastAccessTime, api.OperlogApi.OperationLog)

		//登录用户信息(登录用户对自己信息的维护)
		group.GET("/system/user/getInfo", api.UserApi.GetUserInfo)
		group.GET("/system/user/profile", api.UserApi.GetProfile)
		group.PUT("/system/user/profile", api.UserApi.UpdateProfile)
		group.PUT("/system/user/profile/updatePwd", api.UserApi.UpdatePassword)
		group.POST("/system/user/profile/avatar", api.UserApi.UpdateAvatar)

		//用户信息管理
		group.REST("/system/user", api.UserApi)
		group.GET("/system/user/info", api.UserApi.Info)
		group.PUT("/system/user/resetPwd", api.UserApi.ResetPwdSave)
		group.PUT("/system/user/changeStatus", api.UserApi.ChangeStatus)
		group.GET("/system/user/export", api.UserApi.Export)

		//角色管理
		group.GET("/system/role/getList", api.RoleApi.GetList)
		group.GET("/system/role/get", api.RoleApi.GetInfo)
		group.POST("/system/role/add", api.RoleApi.Add)
		group.PUT("/system/role/edit", api.RoleApi.Edit)
		group.DELETE("/system/role/del", api.RoleApi.Delete)
		group.GET("/system/role/menu", api.RoleApi.GetRoleMenu)
		group.PUT("/system/role/editMenu", api.RoleApi.EditRoleMenu)

		//group.GET("/system/role/channelGroup", api.RoleApi.GetRoleChnGroup)
		//group.PUT("/system/role/editChannelGroup", api.RoleApi.EditRoleChnGroup)
		//group.GET("/system/role/tvWall", api.RoleApi.GetRoleTvWall)
		//group.PUT("/system/role/editTvWall", api.RoleApi.EditRoleTvWall)
		//group.GET("/system/role/matrixScheme", api.RoleApi.GetRoleMatrixScheme)
		//group.PUT("/system/role/editMatrixScheme", api.RoleApi.EditRoleMatrixScheme)
		//group.GET("/system/role/matrixSchemeGroup", api.RoleApi.GetRoleMatrixSchemeGroup)
		//group.PUT("/system/role/editMatrixSchemeGroup", api.RoleApi.EditRoleMatrixSchemeGroup)
		//group.GET("/system/role/meetDispatch", api.RoleApi.GetRoleMeetDispatch)
		//group.PUT("/system/role/editMeetDispatch", api.RoleApi.EditRoleMeetDispatch)
		//group.GET("/system/role/meetGroup", api.RoleApi.GetRoleMeetDispatchGroup)
		//group.PUT("/system/role/editMeetGroup", api.RoleApi.EditRoleMeetDispatchGroup)

		//菜单管理
		group.GET("/system/menu/getTree", api.MenuApi.GetTree)
		group.GET("/system/menu/getList", api.MenuApi.GetMenus)

		//操作日志
		group.REST("/monitor/operlog", api.OperlogApi)
		group.DELETE("/monitor/operlog/clean", api.OperlogApi.Clean)
		group.GET("/monitor/operlog/export", api.OperlogApi.Export)

		//登录日志
		group.REST("/monitor/loginHistory", api.LoginHistoryApi)
		group.DELETE("/monitor/loginHistory/clean", api.LoginHistoryApi.Clean)
		group.GET("/monitor/loginHistory/export", api.LoginHistoryApi.Export)

		//在线用户管理
		group.REST("/monitor/online", api.OnlineUserUserApi)
		group.DELETE("/monitor/online", api.OnlineUserUserApi.ForceLogout)

		//定时任务管理
		group.REST("/monitor/job", api.JobApi)
		group.PUT("/monitor/job/run", api.JobApi.Run)
		group.PUT("/monitor/job/start", api.JobApi.Start)
		group.PUT("/monitor/job/stop", api.JobApi.Stop)

		//系统访问黑白名单
		group.GET("/system/blacklist/getList", api.SysBlackListApi.GetList)
		group.POST("/system/blacklist/add", api.SysBlackListApi.Add)
		group.DELETE("/system/blacklist/delete", api.SysBlackListApi.Delete)
		group.GET("/system/blacklist/mode", api.SysBlackListApi.GetFilterMode)
		group.PUT("/system/blacklist/mode", api.SysBlackListApi.SetFilterMode)

		//系统参数配置
		group.GET("/system/param/getList", api.SysParamApi.GetSysParamCfgList)
		group.PUT("/system/param/edit", api.SysParamApi.EditSysParamCfg)
		//日志存储策略设置
		group.GET("/monitor/log/storageStrategy", api.SysParamApi.GetLogStorageStrategy)
		group.PUT("/monitor/log/storageStrategy", api.SysParamApi.SetLogStorageStrategy)

	})
}

// system模块日志对象初始配置
func loggerInit() {
	//System模块日志对象
	global.LoggerInit(global.LoggerSystem)

	//系统功能-websocket通知日志对象
	global.LoggerInit(global.LoggerWebsocket)
	_ = g.Log(global.LoggerWebsocket).SetLevelStr("notice")

	//Http服务访问记录日志对象初始化
	global.LoggerInit(global.LoggerHttpAccessRecord, fmt.Sprintf("%s_{Ymd}.log", global.LoggerHttpAccessRecord))
	_ = g.Log(global.LoggerHttpAccessRecord).SetLevelStr("notice")
}

//注册定时任务
func crontabRegister() {
	//清理在线用户任务
	if err := task.Add(task.TasksEntity{
		FuncName:  "CleanOnlineUserTask",
		Func:      api.Auth.CleanOnlineUserTask,
		Param:     nil,
		StartTime: 0,
	}); err != nil {
		g.Log().Error(err)
	}

	go func() {
		for {
			time.Sleep(10 * time.Second) //10s清理一次在线用户表
			if appstate.CheckSyncMPUResourcesIsOk() {
				taskInfo := task.GetByName("CleanOnlineUserTask")
				taskInfo.Func()
			}
		}
	}()

	// You can add other task here
}

//wsClientDisConnectHandel
//@summary 登录用户Ws订阅链路断开连接处理接口
//         注意：此接口不能阻塞
//@param1  ctx context.Context "上下文"
//@param2  disConnectClientInfo publish.WSClientInfo "断开连接的客户端信息"
//@return  nil
func wsClientDisConnectHandel(ctx context.Context, disConnectClientInfo publish.WSClientInfo) {
	onlineUser := service.OnlineUser.GetOnlineUserInfoByLoginToken(disConnectClientInfo.Token)
	if onlineUser == nil {
		return
	}
	go func() {
		if disConnectClientInfo.IsCascadeSubChn {
			//级联订阅链路断链开连接,认为客户端登出,删除token缓存及在线用户信息
			api.Auth.DelGTokenCacheUser(disConnectClientInfo.Token)
			err := api.Auth.UserLogOutHandle(disConnectClientInfo.Token)
			if err != nil {
				g.Log().Warningf("客户端断开连接: %v, isCascadeSubChn: %t",
					disConnectClientInfo.WSSubLinkMap, disConnectClientInfo.IsCascadeSubChn)
				g.Log().Error("客户端ws订阅通道已断开,登出失败:", err)
				return
			}
		}

		//Web连接断开,确认是否重连
		api.Auth.CheckOnlineUserWsSubLinkAlive(*onlineUser)
	}()
}

//系统模块初始化
func init() {
	g.Log().Info("system module load.")

	loggerInit()
	routerInit()
	crontabRegister()
	publish.WsPublish.SetWsClientDisConnectHandelFunc(global.ModuleSystem, wsClientDisConnectHandel)
}
