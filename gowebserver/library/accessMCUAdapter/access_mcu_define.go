//======================================================================================
// 模块对外暴露的数据结构定义
// 1.MCU版本枚举、对接MCU状态枚举、MCU连接参数、等数据类型定义
// 2.AccessMCUAdapter接口输入输出参数数据结构定义
//   目的是为对接不同厂商不同版本MCU抽象出通用接口，模块内部封装屏蔽不同版本MCU接口的差异性
//   设计时并没有其它厂家的接口文档做参考，定义的合理性有待商榷
//======================================================================================

package accessMCUAdapter

import (
	"fmt"
)

//EMCUAdapterVer 对接MCU版本枚举
type EMCUAdapterVer uint

//MCUVersions 已支持的MCU版本列表(经过验证的)
var MCUVersions = [...]EMCUAdapterVer{
	EMCUAdapterV5_0,
	EMCUAdapterV5_2,
	EMCUAdapterV6_0,
	EMCUAdapterV7_0,
}

const (
	EMCUAdapterV4_7 EMCUAdapterVer = 1

	EMCUAdapterV5_0 EMCUAdapterVer = 10
	EMCUAdapterV5_2 EMCUAdapterVer = 12

	EMCUAdapterV6_0 EMCUAdapterVer = 20
	EMCUAdapterV6_1 EMCUAdapterVer = 21

	EMCUAdapterV7_0 EMCUAdapterVer = 30

	EMCUAdapterV8_0 EMCUAdapterVer = 40
)

func (x EMCUAdapterVer) String() string {
	switch x {
	case EMCUAdapterV4_7:
		return "keda_mcu_v4.7"

	case EMCUAdapterV5_0:
		return "keda_mcu_v5.0"
	case EMCUAdapterV5_2:
		return "keda_mcu_v5.2"

	case EMCUAdapterV6_0:
		return "keda_mcu_v6.0"
	case EMCUAdapterV6_1:
		return "keda_mcu_v6.1"

	case EMCUAdapterV7_0:
		return "keda_mcu_v7.0"

	case EMCUAdapterV8_0:
		return "keda_mcu_v8.0"

	default:
		return "unknow mcu ver"
	}
}

//MCUServiceConnParam 连接MCU服务器参数信息
type MCUServiceConnParam struct {
	MCUTag              string         //MCU唯一标识
	AdapterVer          EMCUAdapterVer //MCU适配版本
	Ip                  string         //MCU服务地址
	Port                uint           //MCU服务端口
	OAuthConsumerKey    string         //MCU认证软件Key
	OAuthConsumerSecret string         //MCU认证软件Value
	SSLCertFile         string         //SSL证书文件(全路径)
	SSLKeyFile          string         //SSL密钥文件(全路径)
	MCUUsers            []MCUUserInfo  //MCU账户信息
}

func (x *MCUServiceConnParam) String() string {
	strAlias := fmt.Sprintf("[MCUTag]=%s, [AdapterVer]=%s, [Ip]=%s, [Port]=%d, [OAuthConsumerKey]=%s, [OAuthConsumerSecret]=%s, [SSLCertFile]=%s, [SSLKeyFile]=%s, [MCUUsers]=",
		x.MCUTag,
		x.AdapterVer.String(),
		x.Ip,
		x.Port,
		x.OAuthConsumerKey,
		x.OAuthConsumerSecret,
		x.SSLCertFile,
		x.SSLKeyFile)
	for _, user := range x.MCUUsers {
		strAlias += fmt.Sprintf("%s@%s@%s;", user.Username, user.Password, user.UserDomain)
	}
	return strAlias
}

//MCUUserInfo MCU账户信息
type MCUUserInfo struct {
	Username   string //账户
	Password   string //密码
	UserDomain string //用户域
}

//MCULoginCache MCU用户登录缓存信息
type MCULoginCache struct {
	User         *MCUUserInfo
	AccountToken string            //用户token
	UserCookie   string            //用户cookie
	ConnState    EMCUConnStateCode //用户连接状态
}

//EMeetingUpdateTopic 会议更新通知消息类型枚举
type EMeetingUpdateTopic string

