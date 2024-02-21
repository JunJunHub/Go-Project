package global

import (
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

var (
	Version           = "v1.0.0" //软件版本号
	BinBuildTime      string     //编译时初始化: 编译时间
	BinBuildGoVersion string     //编译时初始化: Go版本

	ProcessStartTime *gtime.Time //程序启动时间
)

//PrintVersionInfo 获取版本信息
func PrintVersionInfo() {
	versionInfo := fmt.Sprintf(
		"\n==================================================\n"+
			"[VersionInfo]:\n"+
			"Version=%s\n"+
			"BinBuildTime=%s\n"+
			"BinBuildGoVersion=%s\n"+
			"BinVersionMd5=%s\n"+
			"==================================================\n",
		Version, BinBuildTime, BinBuildGoVersion, gdebug.BinVersionMd5())

	g.Log().Stdout().Notice(versionInfo)
}
