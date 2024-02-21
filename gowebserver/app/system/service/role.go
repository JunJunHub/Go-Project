package service

import (
	"context"
	"fmt"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	defCom "gowebserver/app/common/define"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	comSrv "gowebserver/app/common/service"
	"gowebserver/app/common/utils/page"
	"gowebserver/app/common/utils/valueCheck"
	"gowebserver/app/system/dao"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"gowebserver/app/system/service/publish"
	"strings"
)

//ISysRole 角色管理业务接口
type ISysRole interface {
	SysRoleInit()
	CheckIsDefaultSysRole(roleId int64) bool
	GetListSearch(ctx context.Context, req *define.SysRoleSelectPageReq) (pageInfo *page.Paging, roleList []*model.SysRole, err error)
	Get(ctx context.Context, roleId int64) (roleInfo *model.SysRole, err error)
	Add(ctx context.Context, roleInfo *model.SysRole) (err error)
	Put(ctx context.Context, roleInfo *model.SysRole) (err error)
	Del(ctx context.Context, roleIds []int64) (err error)
	GetRoles(ctx context.Context, roleIds []int64) (roles []model.SysRole, err error)

	GetRoleMenuPolicy(ctx context.Context, roleId int64) (menuIds []int64, err error)
	UpdateRoleMenuPolicy(ctx context.Context, roleId int64, menuIds []int64) (err error)
	AddRoleMenuPolicy(ctx context.Context, roleId int64, menuIds []int64) (err error)
	DelRoleMenuPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error)

	GetRoleChnGroupPolicy(ctx context.Context, roleId int64) (chnGroupIds []string, err error)
	UpdateRoleChnGroupPolicy(ctx context.Context, roleId int64, chnGroupIds []string) (err error)
	AddRoleChnGroupPolicy(ctx context.Context, roleId int64, chnGroupIds []string) (err error)
	DelRoleChnGroupPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error)

	GetRoleTvWallPolicy(ctx context.Context, roleId int64) (tvWallIds []string, err error)
	UpdateRoleTvWallPolicy(ctx context.Context, roleId int64, tvWallIds []string) (err error)
	AddRoleTvWallPolicy(ctx context.Context, roleId int64, tvWallIds []string) (err error)
	DelRoleTvWallPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error)

	GetRoleMatrixSchemePolicy(ctx context.Context, roleId int64) (matrixSchemeIds []int64, err error)
	UpdateRoleMatrixSchemePolicy(ctx context.Context, roleId int64, matrixSchemeIds []int64) (err error)
	AddRoleMatrixSchemePolicy(ctx context.Context, roleId int64, matrixSchemeIds []int64) (err error)
	DelRoleMatrixSchemePolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error)

	GetRoleMatrixSchemeGroupPolicy(ctx context.Context, roleId int64) (matrixSchemeGroupIds []int64, err error)
	UpdateRoleMatrixSchemeGroupPolicy(ctx context.Context, roleId int64, matrixSchemeGroupIds []int64) (err error)
	AddRoleMatrixSchemeGroupPolicy(ctx context.Context, roleId int64, matrixSchemeGroupIds []int64) (err error)
	DelRoleMatrixSchemeGroupPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error)

	GetRoleMeetDispatchPolicy(ctx context.Context, roleId int64) (meetDispatchIds []int64, err error)
	UpdateRoleMeetDispatchPolicy(ctx context.Context, roleId int64, meetDispatchIds []int64) (err error)
	AddRoleMeetDispatchPolicy(ctx context.Context, roleId int64, meetDispatchIds []int64) (err error)
	DelRoleMeetDispatchPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error)

	GetRoleMeetDispatchGroupPolicy(ctx context.Context, roleId int64) (meetDispatchGroupIds []int64, err error)
	UpdateRoleMeetDispatchGroupPolicy(ctx context.Context, roleId int64, meetDispatchGroupIds []int64) (err error)
	AddRoleMeetDispatchGroupPolicy(ctx context.Context, roleId int64, meetDispatchGroupIds []int64) (err error)
	DelRoleMeetDispatchGroupPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error)

	NotifyDataPermissionUpdate(ctx context.Context, roleId int64, permissionType define.EPermissionType)
}