const (
	NotifyConfs         EMeetingUpdateTopic = "confs"           //update:Detail = nil
	NotifyCascades      EMeetingUpdateTopic = "cascades"        //update:Detail = nil
	NotifySpeaker       EMeetingUpdateTopic = "speaker"         //update:Detail = nil
	NotifyChairman      EMeetingUpdateTopic = "chairman"        //update:Detail = nil
	NotifyDualStream    EMeetingUpdateTopic = "dualstream"      //update:Detail = nil
	NotifyConfPoll      EMeetingUpdateTopic = "poll"            //update:Detail = nil
	NotifyChairmanPoll  EMeetingUpdateTopic = "chairmanpoll"    //update:Detail = nil
	NotifyVad           EMeetingUpdateTopic = "vad"             //update:Detail = nil
	NotifyMixs          EMeetingUpdateTopic = "mixs"            //update:Detail = mixId
	NotifyVmps          EMeetingUpdateTopic = "vmps"            //update:Detail = VmpId 		        notify:detail = MCUNotifyVmps
	NotifyMts           EMeetingUpdateTopic = "mts"             //update:Detail = MCUUpdateMts 		    notify:Detail = MCUNotifyAddMts
	NotifyMtVmps        EMeetingUpdateTopic = "mtvmps"          //update:Detail = VmpId
	NotifyMtChnVmps     EMeetingUpdateTopic = "mt_channel_vmps" //update:Detail = MCUUpdateMtChnVmps
	NotifyMtInspections EMeetingUpdateTopic = "inspections"     //update:Detail = MCUUpdateMtInspect 	notify:detail = MCUNotifyInspect
	NotifyMtMix         EMeetingUpdateTopic = "applymix"        //notify:Detail = string(mtId)
	NotifyMtSpeaker     EMeetingUpdateTopic = "applyspeaker"    //notify:Detail = string(mtId)
	NotifyMtChairman    EMeetingUpdateTopic = "applychairman"   //notify:detail = string(mtId)
	NotifyMCUConnState  EMeetingUpdateTopic = "mcuconnstate"    //notify:Detail = MCUConnState
)

//MCUNotifyMessage MCU通知消息体
type MCUNotifyMessage struct {
	Type   EMeetingUpdateTopic  //通知消息类型
	ConfId string               //会议ID
	Method string               //更新方法: update、delete、notify
	Detail interface{}          //通知内容: topic|method不同的对应不同的消息体,详见 EMeetingUpdateTopic 注释
	User   *MCULoginCache       //MCU用户
	Mcu    *MCUServiceConnParam //MCU信息
}

//MCUConnState MCU连接报错
type MCUConnState struct {
	LastErrCode EMCUConnStateCode //错误代码
	McuUsername string            //报错账户(为空则表示整个MCU报错)
}

//EMCUConnStateCode MCU连接报错枚举
type EMCUConnStateCode uint

const (
	EMCUConnOk        EMCUConnStateCode = 0 //连接OK
	EMCUVerErr        EMCUConnStateCode = 1 //版本错误(不支持的MCU版本)
	EMCUNetConnErr    EMCUConnStateCode = 2 //网络不可达
	EMCUSecretInvalid EMCUConnStateCode = 3 //密钥无效
	EMCUUserUnSet     EMCUConnStateCode = 4 //未设置账户信息
	EMCUUserInvalid   EMCUConnStateCode = 5 //无效的账户信息(账户名或密码错误)
)

//MCUUpdateMts	MCU通知会议终端列表变化
type MCUUpdateMts struct {
	MtId      string `json:"mt_id"`      //终端号
	CascadeId string `json:"cascade_id"` //会议级联号，0表示本级会议
}

//MCUUpdateMtInspect MCU通知终端选看变化通知
type MCUUpdateMtInspect struct {
	MtId string `json:"mt_id"` //终端号
	Mode int    `json:"mode"`  //选看模式，1=视频，2=音频
}

//MCUUpdateMtChnVmps MCU通知终端单通道自主画面合成
type MCUUpdateMtChnVmps struct {
	MtId     string `json:"mt_id"`      //终端号
	MtChnIdx string `json:"mt_chn_idx"` //终端通道索引
}

