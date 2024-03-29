// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// SysOperLogDao is the manager for logic model data accessing and custom defined data operations functions management.
type SysOperLogDao struct {
	Table   string            // Table is the underlying table name of the DAO.
	Group   string            // Group is the database configuration group name of current DAO.
	Columns SysOperLogColumns // Columns is the short type for Columns, which contains all the column names of Table for convenient usage.
}

// SysOperLogColumns defines and stores column names for table sys_oper_log.
type SysOperLogColumns struct {
	OperId        string // 日志主键
	Title         string // 模块标题
	BusinessType  string // 业务类型（0其它 1新增 2修改 3删除 4查询 5授权 6导出 7导入 8强退 9清空）
	Method        string // 方法名称
	RequestMethod string // 请求方式
	OperatorType  string // 操作类别（0其它 1前台用户 2级联操作）
	OperName      string // 操作人员
	DeptName      string // 部门名称
	OperUrl       string // 请求url
	OperIp        string // 主机地址
	OperLocation  string // 操作地点
	OperParam     string // 请求参数
	JsonResult    string // 返回参数
	Status        string // 操作结果（0正常 -1异常）
	ErrorMsg      string // 错误消息
	OperTime      string // 操作时间
}

//  sysOperLogColumns holds the columns for table sys_oper_log.
var sysOperLogColumns = SysOperLogColumns{
	OperId:        "oper_id",
	Title:         "title",
	BusinessType:  "business_type",
	Method:        "method",
	RequestMethod: "request_method",
	OperatorType:  "operator_type",
	OperName:      "oper_name",
	DeptName:      "dept_name",
	OperUrl:       "oper_url",
	OperIp:        "oper_ip",
	OperLocation:  "oper_location",
	OperParam:     "oper_param",
	JsonResult:    "json_result",
	Status:        "status",
	ErrorMsg:      "error_msg",
	OperTime:      "oper_time",
}

// NewSysOperLogDao creates and returns a new DAO object for table data access.
func NewSysOperLogDao() *SysOperLogDao {
	return &SysOperLogDao{
		Group:   "default",
		Table:   "sys_oper_log",
		Columns: sysOperLogColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SysOperLogDao) DB() gdb.DB {
	return g.DB(dao.Group)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SysOperLogDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.Table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SysOperLogDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
