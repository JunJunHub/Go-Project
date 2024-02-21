// =====================================
// 显控环境设备相关错误码定义[124000,124999]
// =====================================

package errcode

//[126000,126299] 环境设备控制器相关预留错误码
const (
	codeMpuApplianceCtrlModuleUnknown = iota + codeMpusrvApplianceStart //126000
)

//[126300,126999] 环境设备管理相关错误码
const (
	//环境设备分组名称重复
	codeMpuApplianceSchemeGroupNameDuplicate = iota + codeMpusrvApplianceManageStart //126300
	//环境设备分组非空
	codeMpuApplianceSchemeGroupNoEmpty
	//继电器名称重复
	codeMpuApplianceRelayNameDuplicate
	//继电器开关数目无效
	codeMpuApplianceRelaySwitchOutNumInvalid
	//一个盒子或通道不可被同时绑定多个继电器
	codeMpuApplianceRelayTRxIpOrChannelIdDuplicate
	//继电器已被绑定设备,不支持该操作
	codeMpuApplianceRelayUsed
	//继电器不存在
	codeMpuApplianceRelayNoExist
	//环境设备名称重复
	codeMpuApplianceDevNameDuplicate
	//环境设备名称为空
	codeMpuApplianceDevNameEmpty
	//环境设备未配置控制模块
	codeMpuApplianceDevCtrlModuleEmpty
	//控制模块重复或已被其他环境设备配置使用
	codeMpuApplianceDevCtrlModuleDuplicate
	//环境设备不存在
	codeMpuApplianceDevNoExist
)

var (
	ErrMpuApplianceSchemeGroupNameDuplicate       = NewMpuapsErr(codeMpuApplianceSchemeGroupNameDuplicate, "{#ErrMpuApplianceSchemeGroupNameDuplicate}", nil)
	ErrMpuApplianceSchemeGroupNoEmpty             = NewMpuapsErr(codeMpuApplianceSchemeGroupNoEmpty, "{#ErrMpuApplianceSchemeGroupNoEmpty}", nil)
	ErrMpuApplianceRelayNameDuplicate             = NewMpuapsErr(codeMpuApplianceRelayNameDuplicate, "{#ErrMpuApplianceRelayNameDuplicate}", nil)
	ErrMpuApplianceRelaySwitchOutNumInvalid       = NewMpuapsErr(codeMpuApplianceRelaySwitchOutNumInvalid, "{#ErrMpuApplianceRelaySwitchOutNumInvalid}", nil)
	ErrMpuApplianceRelayTRxIpOrChannelIdDuplicate = NewMpuapsErr(codeMpuApplianceRelayTRxIpOrChannelIdDuplicate, "{#ErrMpuApplianceRelayTRxIpOrChannelIdDuplicate}", nil)
	ErrMpuApplianceRelayUsed                      = NewMpuapsErr(codeMpuApplianceRelayUsed, "{#ErrMpuApplianceRelayUsed}", nil)
	ErrMpuApplianceRelayNoExist                   = NewMpuapsErr(codeMpuApplianceRelayNoExist, "{#ErrMpuApplianceRelayNoExist}", nil)
	ErrMpuApplianceDevNameDuplicate               = NewMpuapsErr(codeMpuApplianceDevNameDuplicate, "{#ErrMpuApplianceDevNameDuplicate}", nil)
	ErrMpuApplianceDevNameEmpty                   = NewMpuapsErr(codeMpuApplianceDevNameEmpty, "{#ErrMpuApplianceDevNameEmpty#}", nil)
	ErrMpuApplianceDevCtrlModuleEmpty             = NewMpuapsErr(codeMpuApplianceDevCtrlModuleEmpty, "{#ErrMpuApplianceDevCtrlModuleEmpty}", nil)
	ErrMpuApplianceDevCtrlModuleDuplicate         = NewMpuapsErr(codeMpuApplianceDevCtrlModuleDuplicate, "{#ErrMpuApplianceDevCtrlModuleDuplicate}", nil)
	ErrMpuApplianceDevNoExist                     = NewMpuapsErr(codeMpuApplianceDevNoExist, "{#ErrMpuApplianceDevNoExist}", nil)
)
