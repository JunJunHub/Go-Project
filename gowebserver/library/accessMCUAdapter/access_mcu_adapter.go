//====================================
// MCU对接适配器
// 封装各版本MCU接口，对上层业务提供统一接口
//====================================

package accessMCUAdapter

import (
	"context"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
)

//AccessMCUAdapter 封装后的MCU平台数据增删改查接口
type AccessMCUAdapter interface {
	GetMCUUserConnState(username string) (connState EMCUConnStateCode)                       //查询MCU账户连接状态查询
	MCUApiSystemGetVersion() (apiLevel uint, err error)                                      //获取MCU版本信息
	MCUApiAMSGetAccount(username string, accountMoid string) (account MCUAccount, err error) //查询MCU账户信息

	MCUApiNMSGetUserDomains(username string) (domains []MCUVMSDomains, err error)                                      //获取用户所属域信息
	MCUApiNMSGetAllTerminalsByUserDomain(username, userDomainMoid string) (meetTerminals []MCUMeetTerminal, err error) //获取指定用户域下的所有会议终端设备

	MCUApiMCGetPersonalTemplates(username string) (templates []MCUMCPersonalTemplateBaseInfo, err error)                          //获取会议模板列表
	MCUApiMCGetOnePersonalTemplateDetail(username, templateId string) (templateDetail MCUMCPersonalTemplateDetailInfo, err error) //获取会议模板详细信息
	MCUApiMCCreateConf(username string, createReq MCUCreateConfReq) (meetingMark MCUMeetingMark, err error)                       //创建会议
	MCUApiMCReleaseConf(username string, confId string, bReleaseSubConf bool) (err error)                                         //结束会议

	MCUApiVCGetConfDetail(username string, confId string) (confDetail MCUMeetingDetail, err error)           //获取会议详情
	MCUApiVCGetConfCascades(username string, confId string) (confCascadeInfo []MCUMeetingCascade, err error) //会议级联信息
	MCUApiVCSetSilence(username string, confId string, value int) (err error)                                //会场静音
	MCUApiVCSetMute(username string, confId string, state int, forceMute int) (err error)                    //会场哑音
	MCUApiVCGetChairman(username string, confId string) (mtId string, err error)                             //获取会议主席
	MCUApiVCSetChairman(username string, confId string, mtId string) (err error)                             //设置会议主席
	MCUApiVCGetSpeaker(username, confId string) (mtId string, forceBroadCase int, err error)                 //获取会议发言人
	MCUApiVCSetSpeaker(username, confId string, mtId string, broadcast int) (err error)                      //设置会议发言人
	MCUApiVCSetForceBroadcast(username, confId string, mode int) (err error)                                 //设置会议是否强制广播
	MCUApiVCGetVadState(username, confId string) (vadState MCUMeetingVadState, err error)                    //查询会议语音激励状态
	MCUApiVCSetVadState(username, confId string, state int, interval uint) (err error)                       //设置会议语音激励状态
	MCUApiVCGetMixs(username, confId string) (meetingMixState MCUMeetingMixState, err error)                 //查询会议混音状态
	MCUApiVCStartMixs(username, confId string, mode int, mtIds []string) (err error)                         //会议开启混音
	MCUApiVCStopMixs(username, confId string, mixId string) (err error)                                      //会议停止混音
	MCUApiVCAddMixMts(username, confId, mixId string, mtIds []string) (err error)                            //添加混音成员
	MCUApiVCDelMixMts(username, confId, mixId string, mtIds []string) (err error)                            //删除混音成员
	MCUApiVCGetVmps(username, confId, vmpId string) (vmpState MCUMeetingVmpParam, err error)                 //查询会议画面合成参数配置
	MCUApiVCStartVmps(username, confId string, vmpParam MCUMeetingVmpParam) (err error)                      //会议开启画面合成
	MCUApiVCUpdateVmps(username, confId string, vmpParam MCUMeetingVmpParam) (err error)                     //更新画面合成参数配置
	MCUApiVCStopVmps(username, confId string) (err error)                                                    //停止画面合成
	MCUApiVCGetCurConfMts(username string, confId string) (terminals []MCUMeetTerminal, err error)           //获取所有本级与会终端
	MCUApiVCGetCascadesMts(username, confId, cascadeId string) (terminals []MCUMeetTerminal, err error)      //获取所有级联与会终端
	MCUApiVCGetMt(username, confId, mtId string) (terminal MCUMeetTerminal, err error)                       //获取与会终端详细信息
	MCUApiVCCallMts(username, confId string, mts []MCUCallMtReq) (err error)                                 //批量呼叫终端
	MCUApiVCDropMts(username, confId string, mtIds []string) (err error)                                     //批量挂断终端
	MCUApiVCMtSetSilence(username, confId string, mtId string, silence int) (err error)                      //设置与会终端静音
	MCUApiVCMtSetMute(username, confId string, mtId string, mute int) (err error)                            //设置与会终端哑音
	MCUApiVCMtSetVolume(username, confId string, mtId string, mode int, volume int) (err error)              //设置与会终端音量

	//模块内部对接各版本会议平台抽象接口
	connect()
	disconnect()
	login(user *MCUUserInfo) (userLoginCache *MCULoginCache, err error)
	heartbeat(user *MCUUserInfo) error
	keepUserAlive()
	subscribeMeeting()
	DebugPrintState()
}

