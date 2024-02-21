// =================================================================================
// This is auto-generated by GoFrame CLI tool only once. Fill this file as you wish.
// =================================================================================

package dao

import (
	"gowebserver/app/system/dao/internal"
)

// sysJobLogDao is the manager for logic model data accessing and custom defined data operations functions management.
// You can define custom methods on it to extend its functionality as you wish.
type sysJobLogDao struct {
	*internal.SysJobLogDao
}

var (
	// SysJobLog is globally public accessible object for table sys_job_log operations.
	SysJobLog = sysJobLogDao{
		internal.NewSysJobLogDao(),
	}
)

// Fill with you ideas below.