// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package model

// SysUserForDao is the golang structure of table sys_user for DAO operations like Where/Data.
type SysUserForDao struct {
	UserId      interface{} `orm:"user_id,primary"`   // 用户id
	LoginName   interface{} `orm:"login_name,unique"` // 登录账号
	UserName    interface{} `orm:"user_name"`         // 用户昵称
	UserType    interface{} `orm:"user_type"`         // 用户类型(0系统管理员 1管理员 2操作员)|20221020实现角色管理功能废弃
	Email       interface{} `orm:"email"`             // 用户邮箱
	PhoneNumber interface{} `orm:"phone_number"`      // 手机号码
	Sex         interface{} `orm:"sex"`               // 用户性别(0保密 1男 2女)
	Avatar      interface{} `orm:"avatar"`            // 头像路径
	Password    interface{} `orm:"password"`          // 密码
	Salt        interface{} `orm:"salt"`              // 加密盐
	Status      interface{} `orm:"status"`            // 帐号状态（0正常 1停用 2未验证）
	LoginIp     interface{} `orm:"login_ip"`          // 最后登陆ip
	LoginDate   interface{} `orm:"login_date"`        // 最后登陆时间
	Remark      interface{} `orm:"remark"`            // 备注信息
	CreateBy    interface{} `orm:"create_by"`         // 创建者
	CreateTime  interface{} `orm:"create_time"`       // 创建时间
	UpdateBy    interface{} `orm:"update_by"`         // 更新者
	UpdateTime  interface{} `orm:"update_time"`       // 更新时间
}

// SysRoleForDao is the golang structure of table sys_role for DAO operations like Where/Data.
type SysRoleForDao struct {
	RoleId        interface{} `orm:"role_id,primary"` // 角色id
	RoleName      interface{} `orm:"role_name"`       // 角色名称
	RoleKey       interface{} `orm:"role_key"`        // 角色标签
	RoleSort      interface{} `orm:"role_sort"`       // 显示顺序
	RoleDataScope interface{} `orm:"role_data_scope"` // 数据范围标识(1：全部数据权限 2：自定数据权限)
	RoleStatus    interface{} `orm:"role_status"`     // 角色状态(0停用 1正常)
	Remark        interface{} `orm:"remark"`          // 备注信息
	CreateBy      interface{} `orm:"create_by"`       // 创建者
	CreateTime    interface{} `orm:"create_time"`     // 创建时间
	UpdateBy      interface{} `orm:"update_by"`       // 更新者
	UpdateTime    interface{} `orm:"update_time"`     // 更新时间
}

// SysUserOnlineForDao is the golang structure of table sys_user_online for DAO operations like Where/Data.
type SysUserOnlineForDao struct {
	Token          interface{} `orm:"token,primary"`    // 用户会话token
	UserId         interface{} `orm:"user_id"`          // 用户id
	LoginName      interface{} `orm:"login_name"`       // 登录账号
	DeptName       interface{} `orm:"dept_name"`        // 部门名称
	Ipaddr         interface{} `orm:"ipaddr"`           // 登录ip地址
	LoginLocation  interface{} `orm:"login_location"`   // 登录地点
	Browser        interface{} `orm:"browser"`          // 浏览器类型
	Os             interface{} `orm:"os"`               // 操作系统
	Status         interface{} `orm:"status"`           // 在线状态on_line在线off_line离线
	StartTimestamp interface{} `orm:"start_timestamp"`  // 创建时间
	LastAccessTime interface{} `orm:"last_access_time"` // 最后访问时间
	ExpireTime     interface{} `orm:"expire_time"`      // 超时时间，单位为分钟
}

// SysLoginHistoryForDao is the golang structure of table sys_login_history for DAO operations like Where/Data.
type SysLoginHistoryForDao struct {
	InfoId        interface{} `orm:"info_id,primary"` // 访问id
	LoginName     interface{} `orm:"login_name"`      // 登录账号
	Ipaddr        interface{} `orm:"ipaddr"`          // 登录ip地址
	LoginLocation interface{} `orm:"login_location"`  // 登录地点
	Browser       interface{} `orm:"browser"`         // 浏览器类型
	Os            interface{} `orm:"os"`              // 操作系统
	Status        interface{} `orm:"status"`          // 登录状态（0成功 其他错误代码）
	Msg           interface{} `orm:"msg"`             // 提示消息
	LoginTime     interface{} `orm:"login_time"`      // 访问时间
}

// SysOperLogForDao is the golang structure of table sys_oper_log for DAO operations like Where/Data.
type SysOperLogForDao struct {
	OperId        interface{} `orm:"oper_id,primary"` // 日志主键
	Title         interface{} `orm:"title"`           // 模块标题
	BusinessType  interface{} `orm:"business_type"`   // 业务类型（0其它 1新增 2修改 3删除 4查询 5授权 6导出 7导入 8强退 9清空）
	Method        interface{} `orm:"method"`          // 方法名称
	RequestMethod interface{} `orm:"request_method"`  // 请求方式
	OperatorType  interface{} `orm:"operator_type"`   // 操作类别（0其它 1前台用户 2级联操作）
	OperName      interface{} `orm:"oper_name"`       // 操作人员
	DeptName      interface{} `orm:"dept_name"`       // 部门名称
	OperUrl       interface{} `orm:"oper_url"`        // 请求url
	OperIp        interface{} `orm:"oper_ip"`         // 主机地址
	OperLocation  interface{} `orm:"oper_location"`   // 操作地点
	OperParam     interface{} `orm:"oper_param"`      // 请求参数
	JsonResult    interface{} `orm:"json_result"`     // 返回参数
	Status        interface{} `orm:"status"`          // 操作结果（0正常 -1异常）
	ErrorMsg      interface{} `orm:"error_msg"`       // 错误消息
	OperTime      interface{} `orm:"oper_time"`       // 操作时间
}

