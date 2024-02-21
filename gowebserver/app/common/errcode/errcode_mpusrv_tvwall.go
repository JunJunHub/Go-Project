// ===================================
// 显控大屏相关错误码定义[121000,121099]
// ===================================

package errcode

//[121000,121099] 大屏相关功能预留错误代码
const (
	codeMpusrvTvWallUnknown = iota + codeMpusrvTvWallStart //121000
)

//[121100,121199] 大屏配置管理错误代码
const (
	codeMpusrvTvWallNameDuplicate         = iota + codeMpusrvTvWallManageStart //121100
	codeMpusrvTvWallRecordingDelForbid                                         //121101
	codeMpusrvTvWallInvalidTvWallId                                            //121102
	codeMpusrvTvWallMultiWinStyleNotExist                                      //121103
)

//[121200,121299] 大屏场景配置管理错误代码
const (
	codeMpusrvTvWallSceneNameDuplicate  = iota + codeMpusrvTvWallSceneStart //121200
	codeMpusrvTvWallSceneInvalidSceneId                                     //121201
	codeMpusrvTvWallSceneSaveErr                                            //121202
	codeMpusrvTvWallSceneLoadFailed                                         //121203
	codeMpusrvTvWallSceneInvalidSceneName
	codeMpusrvTvWallSceneSortIsAlreadyFirst
	codeMpusrvTvWallSceneSortIsAlreadyLast
	codeMpusrvTvWallSceneDefSceneNameUnavailable
)

//[121300,121399] 大屏开窗错误代码
const (
	codeMpusrvTvWallWindowInvalidWinId                                      = iota + codeMpusrvTvWallWindowStart //121300
	codeMpusrvTvWallWindowInvalidSubWinId                                                                        //121301
	codeMpusrvTvWallWindowOpenFailedNoInputSrc                                                                   //121302
	codeMpusrvTvWallWindowOpenFailedCascadeOffline                                                               //121303
	codeMpusrvTvWallWindowOpenFailedInputChnDelete                                                               //121304
	codeMpusrvTvWallWindowOpenFailedInputChnTakeBack                                                             //121305
	codeMpusrvTvWallWindowOpenFailedInputChnOffline                                                              //121306
	codeMpusrvTvWallWindowOpenFailedInputChnUnknown                                                              //121307
	codeMpusrvTvWallWindowOpenFailedCascadeNetInputChnDispatchErr                                                //121308
	codeMpusrvTvWallWindowOpenFailedCascadeAnalogInputAssociationChnUnknown                                      //121309
	codeMpusrvTvWallWindowOpenFailedCascadeAnalogInputDispatchErr                                                //121310
	codeMpusrvTvWallWindowExchangeInputForbidMultiWin                                                            //121311
	codeMpusrvTvWallWindowExchangeInputForbidNotInSameScene                                                      //121312
	codeMpusrvTvWallWindowExchangeFailed                                                                         //121313
	codeMpusrvTvWallWindowInputChnNoPermission                                                                   //121314
	codeMpusrvTvWallWindowOutputChnNoPermission                                                                  //121315
	codeMpusrvTvWallWindowNotSupportAnalogSigChn                                                                 //121316
	codeMpusrvTvWallWindowUpdateFailedNetSigInSync                                                               //121317
	codeMpusrvTvWallWindowUpdateFailedForbidMultiWin
	codeMpusrvTvWallWindowSortIsAlreadyTop
	codeMpusrvTvWallWindowSortIsAlreadyBottom
)