type sysRoleServiceImpl struct{}

var roleService = sysRoleServiceImpl{}

func InterfaceSysRole() ISysRole {
	return &roleService
}

func (s *sysRoleServiceImpl) SysRoleInit() {
	//初始化默认角色信息
	s.createDefaultSysRole()
}

// createDefaultSysRole
// @summary 创建系统默认角色
func (s *sysRoleServiceImpl) createDefaultSysRole() {
	_ = g.Try(func() {
		var roleSysAdmin *model.SysRole //默认系统管理员角色
		var roleOperator *model.SysRole //默认普通操作员角色

		err := dao.SysRole.Ctx(context.TODO()).Where(dao.SysRole.Columns.RoleId, 1).FindScan(&roleSysAdmin)
		valueCheck.ErrIsNil(context.TODO(), err)
		if roleSysAdmin == nil {
			//初始化系统管理员角色信息
			roleSysAdmin = &model.SysRole{
				RoleId:        1,
				RoleName:      "系统管理员",
				RoleKey:       "SysAdmin",
				RoleSort:      0,
				RoleDataScope: 1,
				RoleStatus:    1,
				Remark:        "",
				CreateBy:      "SysAdmin",
				CreateTime:    gtime.Now(),
				UpdateBy:      "",
				UpdateTime:    nil,
			}
			_, err = dao.SysRole.Ctx(context.TODO()).Save(roleSysAdmin)
			valueCheck.ErrIsNil(context.TODO(), err)
		}

		err = dao.SysRole.Ctx(context.TODO()).Where(dao.SysRole.Columns.RoleId, 2).FindScan(&roleOperator)
		valueCheck.ErrIsNil(context.TODO(), err)
		if roleOperator == nil {
			//初始化普通操作员角色信息
			roleOperator = &model.SysRole{
				RoleId:        2,
				RoleName:      "操作员",
				RoleKey:       "Operator",
				RoleSort:      0,
				RoleDataScope: 2,
				RoleStatus:    1,
				Remark:        "",
				CreateBy:      "SysAdmin",
				CreateTime:    gtime.Now(),
				UpdateBy:      "",
				UpdateTime:    nil,
			}
			_, err = dao.SysRole.Ctx(context.TODO()).Save(roleOperator)
			valueCheck.ErrIsNil(context.TODO(), err)
		}
	})
}

// CheckIsDefaultSysRole
// @summary 校验是否为系统默认角色
//          系统内置角色ID范围[1,10]
func (s *sysRoleServiceImpl) CheckIsDefaultSysRole(roleId int64) bool {
	return roleId <= global.RoleIdSystemRetentionMax
}

// CheckRoleNameUnique
// @summary 校验角色名称是否重复
func (s *sysRoleServiceImpl) CheckRoleNameUnique(roleId int64, RoleName string) error {
	err := g.Try(func() {
		count, e := dao.SysRole.Ctx(context.TODO()).
			Where(dao.SysRole.Columns.RoleName, RoleName).
			WhereNot(dao.SysRole.Columns.RoleId, roleId).
			FindCount()
		valueCheck.ErrIsNil(context.TODO(), e, errcode.ErrCommonDbOperationError.Message())
		if count > 0 {
			valueCheck.ErrIsNil(context.TODO(), gerror.NewCode(errcode.ErrSysRoleNameDuplicate))
		}
	})
	return err
}

