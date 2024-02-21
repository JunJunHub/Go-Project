//==================================
// 对接科达会议平台API输入输出参数结构化
//==================================

package accessMCUAdapter

//mcuApiHead 协议头
type mcuApiHead struct {
	mcuApiLevel     uint                 //api版本，默认1
	mcu             *MCUServiceConnParam //mcu地址信息
	mcuAccountToken string               //token
	mcuUserCookie   string               //cookies
}

//mcuSelectPageReq 条件查询
type mcuSelectPageReq struct {
	count int
	start int
	order int
}

//mcuPostReq 请求内容
type mcuPostReq struct {
	AccountToken string      `json:"account_token,omitempty"`
	Params       interface{} `json:"params"`
}

//MCU登录请求
type mcuSystemLoginReq struct {
	AccountToken string `json:"account_token,omitempty"`
	Username     string `json:"username,omitempty"`
	Password     string `json:"password,omitempty"`
}

//mcuVMSDomains 域信息
type mcuVMSDomains struct {
	ParentMoid string `json:"parent_moid"`
	Name       string `json:"name"`
	Moid       string `json:"moid"`
	Type       string `json:"type"`
}

//mcuMCPersonalTemplatesSimple 个人模板简单信息
type mcuMCPersonalTemplatesSimple struct {
	TemplateId     string `json:"template_id,omitempty"`              //模板ID
	Name           string `json:"name,omitempty"`                     //模板名称，最大64字符
	BitRate        int    `json:"bitrate,omitempty"`                  //会议码率
	ClosedConf     int    `json:"closed_conf,omitempty"`              //会议免打扰，0=关闭，1=开启
	SafeConf       int    `json:"safe_conf,omitempty"`                //会议安全，0=公开会议，1=隐藏会议
	EncryptedType  int    `json:"encrypted_type,omitempty"`           //传输加密类型，0=不加密，2=AES加密，3=商密SM4,4=商密SM1
	EncryptedAuth  int    `json:"encrypted_auth,omitempty"`           //终端双向认证，0=关闭，1=开启
	ConfType       int    `json:"conf_type,omitempty"`                //会议类型，0=传统会议，1=端口会议，2=SFU纯转发会议
	EnableRtc      int    `json:"enable_rtc,omitempty"`               //是否允许RTC接入，0=关闭，1=开启
	PublicConf     int    `json:"public_conf,omitempty"`              //是否是来宾会议室
	Mute           int    `json:"mute,omitempty"`                     //初始化哑音
	Silence        int    `json:"silence,omitempty"`                  //初始化静音
	VideoQuality   int    `json:"video_quality,omitempty"`            //视频质量，0=质量优先，1=速度优先
	OneReforming   int    `json:"one_reforming,omitempty"`            //归一重整，0=不使用，1=使用
	EncryptedKey   string `json:"encrypted_key,omitempty"`            //传输加密密钥，最大16字符
	DualMode       int    `json:"dual_mode,omitempty"`                //双流权限，0=发言会场，1=任意会场
	VoiceActivity  int    `json:"voice_activity_detection,omitempty"` //是否开启语言激励，0=否，1=是
	Resolution     int    `json:"resolution,omitempty"`               //分辨率
	MaxJoinMt      int    `json:"max_join_mt,omitempty"`              //最大与会终端数，8=小型8方会议，32=32方会议，64=64方会议，192=大型192方会议
	E164           string `json:"e164,omitempty"`                     //预分配会议号
	AnonymousMt    int    `json:"anonymous_mt,omitempty"`             //是否允许匿名用户免登录，0=不允许，1=允许
	ConfLevel      string `json:"conf_level,omitempty"`               //会议登记0-16，数值越大会议登记越高
	EnableAudience int    `json:"enable_audience,omitempty"`          //是否超大方，0=否，1=是
	ConfProtocol   int    `json:"conf_protocol,omitempty"`            //会议优选协议，0=H323，1-SIP，2=RTC，enable_rtc=1时为RTC，否则为SIP
}

//EMCUAccountType MCU账号类型枚举
type EMCUAccountType int

const (
	EMCUAccountTypeMoid   EMCUAccountType = 1 //moid
	EMCUAccountTypeNoMail EMCUAccountType = 4 //非系统邮箱
	EMCUAccountTypeE164   EMCUAccountType = 5 //e164号
	EMCUAccountTypeTell   EMCUAccountType = 6 //电话
	EMCUAccountTypeIP     EMCUAccountType = 7 //IP地址
	EMCUAccountTypeAlias  EMCUAccountType = 8 //别名@ip(监控前端)
)

//mcuDepartment MCU部门信息
type mcuDepartment struct {
	DepartmentMoid     string `json:"department_moid"`     //部门序号(UUID格式)
	DepartmentName     string `json:"department_name"`     //部门名称
	DepartmentPosition string `json:"department_position"` //职位信息
}

//mcuAccount mcu账户信息
type mcuAccount struct {
	Account        string          `json:"account"`         //账户
	AccountJid     string          `json:"account_jid"`     //文字聊天账户jid
	AccountMoid    string          `json:"account_moid"`    //账户moid
	Binded         int64           `json:"binded"`          //164号是否绑定账号：0-否;1-是;
	DateOfBirth    string          `json:"date_of_birth"`   //出生日期
	E164           string          `json:"e164"`            //e164号
	Email          string          `json:"email"`           //邮箱
	EXTNum         string          `json:"ext_num"`         //电话号|分机号
	Fax            string          `json:"fax"`             //传真
	JobNum         string          `json:"job_num"`         //账户编号|工号
	Limited        int64           `json:"limited"`         //是否为来宾账户: 0-否;1-是;
	Mobile         string          `json:"mobile"`          //手机
	Name           string          `json:"name"`            //账户真实姓名
	OfficeLocation string          `json:"office_location"` //办公地址
	Password       string          `json:"password"`        //密码
	Sex            int64           `json:"sex"`             //性别: 0-女;1-男;
	Departments    []mcuDepartment `json:"departments"`     //账户所属部门信息
}