//MCUNotifyAddMts 添加终端失败通知
//                异步操作通知结果，通知 method 对应 notify
type MCUNotifyAddMts struct {
	ErrorCode      int    `json:"error_code"`       //错误号
	MtAccount      string `json:"mt_account"`       //终端E164、IP或电话号码
	MtAccountType  string `json:"mt_account_type"`  //终端类型，5=E164号码、6=电话、7=IP地址
	OccupyConfName string `json:"occupy_conf_name"` //占用改终端的会议名称，仅当错误号为21410时返回该值
}

//MCUNotifyMixs 开始混音失败通知
type MCUNotifyMixs struct {
	ErrorCode int    `json:"error_code"` //错误号
	MixId     string `json:"mix_id"`     //混音Id
	Mode      int    `json:"mode"`       //混音模式，1=讨论，2=定制混音
}

//MCUNotifyVmps 开始画面合成失败通知
type MCUNotifyVmps struct {
	ErrorCode int    `json:"error_code"` //错误号
	VmpId     string `json:"vmp_id"`     //画面合成Id
}

//MCUNotifyInspect 终端选看失败通知
type MCUNotifyInspect struct {
	ErrorCode int        `json:"error_code"` //错误号
	Mode      int        `json:"mode"`       //选看模式，1=视频，2=音频
	Src       mcuVCMtIns `json:"src"`        //选看源
	Dst       mcuVCMtIns `json:"dst"`        //选看目的，目的终端号
}

//MCUVMSDomains MCU域信息
type MCUVMSDomains struct {
	mcuVMSDomains
}

//MCUMCPersonalTemplateBaseInfo 个人会议模板信息
type MCUMCPersonalTemplateBaseInfo struct {
	mcuMCPersonalTemplatesSimple
}

//MCUMCPersonalTemplateDetailInfo 个人会议模板详细信息
type MCUMCPersonalTemplateDetailInfo struct {
	mcuMCPersonalTemplatesDetail
}

//MCUAccount MCU账户信息
type MCUAccount struct {
	mcuAccount
}

//MCUMeetTerminalSimple 会议终端简单信息
type MCUMeetTerminalSimple struct {
	MtId        string `json:"mt_id,omitempty"`        //终端号，最大48字节
	Name        string `json:"name,omitempty"`         //名称，最大128字节
	Account     string `json:"account,omitempty"`      //账号或别名，最大128字节
	AccountType int64  `json:"account_type,omitempty"` //账号类型，详见账号类型枚举
	Bitrate     int    `json:"bitrate,omitempty"`      //终端呼叫码率，不可超过会议码率
	Protocol    int    `json:"protocol,omitempty"`     //默认AUTO=255
}

//MCUMeetTerminal 会议终端信息
type MCUMeetTerminal struct {
	Account     string `json:"account"`       //终端E164号，IP或电话号码
	AccountType int    `json:"account_type"`  //账号类型，详见终端账号类型枚举
	Alias       string `json:"alias"`         //终端别名，最大128字节
	MtId        string `json:"mt_id"`         //终端号，最大48字节
	Ip          string `json:"ip"`            //终端IP
	Online      int    `json:"online"`        //是否在线
	E164        string `json:"e164"`          //终端E164号
	Type        int    `json:"type"`          //终端类型，详见终端类型枚举
	Bitrate     int    `json:"bitrate"`       //呼叫码率
	ProductId   string `json:"product_id"`    //终端类型型号
	Silence     int    `json:"silence"`       //是否静音
	Mute        int    `json:"mute"`          //是否哑音
	Dual        int    `json:"dual"`          //是否在发送双流
	Mix         int    `json:"mix"`           //是否在混音
	Vmp         int    `json:"vmp"`           //是否在合成
	Inspection  int    `json:"inspection"`    //是否在选看
	Rec         int    `json:"rec"`           //是否在录像
	Poll        int    `json:"poll"`          //是否在轮询
	Upload      int    `json:"upload"`        //是否在上传
	Protocol    int    `json:"protocol"`      //终端优选呼叫协议，0=H323,1=SIP,2=RTC,244=无优选
	CallMode    int    `json:"call_mode"`     //呼叫模式，0=手动，2=自动，3=追呼
	VSndChnNum  int    `json:"v_snd_chn_num"` //终端发送码流数量
	VRcvChnNum  int    `json:"v_rcv_chn_num"` //终端接收码流数量
	//终端详细信息，查询与会终端信息时返回
	VSndChn   []MCUMtChn `json:"v_snd_chn"`  //视频发送通道数组
	VRcvChn   []MCUMtChn `json:"v_rcv_chn"`  //视频接收通道数组
	DVSndChn  []MCUMtChn `json:"dv_snd_chn"` //双流视频发送通道
	DvRcvChn  []MCUMtChn `json:"dv_rcv_chn"` //双流视频接收通道
	ARcvChn   []MCUMtChn `json:"a_rcv_chn"`  //音频接收通道
	ASndChn   []MCUMtChn `json:"a_snd_chn"`  //音频发送通道
	RcvVolume int        `json:"rcv_volume"` //接收音量
	SndVolume int        `json:"snd_volume"` //发送音量
}
type MCUMtChn struct { //会议终端通道信息
	ChnId      int    `json:"chn_id"`     //通道号
	ChnAlias   string `json:"chn_alias"`  //通道别名
	Bitrate    int    `json:"bitrate"`    //双流码率，获取级联终端信息时不返回
	Resolution int    `json:"resolution"` //双流分辨率，获取级联终端信息时不返回
	Format     int    `json:"format"`     //双流格式，详见双流格式枚举，获取级联终端时不返回
	ChnLabel   int    `json:"chn_label"`  //通道标签，详见通道标签枚举
}

