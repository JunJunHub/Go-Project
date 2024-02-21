package accessMCUAdapter

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"gowebserver/app/common/global"
)

const defaultLoggerName = "AccessMCU"

//logger
//@summary 本模块使用的日志对象
func logger() *glog.Logger {
	return g.Log(global.LoggerAccessMCUAdapter)
}

//loggerInit
//@summary 模块内使用日志实例初始化,调试日志可以记录到一个独立的文件中
func loggerInit() {
	global.LoggerInit(global.LoggerAccessMCUAdapter)
}
