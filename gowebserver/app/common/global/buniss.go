// =================================================================================
// 业务类型枚举
// =================================================================================

package global

//ResponseBusinessType 请求上下文中记录业务操作类型的标签
const ResponseBusinessType = "ResponseBusinessType"

type BusinessType int

const (
	BusinessOther       BusinessType = 0 //0其他
	BusinessAdd         BusinessType = 1 //1新增
	BusinessEdit        BusinessType = 2 //2修改
	BusinessDel         BusinessType = 3 //3删除
	BusinessGet         BusinessType = 4 //4查询
	BusinessAuthorize   BusinessType = 5 //5授权
	BusinessExport      BusinessType = 6 //6导出
	BusinessImport      BusinessType = 7 //7导入
	BusinessForceLogout BusinessType = 8 //8强退
	BusinessClean       BusinessType = 9 //9清空
)
