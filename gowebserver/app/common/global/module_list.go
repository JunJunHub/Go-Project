package global

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
)

//Module Http请求上下文记录模块信息
const Module = "Module"

// 模块列表枚举
const (
	ModuleCommon = "common" //通用功能模块

	ModuleSystem       = "system"               //系统功能模块
	ModuleSysWebsocket = "system.websocket.log" //系统功能-Websocket通知日志

	ModuleMpuSrv                 = "mpusrv"           //显控业务模块
	ModuleMpuSrvAccessMPU        = "mpusrv.rpc.log"   //显控业务-对接显控主业务(grpc协议)
	ModuleMpuSrvAccessMCUAdapter = "AccessMCUAdapter" //显控业务-对接MCU适配器

	ModuleCustomization = "customization" //定制功能模块

	// todo add module name
)

// 模块日志对象枚举
// 通过名称获取并使用对应模块日志对象, 例: g.Log(ModuleMpuCfg).Debug("...")
// 新增模块再次添加模块对应的日志对象名, 并修改 CheckLoggerModuleName
const (
	LoggerHttpAccessRecord = "httpGetHistory"
	LoggerCustomization    = ModuleCustomization
	LoggerCommon           = ModuleCommon

	LoggerSystem    = ModuleSystem
	LoggerWebsocket = ModuleSysWebsocket

	LoggerMpuSrv           = ModuleMpuSrv
	LoggerAccessMPU        = ModuleMpuSrvAccessMPU
	LoggerAccessMCUAdapter = ModuleMpuSrvAccessMCUAdapter

	// todo add logger
)

// 日志级别定义
// 按等级: DEBUG < INFO < NOTICE < WARN < ERROR < CRITICAL
// 按模式: ALL、DEV、PROD
const (
	LevelAll  = "all"
	LevelDev  = "dev"
	LevelProd = "Prod"

	LevelDebug    = "Debug"
	LevelInfo     = "Info"
	LevelNotice   = "Notice"
	LevelWarn     = "Warn"
	LevelError    = "Error"
	LevelCritical = "Critical"
)

// CheckLoggerModuleName
// @summary 校验是否已定义对应模块日志对象
// @param1  moduleName string "模块名"
func CheckLoggerModuleName(moduleName string) bool {
	switch moduleName {
	case
		LoggerHttpAccessRecord,
		LoggerCustomization,
		LoggerCommon,
		LoggerSystem,
		LoggerWebsocket,
		LoggerMpuSrv,
		LoggerAccessMPU,
		LoggerAccessMCUAdapter:
		return true

	default:
		return false
	}
}

// InitLoggerDefault
// @summary goframe框架默认日志对象初始化
func InitLoggerDefault() {
	//日志配置参数
	loggerConfig := g.Cfg().GetMap("logger")
	if err := g.Log().SetConfigWithMap(loggerConfig); err != nil {
		g.Log().Line().Error(err)
	}

	//设置日志输出格式
	g.Log().SetFlags(glog.F_ASYNC | glog.F_TIME_STD | glog.F_FILE_SHORT)
	g.Log().SetStackSkip(4)

	//调试模式
	if g.Cfg().GetBool("debug.enable") {
		_ = g.Log().SetLevelStr("debug")
	}
}

// LoggerInit
// @summary 封装自定义日志对象初始化方法
// @summary loggerName string "日志对象Key"
// @summary loggerFile ...string "日志输出文件名"
func LoggerInit(loggerName string, loggerFile ...string) {
	logOutFileName := "mpuaps_{Ymd}.log"
	if len(loggerFile) > 0 {
		logOutFileName = loggerFile[0]
	}

	//日志对象初始化,默认等级为Notice
	err := g.Log(loggerName).SetConfigWithMap(g.Map{
		"Path":     "log/mpuaps",
		"File":     logOutFileName,
		"Prefix":   loggerName,
		"Level":    LevelNotice,
		"Stdout":   false,
		"StStatus": 1,
	})
	if err != nil {
		g.Log().Error(err)
	}

	//日志输出格式
	g.Log(loggerName).SetFlags(glog.F_ASYNC | glog.F_TIME_STD | glog.F_FILE_SHORT)
	g.Log(loggerName).SetStackSkip(4)

	//调试模式
	if g.Cfg().GetBool("debug.enable") {
		_ = g.Log(loggerName).SetLevelStr("debug")
	}
}