// GetListSearch
// @summary 搜索查询角色信息
func (s *sysRoleServiceImpl) GetListSearch(ctx context.Context, req *define.SysRoleSelectPageReq) (pageInfo *page.Paging, roleList []*model.SysRole, err error) {
	if req == nil {
		return nil, nil, gerror.NewCode(errcode.ErrCommonInvalidParameter)
	}

	var total int
	err = g.Try(func() {
		m := dao.SysRole.Ctx(ctx)
		if req.KeyWord != "" {
			//m = m.WhereLike(dao.SysRole.Columns.RoleName, "%"+req.KeyWord+"%")
			req.KeyWord = strings.Replace(req.KeyWord, "'", "''", -1)
			m = m.Where(fmt.Sprintf("locate('%s', %s) > 0", req.KeyWord, dao.SysRole.Columns.RoleName))
		}

		total, err = m.Count()
		valueCheck.ErrIsNil(ctx, err, errcode.ErrCommonDbOperationError.Message())

		if req.PageSize == 0 {
			req.PageSize = total
		}
		if req.OrderByColumn != "" && req.IsAsc != "" {
			req.OrderByColumn = gstr.CaseSnake(req.OrderByColumn)
			m = m.Order(req.OrderByColumn + " " + req.IsAsc)
		} else {
			m = m.OrderAsc(dao.SysRole.Columns.RoleId)
		}

		pageInfo = page.CreatePaging(req.PageNum, req.PageSize, total)
		m = m.Limit(pageInfo.StartNum, pageInfo.PageSize)
		err = m.Scan(&roleList)
		valueCheck.ErrIsNil(ctx, err, errcode.ErrCommonDbOperationError.Message())
	})
	return
}

// Get
// @summary 根据角色ID查询角色信息
func (s *sysRoleServiceImpl) Get(ctx context.Context, roleId int64) (roleInfo *model.SysRole, err error) {
	err = g.Try(func() {
		err = dao.SysRole.Ctx(ctx).Where(dao.SysRole.Columns.RoleId, roleId).FindScan(&roleInfo)
		valueCheck.ErrIsNil(ctx, err, errcode.ErrCommonDbOperationError.Message())

		valueCheck.ValueIsNil(roleInfo, errcode.ErrSysRoleInvalidRoleId.Message())
	})
	return
}

// Add
// @summary 新增角色
func (s *sysRoleServiceImpl) Add(ctx context.Context, roleInfo *model.SysRole) (err error) {
	if roleInfo == nil || roleInfo.RoleName == "" {
		return gerror.NewCode(errcode.ErrCommonInvalidParameter)
	}

	err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		err = g.Try(func() {
			//校验角色名称是否重复
			err = s.CheckRoleNameUnique(roleInfo.RoleId, roleInfo.RoleName)
			valueCheck.ErrIsNil(ctx, err)

			roleIdMax, _ := dao.SysRole.Ctx(ctx).TX(tx).Max(dao.SysRole.Columns.RoleId)
			if roleIdMax <= global.RoleIdSystemRetentionMax {
				roleInfo.RoleId = global.RoleIdSystemRetentionMax + 1
				_, err = dao.SysRole.Ctx(ctx).TX(tx).Insert(roleInfo)
				valueCheck.ErrIsNil(ctx, err, errcode.ErrCommonDbOperationError.Message())
			} else {
				roleId, e := dao.SysRole.Ctx(ctx).TX(tx).InsertAndGetId(roleInfo)
				valueCheck.ErrIsNil(ctx, e, errcode.ErrCommonDbOperationError.Message())
				roleInfo.RoleId = roleId
			}
		})
		return err
	})
	return
}

// Put
// @summary 修改角色信息
func (s *sysRoleServiceImpl) Put(ctx context.Context, roleInfoNew *model.SysRole) (err error) {
	err = g.Try(func() {
		err = s.CheckRoleNameUnique(roleInfoNew.RoleId, roleInfoNew.RoleName)
		valueCheck.ErrIsNil(ctx, err)

		role, e := s.Get(ctx, roleInfoNew.RoleId)
		valueCheck.ErrIsNil(ctx, e)
		valueCheck.ValueIsNil(role, errcode.ErrSysRoleInvalidRoleId.Message())

		//禁止修改系统内置角色信息
		if s.CheckIsDefaultSysRole(roleInfoNew.RoleId) {
			valueCheck.ErrIsNil(ctx, gerror.NewCode(errcode.ErrSysRoleModifyDefaultRoleForbid))
		}

		role.RoleName = roleInfoNew.RoleName
		role.RoleKey = roleInfoNew.RoleKey
		role.RoleSort = roleInfoNew.RoleSort
		role.RoleDataScope = roleInfoNew.RoleDataScope
		role.RoleStatus = roleInfoNew.RoleStatus
		role.Remark = roleInfoNew.Remark
		role.UpdateBy = roleInfoNew.UpdateBy
		role.UpdateTime = roleInfoNew.UpdateTime

		//更新保存数据
		_, e = dao.SysRole.Ctx(ctx).Save(role)
		valueCheck.ErrIsNil(ctx, e, errcode.ErrCommonDbOperationError.Message())
	})
	return
}

