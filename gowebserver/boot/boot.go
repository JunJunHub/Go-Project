package boot

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"runtime"
	"strings"

	"gowebserver/app/common/global"
	"gowebserver/app/common/service/appstate"

	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"github.com/gogf/gf/util/gmode"
	_ "go.uber.org/automaxprocs"
)

//ProfilesActive 使用的配置文件,根据默认配置运行环境[local/dev/prod]指定
var ProfilesActive string

//程序默认配置文件
func initConfig() {
	//获取环境变量(程序运行环境)
	ProfilesActive = os.Getenv("GF_PROFILES_ACTIVE")
	if ProfilesActive == "" {
		ProfilesActive = gconv.String(g.Cfg().Get("profile.active"))
	}
	switch ProfilesActive {
	case "local":
		gmode.SetDevelop()
		g.Cfg().SetFileName("config-local.toml")
	case "dev":
		gmode.SetDevelop()
		g.Cfg().SetFileName("config-dev.toml")
	case "prod":
		gmode.SetProduct()
		g.Cfg().SetFileName("config-prod.toml")
	}
	g.Log().Line().Infof("运行环境：%s[GF_PROFILES_ACTIVE:%s]", ProfilesActive, os.Getenv("GF_PROFILES_ACTIVE"))
}

//数据表初始化
func initDBTable(sqlScriptFile string) error {
	file, err := os.Open(sqlScriptFile)
	if err != nil {
		g.Log().Error(err)
		return err
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			g.Log().Error(err)
		}
	}(file)

	bFindDelimiterCount := 0
	sqlBuf := new(bytes.Buffer)
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		// 跳过注释行和空行
		if len(line) < 1 || line[0] == '-' || line[0] == '#' {
			continue
		}
		sqlBuf.Write(line)

		// 如果发现包含 DELIMITER 关键字, 说明开始声明一个函数, 找到第二个 DELIMITER 表示读完一个函数声明语句
		// 否则判断最后字符是 ; 说明读一条完整的sql语句结束
		if strings.Contains(string(line), "DELIMITER") || strings.Contains(string(line), "delimiter") {
			bFindDelimiterCount += 1
			if bFindDelimiterCount == 2 {
				//执行函数声明语句
				_, err = g.DB().Exec(sqlBuf.String())
				if err != nil {
					g.Log().Errorf("SQL %s Exec Err: %v", sqlBuf.String(), err)
				}
				bFindDelimiterCount = 0
				sqlBuf.Reset()
			}
		} else {
			if bFindDelimiterCount == 0 && line[len(line)-1] == ';' {
				_, err = g.DB().Exec(sqlBuf.String())
				if err != nil {
					if strings.Contains(sqlBuf.String(), "create") || strings.Contains(sqlBuf.String(), "CREATE") {
						//建表语句执行失败,返回数据库初始化错误
						//其它SQL报错只记录错误日志
						g.Log().Errorf("SQL %s Exec Err: %v", sqlBuf.String(), err)
						return err
					}
					g.Log().Warningf("SQL %s Exec Err: %v", sqlBuf.String(), err)
				}
				sqlBuf.Reset()
			}
		}
	}
	return nil
}

//数据库初始化
func initDB() {
	//如果数据库不存在创建数据库
	dbCfgMap := g.Cfg().GetMapStrStr("database")
	if dbCfgMap == nil {
		g.Log().Error("请检查数据库配置")
		appstate.SetMpuApsStatusErr(appstate.EMpuApsStatusDBConnect)
		return
	}
	if _, err := g.DB().Exec(fmt.Sprintf("use %s;", dbCfgMap["name"])); err != nil {
		g.Log().Warning(err)
		g.Log().Info("根据配置创建数据库:", dbCfgMap)

		//数据库参数配置错误 or 配置中指定的数据库名不存在,默认的g.DB()对象不能正常使用,新建一个gdb对象(不指定数据库名),用来创建数据库.
		var dbGroup = "mpuapsBoot"
		var dbCfgNode = gdb.ConfigNode{
			Type: dbCfgMap["type"],
			Host: dbCfgMap["host"],
			Port: dbCfgMap["port"],
			User: dbCfgMap["user"],
			Pass: dbCfgMap["pass"],
		}
		gdb.SetConfigGroup(dbGroup, gdb.ConfigGroup{
			dbCfgNode,
		})
		db, err := gdb.New(dbGroup)
		if err != nil {
			g.Log().Error(err)
		}

		if _, err := db.Exec(fmt.Sprintf("create database %s;", dbCfgMap["name"])); err != nil {
			g.Log().Error(err)
			g.Log().Warning("请确认数据库配置!")
			appstate.SetMpuApsStatusErr(appstate.EMpuApsStatusDBConnect)
			return
		}
	}
	appstate.SetMpuApsStatusOk(appstate.EMpuApsStatusDBConnect)

	//表初始化
	if appstate.CheckDBConnectIsOk() {
		sqlCfg := gconv.String(g.Cfg().Get("database.sql"))
		sqlScriptFiles := strings.Split(sqlCfg, ";")
		g.Log().Noticef("sqlScriptFiles: %s, fileNum: %d", sqlCfg, len(sqlScriptFiles))
		for _, value := range sqlScriptFiles {
			g.Log().Noticef("run sqlScriptFile: %s", value)
			if err := initDBTable(value); err != nil {
				g.Log().Error(err, value, "数据库脚本执行失败!")
				appstate.SetMpuApsStatusErr(appstate.EMpuApsStatusDBInit)
				return
			}
		}
		appstate.SetMpuApsStatusOk(appstate.EMpuApsStatusDBInit)
	}
}

//程序初始化
func init() {
	//设置系统时区
	if err := gtime.SetTimeZone("Asia/Shanghai"); err != nil {
		g.Log().Warning("SetTimeZone(\"Asia/Shanghai\") Failed!", err)
	}

	//程序启动时间
	global.ProcessStartTime = gtime.Now()

	//初始化配置信息
	initConfig()

	//初始化默认日志对象配置
	global.InitLoggerDefault()

	//打印版本信息
	global.PrintVersionInfo()

	//打印GOMAXPROCS数量
	g.Log().Noticef("GOMAXPROCS：%v", runtime.GOMAXPROCS(-1))

	//数据库初始化
	initDB()
}
