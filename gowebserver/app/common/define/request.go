package define

// SelectPageReq 分页查询请求
type SelectPageReq struct {
	BeginTime     string `p:"beginTime"`     //时间范围-起始		非必须(查询条件)
	EndTime       string `p:"endTime"`       //时间范围-结束		非必须(查询条件)
	PageNum       int    `p:"pageNum"`       //当前页码			必须
	PageSize      int    `p:"pageSize"`      //每页数			必须
	OrderByColumn string `p:"orderByColumn"` //排序字段			非必须
	IsAsc         string `p:"isAsc"`         //排序方式			非必须,OrderByColumn不为空时有效,"ASC"升序 "DESC"降序
}
