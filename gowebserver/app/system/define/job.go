package define

import (
	"gowebserver/app/common/define"
	"gowebserver/app/system/model"
)

// JobApiSelectPageReq 分页查询定时任务请求参数
type JobApiSelectPageReq struct {
	JobName        string `p:"jobName"`        //任务名称(模糊匹配)
	JobGroup       string `p:"jobGroup"`       //任务分组
	InvokeTarget   string `p:"invokeTarget"`   //调用目标字符串(模糊匹配)
	CronExpression string `p:"cronExpression"` //cron执行表达式(模糊匹配)
	MisfirePolicy  string `p:"misfirePolicy"`  //计划执行错误策略("1"立即执行 "2"执行一次 "3"放弃执行)
	Concurrent     string `p:"concurrent"`     //是否并发执行("0"允许 "1"禁止)
	Status         string `p:"status"`         //状态("0"正常 "1"暂停)
	define.SelectPageReq
}

// JobApiCreateReq 新增定时任务请求参数
type JobApiCreateReq struct {
	JobName        string `p:"jobName" v:"required#任务名称不能为空"`
	JobParams      string `p:"jobParams"`
	JobGroup       string `p:"jobGroup" `
	InvokeTarget   string `p:"invokeTarget" v:"required#执行方法不能为空"`
	CronExpression string `p:"cronExpression" v:"required#任务表达式不能为空"`
	MisfirePolicy  string `p:"misfirePolicy" v:"required#执行策略"` //计划执行错误策略("1"立即执行 "2"执行一次 "3"放弃执行)
	Concurrent     string `p:"concurrent" v:"required#是否并发执行"`  //否并发执行("0"允许 "1"禁止)
	Status         string `p:"status" v:"required#状态(0正常 1暂停)不能为空"`
	Remark         string `p:"remark" `
}

// JobApiUpdateReq 修改定时任务请求参数
type JobApiUpdateReq struct {
	JobId int64 `p:"jobId" v:"required#主键ID不能为空"`
	JobApiCreateReq
}

// JobApiDeleteReq 删除请求
type JobApiDeleteReq struct {
	Ids string `p:"ids"  v:"required#请选择要删除的数据记录"`
}

// JobServiceList 查询定时任务列表返回数据结构
type JobServiceList struct {
	List  []model.SysJob `json:"list"`
	Page  int            `json:"page"`
	Size  int            `json:"size"`
	Total int            `json:"total"`
}
