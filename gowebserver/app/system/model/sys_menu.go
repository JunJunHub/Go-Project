// ==========================================================================
// 菜单相关数据结构定义
// ==========================================================================

package model

//SysMenuTree 菜单资源树
type SysMenuTree struct {
	*SysMenu
	ChildTree []*SysMenuTree `json:"childTree"`
}
