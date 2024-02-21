// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// SysJobDao is the manager for logic model data accessing and custom defined data operations functions management.
type SysJobDao struct {
	Table   string        // Table is the underlying table name of the DAO.
	Group   string        // Group is the database configuration group name of current DAO.
	Columns SysJobColumns // Columns is the short type for Columns, which contains all the column names of Table for convenient usage.
}

// SysJobColumns defines and stores column names for table sys_job.
type SysJobColumns struct {
	JobId          string // 任务id
	JobName        string // 任务名称
	JobParams      string // 任务参数
	JobGroup       string // 任务组名
	InvokeTarget   string // 调用目标字符串
	CronExpression string // cron执行表达式
	MisfirePolicy  string // 计划执行策略（1多次执行 2执行一次）
	Concurrent     string // 是否并发执行（0允许 1禁止）
	Status         string // 状态（0正常 1暂停）
	Remark         string // 备注信息
	CreateBy       string // 创建者
	CreateTime     string // 创建时间
	UpdateBy       string // 更新者
	UpdateTime     string // 更新时间
}

//  sysJobColumns holds the columns for table sys_job.
var sysJobColumns = SysJobColumns{
	JobId:          "job_id",
	JobName:        "job_name",
	JobParams:      "job_params",
	JobGroup:       "job_group",
	InvokeTarget:   "invoke_target",
	CronExpression: "cron_expression",
	MisfirePolicy:  "misfire_policy",
	Concurrent:     "concurrent",
	Status:         "status",
	Remark:         "remark",
	CreateBy:       "create_by",
	CreateTime:     "create_time",
	UpdateBy:       "update_by",
	UpdateTime:     "update_time",
}

// NewSysJobDao creates and returns a new DAO object for table data access.
func NewSysJobDao() *SysJobDao {
	return &SysJobDao{
		Group:   "default",
		Table:   "sys_job",
		Columns: sysJobColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SysJobDao) DB() gdb.DB {
	return g.DB(dao.Group)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SysJobDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.Table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SysJobDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