//mcuMCMt 终端信息
type mcuMCMt struct {
	MtId        string          `json:"mt_id,omitempty"`        //终端号，最大48字节
	Name        string          `json:"name,omitempty"`         //名称，最大128字节
	Account     string          `json:"account,omitempty"`      //账号或别名，最大128字节
	AccountType EMCUAccountType `json:"account_type,omitempty"` //账号类型，详见账号类型枚举
	Bitrate     int             `json:"bitrate,omitempty"`      //终端呼叫码率，不可超过会议码率
	Protocol    int             `json:"protocol,omitempty"`     //默认AUTO=255
}

//mcuMCMix 混音参数
type mcuMCMix struct {
	Mode    int       `json:"mode"`    //混音模式，1=智能混音、2=定制混音
	Members []mcuMCMt `json:"members"` //制定混音时的混音成员列表
}

//mcuMCVideoFormat 视频格式
type mcuMCVideoFormat struct {
	Format     int `json:"format,omitempty"`     //主视频格式，1=MPEG、2=H.261、3=H.263、4=H.264_HP、5=H.264_BP、6=H.265、H.263+
	Resolution int `json:"resolution,omitempty"` //主视频分辨率，1-QCIF/2=CIF/3=4CIF/12=720P/13=1080P/14=WCIF/15=W4CIF/16=4K
	Frame      int `json:"frame,omitempty"`      //帧率
	Bitrate    int `json:"bitrate,omitempty"`    //码率
}

//画面风格枚举
const (
	VmpStyleOne           = 1  //一画面全屏
	VmpStyleVTwo          = 2  //两画面：左右分
	VmpStyleHTwo          = 3  //两画面: 一大全屏，一小右下
	VmpStyleThree         = 4  //三画面：等大，一上二下
	VmpStyleFour          = 5  //四画面：等大，二行二列
	VmpStyleSix           = 6  //六画面：一大五小，1大左上，2小右上(2行1列)，3小下(1行3列)
	VmpStyleEight         = 7  //八画面: 1大7小，1大左上，3小右上(3行1列)，4小下(1行4列)
	VmpStyleNine          = 8  //九画面：等大，3行3列
	VmpStyleTen           = 9  //十画面(左2右10)
	VmpStyleThirteen      = 10 //十三画面: 1大12小，4小上(1行4列)，2小左中(2行1列)，1大中中，2小右中(2行1列), 4小下(1行4列)
	VmpStyleSixteen       = 11 //十六画面: 16等分，4x4;
	VmpStyleSpecFour      = 12 //特殊四画面: 1大3小，1大左，3小右(3行1列)
	VmpStyleSeven         = 13 //七画面: 3大4小，2大上(1行2列)，1大左下，4小右下(2行2列)
	VmpStyleTwenty        = 14 //二十画面: 2大18小，2大上(1行2列)，18小下(3行6列)
	VmpStyleTenH          = 15 //水平分割的十画面(上2下8)
	VmpStyleSixL2upS4down = 16 //特殊六画面(上2下4)
	VmpStyleFourteen      = 17 //十四画面: 2大12小，2大左上(1行2列)，2小右上(2行1列)，10小下(2行5列)
	VmpStyleTenM          = 18 //十画面(上4中2下4): 2大8小，4小上(1行4列)，2大中(1行2列)，4小下(1行4列)
	VmpStyleThirteenM     = 19 //十三画面(一大在中间): 1大12小，4小上(1行4列)，2小左中(2行1列)，1大中中，2小右中(2行1列), 4小下(1行4列)
	VmpStyleFifteen       = 20 //十五画面: 3大12小，3大上(1行3列)，12小下(2行6列)
	VmpStyleSixDivide     = 21 //六画面(等分)
	VmpStyleAuto          = 22 //三画面: 1大2小，1大左，2小右(2行1列)
	VmpStyleVThree        = 23 //三画面：等大，1左，2右（2行1列）

	VmpStyle1U3DFour   = 26 //四画面：1大3小，1大上，3小下（1行3列）
	VmpStyleTwentyFive = 27 //二十五画面：等大，5行5列

	VmpStyle1L4RFive    = 34 //五画面：1大4小，1大左，4小右(4行1列)
	VmpStyle1U4DFive    = 35 //五画面：1大4小，1大上，4小下（1行4列）
	VmpStyle1U10DEleven = 38 //十一画面：1大10小，1大上，10小下(2行5列)
	VmpStyle3B9LTwelve  = 39 //十二画面：3大9小，2大上(1行2列)，1大左下，9小右下(3行3列)

	VmpStyle1B14LSeventeen = 46 //十七画面：1大16小，1大左上，6小右上(3行2列)，10小下(2行5列)
	VmpStyle6B12LEighteen  = 48 //十八画面：6大12小，6小上(1行6列)，6大居中(2行3列)，6小下(1行6列)
	VmpStyle2B18LNineteen  = 51 //十九画面：2大17小，2大左上(1行2列)，2小右上(2行1列)，15小下(3行5列)
	VmpStyle1B20LTwentyOne = 54 //二十一画面：1大20小，6小上(1行6列)，4小左中(4行1列)，1大中中，4小右中(4行1列)，6小下(1行6列)
	VmpStyle1B21LTwentyTwo = 56 //二十二画面：1大21小，1大左上，6小右上(2行3列)，15小下(3行5列)
	VmpStyle4B20LTwentyTwo = 59 //二十四画面：4大20小，6小上(1行6列)，4小左中(4行1列)，4大中中(2行2列)，4小右中(4行1列)，6小下(1行6列)

	VmpStyleRUTwo        = 61 //两画面：1大1小，1大全屏，1小右上
	VmpStyleLUTwo        = 62 //两画面：1大1小，1大全屏，1小左上
	VmpStyleLDTwo        = 63 //两画面：1大1小，1大全屏，1小左下
	VmpStyle1L5RSix      = 64 //六画面：1大5小，1大左，5小右(5行1列)，3小下(1行3列)
	VmpStyle3U8DEleven   = 65 //十一画面：3大8小，3大上(1行3列)，8小下(2行4列)
	VmpStyle2B10LTwelve  = 66 //十二画面：2大10小，2大上(1行2列)，10小下(2行5列)
	VmpStyle1B14LFifteen = 67 //十五画面：1大14小，1大左上，8小右上(4行2列)，6小下(1行6列)

	VmpStyleNone = 0xff //不支持的格式
)