//MCUMeetingCascade 会议级联信息
type MCUMeetingCascade struct {
	ConfId    string              `json:"conf_id,omitempty"`    //会议号，最大48字节
	Name      string              `json:"name,omitempty"`       //会议名称，最大128字节
	CascadeId string              `json:"cascade_id,omitempty"` //级联号，最大48字节，0标识本本级会议
	MtId      string              `json:"mt_id,omitempty"`      //终端号，最大48字节
	Cascades  []MCUMeetingCascade `json:"cascades,omitempty"`   //级联数组
}

//MCUMeetingDetail 会议详情
type MCUMeetingDetail struct {
	Name               string                `json:"name"`                     //会议名称，最大128字节
	UserDomainMoid     string                `json:"user_domain_moid"`         //会议所属的用户域moid
	UserDomainName     string                `json:"user_domain_name"`         //会议所属的用户域名称
	MeetingRoomName    string                `json:"meeting_room_name"`        //虚拟会议室名称
	MachineRoomMoid    string                `json:"machine_room_moid"`        //会议所属机房moid
	ConfId             string                `json:"conf_id"`                  //会议号，最大48字节
	ConfLevel          int                   `json:"conf_level"`               //会议级别，1-16，值越小级别越高
	EnableRtc          int                   `json:"enable_rtc"`               //是否允许RTC接入
	ConfType           int                   `json:"conf_type"`                //会议类型，0=传统，1=端口，2=SFU纯转发会议
	ConfCategory       int                   `json:"conf_category"`            //会议类别，0=多点会议，1=调度会议，2=点对点
	WaterMark          int                   `json:"watermark"`                //是否开启视频水印
	StartTime          string                `json:"start_time"`               //会议开始时间
	EndTime            string                `json:"end_time"`                 //会议结束时间
	Duration           int                   `json:"duration"`                 //会议时长，0=永久会议
	BitRate            int                   `json:"bit_rate"`                 //会议码率
	ClosedConf         int                   `json:"closed_conf"`              //会议免打扰
	SafeConf           int                   `json:"safe_conf"`                //会议安全，0=公开会议，1=隐藏会议
	EncryptedType      int                   `json:"encrypted_type"`           //传输加密类型，0=不加密，2=AES加密，3=商密（SM4），4=商密（SM1）
	EncryptedAuth      int                   `json:"encrypted_auth"`           //终端双向认证，0=关闭，1=开启
	Mute               int                   `json:"mute"`                     //初始化哑音
	MultiStream        int                   `json:"multi_stream"`             //下级会议中的多流终端回传是否开启
	MuteFilter         int                   `json:"mute_filter"`              //是否开启全场哑音除外
	ForceMute          int                   `json:"force_mute"`               //全场哑音下是否禁止终端取消自身哑音
	Silence            int                   `json:"silence"`                  //初始化静音
	VideoQuality       int                   `json:"video_quality"`            //视频质量，0=质量优先，1=速度优先
	EncryptedKey       string                `json:"encrypted_key"`            //传输加密AES加密秘钥，最大16字节
	DualMode           int                   `json:"dual_mode"`                //双流权限，0=发言会场，1=任意会场
	PublicConf         int                   `json:"public_conf"`              //是否是公共会议室
	AutoEnd            int                   `json:"auto_end"`                 //会议中无终端，是否自动结会
	PreResource        int                   `json:"preoccupy_resource"`       //预占资源模式，0=不预占模式，1=预占模式，2=SFU纯转发模式
	MaxJoinMt          int                   `json:"max_join_mt"`              //最大与会终端数，8=小型8方会议，32=32方会议，64=64方会议，192=大型192方会议
	ForceBroadcast     int                   `json:"force_broadcast"`          //是否强制广播
	FecMode            int                   `json:"fec_mode"`                 //是否强制广播
	VoiceActivity      int                   `json:"voice_activity_detection"` //是否开启语音激励
	VacInterval        int                   `json:"vacinterval"`              //语音激励敏感度，单位秒
	CallTimes          int                   `json:"call_times"`               //呼叫次数
	CallInterval       int                   `json:"call_interval"`            //呼叫间隔，单位秒
	CallMode           int                   `json:"call_mode"`                //呼叫模式，0=手动，2=定时
	CascadeMode        int                   `json:"cascade_mode"`             //级联模式，0=简单级联，1=合并级联
	CascadeUpload      int                   `json:"cascade_upload"`           //是否级联上传
	CascadeReturn      int                   `json:"cascade_return"`           //是否级联回传
	CascadeReturnPara  int                   `json:"cascade_return_para"`      //级联回传带宽参数
	VmpEnable          int                   `json:"vmp_enable"`               //是否在合成
	MixEnable          int                   `json:"mix_enable"`               //是否在混音
	PollEnable         int                   `json:"poll_enable"`              //是否在轮询
	RecEnable          int                   `json:"rec_enable"`               //是否开启录像
	NeedPassword       int                   `json:"need_password"`            //是否需要密码
	OneReforming       int                   `json:"one_reforming"`            //归一重整
	DoubleFlow         int                   `json:"doubleflow"`               //成为发言人后立即发起内容共享
	EnableAudience     int                   `json:"enable_audience"`          //是否超大方
	ConfProtocol       int                   `json:"conf_protocol"`            //会议优选协议，0=H323,1=SIP,2=RTC
	PlatformId         string                `json:"platform_id"`              //创会平台moid
	SuperiorCascade    int                   `json:"superior_cascade"`         //是否有上级级联会议室
	SubordinateCascade int                   `json:"subordinate_cascade"`      //是否有下级会议室
	MeetingId          string                `json:"meeting_id"`               //会议ID
	VideoFormats       []MCUVideoFormat      `json:"video_formats"`            //主视频格式列表
	Creator            MCUMeetTerminalSimple `json:"creator"`                  //会议发起者
}
type MCUVideoFormat struct { //视频格式
	Format     int `json:"format"`     //主视频格式，1=MPEG、2=H.261、3=H.263、4=H.264_HP、5=H.264_BP、6=H.265、H.263+
	Resolution int `json:"resolution"` //主视频分辨率，1-QCIF/2=CIF/3=4CIF/12=720P/13=1080P/14=WCIF/15=W4CIF/16=4K
	Frame      int `json:"frame"`      //帧率
	Bitrate    int `json:"bitrate"`    //码率
}

