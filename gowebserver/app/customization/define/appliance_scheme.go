package define

import "gowebserver/app/common/define"

type MPUApplianceCtrlScheme struct {
	GroupID    int64  `json:"groupId,omitempty"`   // 所属调度预案分组ID,不属于任何分组为0
	GroupName  string `json:"groupName,omitempty"` // 预案所属分组名
	LoadState  int64  `json:"loadState"`           // 预案加载状态
	LoadTimeMS int64  `json:"loadTimeMs"`          // 矩阵预案启动时刻(毫秒时间戳)
	SchemeID   int64  `json:"schemeId"`            // 矩阵调度预案ID
	SchemeName string `json:"schemeName"`          // 矩阵调度预案名称
}

// MPUApplianceCtrlSchemeSelectPageReq 分页查询请求参数
type MPUApplianceCtrlSchemeSelectPageReq struct {
	GroupId int64  `p:"groupId"` // description:"预案分组Id"
	Keyword string `p:"keyword"` //description:"预案名称关键字"
	define.SelectPageReq
}

type MPUApplianceCtrlSchemeList struct {
	List  []MPUApplianceCtrlScheme `json:"list,omitempty"`
	Page  int                      `json:"page"`
	Size  int                      `json:"size"`
	Total int                      `json:"total"`
}

type MPUApplianceCtrlParam struct {
	SchemeID        int64       `json:"schemeId"`                  // 对应设备调度预案ID
	DevID           string      `json:"devId"`                     // 对应设备ID
	DevType         int64       `json:"devType"`                   // 设备类型：1=灯 2=空调 3=窗帘
	CtrlCmd         string      `json:"ctrlCmd"`                   // 控制指令；各类设备的控制指令格式另行说明；
	DevCapabilities interface{} `json:"devCapabilities,omitempty"` // 设备能力信息，根据设备能力展示控制按钮；
}
type MPUApplianceCtrlParamList struct {
	List  []MPUApplianceCtrlParam `json:"list"`
	Page  int                     `json:"page"`
	Size  int                     `json:"size"`
	Total int                     `json:"total"`
}

type MPULampCapabilities struct {
	BSupportBrillianceCtrl int64 `json:"bSupportBrillianceCtrl"` // 亮度调节：0不支持 1支持
	BSupportOnOff          int64 `json:"bSupportOnOff"`          // 开关控制：0不支持 1支持
}

type MPUAirConditionerCapabilities struct {
	BSupportModes                 BSupportModes `json:"bSupportModes"`                 // 支持的模式
	BSupportOnOff                 int64         `json:"bSupportOnOff"`                 // 开关控制：0不支持 1支持
	BSupportTemperatureAdjustment int64         `json:"bSupportTemperatureAdjustment"` // 温度调节：0不支持 1支持
	BSupportWindSpeedAdjustment   int64         `json:"bSupportWindSpeedAdjustment"`   // 风速调节：0不支持 1支持
}
type BSupportModes struct {
	BSupportAutomaticMode       int64 `json:"bSupportAutomaticMode"`       // 自动模式：0不支持 1支持
	BSupportCoolMode            int64 `json:"bSupportCoolMode"`            // 制冷模式：0不支持 1支持
	BSupportHeatMode            int64 `json:"bSupportHeatMode"`            // 制热模式：0不支持 1支持
	BSupportAirDistributionMode int64 `json:"bSupportAirDistributionMode"` // 送风模式：0不支持 1支持
}

type MPUCurtainCapabilities struct {
	BSupportOnOff int64 `json:"bSupportOnOff"` // 开关控制：0不支持 1支持
}
