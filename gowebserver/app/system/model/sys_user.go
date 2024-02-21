package model

//SysUserExtend 用户信息扩展(用户角色)
type SysUserExtend struct {
	SysUser
	Roles []SysRole `json:"roles"`
}
