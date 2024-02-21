// =====================================
// 显控矩阵调度相关错误码定义[124000,124999]
// =====================================

package errcode

//[124000,124099] 矩阵调度相关功能预留错误代码
const (
	codeMpuMatrixDispatchUnknown = iota + codeMpusrvMatrixDispatchStart //124000
)

//[124100,124299] 矩阵调度相关功能错误代码
const (
	codeMpuMatrixDispatchSchemeGroupNameDuplicate       = iota + codeMpusrvMatrixDispatchManageStart //124100
	codeMpuMatrixDispatchSchemeGroupNotExist                                                         //124101
	codeMpuMatrixDispatchSchemeNameDuplicate                                                         //124102
	codeMpuMatrixDispatchSchemeNotExist                                                              //124103
	codeMpuMatrixDispatchSchemeLoadingEditForbid                                                     //124104
	codeMpuMatrixDispatchSchemeLoadingDelForbid                                                      //124105
	codeMpuMatrixDispatchSchemeUnLoading                                                             //124106
	codeMpuMatrixDispatchSchemeDispatchDestChnDuplicate                                              //124107

	codeMpuMatrixDispatchNotExist                       //124108
	codeMpuMatrixDispatchDestChnOccupied                //124109
	codeMpuMatrixDispatchDestChnNotExist                //124110
	codeMpuMatrixDispatchDestChnIsNull                  //124111
	codeMpuMatrixDispatchDestChnNotSupportCascade       //124112
	codeMpuMatrixDispatchSrcChnNotExist                 //124113
	codeMpuMatrixDispatchSrcDispatchFailed              //124114
	codeMpuMatrixDispatchVideoNotSupportDestChnType     //124115
	codeMpuMatrixDispatchVideoNotSupportSrcChnType      //124116
	codeMpuMatrixDispatchAudioNotSupportDestChnType     //124117
	codeMpuMatrixDispatchAudioNotSupportSrcChnType      //124118
	codeMpuMatrixDispatchDefaultSrcCfgNotExist          //124119
	codeMpuMatrixDispatchSwitchFailed                   //124120
	codeMpuMatrixDispatchParameterInvalid               //124121
	codeMpuMatrixDispatchSchemeGroupMemNotNullDelForbid //124122
	codeMpuMatrixDispatchDestChnModify                  //124123
	codeMpuMatrixDispatchTypeModify                     //124124
	codeMpuMatrixDispatchEnabled                        //124125
	codeMpuMatrixDispatchDisabled                       //124126
	codeMpuMatrixDispatchTvWallSrcEmpty                 //124127
	codeMpuMatrixDispatchTvWallSrcModify                //124128
	codeMpuMatrixDispatchTvWallNotSupport               //124129
	codeMpuMatrixDispatchHasData                        //124130
	codeMpuMatrixDispatchDstChnNoPermission             //124131
	codeMpuMatrixDispatchSrcChnNoPermission             //124132
	codeMpuMatrixDispatchNoPermissionEditForbid         //124133

	codeMpuMatrixDispatchFailed
)

//[124300,124399] 矩阵调度音视频源绑定配置错误码
const (
	_ = iota + codeMpusrvMatrixAVChnBindCfgStart //124200
	//配置不存在
	codeMpuMatrixAVChnBindCfgNotExist
	//配置名称重复
	codeMpuMatrixAVChnBindCfgAliasDuplicate
	//重复绑定音视频输入源
	codeMpuMatrixAVChnBindDuplicate
	//音频输入源不在线
	codeMpuMatrixAVChnBindCfgAudioChnOffline
	//视频输入源不在线
	codeMpuMatrixAVChnBindCfgVideoChnOffline
	//音视频输入源不在线
	codeMpuMatrixAVChnBindCfgChnOffline
	//输入源无权限[音频或视频源无权限即配置无权限]
	codeMpuMatrixAVChnBindCfgNoPermission
)