//台标位置枚举
const (
	PosLeftUp    = 0 //左上角
	PosLeftDown  = 1 //左下角
	PosRightUp   = 2 //右上角
	PosRightDown = 3 //右下角
	PosBottom    = 4 //底部中间
)

//画面合成跟随类型枚举
const (
	VMMemberTypeMcsSpec  = 1 //会控指定
	VMMemberTypeSpeaker  = 2 //发言人跟随
	VMMemberTypeChairman = 3 //管理方跟随（主席跟随）
	VMMemberTypeConfPoll = 4 //会议轮训跟随
	VMMemberTypeChnPoll  = 6 //单通道轮询
	VMMemberTypeShare    = 7 //内容共享跟随
)

//EMCUVCPollType MCU会控轮询模式枚举
type EMCUVCPollType int

const (
	VCPollTypeVid            = 1  //视频轮询
	VCPollTypeAVBoth         = 3  //音视频轮询
	VCPollTypeVidChn         = 6  //轮询视频：只给回传通道
	VCPollTypeSpeaker        = 7  //轮询发言人：只给回传通道
	VCPollTypeVidBrdcast     = 8  //轮询视频：给回传通道和本地广播
	VCPollTypeSpeakBrdcast   = 9  //轮询发言人：给回传通道和本地广播
	VCPollTypeVidChairman    = 10 //主席视频轮询
	VCPollTypeAVBothChairman = 11 //主席音视频轮询
)

//轮询状态枚举
const (
	VCPollStateNone   = 0 //未轮询
	VCPollStateNormal = 1 //正常轮询
	VCPollStatePause  = 2 //轮询暂停
)

//mcuMCPoll	轮询参数
type mcuMCPoll struct {
	Mode      int       `json:"mode"`        //轮询模式，详见轮询模式枚举
	State     int       `json:"state"`       //轮询状态，详见轮询状态枚举
	Num       int       `json:"num"`         //轮询次数，0=无限次轮询
	KeepTime  int       `json:"keep_time"`   //轮询间隔时间，单位秒
	PollIndex int       `json:"poll_index"`  //正在轮训的终端顺序索引，为0时表示当前轮训终端无效
	CurPollMt mcuMCMt   `json:"cur_poll_mt"` //当前轮训的终端，只读mtId
	Members   []mcuMCMt `json:"members"`     //轮询成员列表
}

//mcuMCRecorder 录像参数
type mcuMCRecorder struct {
	PublishMode  int    `json:"publish_mode,omitempty"`  //发布模式，0=不发布，1=发布
	DualStream   int    `json:"dual_stream,omitempty"`   //是否内容共享录像
	Anonymous    int    `json:"anonymous,omitempty"`     //是否支持免登录观看直播
	RecorderMode int    `json:"recorder_mode,omitempty"` //录像模式，1=录像，2=直播，3=录像+直播
	VrsId        string `json:"vrs_id,omitempty"`        //VRS的moid
	LivePassword string `json:"live_password,omitempty"` //直播密码，6位字母与数字的组合
}

//mcuMCDcs		数据协作
type mcuMCDcs struct {
	Mode int `json:"mode,omitempty"` //数据协作模式，0=关闭，1=管理方控制；2=自由协作
}

//mcuMCAI
type mcuMCAI struct {
	Record   int `json:"record,omitempty"`    //会议纪要，1=开启，2=关闭
	Sign     int `json:"sign,omitempty"`      //会议签到，1=开启，2=关闭
	SubTitle int `json:"sub_title,omitempty"` //同声字幕，1=开启，2=关闭
}

//mcuMCVmpMtNameStyle 合成画面显示[会议终端名称]参数配置
type mcuMCVmpMtNameStyle struct {
	FontSize  int    `json:"font_size"`  //台标字体大小，0=小，1=中，2=大
	FontColor string `json:"font_color"` //台标字体三原色，#RGB格式，十六进制，默认白色#FFFFFF
	Position  int    `json:"position"`   //台标位置，详见 台标位置枚举
}