// Del
// @summary 删除角色
func (s *sysRoleServiceImpl) Del(ctx context.Context, roleIds []int64) (err error) {
	err = g.DB().Transaction(ctx, func(ctx context.Context, tx *gdb.TX) error {
		err = g.Try(func() {
			for _, roleId := range roleIds {
				roleInfo, e := s.Get(ctx, roleId)
				valueCheck.ErrIsNil(ctx, e)
				valueCheck.ValueIsNil(roleInfo, errcode.ErrSysRoleInvalidRoleId.Message())

				//禁止删除系统内置角色
				if s.CheckIsDefaultSysRole(gconv.Int64(roleId)) {
					valueCheck.ErrIsNil(ctx, gerror.NewCode(errcode.ErrSysRoleDelDefaultRoleForbid))
				}

				//禁止删除已关联了用户的角色
				var roleBindUserIds []int64
				roleBindUserIds, err = s.GetRoleBindUserIds(ctx, roleId)
				valueCheck.ErrIsNil(ctx, err)
				if len(roleBindUserIds) > 0 {
					valueCheck.ErrIsNil(ctx, gerror.NewCode(errcode.ErrSysRoleDelInUseRoleForbid))
				}

				//删除角色所有数据权限规则
				_ = s.DelRoleMenuPolicy(ctx, 0, gconv.String(roleId))
				_ = s.DelRoleChnGroupPolicy(ctx, 0, gconv.String(roleId))
				_ = s.DelRoleTvWallPolicy(ctx, 0, gconv.String(roleId))
				_ = s.DelRoleMatrixSchemePolicy(ctx, 0, gconv.String(roleId))

				_, err = dao.SysRole.Ctx(ctx).TX(tx).Where(dao.SysRole.Columns.RoleId, roleId).Unscoped().Delete()
				valueCheck.ErrIsNil(ctx, err, errcode.ErrCommonDbOperationError.Message())
			}
		})
		return err
	})
	return
}

