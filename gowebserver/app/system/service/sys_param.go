package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/common/utils/page"
	"gowebserver/app/common/utils/valueCheck"
	"gowebserver/app/system/dao"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"strings"
)

//ISysParamMgr 系统参数配置管理接口
type ISysParamMgr interface {
	InitSysParam()
	AddSysParamCfg(ctx context.Context, sysParamCfg *model.SysParam) (err error)
	GetSysParamCfg(ctx context.Context, key string) (sysCfg *model.SysParam, err error)
	SetSysParamCfg(ctx context.Context, key, value string) (err error)

	GetListSearch(ctx context.Context, req *define.SysParamSelectReq) (pageInfo *page.Paging, sysParamList []*model.SysParam, err error)

	GetLogStorageStrategy(ctx context.Context, logType define.ELogType) (logStorageCfg *define.SysParamLogStorageStrategy, err error)
}

type sysParamMgrImpl struct{}

var sysParamMgrServer = sysParamMgrImpl{}

func InterfaceSysParamMgr() ISysParamMgr {
	return &sysParamMgrServer
}

//InitSysParam
//@summary 初始化系统参数信息
func (s *sysParamMgrImpl) InitSysParam() {
	//系统用户相关参数默认值
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "系统账户密码长度最小值",
		ParamKey:   global.SysParamKey_UserPasswordLengthMin,
		ParamValue: "8",
		ParamType:  1,
		Remark:     "用户口令长度最小值,参数不小于8",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "系统账户密码长度最大值",
		ParamKey:   global.SysParamKey_UserPasswordLengthMax,
		ParamValue: "16",
		ParamType:  1,
		Remark:     "用户口令长度最大值,参数不小于8",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "用户登录访问最大连接数",
		ParamKey:   global.SysParamKey_UserLoginAccessConnectionsNumMax,
		ParamValue: "50", //默认值50
		ParamType:  1,
		Remark:     "用户登录访问最大连接数,参数大于0,小于100",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})

	//todo

	//系统上传下载参数默认值
	//todo

	//系统黑白名单模式默认值
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "系统访问黑白名单过滤模式",
		ParamKey:   global.SysParamKey_AccessingBlacklistFilterType,
		ParamValue: "0",
		ParamType:  1,
		Remark:     "0=不启用 1=黑名单模式 2=白名单模式",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})

	//告警参数配置默认值
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "普通告警声音启用状态",
		ParamKey:   global.SysParamKey_AlarmMinorLevelSoundOn,
		ParamValue: "1",
		ParamType:  1,
		Remark:     "0=不启用 1=启用",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "严重告警声音启用状态",
		ParamKey:   global.SysParamKey_AlarmSeverityLevelSoundOn,
		ParamValue: "1",
		ParamType:  1,
		Remark:     "0=不启用 1=启用",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "普通告警声音音调",
		ParamKey:   global.SysParamKey_AlarmMinorLevelAlarmTone,
		ParamValue: "告警提示音1.wav",
		ParamType:  1,
		Remark:     "告警提示音文件名称",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "严重告警声音音调",
		ParamKey:   global.SysParamKey_AlarmSeverityLevelAlarmTone,
		ParamValue: "告警提示音2.wav",
		ParamType:  1,
		Remark:     "告警提示音文件名称",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})

	//日志存储策略默认值
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "日志最大存储条数",
		ParamKey:   global.SysParamKey_LogStorageNumMax,
		ParamValue: "1000000", //100万条
		ParamType:  1,
		Remark:     "整形参数,限制日志最多存储多少条",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "告警日志存储策略",
		ParamKey:   global.SysParamKey_LogStorageDayMax,
		ParamValue: "730", //730天
		ParamType:  1,
		Remark:     "整形参数,限制日志最多存储多少天",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "告警日志存储策略",
		ParamKey:   global.SysParamKey_AlarmLogStorageStrategy,
		ParamValue: "{\"logMaxNum\":200000,\"logMaxDay\":60}",
		ParamType:  1,
		Remark:     "json格式参数",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "设备上下线日志存储策略",
		ParamKey:   global.SysParamKey_DeviceOnlineLogStorageStrategy,
		ParamValue: "{\"logMaxNum\":200000,\"logMaxDay\":60}",
		ParamType:  1,
		Remark:     "json格式参数",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "操作日志存储策略",
		ParamKey:   global.SysParamKey_OperLogStorageStrategy,
		ParamValue: "{\"logMaxNum\":200000,\"logMaxDay\":60}",
		ParamType:  1,
		Remark:     "json格式参数",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})
	_ = s.AddSysParamCfg(context.TODO(), &model.SysParam{
		ParamName:  "登录登出日志存储策略",
		ParamKey:   global.SysParamKey_LoginHistoryLogStorageStrategy,
		ParamValue: "{\"logMaxNum\":200000,\"logMaxDay\":60}",
		ParamType:  1,
		Remark:     "json格式参数",
		CreateBy:   "SysAdmin",
		CreateTime: gtime.Now(),
	})
}

