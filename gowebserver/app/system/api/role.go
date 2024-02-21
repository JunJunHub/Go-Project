package api

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/gconv"
	"gowebserver/app/common/api"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/common/utils"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"gowebserver/app/system/service"
	"strings"
)

var RoleApi = new(roleApi)

type roleApi struct {
	api.BaseController
}

func (a *roleApi) Init(r *ghttp.Request) {
	a.Module = "角色管理"
	r.SetCtxVar(global.Module, a.Module)
}

// GetList
// @summary 查询角色列表
func (a *roleApi) GetList(r *ghttp.Request) {
	var req *define.SysRoleSelectPageReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}

	//分页查询角色列表
	pageInfo, roleList, err := service.InterfaceSysRole().GetListSearch(r.Context(), req)
	if err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}

	//请求返回信息
	var resp define.SysRoleSearchList
	if err = gconv.Struct(roleList, &resp.List); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInternalError, err))
	}
	resp.Total = pageInfo.Total
	resp.Page = pageInfo.PageNum
	resp.Size = pageInfo.PageSize

	//默认角色名称翻译成英文
	if "en" == utils.ParseAcceptLanguage(r.GetHeader("Accept-Language")) {
		for idx, role := range resp.List {
			if service.InterfaceSysRole().CheckIsDefaultSysRole(role.RoleID) {
				switch role.RoleName {
				case "系统管理员":
					resp.List[idx].RoleName = "Administrator"
				case "操作员":
					resp.List[idx].RoleName = "Operator"
				}
			}
		}
	}
	a.RespJsonExit(r, nil, resp)
}

// GetInfo
// @summary 查询角色信息
func (a *roleApi) GetInfo(r *ghttp.Request) {
	roleId := r.GetInt64("roleId")
	if roleId == 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleInvalidRoleId))
	}

	//查询角色信息
	roleInfo, err := service.InterfaceSysRole().Get(r.Context(), roleId)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if roleInfo == nil {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleInvalidRoleId))
	}

	//请求返回信息
	var resp define.SysRoleInfo
	if err = gconv.Struct(roleInfo, &resp); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInternalError, err))
	}
	a.RespJsonExit(r, nil, resp)
}

// Add
// @summary 新增角色信息
func (a *roleApi) Add(r *ghttp.Request) {
	var req *define.SysRoleAddReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}
	if req.RoleDataScope == 0 {
		req.RoleDataScope = 2
	}

	//角色信息
	var roleInfo model.SysRole
	_ = gconv.Struct(req, &roleInfo)
	roleInfo.CreateBy = service.Context.GetUserName(r.GetCtx())
	roleInfo.CreateTime = gtime.Now()
	if err := service.InterfaceSysRole().Add(r.GetCtx(), &roleInfo); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}

	a.RespJsonExit(r, nil, roleInfo)
}

// Edit
// @summary 编辑角色信息
func (a *roleApi) Edit(r *ghttp.Request) {
	var req *define.SysRoleEditReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}

	//角色信息
	roleInfoOld, err := service.InterfaceSysRole().Get(r.GetCtx(), req.RoleID)
	if err != nil {
		a.RespJsonExit(r, err)
	}

	var roleInfoNew model.SysRole
	_ = gconv.Struct(req, &roleInfoNew)
	roleInfoNew.UpdateBy = service.Context.GetUserName(r.GetCtx())
	roleInfoNew.UpdateTime = gtime.Now()

	if err = service.InterfaceSysRole().Put(r.GetCtx(), &roleInfoNew); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}

	if roleInfoOld.RoleName != roleInfoNew.RoleName {
		//更新角色绑定输出通道资源分组名称
		//mpuSrv.MpuChannelGroup.UpdateRoleBindOutputChnGroup(r.GetCtx(), &roleInfoNew, roleInfoOld.RoleName, false)
	}
	a.RespJsonExit(r, nil)
}

