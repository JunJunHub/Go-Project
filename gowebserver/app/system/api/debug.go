package api

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"gowebserver/app/common/api"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/system/service"
	"gowebserver/app/system/service/publish"
)

var DebugApi = new(debugApi)

type debugApi struct {
	api.BaseController
}

func (a *debugApi) Init(r *ghttp.Request) {
	a.Module = "Debug调试"
	r.SetCtxVar(global.Module, a.Module)
}

// RemoveLoginLockUser
// @summary 移除登录锁定账户
// @tags 	DEBUG
// @Param   loginName query string true "登录名"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/removeLoginLockUser [GET]
func (a *debugApi) RemoveLoginLockUser(r *ghttp.Request) {
	loginName := r.GetString("loginName")
	if loginName == "" {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter,
			errcode.ErrCommonInvalidParameter.Message()+"登录名为空"))
	}
	service.User.LoginUnLock(loginName)
	a.RespJsonExit(r, nil)
}

// CheckPasswordRight
// @summary 校验密码是否正确
// @tags    DEBUG
// @Param   loginName query string true "登录名"
// @Param   dbEncryptPass query string true "数据库加密密码"
// @Param   dbSalt query string true "加密盐"
// @Param   inputPass query string true "输入密码"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/checkPasswordRight [GET]
func (a *debugApi) CheckPasswordRight(r *ghttp.Request) {
	loginName := r.GetString("loginName")
	dbEncryptPass := r.GetString("dbEncryptPass")
	dbSalt := r.GetString("dbSalt")
	inputPass := r.GetString("inputPass")
	if service.User.CheckPasswordRight(loginName, dbEncryptPass, dbSalt, inputPass) {
		a.RespJsonExit(r, nil)
	} else {
		a.RespJsonExit(r, nil, "密码错误")
	}
}

// PrintOnlineUserInfo
// @summary 打印在线用户信息
// @tags 	DEBUG
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/printOnlineUserInfo [GET]
func (a *debugApi) PrintOnlineUserInfo(r *ghttp.Request) {
	Auth.CleanOnlineUserTask()
	onlineUsers := service.OnlineUser.GetAllOnlineUser()
	if len(onlineUsers) == 0 {
		g.Log().Info("未查到在线用户信息")
	}
	for _, userInfo := range onlineUsers {
		g.Log().Info(userInfo)
	}
	a.RespJsonExit(r, nil)
}

// PrintWsSubLinkInfo
// @summary 打印ws订阅链路信息
// @tags 	DEBUG
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/printWsSubLinkInfo [GET]
func (a *debugApi) PrintWsSubLinkInfo(r *ghttp.Request) {
	publish.WsPublish.DebugPrintWsSubLinkInfo()
	a.RespJsonExit(r, nil)
}

// PrintUserPermissions
// @summary 打印用户权限规则
// @tags    DEBUG
// @Param   userId query int true "用户ID"
// @produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/debug/printUserPermissions [GET]
func (a *debugApi) PrintUserPermissions(r *ghttp.Request) {
	userId := r.GetInt64("userId")
	roles, _ := service.User.GetUserRoles(r.GetCtx(), userId)
	g.Log().Noticef("User Roles ========================================")
	for idx, role := range roles {
		g.Log().Noticef("Role[%d]: %v", idx, role)
	}

	g.Log().Noticef("User Menus ========================================")
	for _, role := range roles {
		menuList, _ := service.InterfaceSysMenu().GetRoleMenus(r.GetCtx(), role.RoleId)
		for _, menu := range menuList {
			g.Log().Noticef("RoleId[%d] MenuInfo: %v", role.RoleId, menu)
		}
	}

	/*
		g.Log().Noticef("User ChnGroup =====================================")
		for _, role := range roles {
			groupIds, _ := mpuSrv.MpuChannelGroup.GetRoleChnGroupIds(r.GetCtx(), role.RoleId)
			groupList, _ := mpuSrv.MpuChannelGroup.GetChnGroupListByChnGroupIds(r.GetCtx(), groupIds)
			for _, group := range groupList {
				g.Log().Noticef("RoleId[%d] ChnGroupInfo %v", role.RoleId, group)
			}
		}

		g.Log().Noticef("User TvWall =======================================")
		for _, role := range roles {
			tvWallIds, _ := mpuSrv.MpuTvWall.GetRoleTvWallIds(r.GetCtx(), role.RoleId)
			tvWallList, _ := mpuSrv.MpuTvWall.GetTvWallListByTvWallIds(r.GetCtx(), tvWallIds)
			for _, wall := range tvWallList {
				g.Log().Noticef("RoleId[%d] TvWallInfo %v", role.RoleId, wall)
			}
		}

		g.Log().Noticef("User MatrixScheme =================================")
		for _, role := range roles {
			schemeIds, _ := mpuSrv.InterfaceMpuMatrixDispatchScheme().GetRoleMatrixSchemeIds(r.GetCtx(), role.RoleId)
			schemeList, _ := mpuSrv.InterfaceMpuMatrixDispatchScheme().GetMatrixSchemeListByMatrixSchemeIds(r.GetCtx(), schemeIds)
			for _, scheme := range schemeList {
				g.Log().Noticef("RoleId[%d] MatrixSchemeInfo %v", role.RoleId, scheme)
			}
		}

		g.Log().Noticef("User MatrixSchemeGroup ============================")
		for _, role := range roles {
			schemeGroupIds, _ := mpuSrv.InterfaceMpuMatrixDispatchSchemeGroup().GetRoleMatrixSchemeGroupIds(r.GetCtx(), role.RoleId)
			schemeGroupList, _ := mpuSrv.InterfaceMpuMatrixDispatchSchemeGroup().GetMatrixSchemeGroupListByGroupIds(r.GetCtx(), schemeGroupIds)
			for _, schemeGroup := range schemeGroupList {
				g.Log().Noticef("RoleId[%d] MatrixSchemeGroupInfo %v", role.RoleId, schemeGroup)
			}
		}
	*/
	a.RespJsonExit(r, nil)
}