//GetListSearch
//@summary 分页查询配置信息
func (s *sysParamMgrImpl) GetListSearch(ctx context.Context, req *define.SysParamSelectReq) (pageInfo *page.Paging, sysParamList []*model.SysParam, err error) {
	var total int
	err = g.Try(func() {
		m := dao.SysParam.Ctx(ctx)
		if req.Keyword != "" {
			req.Keyword = strings.Replace(req.Keyword, "'", "''", -1)
			m = m.Where(fmt.Sprintf("locate('%s', %s) > 0", req.Keyword, dao.SysParam.Columns.ParamName)).
				WhereOr(fmt.Sprintf("locate('%s', %s) > 0", req.Keyword, dao.SysParam.Columns.ParamKey))
		}

		total, err = m.Count()
		valueCheck.ErrIsNil(ctx, err)

		if req.PageSize == 0 {
			req.PageSize = total
		}
		if req.OrderByColumn != "" && req.IsAsc != "" {
			req.OrderByColumn = gstr.CaseSnake(req.OrderByColumn)
			m = m.Order(req.OrderByColumn + " " + req.IsAsc)
		} else {
			m = m.OrderAsc(dao.SysParam.Columns.ParamKey)
		}

		pageInfo = page.CreatePaging(req.PageNum, req.PageSize, total)
		m = m.Limit(pageInfo.StartNum, pageInfo.PageSize)
		err = m.Scan(&sysParamList)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

//AddSysParamCfg
//@summary 添加系统参数配置
func (s *sysParamMgrImpl) AddSysParamCfg(ctx context.Context, sysParamCfg *model.SysParam) (err error) {
	//校验key是否重复
	if sysCfg, _ := s.GetSysParamCfg(ctx, sysParamCfg.ParamKey); sysCfg != nil {
		return gerror.NewCode(errcode.ErrCommonInvalidParameter, "系统参数key重复")
	}

	//校验值是否合法
	if !s.CheckSysParamValue(sysParamCfg.ParamKey, sysParamCfg.ParamValue) {
		return gerror.NewCode(errcode.ErrCommonInvalidParameter)
	}

	//保存参数配置
	sysParamCfg.ParamId, err = dao.SysParam.Ctx(ctx).Data(sysParamCfg).InsertAndGetId()
	return
}

//GetSysParamCfg
//@summary 查询系统参数信息
func (s *sysParamMgrImpl) GetSysParamCfg(ctx context.Context, key string) (sysCfg *model.SysParam, err error) {
	err = g.Try(func() {
		err = dao.SysParam.Ctx(ctx).Where(dao.SysParam.Columns.ParamKey, key).FindScan(&sysCfg)
		valueCheck.ErrIsNil(ctx, err)
		valueCheck.ValueIsNil(sysCfg, errcode.ErrCommonNotSupported.Message())
	})
	return
}

//SetSysParamCfg
//@summary 更新系统参数信息
func (s *sysParamMgrImpl) SetSysParamCfg(ctx context.Context, key, value string) (err error) {
	//校验值是否合法
	if !s.CheckSysParamValue(key, value) {
		return gerror.NewCode(errcode.ErrCommonInvalidParameter)
	}

	err = g.Try(func() {
		_, err = dao.SysParam.Ctx(ctx).
			Where(dao.SysParam.Columns.ParamKey, key).
			Data(dao.SysParam.Columns.ParamValue, value).
			Update()
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

//GetLogStorageStrategy
//@summary 查询日志存储配置接口
func (s *sysParamMgrImpl) GetLogStorageStrategy(ctx context.Context, logType define.ELogType) (logStorageCfg *define.SysParamLogStorageStrategy, err error) {
	var cfgInfo *model.SysParam
	switch logType {
	case define.ELogTypeOperlog: //操作日志
		cfgInfo, err = s.GetSysParamCfg(ctx, global.SysParamKey_OperLogStorageStrategy)
	case define.ELogTypeLoginHistory: //登录日志
		cfgInfo, err = s.GetSysParamCfg(ctx, global.SysParamKey_LoginHistoryLogStorageStrategy)
	case define.ELogTypeAlarmLog: //告警日志
		cfgInfo, err = s.GetSysParamCfg(ctx, global.SysParamKey_AlarmLogStorageStrategy)
	case define.ELogTypeDeviceOfflineLog: //设备上下线日志
		cfgInfo, err = s.GetSysParamCfg(ctx, global.SysParamKey_DeviceOnlineLogStorageStrategy)
	default:
		return nil, gerror.NewCode(errcode.ErrCommonInvalidParameter)
	}
	if err != nil {
		return
	}
	if cfgInfo == nil {
		return nil, gerror.NewCode(errcode.ErrCommonInvalidParameter)
	}
	if err = gjson.DecodeTo(cfgInfo.ParamValue, &logStorageCfg); err != nil {
		return
	}
	logStorageCfg.LogType = int(logType)
	return
}

//CheckSysParamValue
//@summary 校验系统参数值是否合法
func (s *sysParamMgrImpl) CheckSysParamValue(key, value string) bool {
	switch key {
	case global.SysParamKey_UserPasswordLengthMin:
	case global.SysParamKey_UserPasswordLengthMax:
	case global.SysParamKey_UserPasswordExpirationTime:
	case global.SysParamKey_UserLoginIncorrectAttemptTimes:
	case global.SysParamKey_UserLoginLockTime:

	case global.SysParamKey_UserLoginAccessConnectionsNumMax:
		{
			connMaxNum := gconv.Int(value)
			if connMaxNum <= 0 || connMaxNum > 100 {
				return false
			}
		}

	case global.SysParamKey_UploadPath:
	case global.SysParamKey_UploadImgType:
	case global.SysParamKey_UploadImgSizeMax:
	case global.SysParamKey_UploadFileType:
	case global.SysParamKey_UploadFileSizeMax:
	case global.SysParamKey_UploadFileCountMaxMinute:
	case global.SysParamKey_DownloadPath:

	case global.SysParamKey_AccessingBlacklistFilterType:
		{
			filterMode := gconv.Int(value)
			if filterMode != 0 && filterMode != 1 && filterMode != 2 {
				return false
			}
		}

	case global.SysParamKey_AlarmMinorLevelSoundOn:
	case global.SysParamKey_AlarmSeverityLevelSoundOn:
	case global.SysParamKey_AlarmMinorLevelAlarmTone:
	case global.SysParamKey_AlarmSeverityLevelAlarmTone:

	case global.SysParamKey_LogStorageDayMax:
		{
			num := gconv.Int(value)
			if num <= 0 {
				return false
			}
		}
	case global.SysParamKey_LogStorageNumMax:
		{
			num := gconv.Int(value)
			if num <= 0 {
				return false
			}
		}

	case global.SysParamKey_AlarmLogStorageStrategy, global.SysParamKey_DeviceOnlineLogStorageStrategy,
		global.SysParamKey_OperLogStorageStrategy, global.SysParamKey_LoginHistoryLogStorageStrategy:
		{
			var param define.SysParamLogStorageStrategy
			if err := gjson.DecodeTo(value, &param); err != nil {
				return false
			}
			var (
				logStorageDayMax int
				logStorageNumMax int
			)
			if cfgParam, err := s.GetSysParamCfg(context.TODO(), global.SysParamKey_LogStorageDayMax); err != nil {
				break
			} else {
				logStorageDayMax = gconv.Int(cfgParam.ParamValue)
			}
			if cfgParam, err := s.GetSysParamCfg(context.TODO(), global.SysParamKey_LogStorageNumMax); err != nil {
				break
			} else {
				logStorageNumMax = gconv.Int(cfgParam.ParamValue)
			}
			if param.LogMaxNum > logStorageNumMax || param.LogMaxNum < 0 ||
				param.LogMaxDay > logStorageDayMax || param.LogMaxDay < 0 {
				return false
			}
		}

	default:
		//未知参数无限制(前端自定义的参数)
		return true
	}

	//系统内置参数校验通过
	return true
}
