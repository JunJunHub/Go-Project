package api

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"gowebserver/app/common/api"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/system/define"
	"gowebserver/app/system/service"
)

var JobApi = new(jobApi)

type jobApi struct {
	api.BaseController
}

func (a *jobApi) Init(r *ghttp.Request) {
	a.Module = "定时任务管理"
	r.SetCtxVar(global.Module, a.Module)
}

// Get 分页获取定时任务
// @summary 分页获取定时任务配置
// @tags 	定时任务管理
// @Param   authorization header string true "Bearer Token"
// @Param   jobName         query string false "任务名称"
// @Param   jobGroup        query string false "任务分组"
// @Param   invokeTarget    query string false "调用方法"
// @Param   cronExpression  query string false "cron表达式"
// @Param   misfirePolicy   query string false "执行策略"
// @Param   concurrent      query string false "是否并发执行"
// @Param   status          query string false "任务状态: "0"启动 "1"暂停"
// @Param   pageNum         query int    true  "当前页码"
// @Param   pageSize        query int    true  "每页展示条数"
// @Param   orderByColumn   query string false "排序字段"
// @Param   isAsc           query string false "排序方式: "ASC"升 "DESC"降, 当orderByColumn不为空时有效"
// @Produce json
// @Success 200 {object} response.Response{data=define.JobServiceList} "用户信息列表"
// @Failure 500 {object} response.Response "请求参数错误"
// @Router 	/monitor/job [GET]
func (a *jobApi) Get(r *ghttp.Request) {
	var req *define.JobApiSelectPageReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	result := service.Job.GetList(req)
	a.RespJsonExit(r, nil, result)
}

// Post 新增定时任务
// @summary 新增定时任务
// @tags 	定时任务管理
// @Router 	/monitor/job [POST]
func (a *jobApi) Post(r *ghttp.Request) {
	var req *define.JobApiCreateReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	id, err := service.Job.Create(r.Context(), req)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if id <= 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonDbOperationError))
	}
	a.RespJsonExit(r, nil, id)
}

// Put 修改定时任务信息
// @summary 修改定时任务信息
// @tags 	定时任务管理
// @Router 	/monitor/job [PUT]
func (a *jobApi) Put(r *ghttp.Request) {
	var req *define.JobApiUpdateReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	rs, err := service.Job.Update(r.Context(), req)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if rs <= 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonDbOperationError))
	}
	a.RespJsonExit(r, nil, rs)
}

// Delete 删除定时任务
// @summary 删除定时任务
// @tags 	定时任务管理
// @Router 	/monitor/job [DELETE]
func (a *jobApi) Delete(r *ghttp.Request) {
	var req *define.JobApiDeleteReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidParameter, err))
	}
	rs := service.Job.Delete(req.Ids)
	if rs > 0 {
		a.RespJsonExit(r, nil, rs)
	} else {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonDbOperationError))
	}
}

// Export 导出定时任务配置
// @summary 导出定时任务配置
// @tags 	定时任务管理
// @Router 	/monitor/job/export [GET]
func (a *jobApi) Export(r *ghttp.Request) {

}

// Run 执行定时任务(执行一次)
// @summary 执行定时任务(执行一次)
// @tags 	定时任务管理
// @Router 	/monitor/job/run [PUT]
func (a *jobApi) Run(r *ghttp.Request) {
	jobId := r.GetInt64("jobId")
	if jobId <= 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}
	job, err := service.Job.GetJobInfoByJobId(jobId)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if job == nil {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter,
			errcode.ErrCommonInvalidParameter.Message()+"任务不存在"))
	}
	job.MisfirePolicy = "0"
	job.CronExpression = "* * * * * *"
	err = service.Job.Start(job)
	if err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	} else {
		a.RespJsonExit(r, nil)
	}
}

// Start 启动定时任务
// @summary 启动定时任务
// @tags 	定时任务管理
// @Router 	/monitor/job/start [PUT]
func (a *jobApi) Start(r *ghttp.Request) {
	jobId := r.GetInt64("jobId")
	if jobId <= 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}

	job, err := service.Job.GetJobInfoByJobId(jobId)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if job == nil {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter,
			errcode.ErrCommonInvalidParameter.Message()+"任务不存在"))
	}
	err = service.Job.Start(job)
	if err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	} else {
		a.RespJsonExit(r, nil)
	}
}

// Stop 停止定时任务
// @summary 停止定时任务
// @tags 	定时任务管理
// @Router 	/monitor/job/stop [PUT]
func (a *jobApi) Stop(r *ghttp.Request) {
	jobId := r.GetInt64("jobId")
	if jobId <= 0 {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
	}

	job, err := service.Job.GetJobInfoByJobId(jobId)
	if err != nil {
		a.RespJsonExit(r, err)
	}
	if job == nil {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter,
			errcode.ErrCommonInvalidParameter.Message()+"任务不存在"))
	}
	err = service.Job.Stop(job)
	if err != nil {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonOperationFailed, err.Error()))
	} else {
		a.RespJsonExit(r, nil)
	}
}