// Delete
// @summary 删除角色信息
func (a *roleApi) Delete(r *ghttp.Request) {
	var req *define.SysRoleDeleteReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}

	params := strings.Split(req.RoleIds, ",")
	delRoleIds := make([]int64, 0, len(params))
	for _, roleId := range params {
		//校验删除角色信息
		roleInfo, err := service.InterfaceSysRole().Get(r.GetCtx(), gconv.Int64(roleId))
		if err != nil {
			a.RespJsonExit(r, err)
		}
		if roleInfo == nil {
			a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleInvalidRoleId))
		}
		if service.InterfaceSysRole().CheckIsDefaultSysRole(gconv.Int64(roleId)) {
			a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleDelDefaultRoleForbid))
		}
		delRoleIds = append(delRoleIds, gconv.Int64(roleId))

		//更新角色绑定输出通道资源分组状态为可删除
		//mpuSrv.MpuChannelGroup.UpdateRoleBindOutputChnGroup(r.GetCtx(), roleInfo, roleInfo.RoleName, true)
	}

	if err := service.InterfaceSysRole().Del(r.GetCtx(), delRoleIds); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}
	a.RespJsonExit(r, nil)
}

// GetRoleMenu
// @summary 获取角色菜单权限
func (a *roleApi) GetRoleMenu(r *ghttp.Request) {
	roleId := r.GetInt64("roleId")
	if roleId == 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}

	var resp define.SysRoleMenuList
	if menuIds, err := service.InterfaceSysRole().GetRoleMenuPolicy(r.GetCtx(), roleId); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	} else {
		resp.RoleID = roleId

		//过滤无效规则
		validMenuList, _ := service.InterfaceSysMenu().GetMenuListByMenuIds(r.GetCtx(), menuIds)
		for _, menu := range validMenuList {
			resp.MenuIDS = append(resp.MenuIDS, menu.Id)
		}
		a.RespJsonExit(r, nil, resp)
	}
}

// EditRoleMenu
// @summary 编辑角色菜单权限
func (a *roleApi) EditRoleMenu(r *ghttp.Request) {
	var req *define.SysRoleMenuList
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}

	//校验角色信息
	roleInfo, err := service.InterfaceSysRole().Get(r.GetCtx(), req.RoleID)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if roleInfo.RoleDataScope == 1 {
		//拥有全部数据权限,不需要编辑权限规则
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleNotSupportedCustomPermissions))
	}

	//判断菜单数据是否有效
	for idx1 := range req.MenuIDS {
		for idx2 := range req.MenuIDS {
			if idx1 != idx2 && req.MenuIDS[idx1] == req.MenuIDS[idx2] {
				//参数中包含重复绑定菜单数据
				a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindMenuDuplicate))
			}
		}
	}
	menuList, _ := service.InterfaceSysMenu().GetMenuListByMenuIds(r.GetCtx(), req.MenuIDS)
	if len(menuList) != len(req.MenuIDS) {
		//参数中包含无效的菜单ID
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindMenuInvalid))
	}

	//更新角色绑定菜单的规则
	if err = service.InterfaceSysRole().UpdateRoleMenuPolicy(r.GetCtx(), req.RoleID, req.MenuIDS); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}
	service.InterfaceSysRole().NotifyDataPermissionUpdate(r.GetCtx(), req.RoleID, define.EPermissionTypeMenu)

	//通知冗余机箱刷新权限规则
	//mpuSrv.InterfaceMpuBoxRedundancy().NotifyPeerBoxUpdateCasbinRule()

	a.RespJsonExit(r, nil)
}