// GetRoles
// @summary 批量查询角色信息
func (s *sysRoleServiceImpl) GetRoles(ctx context.Context, roleIds []int64) (roles []model.SysRole, err error) {
	err = g.Try(func() {
		err = dao.SysRole.Ctx(ctx).WhereIn(dao.SysRole.Columns.RoleId, roleIds).FindScan(&roles)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// GetRoleBindUserIds
// @summary 获取角色关联用户Id
func (s *sysRoleServiceImpl) GetRoleBindUserIds(ctx context.Context, roleId int64) (userIds []int64, err error) {
	err = g.Try(func() {
		var roleInfo *model.SysRole
		roleInfo, err = s.Get(ctx, roleId)
		valueCheck.ErrIsNil(ctx, err)
		valueCheck.ValueIsNil(roleInfo, errcode.ErrSysRoleInvalidRoleId.Message())

		//获取已关联用户ID
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := enforcer.GetFilteredGroupingPolicy(1, gconv.String(roleId))
		userIds = make([]int64, len(casbinRules))
		for k, v := range casbinRules {
			userIds[k] = gconv.Int64(v[0])
		}
	})
	return
}

// GetRoleMenuPolicy
// @summary 获取角色关联菜单规则
func (s *sysRoleServiceImpl) GetRoleMenuPolicy(ctx context.Context, roleId int64) (menuIds []int64, err error) {
	err = g.Try(func() {
		var roleInfo *model.SysRole
		roleInfo, err = s.Get(ctx, roleId)
		valueCheck.ErrIsNil(ctx, err)
		valueCheck.ValueIsNil(roleInfo, errcode.ErrSysRoleInvalidRoleId.Message())
		if roleInfo.RoleDataScope == 1 {
			//拥有全部数据权限
			menuIds, err = InterfaceSysMenu().GetAllMenuIds(ctx)
			return
		}

		//获取已分配菜单数据权限
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := enforcer.GetFilteredNamedPolicy(global.CasbinPTypeMenu, 0, gconv.String(roleId))
		menuIds = make([]int64, len(casbinRules))
		for k, v := range casbinRules {
			menuIds[k] = gconv.Int64(v[1])
		}
	})
	return
}

// UpdateRoleMenuPolicy
// @summary 更新角色关联菜单规则
//          tips: 先删除旧规则,再添加新规则
func (s *sysRoleServiceImpl) UpdateRoleMenuPolicy(ctx context.Context, roleId int64, menuIds []int64) (err error) {
	err = g.Try(func() {
		err = s.DelRoleMenuPolicy(ctx, 0, gconv.String(roleId))
		valueCheck.ErrIsNil(ctx, err)

		err = s.AddRoleMenuPolicy(ctx, roleId, menuIds)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// DelRoleMenuPolicy
// @summary 删除角色关联菜单规则
func (s *sysRoleServiceImpl) DelRoleMenuPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)

		_, e = enforcer.RemoveFilteredNamedPolicy(global.CasbinPTypeMenu, fieldIndex, fieldValue)
		valueCheck.ErrIsNil(ctx, e)
	})
	return
}

// AddRoleMenuPolicy
// @summary 添加角色关联菜单规则
func (s *sysRoleServiceImpl) AddRoleMenuPolicy(ctx context.Context, roleId int64, menuIds []int64) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := gconv.Strings(menuIds)
		for _, v := range casbinRules {
			_, err = enforcer.AddNamedPolicy(global.CasbinPTypeMenu, gconv.String(roleId), v, "All")
			valueCheck.ErrIsNil(ctx, err)
		}
	})
	return
}

// GetRoleChnGroupPolicy
// @summary 获取角色关联通道分组资源规则
func (s *sysRoleServiceImpl) GetRoleChnGroupPolicy(ctx context.Context, roleId int64) (chnGroupIds []string, err error) {
	err = g.Try(func() {
		//获取已分配通道分组资源权限
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := enforcer.GetFilteredNamedPolicy(global.CasbinPTypeChnGroup, 0, gconv.String(roleId))
		chnGroupIds = make([]string, len(casbinRules))
		for k, v := range casbinRules {
			chnGroupIds[k] = v[1]
		}
	})
	return
}

// UpdateRoleChnGroupPolicy
// @summary 更新角色关联通道分组资源规则
//          tips: 先删除旧规则,再添加新规则
func (s *sysRoleServiceImpl) UpdateRoleChnGroupPolicy(ctx context.Context, roleId int64, chnGroupIds []string) (err error) {
	err = g.Try(func() {
		err = s.DelRoleChnGroupPolicy(ctx, 0, gconv.String(roleId))
		valueCheck.ErrIsNil(ctx, err)

		err = s.AddRoleChnGroupPolicy(ctx, roleId, chnGroupIds)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// DelRoleChnGroupPolicy
// @summary 删除角色关联通道分组资源规则
func (s *sysRoleServiceImpl) DelRoleChnGroupPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)

		_, e = enforcer.RemoveFilteredNamedPolicy(global.CasbinPTypeChnGroup, fieldIndex, fieldValue)
		valueCheck.ErrIsNil(ctx, e)
	})
	return
}

// AddRoleChnGroupPolicy
// @summary 添加角色关联通道分组资源规则
func (s *sysRoleServiceImpl) AddRoleChnGroupPolicy(ctx context.Context, roleId int64, chnGroupIds []string) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := gconv.Strings(chnGroupIds)
		for _, v := range casbinRules {
			_, err = enforcer.AddNamedPolicy(global.CasbinPTypeChnGroup, gconv.String(roleId), v, "All")
			valueCheck.ErrIsNil(ctx, err)
		}
	})
	return
}

