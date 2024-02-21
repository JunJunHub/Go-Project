// =================================================================================
// 系统用户管理API
// =================================================================================

package api

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"gowebserver/app/common/api"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/common/utils"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"gowebserver/app/system/service"
)

// UserApi 用户管理API对象
var UserApi = new(userApi)

type userApi struct {
	api.BaseController
}

func (a *userApi) Init(r *ghttp.Request) {
	a.Module = "用户管理"
	r.SetCtxVar(global.Module, a.Module)
}

//  ========================系统用户信息管理,一般系统管理员才有权限进行操作=========================

// Get 分页获取用户列表
// @summary 用户列表分页获取
// @tags 	用户管理
// @Param   authorization header string true "Bearer Token"
// @Param   loginName   query string false "用户名称(模糊匹配)"
// @Param   status      query string false "用户状态: 0正常 1停用"
// @Param   email       query string false "邮箱(模糊匹配)"
// @Param   phoneNumber query string false "电话号码(模糊匹配)"
// @Param   beginTime   query string false "用户创建时间范围-起始"
// @Param   endTime     query string false "用户创建时间范围-结束"
// @Param   deptId      query int64  false "部门ID"
// @Param   pageNum     query int    true  "当前页码"
// @Param   pageSize    query int    true  "每页展示条数"
// @Param   sortName    query string false "排序字段"
// @Param   sortOrder   query string false "排序方式: "ASC"升 "DESC"降, 当sortName不为空时有效"
// @Produce json
// @Success 200 {object} response.Response{data=define.UserServiceList} "用户信息列表"
// @Failure 500 {object} response.Response "请求参数错误"
// @Router 	/system/user [GET]
func (a *userApi) Get(r *ghttp.Request) {
	var req *define.UserApiSelectPageReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	resp := service.User.GetList(r.Context(), req)

	//默认账户昵称翻译成英文
	if "en" == utils.ParseAcceptLanguage(r.GetHeader("Accept-Language")) {
		for idx, user := range resp.List {
			if user.UserId == 1 {
				resp.List[idx].UserName = "system administrator"
			}
		}
	}
	a.RespJsonExit(r, nil, resp)
}

// Post 新增用户数据[POST]
// @summary 新增用户数据
// @tags 	用户管理
// @Accept  json
// @Param   Authorization header string true "Bearer Token"
// @Param   UserApiCreateReq body define.UserApiCreateReq{} true "新增用户请求参数"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/system/user [POST]
func (a *userApi) Post(r *ghttp.Request) {
	var req *define.UserApiCreateReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	uid, err := service.User.Create(r.Context(), req)
	if err != nil || uid <= 0 {
		a.RespJsonExit(r, err)
	}

	//通知冗余机箱更新用户角色信息
	//mpuSrv.InterfaceMpuBoxRedundancy().NotifyPeerBoxUpdateCasbinRule()

	a.RespJsonExit(r, nil, uid)
}

// Put 修改用户数据[PUT]
// @summary 修改用户数据
// @tags 	用户管理
// @Accept  json
// @Param   Authorization header string true "Bearer Token"
// @Param   UserApiUpdateReq body define.UserApiUpdateReq{} true "修改用户请求参数"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/system/user [PUT]
func (a *userApi) Put(r *ghttp.Request) {
	var req *define.UserApiUpdateReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	uid, err := service.User.Update(r.Context(), req)
	if err != nil || uid <= 0 {
		a.RespJsonExit(r, err)
	}

	//通知冗余机箱更新用户角色信息
	//mpuSrv.InterfaceMpuBoxRedundancy().NotifyPeerBoxUpdateCasbinRule()

	a.RespJsonExit(r, nil, uid)
}

// Delete 删除数据[DELETE]
// @summary 删除用户数据
// @tags 	用户管理
// @Param   Authorization header string true "Bearer Token"
// @Param   ids query string true "删除用户ID列表,多个用户ID用 , 分割. 例：uid1,uid2,uid3"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/system/user [DELETE]
func (a *userApi) Delete(r *ghttp.Request) {
	var req *define.UserApiDeleteReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	delNum, err := service.User.Delete(r.GetCtx(), req.Ids)
	if err != nil {
		a.RespJsonExit(r, err)
	}

	/*
		//释放用户相关的收藏夹及收藏通道数据
		idArr := convert.ToInt64Array(req.Ids, ",")
		for _, uid := range idArr {
			mpuSrv.MpuChannelFavorites.DelFavoritesDataByUserId(uid)
			mpuSrv.MpuChannelFavoritesAudio.DelFavoritesDataByUserId(uid)
		}

		//通知冗余机箱更新用户角色信息
		mpuSrv.InterfaceMpuBoxRedundancy().NotifyPeerBoxUpdateCasbinRule()
	*/
	a.RespJsonExit(r, nil, delNum)
}

// Info 查看用户详情
// @summary 查看用户详细数据(角色权限)
// @tags 	用户管理
// @Param   Authorization header string true "Bearer Token"
// @Param   userId query int64 true "用户id"
// @Produce json
// @Success 200 {object} response.Response{data=model.SysUserExtend} "用户详细信息"
// @Failure 500 {object} response.Response "错误信息"
// @Router 	/system/user/info [GET]
func (a *userApi) Info(r *ghttp.Request) {
	uid := r.GetInt64("userId")

	var userInfo model.SysUserExtend
	user, _ := service.User.GetInfo(r.GetCtx(), uid)
	if user == nil {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysUserInvalidUid))
	}
	userRoles, _ := service.User.GetUserRoles(r.Context(), user.UserId)
	userInfo.SysUser = *user
	userInfo.Roles = userRoles
	a.RespJsonExit(r, nil, userInfo)
}