//mcuMCVmp 画面合成参数
type mcuMCVmp struct {
	Mode        int                 `json:"mode"`          //画面合成模式，1=定制，2=自动，3=自动合成批量轮询，4=定制合成批量轮询
	Layout      int                 `json:"layout"`        //画面合成风格，详见 画面风格枚举
	ExceptSelf  int                 `json:"except_self"`   //是否开启自动画面和成n-1模式
	VoiceHint   int                 `json:"voice_hint"`    //是否识别声音来源
	Broadcast   int                 `json:"broadcast"`     //是否向终端广播
	ShowMtName  int                 `json:"show_mt_name"`  //是否显示别名
	MtNameStyle mcuMCVmpMtNameStyle `json:"mt_name_style"` //画面合成台标参数
	Members     []mcuMCVmpMember    `json:"members"`       //画面合成成员
	Poll        mcuMCPoll           `json:"poll"`          //画面合成批量轮训信息，当mode为3或4时返回
}
type mcuMCVmpMember struct { //画面合成成员参数
	Account     string    `json:"account"`      //账号，最大128字节
	AccountType int       `json:"account_type"` //账号类型，详见账号类型枚举
	MemberType  int       `json:"member_type"`  //跟随类型，详见画面合成跟随类型枚举
	ChnIdx      int       `json:"chn_idx"`      //在画面合成中的位置，从0开始
	MtId        string    `json:"mt_id"`        //成员终端号
	MtChnIdx    int       `json:"mt_chn_idx"`   //终端视频通道号，0=自动选择，仅会控指定时有效
	Poll        mcuMCPoll `json:"poll"`         //单通道轮询设置
}

//mcuMCPersonalTemplatesDetail 会议模板详情
type mcuMCPersonalTemplatesDetail struct {
	TemplateId         string             `json:"template_id,omitempty"`
	Name               string             `json:"name,omitempty"`                     //模板名称，最大64字符
	BitRate            int                `json:"bitrate,omitempty"`                  //会议码率
	ClosedConf         int                `json:"closed_conf,omitempty"`              //会议免打扰，0=关闭，1=开启
	SafeConf           int                `json:"safe_conf,omitempty"`                //会议安全，0=公开会议，1=隐藏会议
	Password           string             `json:"password,omitempty"`                 //会议密码，最大32字符
	EncryptedType      int                `json:"encrypted_type,omitempty"`           //传输加密类型，0=不加密，2=AES加密，3=商密SM4,4=商密SM1
	EncryptedAuth      int                `json:"encrypted_auth,omitempty"`           //终端双向认证，0=关闭，1=开启
	ConfType           int                `json:"conf_type,omitempty"`                //会议类型，0=传统会议，1=端口会议，2=SFU纯转发会议
	EnableRtc          int                `json:"enable_rtc,omitempty"`               //是否允许RTC接入，0=关闭，1=开启
	CallMode           int                `json:"call_mode,omitempty"`                //呼叫模式，0=手动，2=自动
	CallTimes          int                `json:"call_times,omitempty"`               //呼叫次数
	CallInterval       int                `json:"call_interval,omitempty"`            //呼叫间隔，单位秒
	Mute               int                `json:"mute,omitempty"`                     //初始化哑音
	Silence            int                `json:"silence,omitempty"`                  //初始化静音
	VideoQuality       int                `json:"video_quality,omitempty"`            //视频质量，0=质量优先，1=速度优先
	OneReforming       int                `json:"one_reforming,omitempty"`            //归一重整，0=不使用，1=使用
	EncryptedKey       string             `json:"encrypted_key,omitempty"`            //传输加密密钥，最大16字符
	DualMode           int                `json:"dual_mode,omitempty"`                //双流权限，0=发言会场，1=任意会场
	DoubleFlow         int                `json:"doubleflow,omitempty"`               //成为发言人后立即发起内容共享
	VoiceActivity      int                `json:"voice_activity_detection,omitempty"` //是否开启语言激励，0=否，1=是
	VacInterval        int                `json:"vacinterval,omitempty"`              //语音激励敏感度，单位秒，支持3、5、15、30
	CascadeMode        int                `json:"cascade_mode,omitempty"`             //级联模式，0=简单级联、1=合并级联
	CascadeUpload      int                `json:"cascade_upload,omitempty"`           //是否级联上传
	CascadeReturn      int                `json:"cascade_return,omitempty"`           //是否级联回传
	CascadeReturnPara  int                `json:"cascade_return_para,omitempty"`      //级联回传带宽参数
	PublicConf         int                `json:"public_conf,omitempty"`              //是否是来宾会议室
	MaxJoinMt          int                `json:"max_join_mt,omitempty"`              //最大与会终端数，8=小型8方会议，32=32方会议，64=64方会议，192=大型192方会议
	AutoEnd            int                `json:"auto_end,omitempty"`                 //会议中无终端时，是否自动结会
	PreResource        int                `json:"preoccupy_resource,omitempty"`       //预占资源模式，0=不预占模式，1=预占模式，2=SFU纯转发模式
	FecMode            int                `json:"fec_mode,omitempty"`                 //FEC开关
	MuteFilter         int                `json:"mute_filter,omitempty"`              //是否开启全场哑音例外，0=不对任何人例外，1=对发言方和管理方例外
	Duration           int                `json:"duration,omitempty"`                 //会议时长，0是永久会议
	E164               string             `json:"e164,omitempty"`                     //预分配会议号
	AnonymousMt        int                `json:"anonymous_mt,omitempty"`             //是否允许匿名用户免登录，0=不允许，1=允许
	ConfLevel          string             `json:"conf_level,omitempty"`               //会议登记0-16，数值越大会议登记越高
	EnableAudience     int                `json:"enable_audience,omitempty"`          //是否超大方，0=否，1=是
	ConfProtocol       int                `json:"conf_protocol,omitempty"`            //会议优选协议，0=H323，1-SIP，2=RTC，enable_rtc=1时为RTC，否则为SIP
	MultiStream        int                `json:"multi_stream,omitempty"`             //是否开启多流
	Watermark          int                `json:"watermark,omitempty"`                //是否开启视频水印
	Speaker            mcuMCMt            `json:"speaker,omitempty"`                  //发言人
	Chairman           mcuMCMt            `json:"chairman,omitempty"`                 //主席
	Mix                mcuMCMix           `json:"mix,omitempty"`                      //混音信息
	VideoFormats       []mcuMCVideoFormat `json:"video_formats,omitempty"`            //主视频格式列表
	AudioFormats       []int              `json:"audio_formats,omitempty"`            //音频格式列表
	InviteMembers      []mcuMCMt          `json:"invite_members,omitempty"`           //参会成员
	Vmp                mcuMCVmp           `json:"vmp,omitempty"`                      //画面和成设置
	Vips               []mcuMCMt          `json:"vips,omitempty"`                     //vip成员列表
	Poll               mcuMCPoll          `json:"poll,omitempty"`                     //轮询设置
	Recorder           mcuMCRecorder      `json:"recorder,omitempty"`                 //录像设置
	KeepCallingMembers []mcuMCMt          `json:"keep_calling_members,omitempty"`     //追呼成员数组
	Dcs                mcuMCDcs           `json:"dcs,omitempty"`                      //数据协作
	AI                 mcuMCAI            `json:"ai,omitempty"`
}