var (
	ErrMpuMatrixDispatchUnknown                        = NewMpuapsErr(codeMpuMatrixDispatchUnknown, "{#ErrMpuMatrixDispatchUnknown}", nil)
	ErrMpuMatrixDispatchSchemeGroupNameDuplicate       = NewMpuapsErr(codeMpuMatrixDispatchSchemeGroupNameDuplicate, "{#ErrMpuMatrixDispatchSchemeGroupNameDuplicate}", nil)
	ErrMpuMatrixDispatchSchemeGroupNotExist            = NewMpuapsErr(codeMpuMatrixDispatchSchemeGroupNotExist, "{#ErrMpuMatrixDispatchSchemeGroupNotExist}", nil)
	ErrMpuMatrixDispatchSchemeNameDuplicate            = NewMpuapsErr(codeMpuMatrixDispatchSchemeNameDuplicate, "{#ErrMpuMatrixDispatchSchemeNameDuplicate}", nil)
	ErrMpuMatrixDispatchSchemeNotExist                 = NewMpuapsErr(codeMpuMatrixDispatchSchemeNotExist, "{#ErrMpuMatrixDispatchSchemeNotExist}", nil)
	ErrMpuMatrixDispatchSchemeLoadingEditForbid        = NewMpuapsErr(codeMpuMatrixDispatchSchemeLoadingEditForbid, "{#ErrMpuMatrixDispatchSchemeLoadingEditForbid}", nil)
	ErrMpuMatrixDispatchSchemeLoadingDelForbid         = NewMpuapsErr(codeMpuMatrixDispatchSchemeLoadingDelForbid, "{#ErrMpuMatrixDispatchSchemeLoadingDelForbid}", nil)
	ErrMpuMatrixDispatchSchemeUnLoading                = NewMpuapsErr(codeMpuMatrixDispatchSchemeUnLoading, "{#ErrMpuMatrixDispatchSchemeUnLoading}", nil)
	ErrMpuMatrixDispatchSchemeDispatchDestChnDuplicate = NewMpuapsErr(codeMpuMatrixDispatchSchemeDispatchDestChnDuplicate, "{#ErrMpuMatrixDispatchSchemeDispatchDestChnDuplicate}", nil)

	ErrMpuMatrixDispatchNotExist                       = NewMpuapsErr(codeMpuMatrixDispatchNotExist, "{#ErrMpuMatrixDispatchNotExist}", nil)
	ErrMpuMatrixDispatchDestChnOccupied                = NewMpuapsErr(codeMpuMatrixDispatchDestChnOccupied, "{#ErrMpuMatrixDispatchDestChnOccupied}", nil)
	ErrMpuMatrixDispatchDestChnNotExist                = NewMpuapsErr(codeMpuMatrixDispatchDestChnNotExist, "{#ErrMpuMatrixDispatchDestChnNotExist}", nil)
	ErrMpuMatrixDispatchDestChnIsNull                  = NewMpuapsErr(codeMpuMatrixDispatchDestChnIsNull, "{#ErrMpuMatrixDispatchDestChnIsNull}", nil)
	ErrMpuMatrixDispatchDestChnNotSupportCascade       = NewMpuapsErr(codeMpuMatrixDispatchDestChnNotSupportCascade, "{#ErrMpuMatrixDispatchDestChnNotSupportCascade}", nil)
	ErrMpuMatrixDispatchSrcChnNotExist                 = NewMpuapsErr(codeMpuMatrixDispatchSrcChnNotExist, "{#ErrMpuMatrixDispatchSrcChnNotExist}", nil)
	ErrMpuMatrixDispatchSrcDispatchFailed              = NewMpuapsErr(codeMpuMatrixDispatchSrcDispatchFailed, "{#ErrMpuMatrixDispatchSrcDispatchFailed}", nil)
	ErrMpuMatrixDispatchVideoNotSupportDestChnType     = NewMpuapsErr(codeMpuMatrixDispatchVideoNotSupportDestChnType, "{#ErrMpuMatrixDispatchVideoNotSupportDestChnType}", nil)
	ErrMpuMatrixDispatchVideoNotSupportSrcChnType      = NewMpuapsErr(codeMpuMatrixDispatchVideoNotSupportSrcChnType, "{#ErrMpuMatrixDispatchVideoNotSupportSrcChnType}", nil)
	ErrMpuMatrixDispatchAudioNotSupportDestChnType     = NewMpuapsErr(codeMpuMatrixDispatchAudioNotSupportDestChnType, "{#ErrMpuMatrixDispatchAudioNotSupportDestChnType}", nil)
	ErrMpuMatrixDispatchAudioNotSupportSrcChnType      = NewMpuapsErr(codeMpuMatrixDispatchAudioNotSupportSrcChnType, "{#ErrMpuMatrixDispatchAudioNotSupportSrcChnType}", nil)
	ErrMpuMatrixDispatchDefaultSrcCfgNotExist          = NewMpuapsErr(codeMpuMatrixDispatchDefaultSrcCfgNotExist, "{#ErrMpuMatrixDispatchDefaultSrcCfgNotExist}", nil)
	ErrMpuMatrixDispatchSwitchFailed                   = NewMpuapsErr(codeMpuMatrixDispatchSwitchFailed, "{#ErrMpuMatrixDispatchSwitchFailed}", nil)
	ErrMpuMatrixDispatchParameterInvalid               = NewMpuapsErr(codeMpuMatrixDispatchParameterInvalid, "{#ErrMpuMatrixDispatchParameterInvalid}", nil)
	ErrMpuMatrixDispatchSchemeGroupMemNotNullDelForbid = NewMpuapsErr(codeMpuMatrixDispatchSchemeGroupMemNotNullDelForbid, "{#ErrMpuMatrixDispatchSchemeGroupMemNotNullDelForbid}", nil)
	ErrMpuMatrixDispatchDestChnModify                  = NewMpuapsErr(codeMpuMatrixDispatchDestChnModify, "{#ErrMpuMatrixDispatchDestChnModify}", nil)
	ErrMpuMatrixDispatchTypeModify                     = NewMpuapsErr(codeMpuMatrixDispatchTypeModify, "{#ErrMpuMatrixDispatchTypeModify}", nil)
	ErrMpuMatrixDispatchEnabled                        = NewMpuapsErr(codeMpuMatrixDispatchEnabled, "{#ErrMpuMatrixDispatchEnabled}", nil)
	ErrMpuMatrixDispatchDisabled                       = NewMpuapsErr(codeMpuMatrixDispatchDisabled, "{#ErrMpuMatrixDispatchDisabled}", nil)
	ErrMpuMatrixDispatchTvwallSrcEmpty                 = NewMpuapsErr(codeMpuMatrixDispatchTvWallSrcEmpty, "{#ErrMpuMatrixDispatchTvwallSrcEmpty}", nil)
	ErrMpuMatrixDispatchTvwallSrcModify                = NewMpuapsErr(codeMpuMatrixDispatchTvWallSrcModify, "{#ErrMpuMatrixDispatchTvwallSrcModify}", nil)
	ErrMpuMatrixDispatchTvWallNotSupport               = NewMpuapsErr(codeMpuMatrixDispatchTvWallNotSupport, "{#ErrMpuMatrixDispatchTvWallNotSupport}", nil)
	ErrMpuMatrixDispatchHasData                        = NewMpuapsErr(codeMpuMatrixDispatchHasData, "{#ErrMpuMatrixDispatchHasData}", nil)
	ErrMpuMatrixDispatchDstChnNoPermission             = NewMpuapsErr(codeMpuMatrixDispatchDstChnNoPermission, "{#ErrMpuMatrixDispatchDstChnNoPermission}", nil)
	ErrMpuMatrixDispatchSrcChnNoPermission             = NewMpuapsErr(codeMpuMatrixDispatchSrcChnNoPermission, "{#ErrMpuMatrixDispatchSrcChnNoPermission}", nil)
	ErrMpuMatrixDispatchNoPermissionEditForbid         = NewMpuapsErr(codeMpuMatrixDispatchNoPermissionEditForbid, "{#ErrMpuMatrixDispatchNoPermissionEditForbid}", nil)
	ErrMpuMatrixDispatchFailed                         = NewMpuapsErr(codeMpuMatrixDispatchFailed, "{#ErrMpuMatrixDispatchFailed}", nil)
)

