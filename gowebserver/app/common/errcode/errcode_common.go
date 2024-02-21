// =================================
// 通用错误码定义[100000,109999]
// =================================

package errcode

//Http状态码
const (
	codeHttpOk                  = 200
	codeHttpBadRequest          = 400
	codeHttpUnauthorized        = 401
	codeHttpForbidden           = 403
	codeHttpNotFount            = 404
	codeHttpInternalServerError = 500
	codeHttpServiceUnavailable  = 503
	codeHttpGatewayTimeout      = 504
)

//公共错误代码枚举
const (
	codeOk                         = 0
	codeCommonInternalError        = iota + codeCommonStart //100001
	codeCommonValidationFailed                              //100002
	codeCommonDbOperationError                              //100003
	codeCommonInvalidParameter                              //100004
	codeCommonParsingParameter                              //100005
	codeCommonInvalidOperation                              //100006
	codeCommonInvalidConfiguration                          //100007
	codeCommonMissingConfiguration                          //100008
	codeCommonNotImplemented                                //100009
	codeCommonNotSupported                                  //100010
	codeCommonOperationFailed                               //100011
	codeCommonServerBusy                                    //100012
	codeCommonUnknown                                       //100013
	codeCommonInvalidRequest                                //100014
	codeCommonHttpServerStartErr                            //100015
	codeCommonDBConnectErr                                  //100016
	codeCommonDBLoading                                     //100017
	codeCommonSyncMPUResources                              //100018
	codeCommonSyncReduMPURes                                //100019
)

//公共错误码对象
var (
	ErrHttpOk                  = NewMpuapsErr(codeHttpOk, "{#codeHttpOk}", nil)
	ErrHttpBadRequest          = NewMpuapsErr(codeHttpBadRequest, "{#codeHttpBadRequest}", nil)
	ErrHttpUnauthorized        = NewMpuapsErr(codeHttpUnauthorized, "{#codeHttpUnauthorized}", nil)
	ErrHttpForbidden           = NewMpuapsErr(codeHttpForbidden, "{#codeHttpForbidden}", nil)
	ErrHttpNotFount            = NewMpuapsErr(codeHttpNotFount, "{#codeHttpNotFount}", nil)
	ErrHttpInternalServerError = NewMpuapsErr(codeHttpInternalServerError, "{#codeHttpInternalServerError}", nil)
	ErrHttpServiceUnavailable  = NewMpuapsErr(codeHttpServiceUnavailable, "{#codeHttpServiceUnavailable}", nil)
	ErrHttpGatewayTimeout      = NewMpuapsErr(codeHttpGatewayTimeout, "{#codeHttpGatewayTimeout}", nil)

	CodeOk                        = NewMpuapsErr(codeOk, "", nil)
	ErrCommonInternalError        = NewMpuapsErr(codeCommonInternalError, "{#ErrCommonInternalError}", nil)
	ErrCommonValidationFailed     = NewMpuapsErr(codeCommonValidationFailed, "{#ErrCommonValidationFailed}", nil)
	ErrCommonDbOperationError     = NewMpuapsErr(codeCommonDbOperationError, "{#ErrCommonDbOperationError}", nil)
	ErrCommonInvalidParameter     = NewMpuapsErr(codeCommonInvalidParameter, "{#ErrCommonInvalidParameter}", nil)
	ErrCommonParsingParameter     = NewMpuapsErr(codeCommonParsingParameter, "{#ErrCommonParsingParameter}", nil)
	ErrCommonInvalidOperation     = NewMpuapsErr(codeCommonInvalidOperation, "{#ErrCommonInvalidOperation}", nil)
	ErrCommonInvalidConfiguration = NewMpuapsErr(codeCommonInvalidConfiguration, "{#ErrCommonInvalidConfiguration}", nil)
	ErrCommonMissingConfiguration = NewMpuapsErr(codeCommonMissingConfiguration, "{#ErrCommonMissingConfiguration}", nil)
	ErrCommonNotImplemented       = NewMpuapsErr(codeCommonNotImplemented, "{#ErrCommonNotImplemented}", nil)
	ErrCommonNotSupported         = NewMpuapsErr(codeCommonNotSupported, "{#ErrCommonNotSupported}", nil)
	ErrCommonOperationFailed      = NewMpuapsErr(codeCommonOperationFailed, "{#ErrCommonOperationFailed}", nil)
	ErrCommonServerBusy           = NewMpuapsErr(codeCommonServerBusy, "{#ErrCommonServerBusy}", nil)
	ErrCommonUnknown              = NewMpuapsErr(codeCommonUnknown, "{#ErrCommonUnknown}", nil)
	ErrCommonInvalidRequest       = NewMpuapsErr(codeCommonInvalidRequest, "{#ErrCommonInvalidRequest}", nil)
	ErrCommonHttpServerStartErr   = NewMpuapsErr(codeCommonHttpServerStartErr, "{#ErrCommonHttpServerStartErr}", nil)
	ErrCommonDBConnectErr         = NewMpuapsErr(codeCommonDBConnectErr, "{#ErrCommonDBConnectErr}", nil)
	ErrCommonDBLoading            = NewMpuapsErr(codeCommonDBLoading, "{#ErrCommonDBLoading}", nil)
	ErrCommonSyncMPUResources     = NewMpuapsErr(codeCommonSyncMPUResources, "{#ErrCommonSyncMPUResources}", nil)
	ErrCommonSyncReduMPURes       = NewMpuapsErr(codeCommonSyncReduMPURes, "{#ErrCommonSyncReduMPURes}", nil)
)