// MCUNotifyMsgCBFunc MCU通知消息回传函数接口
type MCUNotifyMsgCBFunc func(context.Context, MCUNotifyMessage)

var (
	//MCU适配器实例管理
	//key   = string            MCU唯一标识
	//value = AccessMCUAdapter  MCU对接实例
	instances = gmap.NewStrAnyMap(true)

	//MCU会议更新通知通道
	chnMCUNotifyMsg chan MCUNotifyMessage

	//MCU会议更新通知回传接口
	funcMCUNotifyMsgCB MCUNotifyMsgCBFunc
)

//init
func init() {
	loggerInit()
	internalMCUConfigs.mcuConfigs = make(MCUConfigs)

	//初始化MCU通知消息通道
	if chnMCUNotifyMsg == nil {
		chnMCUNotifyMsg = make(chan MCUNotifyMessage, 1000)
	}
	//MCU通知回传协程
	go mcuNotifyMsgCallBackPoll()
}

//Instance
//@summary 获取会议对接实例,如未创建则自动创建对接实例
//@param1  mcuKey string "MCU配置信息Key，根据此值找到对应的MCU配置参数"
//@return1 adapter AccessMCUAdapter "返回对应MCU实例"
//@return2 err error "报错信息"
func Instance(mcuKey string) (adapter AccessMCUAdapter, err error) {
	v := instances.GetOrSetFuncLock(mcuKey, func() interface{} {
		adapter, err = newAdapter(mcuKey)
		return adapter
	})
	if v != nil {
		return v.(AccessMCUAdapter), nil
	}
	return
}

//InstanceDestroy
//@summary 销毁会议对接实例
//         1、登出MCU
//         2、停止保活协程
//         3、停止订阅协程
//         4、删除实例
//@param1  mcuKey string "MCU配置信息Key，根据此值找到对应的MCU对接实例"
//@return1 adapter AccessMCUAdapter "返回对应版本接口实例"
//@return2 err error "报错信息"
func InstanceDestroy(mcuKey string) {
	mcuAdapter := instances.Get(mcuKey)
	if mcuAdapter != nil {
		mcuAdapter.(AccessMCUAdapter).disconnect()
	}
	instances.Remove(mcuKey)
}

//DebugPrintMCUConnInfo
//@Summary 调试打印MCU连接信息
func DebugPrintMCUConnInfo(mcuKey string) {
	if mcuKey != "" {
		logger().Notice("==============================================")
		mcuAdapter := instances.Get(mcuKey)
		if mcuAdapter != nil {
			mcuAdapter.(AccessMCUAdapter).DebugPrintState()
		} else {
			logger().Warningf("MCU[%s]连接实例不存在", mcuKey)
		}
		logger().Notice("==============================================")
		return
	}

	//打印说有MCU对接状态
	logger().Notice("==============================================")
	for _, key := range instances.Keys() {
		if mcuAdapter, err := Instance(key); err != nil {
			logger().Warning(err)
		} else {
			mcuAdapter.DebugPrintState()
		}
		logger().Notice("==============================================")
	}
}

