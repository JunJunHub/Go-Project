// ================================================================
// 注:
//   1、GoFrame 框架内部保留错误码 < 1000
//   2、mpuserver 错误码 10000 - 99999
//
// 显控应用代理错误码定义规则:
// 1、错误代码为6位int类型[100000,199999]
// 2、首位固定位1,第二位代表大模块划分模块号,第三位表示大模块下的小模块号
//
// 扩展错误码：新增业务错误代码并初始化错误码对象
// ================================================================

package errcode

import (
	"context"
	"fmt"
	"github.com/gogf/gf/frame/g"
)

//各功能模块错误代码范围声明
const (
	//======================================
	//通用错误代码[100000,109999]
	//======================================
	codeCommonStart = 100000

	//======================================
	//系统基本功能相关错误代码[110000,119999]
	//======================================
	codeSysStart = 110000 //[110000,110999] 系统基本功能预留

	//[111000,111999]
	codeSysUserStart       = 111000 //[111000,111099] 系统用户功能预留
	codeSysUserManageStart = 111100 //[111100,111199] 系统用户管理
	codeSysUserLoginStart  = 111200 //[111200,111299] 系统用户登录

	//[112000,112999]
	codeSysRoleStart       = 112000 //[112000,112099] 系统角色管理预留
	codeSysRoleManageStart = 112100 //[112100,112199] 系统角色管理错误码

	//======================================
	//显控业务功能相关错误代码[12000,129999]
	//======================================
	codeMpusrvStart = 120000 //[120000,120999] 显控业务功能预留

	//[121000,121999]
	codeMpusrvTvWallStart           = 121000 //[121000,121099] 大屏相关功能预留
	codeMpusrvTvWallManageStart     = 121100 //[121100,121199] 大屏配置管理
	codeMpusrvTvWallSceneStart      = 121200 //[121200,121299] 大屏场景配置管理
	codeMpusrvTvWallWindowStart     = 121300 //[121300,121399] 大屏开窗
	codeMpusrvTvWallLoopStart       = 121400 //[121400,121499] 轮循配置
	codeMpusrvTvWallLoopChnCfgStart = 121500 //[121500,121599] 轮循通道配置

	//[122000,122999]
	codeMpusrvCascadeStart = 122000 //[122000,122999] 级联功能相关错误码

	//[123000,123999]
	codeMpusrvChannelStart      = 123000 //[123000,123099] 信号源相关错误码
	codeMpusrvChannelGroupStart = 123100 //[123100,123299] 信号源分组管理错误码
	codeMpusrvChannelFavorite   = 123300 //[123300,123399] 信号源收藏管理错误码
	codeMpusrvAssociationChn    = 123400 //[123400,123499] 信号源关联配置

	//[124000,124999]
	codeMpusrvMatrixDispatchStart       = 124000 //[124000,124099] 矩阵调度相关功能预留
	codeMpusrvMatrixDispatchManageStart = 124100 //[124100,124299] 矩阵调度管理错误码起始
	codeMpusrvMatrixAVChnBindCfgStart   = 124200 //[124300,124399] 矩阵丢的音视频源绑定配置

	//[125000,125999]
	codeMpusrvAlarmManageStart = 125000 //[125000,125999] 告警管理功能错误码

	//[126000,126999]
	codeMpusrvApplianceStart       = 126000 //[124600,126299] 环境设备控制器相关预留
	codeMpusrvApplianceManageStart = 126300 //[126300,126999] 环境设备管理错误码起始

	//[127000,127999]
	codeMpusrvMeetManageStart        = 127000 //[127000,127099] 会议管理功能预留错误码
	codeMpusrvMeetPlatformStart      = 127100 //[127100,127199] 会议平台管理错误码起始
	codeMpusrvMeetTerminalStart      = 127200 //[127200,127299] 会议终端管理错误码起始
	codeMpusrvMeetTemplateStart      = 127300 //[127300,127399] 会议模板管理错误码起始
	codeMpusrvMeetDispatchGroupStart = 127400 //[127400,127499] 会议调度分组管理功能错误码起始
	codeMpusrvMeetDispatchStart      = 127500 //[127500,127999] 会议调度管理功能错误码起始

	//todo add modelErrCodeStart
)

var (
	//NewMpuapsErr 定义错误码对象统一存储到 _codes
	_codes = map[int]MpuapsErr{}

	//ErrUndefined 不存储到 _codes 中
	//GetMpuapsErr 根据错误码在 _codes 中未找到对应错误码对象,则返回此错误
	ErrUndefined = MpuapsErr{
		code:    -1,
		message: "unknow err code",
		detail:  nil,
	}
)

// MpuapsErr is an implementer for interface gcode.Code.
type MpuapsErr struct {
	code    int         // Error code, usually an integer.
	message string      // Brief message for this error code.
	detail  interface{} // As type of interface, it is mainly designed as an extension field for error code.
}

// Code returns the integer number of current error code.
func (c MpuapsErr) Code() int {
	return c.code
}

// Message returns the brief message for current error code.
func (c MpuapsErr) Message() string {
	//此处返回的为定义的错误描述 key
	//在 response.RJson 中会将错误信息翻译成对应国家地区的语言
	return c.message
}

// MessageAreaLanguage returns the brief message for current error code.
// @summary 根据传入参数,返回对应的国家地区的错误描述信息, 可选参数 "zh_CN、en"
func (c MpuapsErr) MessageAreaLanguage(language string) string {
	g.I18n().SetLanguage(language)
	return g.I18n().T(context.Background(), fmt.Sprintf(`%s`, c.message))
}

// Detail returns the detailed information of current error code,
// which is mainly designed as an extension field for error code.
func (c MpuapsErr) Detail() interface{} {
	return c.detail
}

// String returns current error code as a string.
func (c MpuapsErr) String() string {
	if c.detail != nil {
		return fmt.Sprintf(`%d:%s %v`, c.code, c.Message(), c.detail)
	}
	if c.Message() != "" {
		return fmt.Sprintf(`%d:%s`, c.code, c.Message())
	}
	return fmt.Sprintf(`%d`, c.code)
}

// NewMpuapsErr
// @summary 错误码对象定义方法
// @param1  code    int "错误代码"
// @param2  message string "错误信息"
// @param3  detail  interface{} "扩展错误信息"
// @return  code    gcode.Code "返回gcode.Code接口对象"
func NewMpuapsErr(code int, message string, detail interface{}) MpuapsErr {
	if _, ok := _codes[code]; ok {
		panic(fmt.Sprintf("错误码 %d 重复,请检查错误码定义[%s,%s]", code, _codes[code].Message(), message))
	}
	mpuapsErr := MpuapsErr{
		code:    code,
		message: message,
		detail:  detail,
	}
	_codes[code] = mpuapsErr
	return mpuapsErr
}

// GetMpuapsErr
// @summary 根据错误代码获取自定义错误码对象
// @param1  code int "错误代码"
// @return  MpuapsErr implementer for interface gcode.Code
func GetMpuapsErr(code int) MpuapsErr {
	if _, ok := _codes[code]; ok {
		return _codes[code]
	}
	if code == 0 {
		return CodeOk
	}
	return ErrUndefined
}

// GetAllMpuapsErr
// @summary 获取所有已定义错误码
func GetAllMpuapsErr() map[int]MpuapsErr {
	return _codes
}
