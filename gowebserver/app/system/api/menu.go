package api

import (
	"github.com/gogf/gf/net/ghttp"
	"gowebserver/app/common/api"
	"gowebserver/app/common/global"
	"gowebserver/app/system/service"
)

var MenuApi = new(menuApi)

type menuApi struct {
	api.BaseController
}

func (a *menuApi) Init(r *ghttp.Request) {
	a.Module = "菜单管理"
	r.SetCtxVar(global.Module, a.Module)
}

func (a *menuApi) GetTree(r *ghttp.Request) {
	tree, err := service.InterfaceSysMenu().GetAllMenuTree(r.GetCtx())
	if err != nil {
		a.RespJsonExit(r, err)
	}
	a.RespJsonExit(r, nil, tree)
}

func (a *menuApi) GetMenus(r *ghttp.Request) {
	roleId := r.GetInt64("roleId")
	menuList, _ := service.InterfaceSysMenu().GetRoleMenus(r.GetCtx(), roleId)
	a.RespJsonExit(r, nil, menuList)
}
