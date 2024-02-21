package define

import (
	"gowebserver/app/common/define"
	"gowebserver/app/system/model"

	"github.com/gogf/gf/net/ghttp"
)

//UserListItem json用户信息列表中,一个用户信息节点
type UserListItem struct {
	UserId      int64  `json:"userId"`
	LoginName   string `json:"loginName"`
	UserName    string `json:"userName"`
	UserType    string `json:"userType"` //2022.10.28 废弃
	RoleId      int64  `json:"roleId"`
	Email       string `json:"email"`
	Phonenumber string `json:"phoneNumber"`
	Sex         string `json:"sex"`
	Avatar      string `json:"avatar"`
	Password    string `json:"password"`
	Salt        string `json:"salt"`
	Status      string `json:"status"`
	LoginIp     string `json:"loginIp"`
	LoginDate   string `json:"loginDate"`
	Remark      string `json:"remark"`
	CreateBy    string `json:"createBy"`
	CreateTime  string `json:"createTime"`
}

// UserApiLoginReq api接口登录参数
type UserApiLoginReq struct {
	UserName     string `p:"userName"`
	Password     string `p:"password"`
	ValidateCode string `p:"validateCode"`
	IdKey        string `p:"idkey"`
	LoginType    int    `p:"loginType"` //0未知 1调度管理页面登录 2配置管理页面登录 3级联登录
}

// UserApiLoginRsp 登陆返回信息
type UserApiLoginRsp struct {
	Token           string               `json:"token"`           // 登录token
	UserInfo        *model.SysUserExtend `json:"userInfo"`        // 登录用户信息
	MenuList        []*model.SysMenu     `json:"menuList"`        // 登录用户用户拥有功能菜单列表
	RedundancyState int                  `json:"redundancyState"` // 显控冗余部署状态: 0未配置 1主控冗余 2机箱冗余 3主控冗余+机箱冗余
}

// UserApiSelectPageReq 查询用户列表请求参数
type UserApiSelectPageReq struct {
	LoginName   string `p:"loginName"`   //登陆名   			非必须(查询条件)
	Status      string `p:"status"`      //状态     			非必须(查询条件)
	Email       string `p:"email"`       //邮箱              非必须(查询条件)
	Phonenumber string `p:"phoneNumber"` //手机号码  	    非必须(查询条件)
	DeptId      int64  `p:"deptId"`      //部门				非必须(查询条件)
	Export      int
	define.SelectPageReq
}

type UserApiCreateBase struct {
	LoginName   string `p:"loginName"`
	UserName    string `p:"userName"`
	Phonenumber string `p:"phoneNumber"`
	Email       string `p:"email"`
	Sex         string `p:"sex"`
	Status      string `p:"status"`
	UserType    string `p:"userType"`
	RoleIds     string `p:"roleIds"`
	Remark      string `p:"remark"`
}

// UserApiCreateReq 新增用户请求参数
type UserApiCreateReq struct {
	UserApiCreateBase
	Password string `p:"password"`
}

// UserApiUpdateReq 修改用户资料请求参数
type UserApiUpdateReq struct {
	UserId   int64  `p:"userId"`
	Password string `p:"password"`
	UserApiCreateBase
}

// UserApiDeleteReq API执行删除内容
type UserApiDeleteReq struct {
	Ids string `p:"ids"`
}

// UserApiResetPwdReq 重置密码请求参数
type UserApiResetPwdReq struct {
	UserId int64 `p:"userId"`
	UserApiReSetPasswordReq
}

// UserApiProfileReq 修改用户资料请求参数
type UserApiProfileReq struct {
	UserName    string `p:"userName"`
	PhoneNumber string `p:"phoneNumber"`
	Email       string `p:"email"`
	Sex         string `p:"sex"`
}

// UserApiReSetPasswordReq 修改密码请求参数
type UserApiReSetPasswordReq struct {
	OldPassword string `p:"oldPassword"`
	NewPassword string `p:"newPassword"`
	Confirm     string `p:"confirm"`
}

// UserApiAdminResetPwdReq 管理员重置密码请求参数
type UserApiAdminResetPwdReq struct {
	UserId   int64  `p:"userId"`
	Password string `p:"password"`
}

// UserApiAvatarUploadReq 头像上传
type UserApiAvatarUploadReq struct {
	AvatarFile *ghttp.UploadFile `json:"avatarFile"` // 上传文件对象
}

// UserApiChangeStatus 修改状态
type UserApiChangeStatus struct {
	UserId int64  `p:"userId"`
	Status string `p:"status"`
}

// UserServiceLoginReq service登录参数
type UserServiceLoginReq struct {
	UserName string `p:"userName"`
	Password string `p:"password"`
}

// UserServiceList 查询列表返回值
type UserServiceList struct {
	List  []UserListItem `json:"list"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Total int            `json:"total"`
}
