// ==========================================================================
// 公共模块
// ==========================================================================

package common

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
	"gowebserver/app/common/api"
	"gowebserver/app/common/global"
	"gowebserver/app/common/middleware"
)

func routerInit() {
	httpServer := g.Server()
	prefix := g.Cfg().GetString("server.Prefix")

	//公共通用接口
	httpServer.Group(prefix, func(group *ghttp.RouterGroup) {
		routerMpuDebug(group)

		// 上下文
		group.Middleware(middleware.CheckMpuApsStatus)
		group.Middleware(middleware.CtxInit)
		group.Middleware(middleware.RequestId)

		// 图形验证码获取
		group.GET("/auth/captchaImage", api.CaptchaApi.CaptchaImage)

		// 文件下载
		group.GET("/download", api.DownloadApi.Download)

		// 文件上传
		group.POST("/upload/upImg", api.UploadApi.UpImg)
		group.POST("/upload/upImgs", api.UploadApi.UpImgs)
		group.POST("/upload/upFile", api.UploadApi.UpFile)
		group.POST("/upload/upFiles", api.UploadApi.UpFiles)
		group.POST("/upload/ckEditorUp", api.UploadApi.CkEditorUp)
	})
}

func routerMpuDebug(routerGroup *ghttp.RouterGroup) *ghttp.RouterGroup {
	//调试接口路由注册
	return routerGroup.Group("", func(group *ghttp.RouterGroup) {
		// g.Log()日志对象打印等级配置
		group.PUT("/debug/setLoggerLevel", api.DebugApi.SetLoggerLevel)
		// 设置数据库操作日志等级
		group.PUT("debug/setDBLoggerLevel", api.DebugApi.SetDBLoggerLevel)
		// 下载调试日志
		group.GET("debug/downloadDebugLogs", api.DebugApi.DownloadDebugLogs)
		// 启用/停用 PProf性能分析服务
		group.PUT("/debug/setPProfServer", api.DebugApi.SetPProfServer)
		// 打印版本信息
		group.GET("/debug/version", api.DebugApi.Version)
		// 获取程序状态机状态
		group.GET("/debug/appStatus", api.DebugApi.GetMpuApsStatus)
		// 查询错误码描述
		group.GET("/debug/queryErrcode", api.DebugApi.GetMpuErrCode)
		// MPU模块重启
		group.ALL("/debug/reboot", api.DebugApi.MPUModuleReboot)
	})
}

// common模块日志对象初始配置
func loggerInit() {
	err := g.Log(global.LoggerCommon).SetConfigWithMap(g.Map{
		"Path":     "log/mpuaps",
		"File":     global.LoggerCommon + "_{Ymd}.log",
		"Prefix":   global.LoggerCommon,
		"Level":    global.LevelNotice,
		"Stdout":   false,
		"StStatus": 1,
	})
	if err != nil {
		g.Log().Error(err)
	}

	// 日志输出格式
	g.Log(global.LoggerCommon).SetFlags(glog.F_ASYNC | glog.F_TIME_DATE | glog.F_TIME_TIME | glog.F_FILE_SHORT | glog.F_CALLER_FN)
}

func init() {
	g.Log().Info("common module load.")

	loggerInit()
	routerInit()
}