// GetRoleTvWallPolicy
// @summary 获取角色关联大屏规则
func (s *sysRoleServiceImpl) GetRoleTvWallPolicy(ctx context.Context, roleId int64) (tvWallIds []string, err error) {
	err = g.Try(func() {
		//获取已分配大屏数据权限
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := enforcer.GetFilteredNamedPolicy(global.CasbinPTypeTvWall, 0, gconv.String(roleId))
		tvWallIds = make([]string, len(casbinRules))
		for k, v := range casbinRules {
			tvWallIds[k] = v[1]
		}
	})
	return
}

// UpdateRoleTvWallPolicy
// @summary 更新角色关联大屏规则
//          tips: 先删除旧规则,再添加新规则
func (s *sysRoleServiceImpl) UpdateRoleTvWallPolicy(ctx context.Context, roleId int64, tvWallIds []string) (err error) {
	err = g.Try(func() {
		err = s.DelRoleTvWallPolicy(ctx, 0, gconv.String(roleId))
		valueCheck.ErrIsNil(ctx, err)

		err = s.AddRoleTvWallPolicy(ctx, roleId, tvWallIds)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// DelRoleTvWallPolicy
// @summary 删除角色关联大屏规则
func (s *sysRoleServiceImpl) DelRoleTvWallPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)

		_, e = enforcer.RemoveFilteredNamedPolicy(global.CasbinPTypeTvWall, fieldIndex, fieldValue)
		valueCheck.ErrIsNil(ctx, e)
	})
	return
}

// AddRoleTvWallPolicy
// @summary 添加角色关联大屏规则
func (s *sysRoleServiceImpl) AddRoleTvWallPolicy(ctx context.Context, roleId int64, tvWallIds []string) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := tvWallIds
		for _, v := range casbinRules {
			_, err = enforcer.AddNamedPolicy(global.CasbinPTypeTvWall, gconv.String(roleId), v, "All")
			valueCheck.ErrIsNil(ctx, err)
		}
	})
	return
}

// GetRoleMatrixSchemePolicy
// @summary 获取角色关联预案资源规则
func (s *sysRoleServiceImpl) GetRoleMatrixSchemePolicy(ctx context.Context, roleId int64) (matrixSchemeIds []int64, err error) {
	err = g.Try(func() {
		//获取已分配预案资源权限
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := enforcer.GetFilteredNamedPolicy(global.CasbinPTypeMatrixScheme, 0, gconv.String(roleId))
		matrixSchemeIds = make([]int64, len(casbinRules))
		for k, v := range casbinRules {
			matrixSchemeIds[k] = gconv.Int64(v[1])
		}
	})
	return
}

// UpdateRoleMatrixSchemePolicy
// @summary 更新角色关联预案资源规则
//          tips: 先删除旧规则,再添加新规则
func (s *sysRoleServiceImpl) UpdateRoleMatrixSchemePolicy(ctx context.Context, roleId int64, matrixSchemeIds []int64) (err error) {
	err = g.Try(func() {
		err = s.DelRoleMatrixSchemePolicy(ctx, 0, gconv.String(roleId))
		valueCheck.ErrIsNil(ctx, err)

		err = s.AddRoleMatrixSchemePolicy(ctx, roleId, matrixSchemeIds)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// DelRoleMatrixSchemePolicy
// @summary 删除角色关联预案资源规则
func (s *sysRoleServiceImpl) DelRoleMatrixSchemePolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)

		_, e = enforcer.RemoveFilteredNamedPolicy(global.CasbinPTypeMatrixScheme, fieldIndex, fieldValue)
		valueCheck.ErrIsNil(ctx, e)
	})
	return
}

// AddRoleMatrixSchemePolicy
// @summary 添加角色关联预案资源规则
func (s *sysRoleServiceImpl) AddRoleMatrixSchemePolicy(ctx context.Context, roleId int64, matrixSchemeIds []int64) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := gconv.Strings(matrixSchemeIds)
		for _, v := range casbinRules {
			_, err = enforcer.AddNamedPolicy(global.CasbinPTypeMatrixScheme, gconv.String(roleId), v, "All")
			valueCheck.ErrIsNil(ctx, err)
		}
	})
	return
}

