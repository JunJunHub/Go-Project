// =====================================
// 显控会议管理相关错误码定义[127000,127999]
// =====================================

package errcode

//[127000,127099] 会议管理相关功能预留错误代码
const (
	_ = iota + codeMpusrvMeetManageStart
)

//[127100,127199] 会议平台管理错误代码
const (
	_ = iota + codeMpusrvMeetPlatformStart
	//会议平台不存在
	codeMpuMeetPlatformNotExist
	//会议平台名称重复
	codeMpuMeetPlatformAliasDuplicate
	//会议平台IP重复
	codeMpuMeetPlatformIPDuplicate
	//不可删除正在使用的会议平台
	codeMpuMeetPlatformInuseDelForbid
	//不可删除正在使用的会议平台账户
	codeMpuMeetPlatformUserInuseDelForbid
	//不可修改正在使用的会议平台账户
	codeMpuMeetPlatformUserInuseUpdateForbid
	//会议平台用户名重复
	codeMpuMeetPlatformUsernameDuplicate
	//会议平台未连接
	codeMpuMeetPlatformUnConnect
)

//[127200,127299]  会议终端管理错误代码
const (
	_ = iota + codeMpusrvMeetTerminalStart
)

//[127300,127399] 会议模板管理错误代码
const (
	_ = iota + codeMpusrvMeetTemplateStart
)

//[127400,127499] 会议调度分组错误代码
const (
	_ = iota + codeMpusrvMeetDispatchGroupStart
	//会议分组名重复
	codeMpuMeetDispatchGroupNameDuplicate
	//会议分组不存在
	codeMpuMeetDispatchGroupNotExist
	//会议分组不为空,禁止删除
	codeMpuMeetDispatchGroupMemNotNullDelForbid
)

//[127500,127999] 会议调度管理错误代码
const (
	_ = iota + codeMpusrvMeetDispatchStart
	//会议调度不存在
	codeMpuMeetDispatchNotExist
	//会议调度名称重复
	codeMpuMeetDispatchNameDuplicate
	//会议调度名称无效
	codeMpuMeetDispatchNameInvalid
	//会议调度绑定矩阵预案不存在
	codeMpuMeetDispatchBindMatrixSchemeNotExist
	//会议调度绑定会议模板不存在
	codeMpuMeetDispatchBindMeetTemplateNotExist
	//矩阵预案已被其它会议调度绑定
	codeMpuMeetDispatchBindMatrixSchemeDuplicate
	//会议调度未绑定MCU会议模板或绑定模板已删除
	codeMpuMeetDispatchUnBindMeetTemplate
	//会议调度关联MCU会议已关闭
	codeMpuMeetDispatchBindMeetingStopped
	//不可编辑加载中的会议调度配置
	codeMpuMeetDispatchLoadingEditForbid
	//不可删除加载中的会议调度配置
	codeMpuMeetDispatchLoadingDelForbid
	//会议调度已启用
	codeMpuMeetDispatchLoading
	//会议调度未启用
	codeMpuMeetDispatchUnLoading
	//会议调度绑定矩阵预案无权限
	codeMpuMeetDispatchBindMatrixSchemeNoPermission
)