//[121400,121499] 轮循相关错误代码
const (
	codeMpusrvTvWallLoopInvalidLoopCfgId       = iota + codeMpusrvTvWallLoopStart //121400
	codeMpusrvTvWallLoopSceneLoopRunning                                          //121401
	codeMpusrvTvWallLoopWinLoopRunning                                            //121402
	codeMpusrvTvWallLoopInvalidLoopMem                                            //121403
	codeMpusrvTvWallLoopSceneStartNoPermission                                    //121404
	codeMpusrvTvWallLoopSceneLoopNotConfig                                        //121405
	codeMpusrvTvWallLoopSceneLoopMemNotEnough                                     //121406
	codeMpusrvTvWallLoopWinLoopNotConfig                                          //121407
	codeMpusrvTvWallLoopWinLoopMemNotEnough                                       //121408
	codeMpusrvTvWallLoopSceneLoopNotRunning                                       //121409
	codeMpusrvTvWallLoopWinLoopNotRunning                                         //121410
	codeMpusrvTvWallLoopRunningOptDisallow
	//未配置大屏批量通道轮循参数
	codeMpusrvTvWallLoopBatchChnLoopNotConfig
	//大屏批量轮循通道个数不足
	codeMpusrvTvWallLoopBatchChnLoopMemNotEnough
	//未启用大屏批量通道轮循
	codeMpusrvTvWallLoopBatchChnLoopNotRunning
	//大屏批量通道轮循占用窗口不存在
	codeMpusrvTvWallLoopBatchChnLoopWinNotExist
)

//[121500,121599] 轮循通道配置相关错误代码
const (
	_ = iota + codeMpusrvTvWallLoopChnCfgStart //121500
	//轮循通道配置不存在
	codeMpuTvWallLoopChnCfgNotExist
	//轮循通道配置别名重复
	codeMpuTvWallLoopChnCfgAliasDuplicate
	//添加通道信息重复或通道不存在
	codeMpuTvWallLoopChnCfgAddChnParamInvalid
	//通道顺序已经排在首位
	codeMpuTvWallLoopChnIsAlreadyTop
	//通道顺序已经排在末尾
	codeMpuTvWallLoopChnIsAlreadyBottom
	//轮循通道配置正在使用中
	codeMpuTvWallLoopChnCfgInUse
	//轮循通道配置达到上限
	codeMpuTvWallLoopChnCfgExceedsLimit
	//轮循通道个数达到上限
	codeMpuTvWallLoopChnNumExceedsLimit
)

