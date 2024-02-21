// ==========================================================================
// 系统调试接口
// ==========================================================================

package api

import (
	"archive/zip"
	"fmt"
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/errors/gerror"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/common/service/appstate"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var DebugApi = new(debugApi)

type debugApi struct {
	BaseController
}

func (a *debugApi) Init(r *ghttp.Request) {
	a.Module = "Debug调试"
	r.SetCtxVar(global.Module, a.Module)
}

// SetLoggerLevel
// @summary 设置各模块日志打印等级
// @tags 	DEBUG
// @Param   moduleName query string true "日志模块: default、system、mpusrv、httpGetHistory"
// @Param   level query string true "日志级别: Debug < Info < Notice < Warn < Error < Critical || All、Dev、Prod"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/setLoggerLevel [PUT]
func (a *debugApi) SetLoggerLevel(r *ghttp.Request) {
	moduleName := r.GetString("moduleName")
	if !global.CheckLoggerModuleName(moduleName) {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter, errcode.ErrCommonInvalidParameter.Message()+"模块名不存在"))
	}

	level := r.GetString("level")
	if level == "" {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter, errcode.ErrCommonInvalidParameter.Message()+"日志级别错误"))
	}

	if err := g.Log(moduleName).SetLevelStr(level); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInternalError, err))
	}

	a.RespJsonExit(r, nil)
}

// SetDBLoggerLevel
// @summary 设置数据库操作日志打印模式
// @tags 	DEBUG
// @Param   level query string true "日志级别: Debug < Info < Notice < Warn < Error < Critical || All、Dev、Prod"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/setDBLoggerLevel [PUT]
func (a *debugApi) SetDBLoggerLevel(r *ghttp.Request) {
	level := r.GetString("level")
	if level == "" {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter, errcode.ErrCommonInvalidParameter.Message()+"日志级别错误"))
	}
	if err := g.DB().GetLogger().SetLevelStr(level); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInternalError, err))
	}
	a.RespJsonExit(r, nil)
}

// DownloadDebugLogs
// @summary 下载mpuaps调试日志
// @tags 	DEBUG
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/downloadDebugLogs [GET]
func (a *debugApi) DownloadDebugLogs(r *ghttp.Request) {
	//压缩日志目录
	curDir, _ := os.Getwd()
	logDirPath := curDir + "/log"
	logZipFilePath := curDir + g.Cfg().GetString("download.downPath") + "/"
	logZipFileName := "logs_" + strconv.FormatInt(time.Now().UnixNano(), 10) + ".zip"
	logZipFile := logZipFilePath + logZipFileName

	//下载路径不存在创建路径
	path, _ := filepath.Split(logZipFilePath)
	_, err := os.Stat(path)
	if err != nil || os.IsNotExist(err) {
		err = os.MkdirAll(path, os.ModePerm)
	}

	//打开压缩文件
	zipFile, err := os.Create(logZipFile)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	//遍历日志目录并将文件添加到 ZIP 文件中
	err = filepath.Walk(logDirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		//忽略目录本身，只压缩目录中的文件
		if !info.Mode().IsRegular() {
			return nil
		}

		//创建一个新的文件头部
		header, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		header.Name = path[len(logDirPath)+1:]

		//将文件头部写入 ZIP 文件
		writer, err := zipWriter.CreateHeader(header)
		if err != nil {
			return err
		}
		//打开要压缩的文件
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		//将文件内容写入 ZIP 文件
		_, err = io.Copy(writer, file)
		return err
	})
	if err != nil {
		a.RespJsonExit(r, err)
	}

	//重定向到文件下载链接
	downloadUrl := "http://" + r.Host + "/mpuaps/v1/download?fileName=" + logZipFileName + "&delete=true"
	r.Response.RedirectTo(downloadUrl)
}