//mcuMCConfMark 会议标识
type mcuMCConfMark struct {
	ConfId          string `json:"conf_id,omitempty"`           //会议号码
	MeetingId       string `json:"meeting_id,omitempty"`        //视频会议唯一id
	MachineRoomMoid string `json:"machine_room_moid,omitempty"` //机房moid
	Description     string `json:"description,omitempty"`       //结果描述
}

//会议状态枚举
const (
	MeetStateOrder   = 1 //预约
	MeetStateStart   = 2 //开始
	MeetStateEnd     = 3 //结束
	MeetStatePending = 6 //待审批
	MeetStateFail    = 7 //审批不通过
)

//mcuVCConf 会议简单信息
type mcuVCConf struct {
	Name               string  `json:"name,omitempty"`                     //会议名称，最大128字节
	UserDomainMoid     string  `json:"user_domain_moid,omitempty"`         //会议所属的用户域moid
	UserDomainName     string  `json:"user_domain_name,omitempty"`         //会议所属的用户域名称
	MeetingRoomName    string  `json:"meeting_room_name,omitempty"`        //虚拟会议室名称
	ConfId             string  `json:"conf_id,omitempty"`                  //会议号，最大48字节
	ConfLevel          int     `json:"conf_level,omitempty"`               //会议级别，1-16，值越小级别越高
	EnableRtc          int     `json:"enable_rtc,omitempty"`               //是否允许RTC接入
	VoiceActivity      int     `json:"voice_activity_detection,omitempty"` //是否开启语音激励
	VacInterval        int     `json:"vacinterval,omitempty"`              //语音激励敏感度，单位秒
	ConfType           int     `json:"conf_type,omitempty"`                //会议类型，0=传统，1=端口，2=SFU纯转发会议
	WaterMark          int     `json:"watermark,omitempty"`                //是否开启视频水印
	ConfCategory       int     `json:"conf_category,omitempty"`            //会议类别，0=多点会议，1=调度会议，2=点对点
	StartTime          string  `json:"start_time,omitempty"`               //会议开始时间
	EndTime            string  `json:"end_time,omitempty"`                 //会议结束时间
	Duration           int     `json:"duration,omitempty"`                 //会议时长，0=永久会议
	BitRate            int     `json:"bit_rate,omitempty"`                 //会议码率
	ClosedConf         int     `json:"closed_conf,omitempty"`              //会议免打扰
	SafeConf           int     `json:"safe_conf,omitempty"`                //会议安全，0=公开会议，1=隐藏会议
	EncryptedType      int     `json:"encrypted_type,omitempty"`           //传输加密类型，0=不加密，2=AES加密，3=商密（SM4），4=商密（SM1）
	EncryptedAuth      int     `json:"encrypted_auth,omitempty"`           //终端双向认证，0=关闭，1=开启
	Mute               int     `json:"mute,omitempty"`                     //初始化哑音
	MultiStream        int     `json:"multi_stream,omitempty"`             //下级会议中的多流终端回传是否开启
	MuteFilter         int     `json:"mute_filter,omitempty"`              //是否开启全场哑音除外
	FecMode            int     `json:"fec_mode,omitempty"`                 //是否强制广播
	AnonymousMt        int     `json:"anonymous_mt,omitempty"`             //是否允许匿名用户免登录，0=不允许，1=允许
	Silence            int     `json:"silence,omitempty"`                  //初始化静音
	AutoEnd            int     `json:"auto_end,omitempty"`                 //会议中无终端，是否自动结会
	PreResource        int     `json:"preoccupy_resource,omitempty"`       //预占资源模式，0=不预占模式，1=预占模式，2=SFU纯转发模式
	VideoQuality       int     `json:"video_quality,omitempty"`            //视频质量，0=质量优先，1=速度优先
	EncryptedKey       string  `json:"encrypted_key,omitempty"`            //传输加密AES加密秘钥，最大16字节
	DualMode           int     `json:"dual_mode,omitempty"`                //双流权限，0=发言会场，1=任意会场
	PublicConf         int     `json:"public_conf,omitempty"`              //是否是公共会议室
	MaxJoinMt          int     `json:"max_join_mt,omitempty"`              //最大与会终端数，8=小型8方会议，32=32方会议，64=64方会议，192=大型192方会议
	ForceBroadcast     int     `json:"force_broadcast,omitempty"`          //是否强制广播
	NeedPassword       int     `json:"need_password,omitempty"`            //是否需要密码
	OneReforming       int     `json:"one_reforming,omitempty"`            //归一重整
	DoubleFlow         int     `json:"doubleflow,omitempty"`               //成为发言人后立即发起内容共享
	OwnConf            int     `json:"own_conf,omitempty"`                 //会议是否和我相关
	EnableAudience     int     `json:"enable_audience,omitempty"`          //是否超大方
	PlatformId         string  `json:"platform_id,omitempty"`              //创会平台moid
	SuperiorCascade    int     `json:"superior_cascade,omitempty"`         //是否有上级级联会议室
	SubordinateCascade int     `json:"subordinate_cascade,omitempty"`      //是否有下级会议室
	MeetingId          string  `json:"meeting_id,omitempty"`               //会议ID
	Creator            mcuMCMt `json:"creator,omitempty"`                  //会议发起者
}