// GetRoleMatrixSchemeGroupPolicy
// @summary 获取角色关联矩阵预案分组规则
func (s *sysRoleServiceImpl) GetRoleMatrixSchemeGroupPolicy(ctx context.Context, roleId int64) (matrixSchemeGroupIds []int64, err error) {
	err = g.Try(func() {
		//获取已分配预案分组资源权限
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := enforcer.GetFilteredNamedPolicy(global.CasbinPTypeMatrixSchemeGroup, 0, gconv.String(roleId))
		matrixSchemeGroupIds = make([]int64, len(casbinRules))
		for k, v := range casbinRules {
			matrixSchemeGroupIds[k] = gconv.Int64(v[1])
		}
	})
	return
}

// UpdateRoleMatrixSchemeGroupPolicy
// @summary 更新角色关联矩阵预案分组规则
//          tips: 先删除旧规则,再添加新规则
func (s *sysRoleServiceImpl) UpdateRoleMatrixSchemeGroupPolicy(ctx context.Context, roleId int64, matrixSchemeGroupIds []int64) (err error) {
	err = g.Try(func() {
		err = s.DelRoleMatrixSchemeGroupPolicy(ctx, 0, gconv.String(roleId))
		valueCheck.ErrIsNil(ctx, err)

		err = s.AddRoleMatrixSchemeGroupPolicy(ctx, roleId, matrixSchemeGroupIds)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// DelRoleMatrixSchemeGroupPolicy
// @summary 删除角色关联预案分组资源规则
func (s *sysRoleServiceImpl) DelRoleMatrixSchemeGroupPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)

		_, e = enforcer.RemoveFilteredNamedPolicy(global.CasbinPTypeMatrixSchemeGroup, fieldIndex, fieldValue)
		valueCheck.ErrIsNil(ctx, e)
	})
	return
}

// AddRoleMatrixSchemeGroupPolicy
// @summary 添加角色关联预案分组资源规则
func (s *sysRoleServiceImpl) AddRoleMatrixSchemeGroupPolicy(ctx context.Context, roleId int64, matrixSchemeGroupIds []int64) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := gconv.Strings(matrixSchemeGroupIds)
		for _, v := range casbinRules {
			_, err = enforcer.AddNamedPolicy(global.CasbinPTypeMatrixSchemeGroup, gconv.String(roleId), v, "All")
			valueCheck.ErrIsNil(ctx, err)
		}
	})
	return
}

// GetRoleMeetDispatchPolicy
// @summary 获取角色关联会议调度资源规则
func (s *sysRoleServiceImpl) GetRoleMeetDispatchPolicy(ctx context.Context, roleId int64) (meetDispatchIds []int64, err error) {
	err = g.Try(func() {
		//获取已分配会议调度资源权限
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := enforcer.GetFilteredNamedPolicy(global.CasbinPTypeMeetDispatch, 0, gconv.String(roleId))
		meetDispatchIds = make([]int64, len(casbinRules))
		for k, v := range casbinRules {
			meetDispatchIds[k] = gconv.Int64(v[1])
		}
	})
	return
}

// UpdateRoleMeetDispatchPolicy
// @summary 更新角色关联会议调度资源规则
//          tips: 先删除旧规则,再添加新规则
func (s *sysRoleServiceImpl) UpdateRoleMeetDispatchPolicy(ctx context.Context, roleId int64, meetDispatchIds []int64) (err error) {
	err = g.Try(func() {
		err = s.DelRoleMeetDispatchPolicy(ctx, 0, gconv.String(roleId))
		valueCheck.ErrIsNil(ctx, err)

		err = s.AddRoleMeetDispatchPolicy(ctx, roleId, meetDispatchIds)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// DelRoleMeetDispatchPolicy
// @summary 删除角色关联会议调度资源规则
func (s *sysRoleServiceImpl) DelRoleMeetDispatchPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)

		_, e = enforcer.RemoveFilteredNamedPolicy(global.CasbinPTypeMeetDispatch, fieldIndex, fieldValue)
		valueCheck.ErrIsNil(ctx, e)
	})
	return
}

