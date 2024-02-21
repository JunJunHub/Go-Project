package service

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/utils/valueCheck"
	"gowebserver/app/system/dao"
	"gowebserver/app/system/model"
	"strings"
)

type ISysMenu interface {
	GetAllMenuTree(ctx context.Context) (menuTree []*model.SysMenuTree, err error)
	GetAllMenuList(ctx context.Context) (menuList []*model.SysMenu, err error)
	GetAllMenuIds(ctx context.Context) (menuIds []int64, err error)
	GetMenuListByMenuIds(ctx context.Context, menuIds []int64) (menuList []*model.SysMenu, err error)
	GetRoleMenus(ctx context.Context, roleId int64) (menuList []*model.SysMenu, err error)
	GetUserMenus(ctx context.Context) (menuList []*model.SysMenu, err error)

	CheckNeedVerifyPermissions(ctx context.Context, restInterface string) bool
	CheckHttpInterfacePermissions(ctx context.Context, restInterface string) bool
}

type sysMenuServiceImpl struct{}

var menuService = sysMenuServiceImpl{}

func InterfaceSysMenu() ISysMenu {
	return &menuService
}

func (s *sysMenuServiceImpl) SysMenuInit() {
}

//BuildMenuTree
//@summary 构造菜单树
func (s *sysMenuServiceImpl) BuildMenuTree(menuPid int64, menuList []*model.SysMenu) []*model.SysMenuTree {
	tree := make([]*model.SysMenuTree, 0, len(menuList))
	for _, value := range menuList {
		if value.Pid == menuPid {
			t := &model.SysMenuTree{
				SysMenu: value,
			}
			child := s.BuildMenuTree(value.Id, menuList)
			if len(child) > 0 {
				t.ChildTree = child
			}
			tree = append(tree, t)
		}
	}
	return tree
}

//GetAllMenuTree
//@summary 获取所有功能菜单树
func (s *sysMenuServiceImpl) GetAllMenuTree(ctx context.Context) (menuTree []*model.SysMenuTree, err error) {
	err = g.Try(func() {
		var menuList []*model.SysMenu

		err = dao.SysMenu.Ctx(ctx).FindScan(&menuList)
		valueCheck.ErrIsNil(ctx, err)

		//构造菜单树
		menuTree = s.BuildMenuTree(0, menuList)
	})
	return
}

//GetAllMenuList
//@summary 获取所有菜单列表
func (s *sysMenuServiceImpl) GetAllMenuList(ctx context.Context) (menuList []*model.SysMenu, err error) {
	err = g.Try(func() {
		err = dao.SysMenu.Ctx(ctx).FindScan(&menuList)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// GetAllMenuIds
// @summary 获取所有菜单ID
func (s *sysMenuServiceImpl) GetAllMenuIds(ctx context.Context) (menuIds []int64, err error) {
	err = g.Try(func() {
		var result []gdb.Value
		result, err = dao.SysMenu.Ctx(ctx).Fields(dao.SysMenu.Columns.Id).FindArray()
		valueCheck.ErrIsNil(ctx, err)
		for _, value := range result {
			menuIds = append(menuIds, value.Int64())
		}
	})
	return
}

//GetMenuListByMenuIds
//@summary 批量获取菜单列表
func (s *sysMenuServiceImpl) GetMenuListByMenuIds(ctx context.Context, menuIds []int64) (menuList []*model.SysMenu, err error) {
	err = g.Try(func() {
		err = dao.SysMenu.Ctx(ctx).WhereIn(dao.SysMenu.Columns.Id, menuIds).FindScan(&menuList)
		valueCheck.ErrIsNil(ctx, err)
	})
	return
}

// GetRoleMenus
// @summary 角色绑定菜单资源
func (s *sysMenuServiceImpl) GetRoleMenus(ctx context.Context, roleId int64) (menuList []*model.SysMenu, err error) {
	err = g.Try(func() {
		roleInfo, _ := InterfaceSysRole().Get(ctx, roleId)
		valueCheck.ValueIsNil(roleInfo, errcode.ErrSysRoleInvalidRoleId.Message())

		menuIds, e := InterfaceSysRole().GetRoleMenuPolicy(ctx, roleId)
		valueCheck.ErrIsNil(ctx, e)

		menuList, err = InterfaceSysMenu().GetMenuListByMenuIds(ctx, menuIds)
	})
	return
}

// GetUserMenus
// @summary 登录用户获取菜单权限
func (s *sysMenuServiceImpl) GetUserMenus(ctx context.Context) (menuList []*model.SysMenu, err error) {
	userInfo := Context.GetUser(ctx)
	roles, _ := User.GetUserRoles(ctx, userInfo.UserId)
	for _, role := range roles {
		var tmpMenuList []*model.SysMenu
		if tmpMenuList, err = s.GetRoleMenus(ctx, role.RoleId); err != nil {
			return nil, err
		}
		menuList = append(menuList, tmpMenuList...)
	}
	return
}

//GetAllNeedVerifyPermissionsInterface
//@summary 获取所有需要校验权限的接口
func (s *sysMenuServiceImpl) GetAllNeedVerifyPermissionsInterface(ctx context.Context) (restInterfaceList []string) {
	allMenuList, err := s.GetAllMenuList(ctx)
	if err != nil {
		g.Log().Error(err)
	}

	for _, menu := range allMenuList {
		tmpRestInterfaceList := strings.Split(menu.RestInterface, "|")
		restInterfaceList = append(restInterfaceList, tmpRestInterfaceList...)
	}
	return
}

//CheckNeedVerifyPermissions
//@summary 校验http接口是否需要校验权限
func (s *sysMenuServiceImpl) CheckNeedVerifyPermissions(ctx context.Context, restInterface string) bool {
	allNeedVerifyInterface := s.GetAllNeedVerifyPermissionsInterface(ctx)
	for idx := range allNeedVerifyInterface {
		if allNeedVerifyInterface[idx] == restInterface {
			return true
		}
	}
	return false
}

func (s *sysMenuServiceImpl) CheckHttpInterfacePermissions(ctx context.Context, restInterface string) bool {
	if s.CheckNeedVerifyPermissions(ctx, restInterface) {
		var bHasPermissions bool
		var userInterfaceList []string

		menuList, _ := s.GetUserMenus(ctx)
		for _, menu := range menuList {
			userInterfaceList = append(userInterfaceList, strings.Split(menu.RestInterface, "|")...)
		}
		for _, interfaceUrl := range userInterfaceList {
			if interfaceUrl == restInterface {
				bHasPermissions = true
				break
			}
		}
		if !bHasPermissions {
			g.Log().Noticef("restInterface no permissions: %s", restInterface)
			return false
		}
	}
	return true
}