//mcuVCConfDetail 视频会议详细信息
type mcuVCConfDetail struct {
	Name               string             `json:"name,omitempty"`                     //会议名称，最大128字节
	UserDomainMoid     string             `json:"user_domain_moid,omitempty"`         //会议所属的用户域moid
	UserDomainName     string             `json:"user_domain_name,omitempty"`         //会议所属的用户域名称
	MeetingRoomName    string             `json:"meeting_room_name,omitempty"`        //虚拟会议室名称
	MachineRoomMoid    string             `json:"machine_room_moid,omitempty"`        //会议所属机房moid
	ConfId             string             `json:"conf_id,omitempty"`                  //会议号，最大48字节
	ConfLevel          int                `json:"conf_level,omitempty"`               //会议级别，1-16，值越小级别越高
	EnableRtc          int                `json:"enable_rtc,omitempty"`               //是否允许RTC接入
	ConfType           int                `json:"conf_type,omitempty"`                //会议类型，0=传统，1=端口，2=SFU纯转发会议
	ConfCategory       int                `json:"conf_category,omitempty"`            //会议类别，0=多点会议，1=调度会议，2=点对点
	WaterMark          int                `json:"watermark,omitempty"`                //是否开启视频水印
	StartTime          string             `json:"start_time,omitempty"`               //会议开始时间
	EndTime            string             `json:"end_time,omitempty"`                 //会议结束时间
	Duration           int                `json:"duration,omitempty"`                 //会议时长，0=永久会议
	BitRate            int                `json:"bit_rate,omitempty"`                 //会议码率
	ClosedConf         int                `json:"closed_conf,omitempty"`              //会议免打扰
	SafeConf           int                `json:"safe_conf,omitempty"`                //会议安全，0=公开会议，1=隐藏会议
	EncryptedType      int                `json:"encrypted_type,omitempty"`           //传输加密类型，0=不加密，2=AES加密，3=商密（SM4），4=商密（SM1）
	EncryptedAuth      int                `json:"encrypted_auth,omitempty"`           //终端双向认证，0=关闭，1=开启
	Mute               int                `json:"mute,omitempty"`                     //初始化哑音
	MultiStream        int                `json:"multi_stream,omitempty"`             //下级会议中的多流终端回传是否开启
	MuteFilter         int                `json:"mute_filter,omitempty"`              //是否开启全场哑音除外
	ForceMute          int                `json:"force_mute,omitempty"`               //全场哑音下是否禁止终端取消自身哑音
	Silence            int                `json:"silence,omitempty"`                  //初始化静音
	VideoQuality       int                `json:"video_quality,omitempty"`            //视频质量，0=质量优先，1=速度优先
	EncryptedKey       string             `json:"encrypted_key,omitempty"`            //传输加密AES加密秘钥，最大16字节
	DualMode           int                `json:"dual_mode,omitempty"`                //双流权限，0=发言会场，1=任意会场
	PublicConf         int                `json:"public_conf,omitempty"`              //是否是公共会议室
	AutoEnd            int                `json:"auto_end,omitempty"`                 //会议中无终端，是否自动结会
	PreResource        int                `json:"preoccupy_resource,omitempty"`       //预占资源模式，0=不预占模式，1=预占模式，2=SFU纯转发模式
	MaxJoinMt          int                `json:"max_join_mt,omitempty"`              //最大与会终端数，8=小型8方会议，32=32方会议，64=64方会议，192=大型192方会议
	ForceBroadcast     int                `json:"force_broadcast,omitempty"`          //是否强制广播
	FecMode            int                `json:"fec_mode,omitempty"`                 //是否强制广播
	VoiceActivity      int                `json:"voice_activity_detection,omitempty"` //是否开启语音激励
	VacInterval        int                `json:"vacinterval,omitempty"`              //语音激励敏感度，单位秒
	CallTimes          int                `json:"call_times,omitempty"`               //呼叫次数
	CallInterval       int                `json:"call_interval,omitempty"`            //呼叫间隔，单位秒
	CallMode           int                `json:"call_mode,omitempty"`                //呼叫模式，0=手动，2=定时
	CascadeMode        int                `json:"cascade_mode,omitempty"`             //级联模式，0=简单级联，1=合并级联
	CascadeUpload      int                `json:"cascade_upload,omitempty"`           //是否级联上传
	CascadeReturn      int                `json:"cascade_return,omitempty"`           //是否级联回传
	CascadeReturnPara  int                `json:"cascade_return_para,omitempty"`      //级联回传带宽参数
	VmpEnable          int                `json:"vmp_enable,omitempty"`               //是否在合成
	MixEnable          int                `json:"mix_enable,omitempty"`               //是否在混音
	PollEnable         int                `json:"poll_enable,omitempty"`              //是否在轮询
	RecEnable          int                `json:"rec_enable,omitempty"`               //是否开启录像
	NeedPassword       int                `json:"need_password,omitempty"`            //是否需要密码
	OneReforming       int                `json:"one_reforming,omitempty"`            //归一重整
	DoubleFlow         int                `json:"doubleflow,omitempty"`               //成为发言人后立即发起内容共享
	EnableAudience     int                `json:"enable_audience,omitempty"`          //是否超大方
	ConfProtocol       int                `json:"conf_protocol,omitempty"`            //会议优选协议，0=H323,1=SIP,2=RTC
	PlatformId         string             `json:"platform_id,omitempty"`              //创会平台moid
	SuperiorCascade    int                `json:"superior_cascade,omitempty"`         //是否有上级级联会议室
	SubordinateCascade int                `json:"subordinate_cascade,omitempty"`      //是否有下级会议室
	MeetingId          string             `json:"meeting_id,omitempty"`               //会议ID
	VideoFormats       []mcuMCVideoFormat `json:"video_formats,omitempty"`            //主视频格式列表
	Creator            mcuMCMt            `json:"creator,omitempty"`                  //会议发起者
}

