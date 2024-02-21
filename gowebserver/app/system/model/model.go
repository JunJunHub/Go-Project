// =================================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package model

import (
	"github.com/gogf/gf/os/gtime"
)

// SysUser is the golang structure for table sys_user.
type SysUser struct {
	UserId      int64       `orm:"user_id,primary"   json:"userId"      description:"用户id"`
	LoginName   string      `orm:"login_name,unique" json:"loginName"   description:"登录账号"`
	UserName    string      `orm:"user_name"         json:"userName"    description:"用户昵称"`
	UserType    string      `orm:"user_type"         json:"userType"    description:"用户类型(0系统管理员 1管理员 2操作员)|20221020实现角色管理功能废弃"`
	Email       string      `orm:"email"             json:"email"       description:"用户邮箱"`
	PhoneNumber string      `orm:"phone_number"      json:"phoneNumber" description:"手机号码"`
	Sex         string      `orm:"sex"               json:"sex"         description:"用户性别(0保密 1男 2女)"`
	Avatar      string      `orm:"avatar"            json:"avatar"      description:"头像路径"`
	Password    string      `orm:"password"          json:"password"    description:"密码"`
	Salt        string      `orm:"salt"              json:"salt"        description:"加密盐"`
	Status      string      `orm:"status"            json:"status"      description:"帐号状态（0正常 1停用 2未验证）"`
	LoginIp     string      `orm:"login_ip"          json:"loginIp"     description:"最后登陆ip"`
	LoginDate   *gtime.Time `orm:"login_date"        json:"loginDate"   description:"最后登陆时间"`
	Remark      string      `orm:"remark"            json:"remark"      description:"备注信息"`
	CreateBy    string      `orm:"create_by"         json:"createBy"    description:"创建者"`
	CreateTime  *gtime.Time `orm:"create_time"       json:"createTime"  description:"创建时间"`
	UpdateBy    string      `orm:"update_by"         json:"updateBy"    description:"更新者"`
	UpdateTime  *gtime.Time `orm:"update_time"       json:"updateTime"  description:"更新时间"`
}

// SysRole is the golang structure for table sys_role.
type SysRole struct {
	RoleId        int64       `orm:"role_id,primary" json:"roleId"        description:"角色id"`
	RoleName      string      `orm:"role_name"       json:"roleName"      description:"角色名称"`
	RoleKey       string      `orm:"role_key"        json:"roleKey"       description:"角色标签"`
	RoleSort      int         `orm:"role_sort"       json:"roleSort"      description:"显示顺序"`
	RoleDataScope int         `orm:"role_data_scope" json:"roleDataScope" description:"数据范围标识(1：全部数据权限 2：自定数据权限)"`
	RoleStatus    int         `orm:"role_status"     json:"roleStatus"    description:"角色状态(0停用 1正常)"`
	Remark        string      `orm:"remark"          json:"remark"        description:"备注信息"`
	CreateBy      string      `orm:"create_by"       json:"createBy"      description:"创建者"`
	CreateTime    *gtime.Time `orm:"create_time"     json:"createTime"    description:"创建时间"`
	UpdateBy      string      `orm:"update_by"       json:"updateBy"      description:"更新者"`
	UpdateTime    *gtime.Time `orm:"update_time"     json:"updateTime"    description:"更新时间"`
}

// SysUserOnline is the golang structure for table sys_user_online.
type SysUserOnline struct {
	Token          string      `orm:"token,primary"    json:"token"          description:"用户会话token"`
	UserId         int64       `orm:"user_id"          json:"userId"         description:"用户id"`
	LoginName      string      `orm:"login_name"       json:"loginName"      description:"登录账号"`
	DeptName       string      `orm:"dept_name"        json:"deptName"       description:"部门名称"`
	Ipaddr         string      `orm:"ipaddr"           json:"ipaddr"         description:"登录ip地址"`
	LoginLocation  string      `orm:"login_location"   json:"loginLocation"  description:"登录地点"`
	Browser        string      `orm:"browser"          json:"browser"        description:"浏览器类型"`
	Os             string      `orm:"os"               json:"os"             description:"操作系统"`
	Status         string      `orm:"status"           json:"status"         description:"在线状态on_line在线off_line离线"`
	StartTimestamp *gtime.Time `orm:"start_timestamp"  json:"startTimestamp" description:"创建时间"`
	LastAccessTime *gtime.Time `orm:"last_access_time" json:"lastAccessTime" description:"最后访问时间"`
	ExpireTime     int         `orm:"expire_time"      json:"expireTime"     description:"超时时间，单位为分钟"`
}