//MCUMeetingVadState 会议语音激励状态
type MCUMeetingVadState struct {
	VoiceActivityDetection int `json:"voice_activity_detection"` //是否开启语音激励： 0-否；1-是
	VacInterval            int `json:"vacinterval"`              //语音激励敏感度(s)，最小值3s，开启语音激励时有效
}

// MCUMeetingMixState 会议混音状态
type MCUMeetingMixState struct {
	Mode    int64          `json:"mode"`    //混音模式：0-未开混音；1-智能混音；2-定制混音
	Members []MCUMixMember `json:"members"` //混音成员数组
}
type MCUMixMember struct {
	MtId    string `json:"mt_id"`   //成员终端号，最大字符长度：48个字节
	Account string `json:"account"` //成员呼叫别名
}

//MCUCreateConfReq 创建会议请求
type MCUCreateConfReq struct {
	CreateType int    //创建会议类型： 1=及时会议，2=公共模板，3=个人模板，4=根据虚拟会议室创建，5=预约会议提前召开，当前只支持根据个人模板召开会议
	TemplateId string //会议模板ID
	Name       string //会议名称
	Duration   uint   //会议时长：0=永久会议（即手动结束），否则按时长自动结束单位分钟
}

//MCUMeetingMark MCU会议标识(创建会议返回参数)
type MCUMeetingMark struct {
	ConfId          string `json:"conf_id"`           //会议号码
	MeetingId       string `json:"meeting_id"`        //视频会议唯一id
	MachineRoomMoid string `json:"machine_room_moid"` //机房moid
	Description     string `json:"description"`       //结果描述
}