var (
	ErrMpuMatrixAVChnBindCfgNotExist        = NewMpuapsErr(codeMpuMatrixAVChnBindCfgNotExist, "{#ErrMpuMatrixAVChnBindCfgNotExist}", nil)
	ErrMpuMatrixAVChnBindCfgAliasDuplicate  = NewMpuapsErr(codeMpuMatrixAVChnBindCfgAliasDuplicate, "{#ErrMpuMatrixAVChnBindCfgAliasDuplicate}", nil)
	ErrMpuMatrixAVChnBindDuplicate          = NewMpuapsErr(codeMpuMatrixAVChnBindDuplicate, "{#ErrMpuMatrixAVChnBindDuplicate}", nil)
	ErrMpuMatrixAVChnBindCfgAudioChnOffline = NewMpuapsErr(codeMpuMatrixAVChnBindCfgAudioChnOffline, "{#ErrMpuMatrixAVChnBindCfgAudioChnOffline}", nil)
	ErrMpuMatrixAVChnBindCfgVideoChnOffline = NewMpuapsErr(codeMpuMatrixAVChnBindCfgVideoChnOffline, "{#ErrMpuMatrixAVChnBindCfgVideoChnOffline}", nil)
	ErrMpuMatrixAVChnBindCfgChnOffline      = NewMpuapsErr(codeMpuMatrixAVChnBindCfgChnOffline, "{#ErrMpuMatrixAVChnBindCfgChnOffline}", nil)
	ErrMpuMatrixAVChnBindCfgNoPermission    = NewMpuapsErr(codeMpuMatrixAVChnBindCfgNoPermission, "{#ErrMpuMatrixAVChnBindCfgNoPermission}", nil)
)