// SysLoginHistory is the golang structure for table sys_login_history.
type SysLoginHistory struct {
	InfoId        int64       `orm:"info_id,primary" json:"infoId"        description:"访问id"`
	LoginName     string      `orm:"login_name"      json:"loginName"     description:"登录账号"`
	Ipaddr        string      `orm:"ipaddr"          json:"ipaddr"        description:"登录ip地址"`
	LoginLocation string      `orm:"login_location"  json:"loginLocation" description:"登录地点"`
	Browser       string      `orm:"browser"         json:"browser"       description:"浏览器类型"`
	Os            string      `orm:"os"              json:"os"            description:"操作系统"`
	Status        string      `orm:"status"          json:"status"        description:"登录状态（0成功 其他错误代码）"`
	Msg           string      `orm:"msg"             json:"msg"           description:"提示消息"`
	LoginTime     *gtime.Time `orm:"login_time"      json:"loginTime"     description:"访问时间"`
}

// SysOperLog is the golang structure for table sys_oper_log.
type SysOperLog struct {
	OperId        int64       `orm:"oper_id,primary" json:"operId"        description:"日志主键"`
	Title         string      `orm:"title"           json:"title"         description:"模块标题"`
	BusinessType  int         `orm:"business_type"   json:"businessType"  description:"业务类型（0其它 1新增 2修改 3删除 4查询 5授权 6导出 7导入 8强退 9清空）"`
	Method        string      `orm:"method"          json:"method"        description:"方法名称"`
	RequestMethod string      `orm:"request_method"  json:"requestMethod" description:"请求方式"`
	OperatorType  int         `orm:"operator_type"   json:"operatorType"  description:"操作类别（0其它 1前台用户 2级联操作）"`
	OperName      string      `orm:"oper_name"       json:"operName"      description:"操作人员"`
	DeptName      string      `orm:"dept_name"       json:"deptName"      description:"部门名称"`
	OperUrl       string      `orm:"oper_url"        json:"operUrl"       description:"请求url"`
	OperIp        string      `orm:"oper_ip"         json:"operIp"        description:"主机地址"`
	OperLocation  string      `orm:"oper_location"   json:"operLocation"  description:"操作地点"`
	OperParam     string      `orm:"oper_param"      json:"operParam"     description:"请求参数"`
	JsonResult    string      `orm:"json_result"     json:"jsonResult"    description:"返回参数"`
	Status        int         `orm:"status"          json:"status"        description:"操作结果（0正常 -1异常）"`
	ErrorMsg      string      `orm:"error_msg"       json:"errorMsg"      description:"错误消息"`
	OperTime      *gtime.Time `orm:"oper_time"       json:"operTime"      description:"操作时间"`
}

// SysMenu is the golang structure for table sys_menu.
type SysMenu struct {
	Id            int64       `orm:"id,primary"     json:"id"            description:"菜单id"`
	Pid           int64       `orm:"pid"            json:"pid"           description:"父菜单id"`
	Title         string      `orm:"title"          json:"title"         description:"菜单标题"`
	TitleTag      string      `orm:"title_tag"      json:"titleTag"      description:"菜单标签"`
	Type          uint        `orm:"type"           json:"type"          description:"菜单类型(0目录 1菜单 2按钮)"`
	Component     string      `orm:"component"      json:"component"     description:"前端组件路径"`
	Icon          string      `orm:"icon"           json:"icon"          description:"菜单图标"`
	Router        string      `orm:"router"         json:"router"        description:"后端路由规则"`
	RestInterface string      `orm:"rest_interface" json:"restInterface" description:"请求接口,一个菜单功能项可对应多个接口,多个接口 | 分隔"`
	Weigh         int         `orm:"weigh"          json:"weigh"         description:"菜单权重(菜单显示顺序)"`
	IsHide        uint        `orm:"is_hide"        json:"isHide"        description:"是否隐藏: 0不隐藏 1隐藏"`
	Remark        string      `orm:"remark"         json:"remark"        description:"备注"`
	CreateBy      string      `orm:"create_by"      json:"createBy"      description:"创建者"`
	CreateTime    *gtime.Time `orm:"create_time"    json:"createTime"    description:"创建时间"`
	UpdateBy      string      `orm:"update_by"      json:"updateBy"      description:"更新者"`
	UpdateTime    *gtime.Time `orm:"update_time"    json:"updateTime"    description:"更新时间"`
}

