package page

import (
	"math"
)

type Paging struct {
	PageNum   int //当前页
	PageSize  int //每页条数
	Total     int //总条数
	PageCount int //总页数
	StartNum  int //起始行
}

// CreatePaging 创建分页
func CreatePaging(pageNum, pageSize, total int) *Paging {
	if pageNum < 1 {
		pageNum = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}

	pageCount := math.Ceil(float64(total) / float64(pageSize))
	startNum := pageSize * (pageNum - 1)
	paging := new(Paging)
	paging.PageNum = pageNum
	paging.PageSize = pageSize
	paging.Total = total
	paging.PageCount = int(pageCount)
	paging.StartNum = startNum
	return paging
}
