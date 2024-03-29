// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"gowebserver/app/system/dao/internal"
)

// sysBlackListDao is the manager for logic model data accessing and custom defined data operations functions management.
// You can define custom methods on it to extend its functionality as you wish.
type sysBlackListDao struct {
	*internal.SysBlackListDao
}

var (
	// SysBlackList is globally public accessible object for table sys_black_list operations.
	SysBlackList = sysBlackListDao{
		internal.NewSysBlackListDao(),
	}
)

// Fill with you ideas below.
