package define

import (
	"gowebserver/app/common/define"
	"gowebserver/app/system/model"

	"github.com/gogf/gf/os/gtime"
)

//OnlineUserApiSelectPageReq 分页请求参数
type OnlineUserApiSelectPageReq struct {
	Token          string      `p:"token"`          //用户会话id
	LoginName      string      `p:"loginName"`      //登录账号
	DeptName       string      `p:"deptName"`       //部门名称
	Ipaddr         string      `p:"ipaddr"`         //登录IP地址
	LoginLocation  string      `p:"loginLocation"`  //登录地点
	Browser        string      `p:"browser"`        //浏览器类型
	Os             string      `p:"os"`             //操作系统
	Status         string      `p:"status"`         //在线状态online在线offline离线
	StartTimestamp *gtime.Time `p:"startTimestamp"` //session创建时间
	LastAccessTime *gtime.Time `p:"lastAccessTime"` //session最后访问时间
	ExpireTime     int         `p:"expireTime"`     //超时时间，单位为分钟
	define.SelectPageReq
}

//OnlineUserServiceList 查询列表返回值
type OnlineUserServiceList struct {
	List  []model.SysUserOnline `json:"list"`
	Page  int                   `json:"page"`
	Size  int                   `json:"size"`
	Total int                   `json:"total"`
}

//EForceLogoutType 强制登出类型
type EForceLogoutType int

const (
	EForceLogOut_TokenInValid        EForceLogoutType = 1 //Token过期
	EForceLogOut_UnConnWsSunLink     EForceLogoutType = 2 //未连接Ws订阅服务
	EForceLogOut_UnConnRedundancyBox EForceLogoutType = 3 //未连接冗余机箱
	EForceLogOut_PermissionChange    EForceLogoutType = 4 //权限变更
)

//OnlineUserForceLogoutNotify 强制登出通知消息
type OnlineUserForceLogoutNotify struct {
	Token  string           `json:"token"`
	Reason EForceLogoutType `json:"reason"` //emReduBox
}