// SysMenuForDao is the golang structure of table sys_menu for DAO operations like Where/Data.
type SysMenuForDao struct {
	Id            interface{} `orm:"id,primary"`     // 菜单id
	Pid           interface{} `orm:"pid"`            // 父菜单id
	Title         interface{} `orm:"title"`          // 菜单标题
	TitleTag      interface{} `orm:"title_tag"`      // 菜单标签
	Type          interface{} `orm:"type"`           // 菜单类型(0目录 1菜单 2按钮)
	Component     interface{} `orm:"component"`      // 前端组件路径
	Icon          interface{} `orm:"icon"`           // 菜单图标
	Router        interface{} `orm:"router"`         // 后端路由规则
	RestInterface interface{} `orm:"rest_interface"` // 请求接口,一个菜单功能项可对应多个接口,多个接口 | 分隔
	Weigh         interface{} `orm:"weigh"`          // 菜单权重(菜单显示顺序)
	IsHide        interface{} `orm:"is_hide"`        // 是否隐藏: 0不隐藏 1隐藏
	Remark        interface{} `orm:"remark"`         // 备注
	CreateBy      interface{} `orm:"create_by"`      // 创建者
	CreateTime    interface{} `orm:"create_time"`    // 创建时间
	UpdateBy      interface{} `orm:"update_by"`      // 更新者
	UpdateTime    interface{} `orm:"update_time"`    // 更新时间
}

// SysParamForDao is the golang structure of table sys_param for DAO operations like Where/Data.
type SysParamForDao struct {
	ParamId    interface{} `orm:"param_id,primary"` // 参数配置id
	ParamName  interface{} `orm:"param_name"`       // 参数名称
	ParamKey   interface{} `orm:"param_key"`        // 参数键名
	ParamValue interface{} `orm:"param_value"`      // 参数值
	ParamType  interface{} `orm:"param_type"`       // 参数类型: 1系统内置 2非系统内置参数
	Remark     interface{} `orm:"remark"`           // 备注信息
	CreateBy   interface{} `orm:"create_by"`        // 创建者
	CreateTime interface{} `orm:"create_time"`      // 创建时间
	UpdateBy   interface{} `orm:"update_by"`        // 更新者
	UpdateTime interface{} `orm:"update_time"`      // 更新时间
}

// SysBlackListForDao is the golang structure of table sys_black_list for DAO operations like Where/Data.
type SysBlackListForDao struct {
	Id         interface{} `orm:"id,primary"`  // 配置ID
	Type       interface{} `orm:"type"`        // 1黑名单 2白名单
	IpAddr     interface{} `orm:"ip_addr"`     // 访问IP地址
	CreateBy   interface{} `orm:"create_by"`   // 创建者
	CreateTime interface{} `orm:"create_time"` // 创建时间
	UpdateBy   interface{} `orm:"update_by"`   // 更新者
	UpdateTime interface{} `orm:"update_time"` // 更新时间
}

// SysJobForDao is the golang structure of table sys_job for DAO operations like Where/Data.
type SysJobForDao struct {
	JobId          interface{} `orm:"job_id,primary"`    // 任务id
	JobName        interface{} `orm:"job_name,primary"`  // 任务名称
	JobParams      interface{} `orm:"job_params"`        // 任务参数
	JobGroup       interface{} `orm:"job_group,primary"` // 任务组名
	InvokeTarget   interface{} `orm:"invoke_target"`     // 调用目标字符串
	CronExpression interface{} `orm:"cron_expression"`   // cron执行表达式
	MisfirePolicy  interface{} `orm:"misfire_policy"`    // 计划执行策略（1多次执行 2执行一次）
	Concurrent     interface{} `orm:"concurrent"`        // 是否并发执行（0允许 1禁止）
	Status         interface{} `orm:"status"`            // 状态（0正常 1暂停）
	Remark         interface{} `orm:"remark"`            // 备注信息
	CreateBy       interface{} `orm:"create_by"`         // 创建者
	CreateTime     interface{} `orm:"create_time"`       // 创建时间
	UpdateBy       interface{} `orm:"update_by"`         // 更新者
	UpdateTime     interface{} `orm:"update_time"`       // 更新时间
}

// SysJobLogForDao is the golang structure of table sys_job_log for DAO operations like Where/Data.
type SysJobLogForDao struct {
	JobLogId      interface{} `orm:"job_log_id,primary"` // 任务日志id
	JobName       interface{} `orm:"job_name"`           // 任务名称
	JobGroup      interface{} `orm:"job_group"`          // 任务组名
	InvokeTarget  interface{} `orm:"invoke_target"`      // 调用目标字符串
	JobMessage    interface{} `orm:"job_message"`        // 日志信息
	Status        interface{} `orm:"status"`             // 执行状态（0正常 1失败）
	ExceptionInfo interface{} `orm:"exception_info"`     // 异常信息
	CreateTime    interface{} `orm:"create_time"`        // 创建时间
}
