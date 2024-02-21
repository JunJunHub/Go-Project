// ========================================================
// 测试casbin单例模式管理多类数据权限(多类数据权限策略存储在同一张表)
// ========================================================

package service

import (
	"context"
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gogf/gf/frame/g"
	"testing"
)

//数据权限类型
const EnforceMenu string = ""
const EnforceAVChannelGroup string = "2"
const EnforceTvwall string = "3"
const EnforceMatrixScheme string = "4"

func TestEnforcerCheck(t *testing.T) {
	g.Cfg().SetFileName("config-local.toml")

	enforcer, err := CasbinEnforcer(context.TODO())
	if err != nil {
		t.Log(err)
		return
	}

	//测试录入用户角色绑定关系(userId,roleId)
	//user
	//  userId 99
	//role
	//  roleSysAdmin 1
	//接口描述不会重复添加,实测会重复???
	_, err = enforcer.AddGroupingPolicy(
		fmt.Sprintf("u_%d", 99),
		fmt.Sprintf("r_%d", 1))
	if err != nil {
		t.Log(err)
	}

	//测试录入角色菜单权限(权限标识 All)
	_, err = enforcer.AddNamedPolicy(
		fmt.Sprintf("p%s", EnforceMenu),
		fmt.Sprintf("r_%d", 1),
		fmt.Sprintf("Menu_%d", 10086),
		"All")
	if err != nil {
		t.Log(err)
		return
	}

	//测试录入角色音视频通道数据权限(权限标识 All)
	_, err = enforcer.AddNamedPolicy(
		fmt.Sprintf("p%s", EnforceAVChannelGroup),
		fmt.Sprintf("r_%d", 1),
		fmt.Sprintf("ChnGroup_%d", 10086),
		"All")
	if err != nil {
		t.Log(err)
		return
	}

	//测试录入角色大屏数据权限(权限标识 All)
	_, err = enforcer.AddNamedPolicy(
		fmt.Sprintf("p%s", EnforceTvwall),
		fmt.Sprintf("r_%d", 1),
		fmt.Sprintf("TvWall_%d", 10086),
		"All")
	if err != nil {
		t.Log(err)
		return
	}

	//测试录入角色矩阵调度预案权限(权限标识 All)
	_, err = enforcer.AddNamedPolicy(
		fmt.Sprintf("p%s", EnforceMatrixScheme),
		fmt.Sprintf("r_%d", 1),
		fmt.Sprintf("MatrixScheme_%d", 10086),
		"All")
	if err != nil {
		t.Log(err)
		return
	}

	//校验菜单数据权限(sub=userId obj=menuId act={})
	hasAccess := false
	enforceCtx := casbin.NewEnforceContext(EnforceMenu)
	hasAccess, err = enforcer.Enforce(enforceCtx, "u_99", "Menu_10086", "All")
	if err != nil {
		t.Log(err)
	}
	if !hasAccess {
		t.Log("没有菜单权限")
	} else {
		t.Log("拥有菜单权限")
	}

	//校验音视频通道数据权限(sub=userId obj=channelId act={})
	enforceCtx = casbin.NewEnforceContext(EnforceAVChannelGroup)
	hasAccess, err = enforcer.Enforce(enforceCtx, "u_99", "ChnGroup_10086", "All")
	if err != nil {
		t.Log(err)
	}
	if !hasAccess {
		t.Log("没有通道资源组权限")
	} else {
		t.Log("拥有通道资源组权限")
	}

	//使用音视频数据权限匹配规则和策略判断菜单数据权限应该是没权限
	hasAccess, err = enforcer.Enforce(enforceCtx, "u_99", "Menu_10086", "All")
	if err != nil {
		t.Log(err)
	}
	if !hasAccess {
		t.Log("没有菜单权限")
	} else {
		t.Log("拥有菜单权限")
	}
}
