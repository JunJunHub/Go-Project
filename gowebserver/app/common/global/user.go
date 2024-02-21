package global

const UserStatusOnline = "online"          //用户状态: 在线
const UserStatusOffline = "offline"        //用户状态: 离线
const UserPasswdErrCountMax = 5            //用户登录最大尝试次数
const UserLockTime = 30                    //用户锁定时间(分钟)
const UserLoginKeepaliveIntervalTimeS = 15 //登录用户保活探测包发送时间间隔(秒)

type ELoginType int

const (
	ELoginTypeUnknow            ELoginType = 0 //默认(调度管理页面)
	ELoginTypeManagerClient     ELoginType = 1 //调度管理页面   默认不支持多端登录，同一个账户只能登录一次
	ELoginTypeConfigClient      ELoginType = 2 //配置管理页面   默认不支持多端登录，同一个账户只能登录一次
	ELoginTypeCascade           ELoginType = 3 //级联登录      默认不支持多端登录，同一个账户只能登录一次
	ELoginTypeBoxRedundancy     ELoginType = 4 //机箱冗余登录   默认不支持多端登录，同一个账户只能登录一次
	ELoginTypeThirdPartyProgram ELoginType = 5 //第三方进程登录 默认支持同一个账户多次登录，不抢占
)

//AliasString 获取登录方式描述
func (x ELoginType) AliasString() string {
	var strAlias string
	switch x {
	case ELoginTypeUnknow:
		strAlias = "默认"
	case ELoginTypeManagerClient:
		strAlias = "调度管理客户端"
	case ELoginTypeConfigClient:
		strAlias = "配置管理客户端"
	case ELoginTypeCascade:
		strAlias = "级联登录"
	case ELoginTypeBoxRedundancy:
		strAlias = "冗余机箱登录"
	case ELoginTypeThirdPartyProgram:
		strAlias = "第三方登录"
	default:
		return "未知"
	}
	return strAlias
}