var (
	ErrMpuMeetPlatformNotExist              = NewMpuapsErr(codeMpuMeetPlatformNotExist, "{#ErrMpuMeetPlatformNotExist}", nil)
	ErrMpuMeetPlatformAliasDuplicate        = NewMpuapsErr(codeMpuMeetPlatformAliasDuplicate, "{#ErrMpuMeetPlatformAliasDuplicate}", nil)
	ErrMpuMeetPlatformIPDuplicate           = NewMpuapsErr(codeMpuMeetPlatformIPDuplicate, "{#ErrMpuMeetPlatformIPDuplicate}", nil)
	ErrMpuMeetPlatformInuseDelForbid        = NewMpuapsErr(codeMpuMeetPlatformInuseDelForbid, "{#ErrMpuMeetPlatformInuseDelForbid}", nil)
	ErrMpuMeetPlatformUserInuseDelForbid    = NewMpuapsErr(codeMpuMeetPlatformUserInuseDelForbid, "{#ErrMpuMeetPlatformUserInuseDelForbid}", nil)
	ErrMpuMeetPlatformUserInuseUpdateForbid = NewMpuapsErr(codeMpuMeetPlatformUserInuseUpdateForbid, "{#ErrMpuMeetPlatformUserInuseUpdateForbid}", nil)
	ErrMpuMeetPlatformUsernameDuplicate     = NewMpuapsErr(codeMpuMeetPlatformUsernameDuplicate, "{#ErrMpuMeetPlatformUsernameDuplicate}", nil)
	ErrMpuMeetPlatformUnConnect             = NewMpuapsErr(codeMpuMeetPlatformUnConnect, "{#ErrMpuMeetPlatformUnConnect}", nil)

	ErrMpuMeetDispatchGroupNameDuplicate       = NewMpuapsErr(codeMpuMeetDispatchGroupNameDuplicate, "{#ErrMpuMeetDispatchGroupNameDuplicate}", nil)
	ErrMpuMeetDispatchGroupNotExist            = NewMpuapsErr(codeMpuMeetDispatchGroupNotExist, "{#ErrMpuMeetDispatchGroupNotExist}", nil)
	ErrMpuMeetDispatchGroupMemNotNullDelForbid = NewMpuapsErr(codeMpuMeetDispatchGroupMemNotNullDelForbid, "{#ErrMpuMeetDispatchGroupMemNotNullDelForbid}", nil)

	ErrMpuMeetDispatchNotExist                     = NewMpuapsErr(codeMpuMeetDispatchNotExist, "{#ErrMpuMeetDispatchNotExist}", nil)
	ErrMpuMeetDispatchNameDuplicate                = NewMpuapsErr(codeMpuMeetDispatchNameDuplicate, "{#ErrMpuMeetDispatchNameDuplicate}", nil)
	ErrMpuMeetDispatchNameInvalid                  = NewMpuapsErr(codeMpuMeetDispatchNameInvalid, "{#ErrMpuMeetDispatchNameInvalid}", nil)
	ErrMpuMeetDispatchBindMatrixSchemeNotExist     = NewMpuapsErr(codeMpuMeetDispatchBindMatrixSchemeNotExist, "{#ErrMpuMeetDispatchBindMatrixSchemeNotExist}", nil)
	ErrMpuMeetDispatchBindMeetTemplateNotExist     = NewMpuapsErr(codeMpuMeetDispatchBindMeetTemplateNotExist, "{#ErrMpuMeetDispatchBindMeetTemplateNotExist}", nil)
	ErrMpuMeetDispatchBindMatrixSchemeDuplicate    = NewMpuapsErr(codeMpuMeetDispatchBindMatrixSchemeDuplicate, "{#ErrMpuMeetDispatchBindMatrixSchemeDuplicate}", nil)
	ErrMpuMeetDispatchUnBindMeetTemplate           = NewMpuapsErr(codeMpuMeetDispatchUnBindMeetTemplate, "{#ErrMpuMeetDispatchUnBindMeetTemplate}", nil)
	ErrMpuMeetDispatchBindMeetingStopped           = NewMpuapsErr(codeMpuMeetDispatchBindMeetingStopped, "{#ErrMpuMeetDispatchBindMeetingStopped}", nil)
	ErrMpuMeetDispatchLoadingEditForbid            = NewMpuapsErr(codeMpuMeetDispatchLoadingEditForbid, "{#ErrMpuMeetDispatchLoadingEditForbid}", nil)
	ErrMpuMeetDispatchLoadingDelForbid             = NewMpuapsErr(codeMpuMeetDispatchLoadingDelForbid, "{#ErrMpuMeetDispatchLoadingDelForbid}", nil)
	ErrMpuMeetDispatchLoading                      = NewMpuapsErr(codeMpuMeetDispatchLoading, "{#ErrMpuMeetDispatchLoading}", nil)
	ErrMpuMeetDispatchUnLoading                    = NewMpuapsErr(codeMpuMeetDispatchUnLoading, "{#ErrMpuMeetDispatchUnLoading}", nil)
	ErrMpuMeetDispatchBindMatrixSchemeNoPermission = NewMpuapsErr(codeMpuMeetDispatchBindMatrixSchemeNoPermission, "{#ErrMpuMeetDispatchBindMatrixSchemeNoPermission}", nil)
)
