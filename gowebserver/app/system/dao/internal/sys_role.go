// ==========================================================================
// Code generated by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"
	"github.com/gogf/gf/database/gdb"
	"github.com/gogf/gf/frame/g"
)

// SysRoleDao is the manager for logic model data accessing and custom defined data operations functions management.
type SysRoleDao struct {
	Table   string         // Table is the underlying table name of the DAO.
	Group   string         // Group is the database configuration group name of current DAO.
	Columns SysRoleColumns // Columns is the short type for Columns, which contains all the column names of Table for convenient usage.
}

// SysRoleColumns defines and stores column names for table sys_role.
type SysRoleColumns struct {
	RoleId        string // 角色id
	RoleName      string // 角色名称
	RoleKey       string // 角色标签
	RoleSort      string // 显示顺序
	RoleDataScope string // 数据范围标识(1：全部数据权限 2：自定数据权限)
	RoleStatus    string // 角色状态(0停用 1正常)
	Remark        string // 备注信息
	CreateBy      string // 创建者
	CreateTime    string // 创建时间
	UpdateBy      string // 更新者
	UpdateTime    string // 更新时间
}

//  sysRoleColumns holds the columns for table sys_role.
var sysRoleColumns = SysRoleColumns{
	RoleId:        "role_id",
	RoleName:      "role_name",
	RoleKey:       "role_key",
	RoleSort:      "role_sort",
	RoleDataScope: "role_data_scope",
	RoleStatus:    "role_status",
	Remark:        "remark",
	CreateBy:      "create_by",
	CreateTime:    "create_time",
	UpdateBy:      "update_by",
	UpdateTime:    "update_time",
}

// NewSysRoleDao creates and returns a new DAO object for table data access.
func NewSysRoleDao() *SysRoleDao {
	return &SysRoleDao{
		Group:   "default",
		Table:   "sys_role",
		Columns: sysRoleColumns,
	}
}

// DB retrieves and returns the underlying raw database management object of current DAO.
func (dao *SysRoleDao) DB() gdb.DB {
	return g.DB(dao.Group)
}

// Ctx creates and returns the Model for current DAO, It automatically sets the context for current operation.
func (dao *SysRoleDao) Ctx(ctx context.Context) *gdb.Model {
	return dao.DB().Model(dao.Table).Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rollbacks the transaction and returns the error from function f if it returns non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function f
// as it is automatically handled by this function.
func (dao *SysRoleDao) Transaction(ctx context.Context, f func(ctx context.Context, tx *gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
