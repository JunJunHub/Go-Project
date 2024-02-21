// ==========================================
// 系统内置参数key定义
// 定义需要支持配置的系统参数
// 例: 登录尝试次数 登录锁定时间 密码长度 密码有效期
// ==========================================

package global

const (
	SysParamKey_UploadPath               = "sys.uploadFile.uploadPath"               //文件上传路径
	SysParamKey_UploadImgType            = "sys.uploadFile.uploadImageType"          //支持上传图片格式
	SysParamKey_UploadImgSizeMax         = "sys.uploadFile.uploadImageSizeMax"       //支持上传图片大小
	SysParamKey_UploadFileType           = "sys.uploadFile.uploadFileType"           //支持上传文件格式
	SysParamKey_UploadFileSizeMax        = "sys.uploadFile.uploadFileSizeMax"        //支持上传文件大小
	SysParamKey_UploadFileCountMaxMinute = "sys.uploadFile.uploadFileCountMaxMinute" //一分钟内上传文件个数
	SysParamKey_DownloadPath             = "sys.downloadFile.downloadPath"           //文件下载路径

	SysParamKey_UserPasswordLengthMin            = "sys.user.userPasswordLengthMin"            //用户口令长度最小值
	SysParamKey_UserPasswordLengthMax            = "sys.user.userPasswordLengthMax"            //用户口令长度最大值
	SysParamKey_UserPasswordExpirationTime       = "sys.user.userPasswordExpirationTime"       //用户口令过期时间
	SysParamKey_UserLoginIncorrectAttemptTimes   = "sys.user.userLoginIncorrectAttemptTimes"   //用户登录失败尝试次数
	SysParamKey_UserLoginLockTime                = "sys.user.userLoginLockTime"                //用户登录锁定时间
	SysParamKey_UserLoginAccessConnectionsNumMax = "sys.user.userLoginAccessConnectionsNumMax" //用户登录访问最大连接数

	SysParamKey_AccessingBlacklistFilterType = "sys.blacklist.accessingBlacklistFilterType" //系统访问黑名单过滤方式

	SysParamKey_AlarmMinorLevelSoundOn      = "mpusrv.alarm.minorLevelSoundOn"      //普通告警声音启用禁用状态
	SysParamKey_AlarmSeverityLevelSoundOn   = "mpusrv.alarm.severityLevelSoundOn"   //严重告警声音启用禁用状态
	SysParamKey_AlarmMinorLevelAlarmTone    = "mpusrv.alarm.minorLevelAlarmTone"    //普通告警提示声音
	SysParamKey_AlarmSeverityLevelAlarmTone = "mpusrv.alarm.severityLevelAlarmTone" //严重告警提示声音

	SysParamKey_LogStorageNumMax               = "mpusrv.log.storageNumMax"                  //日志最大存储条数
	SysParamKey_LogStorageDayMax               = "mpusrv.log.storageDayMax"                  //日志最大存储天数
	SysParamKey_AlarmLogStorageStrategy        = "mpusrv.log.alarmLogStorageStrategy"        //告警日志存储策略
	SysParamKey_DeviceOnlineLogStorageStrategy = "mpusrv.log.deviceOnlineLogStorageStrategy" //设备上下线日志存储策略
	SysParamKey_OperLogStorageStrategy         = "mpusrv.log.operLogStorageStrategy"         //操作日志存储策略
	SysParamKey_LoginHistoryLogStorageStrategy = "mpusrv.log.loginHistoryLogStorageStrategy" //登录登出日志存储策略
)
