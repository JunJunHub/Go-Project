package define

import (
	"gowebserver/app/common/define"
	"gowebserver/app/system/model"
)

type SysParamEditReq struct {
	ParamID    int64  `json:"paramId"`    // 参数配置id
	ParamKey   string `json:"ParamKey"`   // 参数键名
	ParamValue string `json:"ParamValue"` // 参数值
}

type SysParamSelectReq struct {
	Keyword string `p:"keyword"` // 搜索查询关键字
	define.SelectPageReq
}

type SysParamSelectRes struct {
	List  []model.SysParam `json:"list"`  // 系统参数信息列表
	Page  int              `json:"page"`  // 页码
	Size  int              `json:"size"`  // 每页条数
	Total int              `json:"total"` // 总数
}

type SysParamLogStorageStrategy struct {
	LogType   int `json:"logType"`   //日志类别：1、操作日志 2、登录日志 3、告警日志 4、设备上下线日志
	LogMaxNum int `json:"logMaxNum"` //最大条数，单位条，0表示不根据条数限制
	LogMaxDay int `json:"logMaxDay"` //最大天数，单位天，0表示不根据天数限制
}

//ELogType 日志类型
type ELogType int

const (
	ELogTypUnknow            ELogType = 0 //未知
	ELogTypeOperlog          ELogType = 1 //操作日志
	ELogTypeLoginHistory     ELogType = 2 //登录日志
	ELogTypeAlarmLog         ELogType = 3 //告警日志
	ELogTypeDeviceOfflineLog ELogType = 4 //设备离线日志
)