/*
// GetRoleChnGroup
// @summary 查询角色通道资源权限
func (a *roleApi) GetRoleChnGroup(r *ghttp.Request) {
	roleId := r.GetInt64("roleId")
	if roleId == 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}

	var resp define.SysRoleChnGroupList
	resp.RoleID = roleId
	chnGroupIDS, err := mpuSrv.MpuChannelGroup.GetRoleChnGroupIds(r.GetCtx(), roleId)
	if err != nil {
		a.RespJsonExit(r, err)
	}

	//过滤无效规则
	validChnGroupList, _ := mpuSrv.MpuChannelGroup.GetChnGroupListByChnGroupIds(r.GetCtx(), chnGroupIDS)
	for _, group := range validChnGroupList {
		resp.ChnGroupIDS = append(resp.ChnGroupIDS, group.GroupId)
	}
	a.RespJsonExit(r, nil, resp)
}

// EditRoleChnGroup
// @summary 编辑角色通道资源权限
func (a *roleApi) EditRoleChnGroup(r *ghttp.Request) {
	var req *define.SysRoleChnGroupList
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}

	//校验角色信息
	roleInfo, err := service.InterfaceSysRole().Get(r.GetCtx(), req.RoleID)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if roleInfo.RoleDataScope == 1 {
		//拥有全部数据权限,不需要编辑权限规则
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleNotSupportedCustomPermissions))
	}

	//判断通道分组数据是否有效
	for idx1 := range req.ChnGroupIDS {
		for idx2 := range req.ChnGroupIDS {
			if idx1 != idx2 && req.ChnGroupIDS[idx1] == req.ChnGroupIDS[idx2] {
				//参数中包含重复绑定通道分组数据
				a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindChnGroupDuplicate))
			}
		}
	}
	chnGroupList, _ := mpuSrv.MpuChannelGroup.GetChnGroupListByChnGroupIds(r.GetCtx(), req.ChnGroupIDS)
	if len(chnGroupList) != len(req.ChnGroupIDS) {
		//参数中包含无效的通道分组ID
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindChnGroupInvalid))
	}

	//更新角色绑定通道分组资源规则
	if err = service.InterfaceSysRole().UpdateRoleChnGroupPolicy(r.GetCtx(), req.RoleID, req.ChnGroupIDS); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}
	service.InterfaceSysRole().NotifyDataPermissionUpdate(r.GetCtx(), req.RoleID, define.EPermissionTypeChnGroup)

	//通知冗余机箱刷新权限规则
	mpuSrv.InterfaceMpuBoxRedundancy().NotifyPeerBoxUpdateCasbinRule()

	a.RespJsonExit(r, nil)
}

// GetRoleTvWall
// @summary 查询角色大屏资源权限
func (a *roleApi) GetRoleTvWall(r *ghttp.Request) {
	roleId := r.GetInt64("roleId")
	if roleId == 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}

	var resp define.SysRoleTvWallList
	resp.RoleID = roleId
	tvWallIds, err := mpuSrv.MpuTvWall.GetRoleTvWallIds(r.GetCtx(), roleId)
	if err != nil {
		a.RespJsonExit(r, err)
	}

	//过滤无效规则
	validTvWallList, _ := mpuSrv.MpuTvWall.GetTvWallListByTvWallIds(r.GetCtx(), tvWallIds)
	for _, wall := range validTvWallList {
		resp.TvWallIDS = append(resp.TvWallIDS, wall.TvWallId)
	}
	a.RespJsonExit(r, nil, resp)
}

// EditRoleTvWall
// 编辑角色大屏资源权限
func (a *roleApi) EditRoleTvWall(r *ghttp.Request) {
	var req *define.SysRoleTvWallList
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}

	//校验参数
	if err := a.checkRoleBindTvWallReq(r.GetCtx(), req); err != nil {
		a.RespJsonExit(r, err)
	}

	//更新角色绑定大屏资源规则
	if err := service.InterfaceSysRole().UpdateRoleTvWallPolicy(r.GetCtx(), req.RoleID, req.TvWallIDS); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}
	service.InterfaceSysRole().NotifyDataPermissionUpdate(r.GetCtx(), req.RoleID, define.EPermissionTypeTvWall)

	//通知冗余机箱刷新权限规则
	mpuSrv.InterfaceMpuBoxRedundancy().NotifyPeerBoxUpdateCasbinRule()

	a.RespJsonExit(r, nil)
}

// GetRoleMatrixScheme
// @summary 查询角色绑定矩阵预案资源
func (a *roleApi) GetRoleMatrixScheme(r *ghttp.Request) {
	roleId := r.GetInt64("roleId")
	if roleId == 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}

	var resp define.SysRoleMatrixSchemeList
	resp.RoleID = roleId
	matrixSchemeIds, err := mpuSrv.InterfaceMpuMatrixDispatchScheme().GetRoleMatrixSchemeIds(r.GetCtx(), roleId)
	if err != nil {
		a.RespJsonExit(r, err)
	}

	//过滤无效规则
	validSchemeList, _ := mpuSrv.InterfaceMpuMatrixDispatchScheme().GetMatrixSchemeListByMatrixSchemeIds(r.GetCtx(), matrixSchemeIds)
	for _, scheme := range validSchemeList {
		resp.MatrixSchemeIDS = append(resp.MatrixSchemeIDS, scheme.SchemeId)
	}
	a.RespJsonExit(r, nil, resp)
}

// EditRoleMatrixScheme
// @summary 编辑角色绑定矩阵预案资源
func (a *roleApi) EditRoleMatrixScheme(r *ghttp.Request) {
	var req *define.SysRoleMatrixSchemeList
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}

	//校验角色信息
	roleInfo, err := service.InterfaceSysRole().Get(r.GetCtx(), req.RoleID)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if roleInfo.RoleDataScope == 1 {
		//拥有全部数据权限,不需要编辑权限规则
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleNotSupportedCustomPermissions))
	}

	//判断矩阵预案资源是否有效
	for idx1 := range req.MatrixSchemeIDS {
		for idx2 := range req.MatrixSchemeIDS {
			if idx1 != idx2 && req.MatrixSchemeIDS[idx1] == req.MatrixSchemeIDS[idx2] {
				//参数中包含重复绑定预案数据
				a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindMatrixSchemeDuplicate))
			}
		}
	}
	matrixSchemeList, _ := mpuSrv.InterfaceMpuMatrixDispatchScheme().GetMatrixSchemeListByMatrixSchemeIds(r.GetCtx(), req.MatrixSchemeIDS)
	if len(matrixSchemeList) != len(req.MatrixSchemeIDS) {
		//参数中包含无效的预案资源ID
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindMatrixSchemeInvalid))
	}

	//更新角色绑定矩阵预案资源规则
	if err = service.InterfaceSysRole().UpdateRoleMatrixSchemePolicy(r.GetCtx(), req.RoleID, req.MatrixSchemeIDS); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}

	//自动分配关联大屏、输出通道权限
	var usedVideoOutputChnList []string
	var usedAudioOutputChnList []string
	var usedTvWallIds []string
	//获取已分配矩阵预案使用的音视频输出通道及大屏资源
	matrixDispatchList, _ := mpuSrv.InterfaceMpuMatrixDispatch().GetListBySchemeIds(r.GetCtx(), req.MatrixSchemeIDS)
	for _, dispatchInfo := range matrixDispatchList {
		//预案使用的大屏资源
		if dispatchInfo.DispatchType == int(defMpu.EMatrixDispatchTvWall) {
			var bDuplicationWallId bool
			for _, wallId := range usedTvWallIds {
				if wallId == dispatchInfo.DstChannelId {
					bDuplicationWallId = true
					break
				}
			}
			if !bDuplicationWallId {
				usedTvWallIds = append(usedTvWallIds, dispatchInfo.DstChannelId)
			}
		}

		//预案使用的视频输出通道
		if dispatchInfo.DispatchType == int(defMpu.EMatrixDispatchVideo) || dispatchInfo.DispatchType == int(defMpu.EMatrixDispatchVideoAndAudio) {
			usedVideoOutputChnList = append(usedVideoOutputChnList, dispatchInfo.DstChannelId)
		}

		//预案使用的音频输出通道
		if dispatchInfo.DispatchType == int(defMpu.EMatrixDispatchAudio) || dispatchInfo.DispatchType == int(defMpu.EMatrixDispatchAudioMixer) {
			usedAudioOutputChnList = append(usedAudioOutputChnList, dispatchInfo.DstChannelId)
		}
	}
	//自动给角色绑定输出通道资源权限
	roleVideoChnGroup, roleAudioChnGroup := mpuSrv.MpuChannelGroup.GetRoleBindOutputChnGroup(r.GetCtx(), roleInfo, true)
	if roleVideoChnGroup != nil {
		err = mpuSrv.MpuChannelGroup.AddGroupMember(&defMpu.MPUChannelGroupMemberCfgReq{
			GroupId:       roleVideoChnGroup.GroupId,
			ChannelIdList: usedVideoOutputChnList,
		}, true)
		if err != nil {
			g.Log().Error(err)
		}
	}
	if roleAudioChnGroup != nil {
		err = mpuSrv.MpuChannelGroup.AddGroupMember(&defMpu.MPUChannelGroupMemberCfgReq{
			GroupId:       roleAudioChnGroup.GroupId,
			ChannelIdList: usedAudioOutputChnList,
		}, true)
		if err != nil {
			g.Log().Error(err)
		}
	}

	//自动给角色分配预案使用的大屏资源权限
	err = a.checkRoleBindTvWallReq(r.GetCtx(), &define.SysRoleTvWallList{
		RoleID:    req.RoleID,
		TvWallIDS: usedTvWallIds,
	})
	if err == nil && len(usedTvWallIds) > 0 {
		_ = service.InterfaceSysRole().AddRoleTvWallPolicy(r.GetCtx(), roleInfo.RoleId, usedTvWallIds)
	}

	//触发权限更新通知
	service.InterfaceSysRole().NotifyDataPermissionUpdate(r.GetCtx(), req.RoleID, define.EPermissionTypeMatrixScheme)

	//通知冗余机箱刷新权限规则
	mpuSrv.InterfaceMpuBoxRedundancy().NotifyPeerBoxUpdateCasbinRule()

	a.RespJsonExit(r, nil)
}

// GetRoleMatrixSchemeGroup
// @summary 查询角色绑定矩阵预案分组资源
func (a *roleApi) GetRoleMatrixSchemeGroup(r *ghttp.Request) {
	roleId := r.GetInt64("roleId")
	if roleId == 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}

	var resp define.SysRoleMatrixSchemeGroupList
	resp.RoleID = roleId
	matrixSchemeGroupIds, err := mpuSrv.InterfaceMpuMatrixDispatchSchemeGroup().GetRoleMatrixSchemeGroupIds(r.GetCtx(), roleId)
	if err != nil {
		a.RespJsonExit(r, err)
	}

	//过滤无效规则
	validSchemeGroupList, _ := mpuSrv.InterfaceMpuMatrixDispatchSchemeGroup().GetMatrixSchemeGroupListByGroupIds(r.GetCtx(), matrixSchemeGroupIds)
	for _, group := range validSchemeGroupList {
		resp.MatrixSchemeGroupIdS = append(resp.MatrixSchemeGroupIdS, group.GroupId)
	}
	a.RespJsonExit(r, nil, resp)
}

// EditRoleMatrixSchemeGroup
// @summary 编辑角色绑定矩阵预案分组资源
func (a *roleApi) EditRoleMatrixSchemeGroup(r *ghttp.Request) {
	var req *define.SysRoleMatrixSchemeGroupList
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}

	//校验角色信息
	roleInfo, err := service.InterfaceSysRole().Get(r.GetCtx(), req.RoleID)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if roleInfo.RoleDataScope == 1 {
		//拥有全部数据权限,不需要编辑权限规则
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleNotSupportedCustomPermissions))
	}

	//判断矩阵预案资源是否有效
	for idx1 := range req.MatrixSchemeGroupIdS {
		for idx2 := range req.MatrixSchemeGroupIdS {
			if idx1 != idx2 && req.MatrixSchemeGroupIdS[idx1] == req.MatrixSchemeGroupIdS[idx2] {
				//参数中包含重复绑定预案数据
				a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindMatrixSchemeGroupDuplicate))
			}
		}
	}
	matrixSchemeGroupList, _ := mpuSrv.InterfaceMpuMatrixDispatchSchemeGroup().GetMatrixSchemeGroupListByGroupIds(r.GetCtx(), req.MatrixSchemeGroupIdS)
	if len(matrixSchemeGroupList) != len(req.MatrixSchemeGroupIdS) {
		//参数中包含无效的预案资源ID
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindMatrixSchemeGroupInvalid))
	}

	//更新角色绑定矩阵预案分组资源规则
	if err = service.InterfaceSysRole().UpdateRoleMatrixSchemeGroupPolicy(r.GetCtx(), req.RoleID, req.MatrixSchemeGroupIdS); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}
	service.InterfaceSysRole().NotifyDataPermissionUpdate(r.GetCtx(), req.RoleID, define.EPermissionTypeMatrixSchemeGroup)

	//通知冗余机箱刷新权限规则
	mpuSrv.InterfaceMpuBoxRedundancy().NotifyPeerBoxUpdateCasbinRule()

	a.RespJsonExit(r, nil)
}

// GetRoleMeetDispatch
// @summary 查询角色绑定会议调度资源
func (a *roleApi) GetRoleMeetDispatch(r *ghttp.Request) {
	roleId := r.GetInt64("roleId")
	if roleId == 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}

	var resp define.SysRoleMeetDispatchList
	resp.RoleID = roleId
	meetDispatchIds, err := mpuSrv.InterfaceMpuMeetDispatch().GetRoleMeetDispatchIds(r.GetCtx(), roleId)
	if err != nil {
		a.RespJsonExit(r, err)
	}

	//过滤无效规则
	validMeetDspList, _ := mpuSrv.InterfaceMpuMeetDispatch().GetMeetDispatchListByMeetDispatchIds(r.GetCtx(), meetDispatchIds)
	for _, meet := range validMeetDspList {
		resp.MeetDispatchIds = append(resp.MeetDispatchIds, meet.MeetDispatchId)
	}
	a.RespJsonExit(r, nil, resp)
}

// EditRoleMeetDispatch
// @summary 编辑角色绑定会议调度资源
func (a *roleApi) EditRoleMeetDispatch(r *ghttp.Request) {
	var req *define.SysRoleMeetDispatchList
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}

	//校验角色信息
	roleInfo, err := service.InterfaceSysRole().Get(r.GetCtx(), req.RoleID)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if roleInfo.RoleDataScope == 1 {
		//拥有全部数据权限,不需要编辑权限规则
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleNotSupportedCustomPermissions))
	}

	//判断会议调度资源是否有效
	for idx1 := range req.MeetDispatchIds {
		for idx2 := range req.MeetDispatchIds {
			if idx1 != idx2 && req.MeetDispatchIds[idx1] == req.MeetDispatchIds[idx2] {
				//参数中包含重复绑定会议调度数据
				a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindMeetDispatchDuplicate))
			}
		}
	}
	meetDspList, _ := mpuSrv.InterfaceMpuMeetDispatch().GetMeetDispatchListByMeetDispatchIds(r.GetCtx(), req.MeetDispatchIds)
	if len(meetDspList) != len(req.MeetDispatchIds) {
		//参数中包含无效的会议资源ID
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindMeetDispatchInvalid))
	}

	//更新角色绑定会议调度资源规则
	if err = service.InterfaceSysRole().UpdateRoleMeetDispatchPolicy(r.GetCtx(), req.RoleID, req.MeetDispatchIds); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}

	//自动分配关联矩阵预案、大屏、输出通道权限
	var usedMatrixSchemeIds []int64
	var usedVideoOutputChnList []string
	var usedAudioOutputChnList []string
	var usedTvWallIds []string

	//获取会议调度中使用的矩阵预案
	meetDspList, _ = mpuSrv.InterfaceMpuMeetDispatch().GetMeetDispatchListByMeetDispatchIds(r.GetCtx(), req.MeetDispatchIds)
	for _, meet := range meetDspList {
		usedMatrixSchemeIds = append(usedMatrixSchemeIds, meet.MatrixSchemeId)
	}

	//获取矩阵预案中使用的大屏、输出通道
	matrixDispatchList, _ := mpuSrv.InterfaceMpuMatrixDispatch().GetListBySchemeIds(r.GetCtx(), usedMatrixSchemeIds)
	for _, dispatchInfo := range matrixDispatchList {
		//预案使用的大屏资源
		if dispatchInfo.DispatchType == int(defMpu.EMatrixDispatchTvWall) {
			var bDuplicationWallId bool
			for _, wallId := range usedTvWallIds {
				if wallId == dispatchInfo.DstChannelId {
					bDuplicationWallId = true
					break
				}
			}
			if !bDuplicationWallId {
				usedTvWallIds = append(usedTvWallIds, dispatchInfo.DstChannelId)
			}
		}

		//预案使用的视频输出通道
		if dispatchInfo.DispatchType == int(defMpu.EMatrixDispatchVideo) || dispatchInfo.DispatchType == int(defMpu.EMatrixDispatchVideoAndAudio) {
			usedVideoOutputChnList = append(usedVideoOutputChnList, dispatchInfo.DstChannelId)
		}

		//预案使用的音频输出通道
		if dispatchInfo.DispatchType == int(defMpu.EMatrixDispatchAudio) || dispatchInfo.DispatchType == int(defMpu.EMatrixDispatchAudioMixer) {
			usedAudioOutputChnList = append(usedAudioOutputChnList, dispatchInfo.DstChannelId)
		}
	}

	//自动给角色绑定矩阵预案权限
	if len(usedMatrixSchemeIds) > 0 {
		_ = service.InterfaceSysRole().AddRoleMatrixSchemePolicy(r.GetCtx(), roleInfo.RoleId, usedMatrixSchemeIds)
	}

	//自动给角色绑定预案中使用输出通道资源权限
	roleVideoChnGroup, roleAudioChnGroup := mpuSrv.MpuChannelGroup.GetRoleBindOutputChnGroup(r.GetCtx(), roleInfo, true)
	if roleVideoChnGroup != nil {
		err = mpuSrv.MpuChannelGroup.AddGroupMember(&defMpu.MPUChannelGroupMemberCfgReq{
			GroupId:       roleVideoChnGroup.GroupId,
			ChannelIdList: usedVideoOutputChnList,
		}, true)
		if err != nil {
			g.Log().Error(err)
		}
	}
	if roleAudioChnGroup != nil {
		err = mpuSrv.MpuChannelGroup.AddGroupMember(&defMpu.MPUChannelGroupMemberCfgReq{
			GroupId:       roleAudioChnGroup.GroupId,
			ChannelIdList: usedAudioOutputChnList,
		}, true)
		if err != nil {
			g.Log().Error(err)
		}
	}

	//自动给角色分配预案使用的大屏资源权限
	err = a.checkRoleBindTvWallReq(r.GetCtx(), &define.SysRoleTvWallList{
		RoleID:    req.RoleID,
		TvWallIDS: usedTvWallIds,
	})
	if err == nil && len(usedTvWallIds) > 0 {
		_ = service.InterfaceSysRole().AddRoleTvWallPolicy(r.GetCtx(), roleInfo.RoleId, usedTvWallIds)
	}

	//触发权限更新通知
	service.InterfaceSysRole().NotifyDataPermissionUpdate(r.GetCtx(), req.RoleID, define.EPermissionTypeMeetDispatch)

	//通知冗余机箱刷新权限规则
	mpuSrv.InterfaceMpuBoxRedundancy().NotifyPeerBoxUpdateCasbinRule()

	a.RespJsonExit(r, nil)
}

// GetRoleMeetDispatchGroup
// @summary 查询角色绑定会议调度分组资源
func (a *roleApi) GetRoleMeetDispatchGroup(r *ghttp.Request) {
	roleId := r.GetInt64("roleId")
	if roleId == 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}

	var resp define.SysRoleMeetDispatchGroupList
	resp.RoleID = roleId
	meetDispatchGroupIds, err := mpuSrv.InterfaceMpuMeetDispatchGroup().GetRoleMeetDispatchGroupIds(r.GetCtx(), roleId)
	if err != nil {
		a.RespJsonExit(r, err)
	}

	//过滤无效规则
	validMeetGroupList, _ := mpuSrv.InterfaceMpuMeetDispatchGroup().GetMeetDispatchGroupListByGroupIds(r.GetCtx(), meetDispatchGroupIds)
	for _, group := range validMeetGroupList {
		resp.MeetDispatchGroupIds = append(resp.MeetDispatchGroupIds, group.GroupId)
	}
	a.RespJsonExit(r, nil, resp)
}

// EditRoleMeetDispatchGroup
// @summary 编辑角色绑定会议调度分组资源
func (a *roleApi) EditRoleMeetDispatchGroup(r *ghttp.Request) {
	var req *define.SysRoleMeetDispatchGroupList
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}

	//校验角色信息
	roleInfo, err := service.InterfaceSysRole().Get(r.GetCtx(), req.RoleID)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if roleInfo.RoleDataScope == 1 {
		//拥有全部数据权限,不需要编辑权限规则
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleNotSupportedCustomPermissions))
	}

	//判断会议调度分组资源是否有效
	for idx1 := range req.MeetDispatchGroupIds {
		for idx2 := range req.MeetDispatchGroupIds {
			if idx1 != idx2 && req.MeetDispatchGroupIds[idx1] == req.MeetDispatchGroupIds[idx2] {
				//参数中包含重复绑定会议调度分组
				a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindMeetDispatchGroupDuplicate))
			}
		}
	}
	meetDispatchGroupList, _ := mpuSrv.InterfaceMpuMeetDispatchGroup().GetMeetDispatchGroupListByGroupIds(r.GetCtx(), req.MeetDispatchGroupIds)
	if len(meetDispatchGroupList) != len(req.MeetDispatchGroupIds) {
		//参数中包含无效的会议分组资源ID
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrSysRoleBindMeetDispatchGroupInvalid))
	}

	//更新角色绑定会议分组资源规则
	if err = service.InterfaceSysRole().UpdateRoleMeetDispatchGroupPolicy(r.GetCtx(), req.RoleID, req.MeetDispatchGroupIds); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}
	service.InterfaceSysRole().NotifyDataPermissionUpdate(r.GetCtx(), req.RoleID, define.EPermissionTypeMeetDispatchGroup)

	//通知冗余机箱刷新权限规则
	mpuSrv.InterfaceMpuBoxRedundancy().NotifyPeerBoxUpdateCasbinRule()

	a.RespJsonExit(r, nil)
}

// checkRoleBindTvWall
// @summary 校验角色绑定大屏资源请求参数
//          校验通过将大屏使用输出通道资源自动分配给角色
func (a *roleApi) checkRoleBindTvWallReq(ctx context.Context, req *define.SysRoleTvWallList) (err error) {
	//校验角色信息
	roleInfo, e := service.InterfaceSysRole().Get(ctx, req.RoleID)
	if e != nil {
		return e
	}
	if roleInfo.RoleDataScope == 1 {
		//拥有全部数据权限,不需要编辑权限规则
		return gerror.NewCode(errcode.ErrSysRoleNotSupportedCustomPermissions)
	}

	//判断大屏资源是否有效,并记录大屏绑定输出通道列表
	for idx1 := range req.TvWallIDS {
		for idx2 := range req.TvWallIDS {
			if idx1 != idx2 && req.TvWallIDS[idx1] == req.TvWallIDS[idx2] {
				//参数中包含重复绑定大屏数据
				return gerror.NewCode(errcode.ErrSysRoleBindTvWallDuplicate)
			}
		}
	}

	tvWallList, _ := mpuSrv.MpuTvWall.GetTvWallListByTvWallIds(ctx, req.TvWallIDS)
	if len(tvWallList) != len(req.TvWallIDS) {
		//参数中包含无效的大屏资源ID
		return gerror.NewCode(errcode.ErrSysRoleBindTvWallInvalid)
	}

	//获取已分配大屏使用的输出通道
	var outputChnList []string
	for _, tvWall := range tvWallList {
		for _, cell := range tvWall.Cells {
			outputChnList = append(outputChnList, cell.CellChannelId)
		}
	}
	//自动给角色绑定大屏输出通道
	roleVideoChnGroup, _ := mpuSrv.MpuChannelGroup.GetRoleBindOutputChnGroup(ctx, roleInfo, true)
	if roleVideoChnGroup != nil {
		err = mpuSrv.MpuChannelGroup.AddGroupMember(&defMpu.MPUChannelGroupMemberCfgReq{
			GroupId:       roleVideoChnGroup.GroupId,
			ChannelIdList: outputChnList,
		}, true)
		if err != nil {
			g.Log().Error(err)
		}
	}
	return
}
*/
