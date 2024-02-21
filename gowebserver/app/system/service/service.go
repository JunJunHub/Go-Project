// =================================================================================
// 业务逻辑包 service
// =================================================================================

package service

import (
	"gowebserver/app/common/service/appstate"
)

// SystemServiceInit
// @summary 系统业务逻辑初始化
func SystemServiceInit() {
	//等待代理数据库连接并初始化完成
	appstate.WaitDBInit()

	//初始化系统参数信息
	InterfaceSysParamMgr().InitSysParam()

	//初始化默认角色信息
	InterfaceSysRole().SysRoleInit()

	//初始化默认账户信息
	User.SysUserInit()

	//清理在线用户表
	OnlineUser.Clean()
}

func init() {
	go SystemServiceInit()
}