//MCUCallMtReq 呼叫终端请求
type MCUCallMtReq struct {
	MtId       string `json:"mt_id"`       //终端号，最大 48 字节
	ForcedCall int    `json:"forced_call"` //是否强制呼叫，仅支持本级会议的情况，0=不强呼，1=强呼
}

//MCUMeetingVmpParam 会议画面合成参数信息
type MCUMeetingVmpParam struct {
	Mode        int               `json:"mode"`          //画面合成模式，1=定制，2=自动，3=自动合成批量轮询，4=定制合成批量轮询
	Layout      int               `json:"layout"`        //画面合成风格，详见会议接口文档
	ExceptSelf  int               `json:"except_self"`   //是否开启自动画面合成(n-1)模式：0-否；1-是；
	VoiceHint   int               `json:"voice_hint"`    //是否识别声音来源：0-否；1-是；
	Broadcast   int               `json:"broadcast"`     //是否向终端广播：0-否；1-是；
	ShowMtName  int               `json:"show_mt_name"`  //是否显示别名：0-否；1-是；
	MtNameStyle MCUVmpMTNameStyle `json:"mt_name_style"` //画面合成台标参数
	Members     []MCUVmpMember    `json:"members"`       //画面合成成员
	Poll        MCUPollParam      `json:"poll"`          //画面合成批量轮训信息，当mode为3或4时返回
}
type MCUVmpMTNameStyle struct { //画面合成台标参数
	FontSize  int    `json:"font_size"`  //台标字体大小，0=小，1=中，2=大
	FontColor string `json:"font_color"` //台标字体三原色，#RGB格式，十六进制，默认白色#FFFFFF
	Position  int    `json:"position"`   //台标位置，详见 台标位置枚举
}
type MCUVmpMember struct {
	Account     string       `json:"account"`      //账号，最大 128 字节
	AccountType int          `json:"account_type"` //账号类型，详见会议账号类型枚举
	MemberType  int          `json:"member_type"`  //跟随类型，详见会议画面合成跟随类型枚举
	ChnIdx      int          `json:"chn_idx"`      //在画面合成中的位置，从0开始
	MtId        string       `json:"mt_id"`        //成员终端号
	MtChnIdx    int          `json:"mt_chn_idx"`   //终端视频通道号，0=自动选择，仅会控指定时有效
	Poll        MCUPollParam `json:"poll"`         //单通道轮询设置
}

//MCUPollParam 会议轮询参数
type MCUPollParam struct {
	Mode      int                     `json:"mode"`        //轮询模式，详见轮询模式枚举
	State     int                     `json:"state"`       //轮询状态，详见轮询状态枚举
	Num       int                     `json:"num"`         //轮询次数，0=无限次轮询
	KeepTime  int                     `json:"keep_time"`   //轮询间隔时间，单位秒
	PollIndex int                     `json:"poll_index"`  //正在轮训的终端顺序索引，为0时表示当前轮训终端无效
	CurPollMt MCUMeetTerminalSimple   `json:"cur_poll_mt"` //当前轮训的终端，只读mtId
	Members   []MCUMeetTerminalSimple `json:"members"`     //轮询成员列表
}