//大屏相关业务错误码对象
var (
	ErrMpusrvTvWallUnknown               = NewMpuapsErr(codeMpusrvTvWallUnknown, "{#ErrMpusrvTvWallUnknown}", nil)
	ErrMpusrvTvWallNameDuplicate         = NewMpuapsErr(codeMpusrvTvWallNameDuplicate, "{#ErrMpusrvTvWallNameDuplicate}", nil)
	ErrMpusrvTvWallInvalidTvWallId       = NewMpuapsErr(codeMpusrvTvWallInvalidTvWallId, "{#ErrMpusrvTvWallInvalidTvWallId}", nil)
	ErrMpusrvTvWallMultiWinStyleNotExist = NewMpuapsErr(codeMpusrvTvWallMultiWinStyleNotExist, "{#ErrMpusrvTvWallMultiWinStyleNotExist}", nil)

	ErrMpusrvTvWallRecordingDelForbid           = NewMpuapsErr(codeMpusrvTvWallRecordingDelForbid, "{#ErrMpusrvTvWallRecordingDelForbid}", nil)
	ErrMpusrvTvWallSceneNameDuplicate           = NewMpuapsErr(codeMpusrvTvWallSceneNameDuplicate, "{#ErrMpusrvTvWallSceneNameDuplicate}", nil)
	ErrMpusrvTvWallSceneInvalidSceneId          = NewMpuapsErr(codeMpusrvTvWallSceneInvalidSceneId, "{#ErrMpusrvTvWallSceneInvalidSceneId}", nil)
	ErrMpusrvTvWallSceneSaveErr                 = NewMpuapsErr(codeMpusrvTvWallSceneSaveErr, "{#ErrMpusrvTvWallSceneSaveErr}", nil)
	ErrMpusrvTvWallSceneLoadFailed              = NewMpuapsErr(codeMpusrvTvWallSceneLoadFailed, "{#ErrMpusrvTvWallSceneLoadFailed}", nil)
	ErrMpusrvTvWallSceneInvalidSceneName        = NewMpuapsErr(codeMpusrvTvWallSceneInvalidSceneName, "{#ErrMpusrvTvWallSceneInvalidSceneName}", nil)
	ErrMpusrvTvWallSceneSortIsAlreadyFirst      = NewMpuapsErr(codeMpusrvTvWallSceneSortIsAlreadyFirst, "{#ErrMpusrvTvWallSceneSortIsAlreadyFirst}", nil)
	ErrMpusrvTvWallSceneSortIsAlreadyLast       = NewMpuapsErr(codeMpusrvTvWallSceneSortIsAlreadyLast, "{#ErrMpusrvTvWallSceneSortIsAlreadyLast}", nil)
	ErrMpusrvTvWallSceneDefSceneNameUnavailable = NewMpuapsErr(codeMpusrvTvWallSceneDefSceneNameUnavailable, "{#ErrMpusrvTvWallSceneDefSceneNameUnavailable}", nil)

	ErrMpusrvTvWallWindowInvalidWinId                                      = NewMpuapsErr(codeMpusrvTvWallWindowInvalidWinId, "{#ErrMpusrvTvWallWindowInvalidWinId}", nil)
	ErrMpusrvTvWallWindowInvalidSubWinId                                   = NewMpuapsErr(codeMpusrvTvWallWindowInvalidSubWinId, "{#ErrMpusrvTvWallWindowInvalidSubWinId}", nil)
	ErrMpusrvTvWallWindowOpenFailedNoInputSrc                              = NewMpuapsErr(codeMpusrvTvWallWindowOpenFailedNoInputSrc, "{#ErrMpusrvTvWallWindowOpenFailedNoInputSrc}", nil)
	ErrMpusrvTvWallWindowOpenFailedCascadeOffline                          = NewMpuapsErr(codeMpusrvTvWallWindowOpenFailedCascadeOffline, "{#ErrMpusrvTvWallWindowOpenWinFailedCascadeOffline}", nil)
	ErrMpusrvTvWallWindowOpenFailedInputChnDelete                          = NewMpuapsErr(codeMpusrvTvWallWindowOpenFailedInputChnDelete, "{#ErrMpusrvTvWallWindowOpenFailedInputChnDelete}", nil)
	ErrMpusrvTvWallWindowOpenFailedInputChnTakeBack                        = NewMpuapsErr(codeMpusrvTvWallWindowOpenFailedInputChnTakeBack, "{#ErrMpusrvTvWallWindowOpenFailedInputChnTakeBack}", nil)
	ErrMpusrvTvWallWindowOpenFailedInputChnOffline                         = NewMpuapsErr(codeMpusrvTvWallWindowOpenFailedInputChnOffline, "{#ErrMpusrvTvWallWindowOpenFailedInputChnOffline}", nil)
	ErrMpusrvTvWallWindowOpenFailedInputChnUnknown                         = NewMpuapsErr(codeMpusrvTvWallWindowOpenFailedInputChnUnknown, "{#ErrMpusrvTvWallWindowOpenFailedInputChnUnknown}", nil)
	ErrMpusrvTvWallWindowOpenFailedCascadeNetInputChnDispatchErr           = NewMpuapsErr(codeMpusrvTvWallWindowOpenFailedCascadeNetInputChnDispatchErr, "{#ErrMpusrvTvWallWindowOpenFailedCascadeNetInputChnDispatchErr}", nil)
	ErrMpusrvTvWallWindowOpenFailedCascadeAnalogInputAssociationChnUnknown = NewMpuapsErr(codeMpusrvTvWallWindowOpenFailedCascadeAnalogInputAssociationChnUnknown, "{#ErrMpusrvTvWallWindowOpenFailedCascadeAnalogInputAssociationChnUnknown}", nil)
	ErrMpusrvTvWallWindowOpenFailedCascadeAnalogInputDispatchErr           = NewMpuapsErr(codeMpusrvTvWallWindowOpenFailedCascadeAnalogInputDispatchErr, "{#ErrMpusrvTvWallWindowOpenFailedCascadeAnalogInputDispatchErr}", nil)
	ErrMpusrvTvWallWindowExchangeInputForbidMultiWin                       = NewMpuapsErr(codeMpusrvTvWallWindowExchangeInputForbidMultiWin, "{#ErrMpusrvTvWallWindowExchangeInputForbidMultiWin}", nil)
	ErrMpusrvTvWallWindowExchangeInputForbidNotInSameScene                 = NewMpuapsErr(codeMpusrvTvWallWindowExchangeInputForbidNotInSameScene, "{#ErrMpusrvTvWallWindowExchangeInputForbidNotInSameScene}", nil)
	ErrMpusrvTvWallWindowExchangeFailed                                    = NewMpuapsErr(codeMpusrvTvWallWindowExchangeFailed, "{#ErrMpusrvTvWallWindowExchangeFailed}", nil)
	ErrMpusrvTvWallWindowInputChnNoPermission                              = NewMpuapsErr(codeMpusrvTvWallWindowInputChnNoPermission, "{#ErrMpusrvTvWallWindowInputChnNoPermission}", nil)
	ErrMpusrvTvWallWindowOutputChnNoPermission                             = NewMpuapsErr(codeMpusrvTvWallWindowOutputChnNoPermission, "{#ErrMpusrvTvWallWindowOutputChnNoPermission}", nil)
	ErrMpusrvTvWallWindowNotSupportAnalogSigChn                            = NewMpuapsErr(codeMpusrvTvWallWindowNotSupportAnalogSigChn, "{#ErrMpusrvTvWallWindowNotSupportAnalogSigChn}", nil)
	ErrMpusrvTvWallWindowUpdateFailedNetSigInSync                          = NewMpuapsErr(codeMpusrvTvWallWindowUpdateFailedNetSigInSync, "{#ErrMpusrvTvWallWindowUpdateFailedNetSigInSync}", nil)
	ErrMpusrvTvWallWindowUpdateFailedForbidMultiWin                        = NewMpuapsErr(codeMpusrvTvWallWindowUpdateFailedForbidMultiWin, "{#ErrMpusrvTvWallWindowUpdateFailedForbidMultiWin}", nil)
	ErrMpusrvTvWallWindowSortIsAlreadyTop                                  = NewMpuapsErr(codeMpusrvTvWallWindowSortIsAlreadyTop, "{#ErrMpusrvTvWallWindowSortIsAlreadyTop}", nil)
	ErrMpusrvTvWallWindowSortIsAlreadyBottom                               = NewMpuapsErr(codeMpusrvTvWallWindowSortIsAlreadyBottom, "{#ErrMpusrvTvWallWindowSortIsAlreadyBottom}", nil)

	ErrMpusrvTvWallLoopInvalidLoopCfgId         = NewMpuapsErr(codeMpusrvTvWallLoopInvalidLoopCfgId, "{#ErrMpusrvTvWallLoopInvalidLoopCfgId}", nil)
	ErrMpusrvTvWallLoopSceneLoopRunning         = NewMpuapsErr(codeMpusrvTvWallLoopSceneLoopRunning, "{#ErrMpusrvTvWallLoopSceneLoopRunning}", nil)
	ErrMpusrvTvWallLoopWinLoopRunning           = NewMpuapsErr(codeMpusrvTvWallLoopWinLoopRunning, "{#ErrMpusrvTvWallLoopWinLoopRunning}", nil)
	ErrMpusrvTvWallLoopInvalidLoopMem           = NewMpuapsErr(codeMpusrvTvWallLoopInvalidLoopMem, "{#ErrMpusrvTvWallLoopInvalidLoopMem}", nil)
	ErrMpusrvTvWallLoopSceneStartNoPermission   = NewMpuapsErr(codeMpusrvTvWallLoopSceneStartNoPermission, "{#ErrMpusrvTvWallLoopSceneStartNoPermission}", nil)
	ErrMpusrvTvWallLoopSceneLoopNotConfig       = NewMpuapsErr(codeMpusrvTvWallLoopSceneLoopNotConfig, "{#ErrMpusrvTvWallLoopSceneLoopNotConfig}", nil)
	ErrMpusrvTvWallLoopSceneLoopMemNotEnough    = NewMpuapsErr(codeMpusrvTvWallLoopSceneLoopMemNotEnough, "{#ErrMpusrvTvWallLoopSceneLoopMemNotEnough}", nil)
	ErrMpusrvTvWallLoopWinLoopNotConfig         = NewMpuapsErr(codeMpusrvTvWallLoopWinLoopNotConfig, "{#ErrMpusrvTvWallLoopWinLoopNotConfig}", nil)
	ErrMpusrvTvWallLoopWinLoopMemNotEnough      = NewMpuapsErr(codeMpusrvTvWallLoopWinLoopMemNotEnough, "{#ErrMpusrvTvWallLoopWinLoopMemNotEnough}", nil)
	ErrMpusrvTvWallLoopSceneLoopNotRunning      = NewMpuapsErr(codeMpusrvTvWallLoopSceneLoopNotRunning, "{#ErrMpusrvTvWallLoopSceneLoopNotRunning}", nil)
	ErrMpusrvTvWallLoopWinLoopNotRunning        = NewMpuapsErr(codeMpusrvTvWallLoopWinLoopNotRunning, "{#ErrMpusrvTvWallLoopWinLoopNotRunning}", nil)
	ErrMpusrvTvWallLoopRunningOptDisallow       = NewMpuapsErr(codeMpusrvTvWallLoopRunningOptDisallow, "{#ErrMpusrvTvWallLoopRunningOptDisallow}", nil)
	ErrMpusrvTvWallLoopBatchChnLoopNotConfig    = NewMpuapsErr(codeMpusrvTvWallLoopBatchChnLoopNotConfig, "{#ErrMpusrvTvWallLoopBatchChnLoopNotConfig}", nil)
	ErrMpusrvTvWallLoopBatchChnLoopMemNotEnough = NewMpuapsErr(codeMpusrvTvWallLoopBatchChnLoopMemNotEnough, "{#ErrMpusrvTvWallLoopBatchChnLoopMemNotEnough}", nil)
	ErrMpusrvTvWallLoopBatchChnLoopNotRunning   = NewMpuapsErr(codeMpusrvTvWallLoopBatchChnLoopNotRunning, "{#ErrMpusrvTvWallLoopBatchChnLoopNotRunning}", nil)
	ErrMpusrvTvWallLoopBatchChnLoopWinNotExist  = NewMpuapsErr(codeMpusrvTvWallLoopBatchChnLoopWinNotExist, "{#ErrMpusrvTvWallLoopBatchChnLoopWinNotExist}", nil)

	ErrMpuTvWallLoopChnCfgNotExist           = NewMpuapsErr(codeMpuTvWallLoopChnCfgNotExist, "{#ErrMpuTvWallLoopChnCfgNotExist}", nil)
	ErrMpuTvWallLoopChnCfgAliasDuplicate     = NewMpuapsErr(codeMpuTvWallLoopChnCfgAliasDuplicate, "{#ErrMpuTvWallLoopChnCfgAliasDuplicate}", nil)
	ErrMpuTvWallLoopChnCfgAddChnParamInvalid = NewMpuapsErr(codeMpuTvWallLoopChnCfgAddChnParamInvalid, "{#ErrMpuTvWallLoopChnCfgAddChnParamInvalid}", nil)
	ErrMpuTvWallLoopChnIsAlreadyTop          = NewMpuapsErr(codeMpuTvWallLoopChnIsAlreadyTop, "{#ErrMpuTvWallLoopChnIsAlreadyTop}", nil)
	ErrMpuTvWallLoopChnIsAlreadyBottom       = NewMpuapsErr(codeMpuTvWallLoopChnIsAlreadyBottom, "{#ErrMpuTvWallLoopChnIsAlreadyBottom}", nil)
	ErrMpuTvWallLoopChnCfgInUse              = NewMpuapsErr(codeMpuTvWallLoopChnCfgInUse, "{#ErrMpuTvWallLoopChnCfgInUse}", nil)
	ErrMpuTvWallLoopChnCfgExceedsLimit       = NewMpuapsErr(codeMpuTvWallLoopChnCfgExceedsLimit, "{#ErrMpuTvWallLoopChnCfgExceedsLimit}", nil)
	ErrMpuTvWallLoopChnNumExceedsLimit       = NewMpuapsErr(codeMpuTvWallLoopChnNumExceedsLimit, "{#ErrMpuTvWallLoopChnNumExceedsLimit}", nil)
)