// SetPProfServer
// @summary [启用/关闭]服务PProf性能分析
// @tags 	DEBUG
// @Param   enable query bool true "是否启用性能分析服务"
// @Param   port   query int false "pprof性能分析端口"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/setPProfServer [PUT]
func (a *debugApi) SetPProfServer(r *ghttp.Request) {
	const defaultPProfServerName = "pprof-server" //默认pprof性能分析服务名

	pprofPort := r.GetInt("port")
	if pprofPort == 0 {
		pprofPort = 9182
	}
	pprofEnable := r.GetBool("enable")
	if pprofEnable {
		runtime.SetMutexProfileFraction(1) //开启对锁调用的跟踪
		runtime.SetBlockProfileRate(1)     //开启对阻塞操作的跟踪

		// PProd Server
		go ghttp.StartPProfServer(pprofPort)
	} else {
		runtime.SetMutexProfileFraction(0) //开启对锁调用的跟踪
		runtime.SetBlockProfileRate(0)     //开启对阻塞操作的跟踪

		pprofServer := ghttp.GetServer(defaultPProfServerName)
		err := pprofServer.Shutdown()
		if err != nil {
			g.Log().Error("pprofServer Shutdown err:", err)
			a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInternalError, err))
		}
	}
	hostParam := strings.Split(r.Host, ":")
	a.RespJsonExit(r, nil, g.Map{
		"pprofUrl":          fmt.Sprintf("http://%s:%d/debug/pprof", hostParam[0], pprofPort),
		"pprofRouter":       "/debug/pprof /debug/pprof/*action /debug/pprof/cmdline /debug/pprof/profile /debug/pprof/symbol /debug/pprof/trace",
		"goToolParseCPUCmd": fmt.Sprintf("go tool pprof \"http://%s:%d/debug/pprof/profile\"", hostParam[0], pprofPort),
		"goToolParseMEMCmd": fmt.Sprintf("go tool pprof \"http://%s:%d/debug/pprof/heap\"", hostParam[0], pprofPort),
		"gfDoc":             "https://goframe.org/pages/viewpage.action?pageId=1114350",
		"golangDoc":         "blog.golang.org/profiling-go-programs",
	})
}

// Version
// @summary 获取版本信息
// @tags 	DEBUG
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/version [GET]
func (a *debugApi) Version(r *ghttp.Request) {
	global.PrintVersionInfo()
	a.RespJsonExit(r, nil, g.Map{
		"version":           global.Version,
		"buildTime":         global.BinBuildTime,
		"binBuildGoVersion": global.BinBuildGoVersion,
		"binVersionMd5":     gdebug.BinVersionMd5(),
	})
}

// GetMpuApsStatus
// @summary 获取程序状态
// @tags 	DEBUG
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/appStatus [GET]
func (a *debugApi) GetMpuApsStatus(r *ghttp.Request) {
	a.RespJsonExit(r, nil, g.Map{
		"mpuapsStatus": appstate.GetMpuApsStatus(),
		"remark":       "状态码：httpServerStatus: 1, DBConnectStatus: 2, DBInitStatus: 4, SyncMpuResource: 8",
	})
}

// GetMpuErrCode
// @summary 获取错误码描述信息
// @tags    DEBUG
// @Param   errcode  query int  true  "要查询的错误码"
// @Param   printAll query bool false "是否返回所有错误码信息"
// @Param   language query string true "国家地区语言: zh_CN、en"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/queryErrcode [GET]
func (a *debugApi) GetMpuErrCode(r *ghttp.Request) {
	//错误码信息
	type mpuapsErrcodeInfo struct {
		Code int    `json:"code"` //错误代码
		Key  string `json:"key"`  //多语言翻译Key
		Msg  string `json:"msg"`  //错误描述
	}
	//接口返回数据结果
	type respInfo struct {
		QueryCode  mpuapsErrcodeInfo   `json:"queryCode"`
		AllErrcode []mpuapsErrcodeInfo `json:"allCode"`
	}

	var resp respInfo
	queryCode := r.GetInt("errcode")
	bPrintAll := r.GetBool("printAll")
	language := r.GetString("language")
	if language == "" {
		language = "zh_CN"
	}

	//要查询的错误码信息
	queryErr := errcode.GetMpuapsErr(queryCode)
	resp.QueryCode = mpuapsErrcodeInfo{
		Code: queryCode,
		Key:  queryErr.Message(),
		Msg:  queryErr.MessageAreaLanguage(language),
	}

	//返回所有错误码描述
	if bPrintAll {
		allErrCodeMap := errcode.GetAllMpuapsErr()

		var codeIds []int
		for codeId := range allErrCodeMap {
			codeIds = append(codeIds, codeId)
		}
		sort.Ints(codeIds)
		for _, codeId := range codeIds {
			resp.AllErrcode = append(resp.AllErrcode, mpuapsErrcodeInfo{
				Code: allErrCodeMap[codeId].Code(),
				Key:  allErrCodeMap[codeId].Message(),
				Msg:  allErrCodeMap[codeId].MessageAreaLanguage(language),
			})
		}
	}
	a.RespJsonExit(r, nil, resp)
}

// MPUModuleReboot
// @summary MPU模块重启
// @tags    DEBUG
// @Param   moduleId  query int  true  "要重启的模块: 1=mpuaps 2=smmgr 3=mpuserver 4="
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/reboot [POST]
func (a *debugApi) MPUModuleReboot(r *ghttp.Request) {
	moduleId := r.GetInt("moduleId")

	switch moduleId {
	case 1:
		panic("debug kill [mpuaps] and the process will be reboot!")

	default:
		a.RespJsonExit(r, nil, g.Map{
			"Tips": "不支持的操作",
		})
	}
	a.RespJsonExit(r, nil, g.Map{
		"Tips": "执行成功，请稍后...",
	})
}