//mcuVCCascadeConf 级联会议
type mcuVCCascadeConf struct {
	ConfId    string             `json:"conf_id,omitempty"`    //会议号，最大48字节
	Name      string             `json:"name,omitempty"`       //会议名称，最大128字节
	CascadeId string             `json:"cascade_id,omitempty"` //级联号，最大48字节，0标识本本级会议
	MtId      string             `json:"mt_id,omitempty"`      //终端号，最大48字节
	Cascades  []mcuVCCascadeConf `json:"cascades,omitempty"`   //级联数组
}

//终端类型枚举
const (
	VCMtTypeNor        = 1 //普通终端
	VCMtTypeTell       = 3 //电话终端
	VCMtTypeGps        = 5 //卫星电话
	VCMtTypeLowerLevel = 7 //上级会议
	VCMtTypeHighLevel  = 8 //下级会议
)

//双流格式枚举
const (
	VCMtFormatMpeg     = 1 //MPEG
	VCMtFormatH261     = 2 //H.261
	VCMtFormatH263     = 3 //H.263
	VCMtFormatH264HP   = 4 //H.264_HP
	VCMtFormatH264BP   = 5 //H.264_BP
	VCMtFormatH265     = 6 //H.265
	VCMtFormatH263Plus = 7 //H.263+
)

//通道标签枚举
const (
	VCMtLabelAuto   = 0 //自动或未开启
	VCMtLabelLeft   = 1 //左边画面
	VCMtLabelCenter = 2 //中间画面
	VCMtLabelRight  = 3 //右边画面
)

//mcuVCMtChn 终端视频源信息
type mcuVCMtChn struct {
	ChnId      int    `json:"chn_id,omitempty"`     //通道号
	ChnAlias   string `json:"chn_alias,omitempty"`  //通道别名
	Bitrate    int    `json:"bitrate,omitempty"`    //双流码率，获取级联终端信息时不返回
	Resolution int    `json:"resolution,omitempty"` //双流分辨率，获取级联终端信息时不返回
	Format     int    `json:"format,omitempty"`     //双流格式，详见双流格式枚举，获取级联终端时不返回
	ChnLabel   int    `json:"chn_label,omitempty"`  //通道标签，详见通道标签枚举
}

//mcuTerminal 会议终端基本信息
type mcuTerminalBaseInfo struct {
	DomainMoid string `json:"domain_moid"`
	Name       string `json:"name"`
	Moid       string `json:"moid"`
	E164       string `json:"e164"`
	Online     string `json:"online"`
	Ip         string `json:"ip"`
}

//mcuVCMt 会议终端信息
type mcuVCMt struct {
	Account     string `json:"account,omitempty"`       //终端E164号，IP或电话号码
	AccountType int    `json:"account_type,omitempty"`  //账号类型，详见终端账号类型枚举
	Alias       string `json:"alias,omitempty"`         //终端别名，最大128字节
	MtId        string `json:"mt_id,omitempty"`         //终端号，最大48字节
	Ip          string `json:"ip,omitempty"`            //终端IP
	Online      int    `json:"online,omitempty"`        //是否在线
	E164        string `json:"e164,omitempty"`          //终端E164号
	Type        int    `json:"type,omitempty"`          //终端类型，详见终端类型枚举
	Bitrate     int    `json:"bitrate,omitempty"`       //呼叫码率
	ProductId   string `json:"product_id,omitempty"`    //终端类型型号
	Silence     int    `json:"silence,omitempty"`       //是否静音
	Mute        int    `json:"mute,omitempty"`          //是否哑音
	Dual        int    `json:"dual,omitempty"`          //是否在发送双流
	Mix         int    `json:"mix,omitempty"`           //是否在混音
	Vmp         int    `json:"vmp,omitempty"`           //是否在合成
	Inspection  int    `json:"inspection,omitempty"`    //是否在选看
	Rec         int    `json:"rec,omitempty"`           //是否在录像
	Poll        int    `json:"poll,omitempty"`          //是否在轮询
	Upload      int    `json:"upload,omitempty"`        //是否在上传
	Protocol    int    `json:"protocol,omitempty"`      //终端优选呼叫协议，0=H323,1=SIP,2=RTC,244=无优选
	CallMode    int    `json:"call_mode,omitempty"`     //呼叫模式，0=手动，2=自动，3=追呼
	VSndChnNum  int    `json:"v_snd_chn_num,omitempty"` //终端发送码流数量
	VRcvChnNum  int    `json:"v_rcv_chn_num,omitempty"` //终端接收码流数量
	//终端详细信息，查询与会终端信息时返回
	VSndChn   []mcuVCMtChn `json:"v_snd_chn,omitempty"`  //视频发送通道数组
	VRcvChn   []mcuVCMtChn `json:"v_rcv_chn,omitempty"`  //视频接收通道数组
	DVSndChn  []mcuVCMtChn `json:"dv_snd_chn,omitempty"` //双流视频发送通道
	DvRcvChn  []mcuVCMtChn `json:"dv_rcv_chn,omitempty"` //双流视频接收通道
	ARcvChn   []mcuVCMtChn `json:"a_rcv_chn,omitempty"`  //音频接收通道
	ASndChn   []mcuVCMtChn `json:"a_snd_chn,omitempty"`  //音频发送通道
	RcvVolume int          `json:"rcv_volume,omitempty"` //接收音量
	SndVolume int          `json:"snd_volume,omitempty"` //发送音量
}