// SysParam is the golang structure for table sys_param.
type SysParam struct {
	ParamId    int64       `orm:"param_id,primary" json:"paramId"    description:"参数配置id"`
	ParamName  string      `orm:"param_name"       json:"paramName"  description:"参数名称"`
	ParamKey   string      `orm:"param_key"        json:"paramKey"   description:"参数键名"`
	ParamValue string      `orm:"param_value"      json:"paramValue" description:"参数值"`
	ParamType  int         `orm:"param_type"       json:"paramType"  description:"参数类型: 1系统内置 2非系统内置参数"`
	Remark     string      `orm:"remark"           json:"remark"     description:"备注信息"`
	CreateBy   string      `orm:"create_by"        json:"createBy"   description:"创建者"`
	CreateTime *gtime.Time `orm:"create_time"      json:"createTime" description:"创建时间"`
	UpdateBy   string      `orm:"update_by"        json:"updateBy"   description:"更新者"`
	UpdateTime *gtime.Time `orm:"update_time"      json:"updateTime" description:"更新时间"`
}

// SysBlackList is the golang structure for table sys_black_list.
type SysBlackList struct {
	Id         int         `orm:"id,primary"  json:"id"         description:"配置ID"`
	Type       uint        `orm:"type"        json:"type"       description:"1黑名单 2白名单"`
	IpAddr     string      `orm:"ip_addr"     json:"ipAddr"     description:"访问IP地址"`
	CreateBy   string      `orm:"create_by"   json:"createBy"   description:"创建者"`
	CreateTime *gtime.Time `orm:"create_time" json:"createTime" description:"创建时间"`
	UpdateBy   string      `orm:"update_by"   json:"updateBy"   description:"更新者"`
	UpdateTime *gtime.Time `orm:"update_time" json:"updateTime" description:"更新时间"`
}

// SysJob is the golang structure for table sys_job.
type SysJob struct {
	JobId          int64       `orm:"job_id,primary"    json:"jobId"          description:"任务id"`
	JobName        string      `orm:"job_name,primary"  json:"jobName"        description:"任务名称"`
	JobParams      string      `orm:"job_params"        json:"jobParams"      description:"任务参数"`
	JobGroup       string      `orm:"job_group,primary" json:"jobGroup"       description:"任务组名"`
	InvokeTarget   string      `orm:"invoke_target"     json:"invokeTarget"   description:"调用目标字符串"`
	CronExpression string      `orm:"cron_expression"   json:"cronExpression" description:"cron执行表达式"`
	MisfirePolicy  string      `orm:"misfire_policy"    json:"misfirePolicy"  description:"计划执行策略（1多次执行 2执行一次）"`
	Concurrent     string      `orm:"concurrent"        json:"concurrent"     description:"是否并发执行（0允许 1禁止）"`
	Status         string      `orm:"status"            json:"status"         description:"状态（0正常 1暂停）"`
	Remark         string      `orm:"remark"            json:"remark"         description:"备注信息"`
	CreateBy       string      `orm:"create_by"         json:"createBy"       description:"创建者"`
	CreateTime     *gtime.Time `orm:"create_time"       json:"createTime"     description:"创建时间"`
	UpdateBy       string      `orm:"update_by"         json:"updateBy"       description:"更新者"`
	UpdateTime     *gtime.Time `orm:"update_time"       json:"updateTime"     description:"更新时间"`
}

// SysJobLog is the golang structure for table sys_job_log.
type SysJobLog struct {
	JobLogId      int64       `orm:"job_log_id,primary" json:"jobLogId"      description:"任务日志id"`
	JobName       string      `orm:"job_name"           json:"jobName"       description:"任务名称"`
	JobGroup      string      `orm:"job_group"          json:"jobGroup"      description:"任务组名"`
	InvokeTarget  string      `orm:"invoke_target"      json:"invokeTarget"  description:"调用目标字符串"`
	JobMessage    string      `orm:"job_message"        json:"jobMessage"    description:"日志信息"`
	Status        string      `orm:"status"             json:"status"        description:"执行状态（0正常 1失败）"`
	ExceptionInfo string      `orm:"exception_info"     json:"exceptionInfo" description:"异常信息"`
	CreateTime    *gtime.Time `orm:"create_time"        json:"createTime"    description:"创建时间"`
}