//SetMCUNotifyMsgCBFunc
//@summary 设置接收MCU会议更新通知消息回调函数
//@param1  cbFunc MCUNotifyMsgCBFunc
func SetMCUNotifyMsgCBFunc(cbFunc MCUNotifyMsgCBFunc) {
	if funcMCUNotifyMsgCB != nil {
		logger().Warning("重复设置MCU会议更新通知消息回调函数\n", logger().GetStack())
	}
	funcMCUNotifyMsgCB = cbFunc
}

//mcuNotifyMsgCallBackPoll
//@summary MCU通知回传协程
func mcuNotifyMsgCallBackPoll() {
	for msg := range chnMCUNotifyMsg {
		logger().Debug("MCUNotifyMsg:", msg)
		if funcMCUNotifyMsgCB != nil {
			funcMCUNotifyMsgCB(context.TODO(), msg)
		} else {
			logger().Debug("MCU会议更新通知消息回调函数为nil")
		}
	}
}

//newAdapter
//@summary 创建一个会议对接实例
//         tips: 调用此接口将创建一个MCU对接实例
//               调用此接口前需要先配置MCU参数信息
//@param1  mcuKey string "MCU配置信息Key，根据此值找到对应的MCU配置参数"
//@return1 adapter AccessMCUAdapter "返回对应版本接口实例"
//@return2 err error "报错信息"
func newAdapter(mcuKey string) (adapter AccessMCUAdapter, err error) {
	internalMCUConfigs.RLock()
	defer internalMCUConfigs.RUnlock()

	if len(internalMCUConfigs.mcuConfigs) < 1 {
		//MCU配置信息为空
		return nil, gerror.NewCode(gcode.CodeInvalidConfiguration, "mcu configuration is empty, please set the mcu configuration before using")
	}

	if _, ok := internalMCUConfigs.mcuConfigs[mcuKey]; ok {
		//匹配到MCU配置信息
		mcuConnParam := internalMCUConfigs.mcuConfigs[mcuKey]
		core := &McuAdapterCore{
			mcuKey:             mcuKey,
			mcuConnParam:       &mcuConnParam,
			mcuUserLoginCaches: gmap.NewStrAnyMap(true),
		}

		//创建MCU对接实例
		switch mcuConnParam.AdapterVer {
		//科达MCU4.7基于SDK对接
		//case EMCUAdapterV4_7:

		//科达MCU5.0-5.2对接方式相同. 订阅机制为Cometd
		case EMCUAdapterV5_0, EMCUAdapterV5_2:
			adapter = &MCUAdapterV5_2{
				core,
			}

		//科达MCU5.0-6.0对接方式相同. 订阅机制为Cometd
		case EMCUAdapterV6_0:
			adapter = &MCUAdapterV6_0{
				core,
			}

		//科达MCU6.1-8.0对接方式相同. 订阅机制为Mqtt
		case EMCUAdapterV6_1, EMCUAdapterV7_0, EMCUAdapterV8_0:
			adapter = &MCUAdapterV8_0{
				core,
			}

		default:
			return nil, gerror.NewCodef(gcode.CodeInvalidConfiguration,
				`mcu configuration node "%s" is invalid configuration, AdapterVer "%s" not support`, mcuKey, mcuConnParam.AdapterVer)
		}

		//连接MCU
		core.adapter = adapter
		adapter.connect()
		return

	} else {
		//未匹配到MCU配置信息
		return nil, gerror.NewCodef(
			gcode.CodeInvalidConfiguration,
			`mcu configuration node "%s" is not found, did you misspell mcuKey "%s" or miss the mcu configuration?`, mcuKey, mcuKey)
	}
}