// ResetPwdSave 管理员重置密码
// @summary 管理员重置用户密码
// @tags 	用户管理
// @Accept  json
// @Param   Authorization header string true "Bearer Token"
// @Param   UserApiResetPwdReq body define.UserApiResetPwdReq true "用户Id和用户密码"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/system/user/resetPwd [PUT]
func (a *userApi) ResetPwdSave(r *ghttp.Request) {
	var req *define.UserApiResetPwdReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	result, err := service.User.ResetPassword(r.GetCtx(), req)
	if err != nil || !result {
		a.RespJsonExit(r, err)
	} else {
		a.RespJsonExit(r, nil)
	}
}

// ChangeStatus 管理员修改用户状态("status":"0"启用、"status":"1"停用)
// @summary 管理员修改用户状态
// @tags 	用户管理
// @Accept  json
// @Param   Authorization header string true "Bearer Token"
// @Param   UserApiChangeStatus body define.UserApiChangeStatus true "用户状态: "status":"0"启用 "1"停用"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/system/user/changeStatus [PUT]
func (a *userApi) ChangeStatus(r *ghttp.Request) {
	var req *define.UserApiChangeStatus
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	err := service.User.ChangeStatus(req.UserId, req.Status)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	a.RespJsonExit(r, nil)
}

// Export 管理员导出用户数据
// @summary 管理员导出用户数据
// @tags 	用户管理
// @Accept  json
// @Param   Authorization header string true "Bearer Token"
// @Param   loginName   query string false "用户名称(模糊匹配)"
// @Param   status      query string false "用户状态: 0正常 1停用"
// @Param   email       query string false "邮箱(模糊匹配)"
// @Param   phoneNumber query string false "电话号码(模糊匹配)"
// @Param   beginTime   query string false "用户创建时间范围-起始"
// @Param   endTime     query string false "用户创建时间范围-结束"
// @Param   deptId      query int64  false "部门ID"
// @Param   pageNum     query int    true  "当前页码"
// @Param   pageSize    query int    true  "每页展示条数"
// @Param   sortName    query string false "排序字段"
// @Param   sortOrder   query string false "排序方式: "ASC"升 "DESC"降, 当sortName不为空时有效"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/system/user/export [GET]
func (a *userApi) Export(r *ghttp.Request) {
	r.SetCtxVar(global.ResponseBusinessType, global.BusinessExport)
	var req *define.UserApiSelectPageReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	url, err := service.User.Export(req)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	a.RespJsonExit(r, nil, url)
}

// ============ 登录用户(以下API为登录用户对自己信息的修改维护, 不需要验证casbin_rule) ==========

// GetUserInfo 获取登录用户资料
// @summary 登录用户查看自身信息,登录成功之后调用
// @tags 	用户管理
// @Param   Authorization header string true "Bearer Token"
// @Produce json
// @Success 200 {object} response.Response{data=model.SysUserExtend} "用户详细信息"
// @Failure 500 {object} response.Response "错误信息"
// @Router 	/system/user/getInfo [GET]
func (a *userApi) GetUserInfo(r *ghttp.Request) {
	userInfo, err := service.User.GetProfile(r.Context())
	if err != nil {
		a.RespJsonExit(r, err)
	}
	a.RespJsonExit(r, nil, userInfo)
}

// GetProfile 登录用户查看自身详情
// @summary 登录用户查看自身详细信息
// @tags 	用户管理
// @Param   Authorization header string true "Bearer Token"
// @Produce json
// @Success 200 {object} response.Response{data=model.SysUserExtend} "用户详细信息"
// @Failure 500 {object} response.Response "错误信息"
// @Router 	/system/user/profile [GET]
func (a *userApi) GetProfile(r *ghttp.Request) {
	userInfo, err := service.User.GetProfile(r.Context())
	if err != nil || userInfo == nil {
		a.RespJsonExit(r, err)
	}
	if nil == userInfo {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysUserInvalidUid))
	}
	a.RespJsonExit(r, nil, userInfo)
}

// UpdateProfile 登录用户修改个人信息
// @summary 登录用户修改个人信息
// @tags 	用户管理
// @Accept  json
// @Param   Authorization header string true "Bearer Token"
// @Param   UserApiProfileReq body define.UserApiProfileReq true "用户修改信息"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/system/user/profile [PUT]
func (a *userApi) UpdateProfile(r *ghttp.Request) {
	var req *define.UserApiProfileReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	err := service.User.UpdateProfile(r.Context(), req)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	a.RespJsonExit(r, nil)
}

// UpdatePassword 登录用户修改自身密码
// @summary 登录用户修改自身密码
// @tags 	用户管理
// @Accept  json
// @Param   Authorization header string true "Bearer Token"
// @Param   UserApiReSetPasswordReq body define.UserApiReSetPasswordReq true "旧密码和新密码"
// @Produce json
// @Success 200 {object} response.Response "执行结果"
// @Router 	/system/user/profile/updatePwd [PUT]
func (a *userApi) UpdatePassword(r *ghttp.Request) {
	var req *define.UserApiReSetPasswordReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	err := service.User.UpdatePassword(r.Context(), req)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	a.RespJsonExit(r, nil)
}

// UpdateAvatar 登录用户修改头像(测试功能不需要)
func (a *userApi) UpdateAvatar(r *ghttp.Request) {
	file := r.GetUploadFile("avatarFile")
	if file == nil {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}
	avatar, err := service.User.UpdateAvatar(r.Context(), &define.UserApiAvatarUploadReq{AvatarFile: file})
	if err != nil {
		a.RespJsonExit(r, err)
	}
	a.RespJsonExit(r, nil, avatar)
}
