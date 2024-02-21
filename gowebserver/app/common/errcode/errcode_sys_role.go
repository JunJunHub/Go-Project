// ===================================
// 系统角色相关错误码定义[112000,112999]
// ===================================

package errcode

//[112000,112099] 系统角色管理预留错误码
const (
//codeSysRole = iota + codeSysRoleStart //112000
)

//[112100,112199] 系统角色管理相关错误代码
const (
	codeSysRoleInvalidRoleId = iota + codeSysRoleManageStart //112100
	codeSysRoleNameDuplicate                                 //112101
	codeSysRoleNotSupportedCustomPermissions
	codeSysRoleDelDefaultRoleForbid
	codeSysRoleDelInUseRoleForbid
	codeSysRoleModifyDefaultRoleForbid
	codeSysRoleBindMenuDuplicate
	codeSysRoleBindMenuInvalid
	codeSysRoleBindChnGroupDuplicate
	codeSysRoleBindChnGroupInvalid
	codeSysRoleBindTvWallDuplicate
	codeSysRoleBindTvWallInvalid
	codeSysRoleBindMatrixSchemeDuplicate
	codeSysRoleBindMatrixSchemeInvalid
	codeSysRoleBindMatrixSchemeGroupDuplicate
	codeSysRoleBindMatrixSchemeGroupInvalid

	codeSysRoleBindMeetDispatchDuplicate
	codeSysRoleBindMeetDispatchInvalid
	codeSysRoleBindMeetDispatchGroupDuplicate
	codeSysRoleBindMeetDispatchGroupInvalid
)

var (
	ErrSysRoleInvalidRoleId                  = NewMpuapsErr(codeSysRoleInvalidRoleId, "{#ErrSysRoleInvalidRoleId}", nil)
	ErrSysRoleNameDuplicate                  = NewMpuapsErr(codeSysRoleNameDuplicate, "{#ErrSysRoleNameDuplicate}", nil)
	ErrSysRoleNotSupportedCustomPermissions  = NewMpuapsErr(codeSysRoleNotSupportedCustomPermissions, "{#ErrSysRoleNotSupportedCustomPermissions}", nil)
	ErrSysRoleDelDefaultRoleForbid           = NewMpuapsErr(codeSysRoleDelDefaultRoleForbid, "{#ErrSysRoleDelDefaultRoleForbid}", nil)
	ErrSysRoleDelInUseRoleForbid             = NewMpuapsErr(codeSysRoleDelInUseRoleForbid, "{#ErrSysRoleDelInUseRoleForbid}", nil)
	ErrSysRoleModifyDefaultRoleForbid        = NewMpuapsErr(codeSysRoleModifyDefaultRoleForbid, "{#ErrSysRoleModifyDefaultRoleForbid}", nil)
	ErrSysRoleBindMenuDuplicate              = NewMpuapsErr(codeSysRoleBindMenuDuplicate, "{#ErrSysRoleBindMenuDuplicate}", nil)
	ErrSysRoleBindMenuInvalid                = NewMpuapsErr(codeSysRoleBindMenuInvalid, "{#ErrSysRoleBindMenuInvalid}", nil)
	ErrSysRoleBindChnGroupDuplicate          = NewMpuapsErr(codeSysRoleBindChnGroupDuplicate, "{#ErrSysRoleBindChnGroupDuplicate}", nil)
	ErrSysRoleBindChnGroupInvalid            = NewMpuapsErr(codeSysRoleBindChnGroupInvalid, "{#ErrSysRoleBindChnGroupInvalid}", nil)
	ErrSysRoleBindTvWallDuplicate            = NewMpuapsErr(codeSysRoleBindTvWallDuplicate, "{#ErrSysRoleBindTvWallDuplicate}", nil)
	ErrSysRoleBindTvWallInvalid              = NewMpuapsErr(codeSysRoleBindTvWallInvalid, "{#ErrSysRoleBindTvWallInvalid}", nil)
	ErrSysRoleBindMatrixSchemeDuplicate      = NewMpuapsErr(codeSysRoleBindMatrixSchemeDuplicate, "{#ErrSysRoleBindMatrixSchemeDuplicate}", nil)
	ErrSysRoleBindMatrixSchemeInvalid        = NewMpuapsErr(codeSysRoleBindMatrixSchemeInvalid, "{#ErrSysRoleBindMatrixSchemeInvalid}", nil)
	ErrSysRoleBindMatrixSchemeGroupDuplicate = NewMpuapsErr(codeSysRoleBindMatrixSchemeGroupDuplicate, "{#ErrSysRoleBindMatrixSchemeGroupDuplicate}", nil)
	ErrSysRoleBindMatrixSchemeGroupInvalid   = NewMpuapsErr(codeSysRoleBindMatrixSchemeGroupInvalid, "{#ErrSysRoleBindMatrixSchemeGroupInvalid}", nil)

	ErrSysRoleBindMeetDispatchDuplicate      = NewMpuapsErr(codeSysRoleBindMeetDispatchDuplicate, "{#ErrSysRoleBindMeetDispatchDuplicate}", nil)
	ErrSysRoleBindMeetDispatchInvalid        = NewMpuapsErr(codeSysRoleBindMeetDispatchInvalid, "{#ErrSysRoleBindMeetDispatchInvalid}", nil)
	ErrSysRoleBindMeetDispatchGroupDuplicate = NewMpuapsErr(codeSysRoleBindMeetDispatchGroupDuplicate, "{#ErrSysRoleBindMeetDispatchGroupDuplicate}", nil)
	ErrSysRoleBindMeetDispatchGroupInvalid   = NewMpuapsErr(codeSysRoleBindMeetDispatchGroupInvalid, "{#ErrSysRoleBindMeetDispatchGroupInvalid}", nil)
)
