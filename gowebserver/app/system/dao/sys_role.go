// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"gowebserver/app/system/dao/internal"
)

// sysRoleDao is the manager for logic model data accessing and custom defined data operations functions management.
// You can define custom methods on it to extend its functionality as you wish.
type sysRoleDao struct {
	*internal.SysRoleDao
}

var (
	// SysRole is globally public accessible object for table sys_role operations.
	SysRole = sysRoleDao{
		internal.NewSysRoleDao(),
	}
)

// Fill with you ideas below.
