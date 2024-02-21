package define

import (
	"gowebserver/app/common/define"
)

//SysRoleSelectPageReq 分页请求参数
type SysRoleSelectPageReq struct {
	KeyWord string `p:"keyWord"` //查询关键字
	define.SelectPageReq
}

//SysRoleAddReq 新增角色请求
type SysRoleAddReq struct {
	RoleKey       string `json:"roleKey"`       // 角色Key
	RoleName      string `json:"roleName"`      // 角色名称
	RoleSort      int64  `json:"roleSort"`      // 角色排序
	RoleStatus    int64  `json:"roleStatus"`    // 角色状态：0停用 1正常
	RoleDataScope int    `json:"roleDataScope"` // 数据范围标识(1：全部数据权限 2：自定数据权限)
	Remark        string `json:"remark"`        // 备注信息
}

//SysRoleInfo 请求返回角色基本信息
type SysRoleInfo struct {
	RoleID        int64  `json:"roleId"`        // 角色id
	RoleKey       string `json:"roleKey"`       // 角色标签
	RoleName      string `json:"roleName"`      // 角色名称
	RoleSort      int64  `json:"roleSort"`      // 显示顺序
	RoleStatus    int64  `json:"roleStatus"`    // 角色状态(0停用 1正常)
	RoleDataScope int    `json:"roleDataScope"` // 数据范围标识(1：全部数据权限 2：自定数据权限)
	Remark        string `json:"remark"`        // 备注信息
}

//SysRoleEditReq 编辑角色基本信息请求
type SysRoleEditReq struct {
	SysRoleInfo
}

//SysRoleDeleteReq 删除角色信息请求
type SysRoleDeleteReq struct {
	RoleIds string `p:"roleIds"`
}

//SysRoleSearchList 查询列表返回值
type SysRoleSearchList struct {
	List  []SysRoleInfo `json:"list"`
	Page  int           `json:"page"`
	Size  int           `json:"size"`
	Total int           `json:"total"`
}

//SysRoleMenuList 角色菜单资源
type SysRoleMenuList struct {
	RoleID  int64   `json:"roleId"`  // 角色ID
	MenuIDS []int64 `json:"menuIds"` // 菜单资源ID列表
}

//SysRoleChnGroupList 角色通道分组资源
type SysRoleChnGroupList struct {
	RoleID      int64    `json:"roleId"`      // 角色ID
	ChnGroupIDS []string `json:"chnGroupIds"` // 通道分组ID列表
}

//SysRoleTvWallList 角色大屏资源列表
type SysRoleTvWallList struct {
	RoleID    int64    `json:"roleId"`    // 角色ID
	TvWallIDS []string `json:"tvWallIds"` // 大屏Id列表
}

//SysRoleMatrixSchemeList 角色矩阵预案资源列表
type SysRoleMatrixSchemeList struct {
	RoleID          int64   `json:"roleId"`
	MatrixSchemeIDS []int64 `json:"matrixSchemeIds"` // 矩阵预案Id列表
}

//SysRoleMatrixSchemeGroupList 角色矩阵预案分组资源列表
type SysRoleMatrixSchemeGroupList struct {
	RoleID               int64   `json:"roleId"`
	MatrixSchemeGroupIdS []int64 `json:"matrixSchemeGroupIds"` // 矩阵预案分组Id列表
}

//SysRoleMeetDispatchList 角色会议调度资源列表
type SysRoleMeetDispatchList struct {
	RoleID          int64   `json:"roleId"`
	MeetDispatchIds []int64 `json:"meetDispatchIds"` // 会议调度Id列表
}

//SysRoleMeetDispatchGroupList 角色会议调度分组资源列表
type SysRoleMeetDispatchGroupList struct {
	RoleID               int64   `json:"roleId"`
	MeetDispatchGroupIds []int64 `json:"meetDispatchGroupIds"` // 会议调度分组Id列表
}

//EPermissionType 数据权限类型
type EPermissionType int

const (
	EPermissionTypeMenu              EPermissionType = 1
	EPermissionTypeChnGroup          EPermissionType = 2
	EPermissionTypeTvWall            EPermissionType = 3
	EPermissionTypeMatrixScheme      EPermissionType = 4
	EPermissionTypeMatrixSchemeGroup EPermissionType = 5
	EPermissionTypeMeetDispatch      EPermissionType = 6
	EPermissionTypeMeetDispatchGroup EPermissionType = 7
)

// SysRoleDataPermissionNotify 角色数据权限更新通知
type SysRoleDataPermissionNotify struct {
	Type           string          `json:"type"` //EOptType
	PermissionType EPermissionType `json:"permissionType"`
}