// AddRoleMeetDispatchPolicy
// @summary 添加角色关联会议调度资源规则
func (s *sysRoleServiceImpl) AddRoleMeetDispatchPolicy(ctx context.Context, roleId int64, meetDispatchIds []int64) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := gconv.Strings(meetDispatchIds)
		for _, v := range casbinRules {
			_, err = enforcer.AddNamedPolicy(global.CasbinPTypeMeetDispatch, gconv.String(roleId), v, "All")
			valueCheck.ErrIsNil(ctx, err)
		}
	})
	return
}

// GetRoleMeetDispatchGroupPolicy
// @summary 获取角色关联会议分组规则
func (s *sysRoleServiceImpl) GetRoleMeetDispatchGroupPolicy(ctx context.Context, roleId int64) (meetDispatchGroupIds []int64, err error) {
	err = g.Try(func() {
		//获取已分配会议分组资源权限
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := enforcer.GetFilteredNamedPolicy(global.CasbinPTypeMeetDispatchGroup, 0, gconv.String(roleId))
		meetDispatchGroupIds = make([]int64, len(casbinRules))
		for k, v := range casbinRules {
			meetDispatchGroupIds[k] = gconv.Int64(v[1])
		}
	})
	return
}

// UpdateRoleMeetDispatchGroupPolicy
// @summary 更新角色关联会议分组规则
//          tips: 先删除旧规则,再添加新规则
func (s *sysRoleServiceImpl) UpdateRoleMeetDispatchGroupPolicy(ctx context.Context, roleId int64, meetDispatchGroupIds []int64) (err error) {
	err = g.Try(func() {
		err = s.DelRoleMeetDispatchGroupPolicy(ctx, 0, gconv.String(roleId))
		valueCheck.ErrIsNil(ctx, err)

		err = s.AddRoleMeetDispatchGroupPolicy(ctx, roleId, meetDispatchGroupIds)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// DelRoleMeetDispatchGroupPolicy
// @summary 删除角色关联会议分组资源规则
func (s *sysRoleServiceImpl) DelRoleMeetDispatchGroupPolicy(ctx context.Context, fieldIndex int, fieldValue string) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)

		_, e = enforcer.RemoveFilteredNamedPolicy(global.CasbinPTypeMeetDispatchGroup, fieldIndex, fieldValue)
		valueCheck.ErrIsNil(ctx, e)
	})
	return
}

// AddRoleMeetDispatchGroupPolicy
// @summary 添加角色关联会议分组资源规则
func (s *sysRoleServiceImpl) AddRoleMeetDispatchGroupPolicy(ctx context.Context, roleId int64, meetDispatchGroupIds []int64) (err error) {
	err = g.Try(func() {
		enforcer, e := comSrv.CasbinEnforcer(ctx)
		valueCheck.ErrIsNil(ctx, e)
		casbinRules := gconv.Strings(meetDispatchGroupIds)
		for _, v := range casbinRules {
			_, err = enforcer.AddNamedPolicy(global.CasbinPTypeMeetDispatchGroup, gconv.String(roleId), v, "All")
			valueCheck.ErrIsNil(ctx, err)
		}
	})
	return
}

// NotifyDataPermissionUpdate
// @summary 通知在线用户数据权限更新
func (s *sysRoleServiceImpl) NotifyDataPermissionUpdate(ctx context.Context, roleId int64, permissionType define.EPermissionType) {
	onlineUserList := OnlineUser.GetAllOnlineUser()
	for _, user := range onlineUserList {
		if User.CheckUserRole(ctx, user.UserId, roleId) {
			publish.WsPublish.PublishMsgToSpecifiedClient(user.Token, defCom.ENotifyUserPermissionUpdate.String(), define.SysRoleDataPermissionNotify{
				Type:           defCom.EOptModify.String(),
				PermissionType: permissionType,
			})
		}
	}
}