//mcuVCMtSim 会议终端简单信息
type mcuVCMtSim struct {
	Account     string `json:"account,omitempty"`      //终端E164号，IP或电话号码
	AccountType int    `json:"account_type,omitempty"` //账号类型，详见终端账号类型枚举
	MtId        string `json:"mt_id,omitempty"`        //终端号，最大48字节
	Bitrate     int    `json:"bitrate"`                //呼叫码率
	Protocol    int    `json:"protocol"`               //终端优选呼叫协议，0=H323,1=SIP,2=RTC,244=无优选
	ForcedCall  int    `json:"forced_call"`            //是否强制呼叫，仅支持本级会议的情况，0=不强呼，1=强呼
	CallMode    int    `json:"call_mode"`              //呼叫模式，0=手动，2=自动，3=追呼
	CascadesId  string `json:"cascades_id,omitempty"`  //终端所在会议所属的Mcu索引
}

//
type mcuVCMtId struct {
	MtId       string `json:"mt_id,omitempty"`       //终端号，最大48字节
	ForcedCall int    `json:"forced_call,omitempty"` //是否强制呼叫，仅支持本级会议的情况，0=不强呼，1=强呼
}

//mcuVCMtSet 会议终端设置请求，例如指定主席、发言人等
type mcuVCMtSet struct {
	MtId           string `json:"mt_id" p:"mt_id"`                     //终端号，最大48字节
	ForceBroadcase int    `json:"force_broadcast" p:"force_broadcast"` //是否设置发言人强制广播，仅SFU会议有效
}

//云台控制类型枚举值
const (
	VCMtPtzUp         = 1  //上
	VCMtPtzDown       = 2  //下
	VCMtPtzLeft       = 3  //左
	VCMtPtzRight      = 4  //右
	VCMtPtzUpLeft     = 5  //上左
	VCMtPtzUpRight    = 6  //上右
	VCMtPtzDownLeft   = 7  //下左
	VCMtPtzDownRight  = 8  //下右
	VCMtPtzZoomOut    = 9  //视野大
	VCMtPtzZoomIn     = 10 //视野小
	VCMtPtzFocusOut   = 11 //焦距长
	VCMtPtzFocusIn    = 12 //焦距短
	VCMtPtzBrightUp   = 13 //亮度加
	VCMtPtzBrightDown = 14 //亮度减
	VCMtPtzAutoFocus  = 15 //自动调焦
)

//终端的视频源信息
type mcuVCMtVideo struct {
	VideoIdx   int    `json:"video_idx,omitempty"`   //视频源通道号
	VideoAlise string `json:"video_alise,omitempty"` //视频源别名
}

//当前终端的视频源通道信息
type mcuVCMtVideos struct {
	CurVideoIdx int            `json:"cur_video_idx,omitempty"` //当前终端的视频源通道号
	MtVideos    []mcuVCMtVideo `json:"mt_videos,omitempty"`     //终端视频源数组，组多返回10个信号源
}

//短消息参数
type mcuVCSms struct {
	Message   string      `json:"message,omitempty"`    //消息内容，最大1500字节
	RollNum   int         `json:"roll_num,omitempty"`   //滚动次数，1-255，新版本终端255为无限轮询
	RollSpeed int         `json:"roll_speed,omitempty"` //滚动速度，1=慢速，2=中速，3=快速
	Type      int         `json:"type,omitempty"`       //短消息类型，0=自右至左滚动，1=翻页滚动，2=全页滚动
	Mts       []mcuVCMtId `json:"mts,omitempty"`        //接收消息的终端数组
}

//语音激励
type mcuVCVac struct {
	VoiceActivityDetection int `json:"voice_activity_detection"` //是否开启语音激励：0=关闭，1=开启
	VacInterval            int `json:"vacinterval"`              //语音激励敏感度(s)，最小值3s，开启语音激励时有效
}

//选看终端
type mcuVCMtIns struct {
	MtId string `json:"mt_id,omitempty"` //终端号，最大48字节
	Type int    `json:"type,omitempty"`  //选看源类型，1=终端，2=画面合成
}

//选看参数
type mcuVCInspec struct {
	Mode int        `json:"mode,omitempty"` //选看模式，1=视频，2=音频
	Src  mcuVCMtIns `json:"src,omitempty"`  //选看源
	Dst  mcuVCMtIns `json:"dst,omitempty"`  //选看目的，目的终端号
}

// 终端成员列表
type mcuVCMtMember struct {
	Members []mcuVCMtId `json:"members"`
}
