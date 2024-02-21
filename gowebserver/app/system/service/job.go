// ================================================================================
// 定时任务管理规则: 同一个任务方法可以配置多个定时任务, 同一个任务方法同时只能有一个定时任务在执行.
// ================================================================================

package service

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gcron"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"gowebserver/app/common/task"
	"gowebserver/app/common/utils/convert"
	"gowebserver/app/common/utils/page"
	"gowebserver/app/system/dao"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"
	"strings"
)

// Job 定时任务管理服务
var Job = &jobService{}

type jobService struct{}

//CheckJobNameUnique 检查任务名是否重复  false不重复 true重复
func (s *jobService) CheckJobNameUnique(jobName string) bool {
	i, err := dao.SysJob.Ctx(context.Background()).FindCount(g.Map{
		dao.SysJob.Columns.JobName: jobName,
	})
	if err != nil {
		return false
	}
	return i > 0
}

//GetJobInfoByInvokeTarget 检查任务方法是否在有对应的定时任务在执行
func (s *jobService) GetJobInfoByInvokeTarget(invokeTag string) (*model.SysJob, error) {
	record, err := dao.SysJob.Ctx(context.Background()).FindOne(dao.SysJob.Columns.InvokeTarget, invokeTag)
	if err != nil {
		return nil, err
	}
	var job *model.SysJob
	if err = record.Struct(job); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *jobService) GetJobInfoByJobId(id int64) (*model.SysJob, error) {
	record, err := dao.SysJob.Ctx(context.Background()).FindOne(dao.SysJob.Columns.JobId, id)
	if err != nil {
		return nil, err
	}

	var job *model.SysJob
	if err = record.Struct(job); err != nil {
		return nil, err
	}
	return job, nil
}

func (s *jobService) GetList(param *define.JobApiSelectPageReq) *define.JobServiceList {
	m := dao.SysJob.Ctx(context.Background()).As("t")
	if param != nil {
		if param.JobName != "" {
			m = m.Where("t.job_name like ?", "%"+param.JobName+"%")
		}
		if param.JobGroup != "" {
			m = m.Where("t.job_group like = ?", param.JobGroup)
		}
		if param.InvokeTarget != "" {
			m = m.Where("t.invoke_target like = ?", param.InvokeTarget)
		}
		if param.CronExpression != "" {
			m = m.Where("t.cron_expression like = ?", param.CronExpression)
		}
		if param.MisfirePolicy != "" {
			m = m.Where("t.misfire_policy = ?", param.MisfirePolicy)
		}
		if param.Concurrent != "" {
			m = m.Where("t.concurrent = ?", param.Concurrent)
		}
		if param.Status != "" {
			m = m.Where("t.status = ?", param.Status)
		}
		if param.BeginTime != "" {
			m = m.Where("date_format(t.create_time,'%y%m%d') >= date_format(?,'%y%m%d') ", param.BeginTime)
		}
		if param.EndTime != "" {
			m = m.Where("date_format(t.create_time,'%y%m%d') <= date_format(?,'%y%m%d') ", param.EndTime)
		}
	}
	total, err := m.Count()
	if err != nil {
		return nil
	}

	pageInfo := page.CreatePaging(param.PageNum, param.PageSize, total)
	result := &define.JobServiceList{
		Page:  pageInfo.PageNum,
		Size:  pageInfo.PageSize,
		Total: pageInfo.Total,
	}
	m = m.Limit(pageInfo.StartNum, pageInfo.PageSize)
	if param.OrderByColumn != "" && param.IsAsc != "" {
		param.OrderByColumn = gstr.CaseSnake(param.OrderByColumn)
		m = m.Order(param.OrderByColumn + " " + param.IsAsc)
	} else {
		m = m.Order("t.create_time desc")
	}
	if err = m.Scan(&result.List); err != nil {
		return nil
	}
	return result
}

func (s *jobService) Create(ctx context.Context, req *define.JobApiCreateReq) (int64, error) {
	//检查任务名称是否存在
	if s.CheckJobNameUnique(req.JobName) {
		return 0, gerror.New("任务名称已经存在")
	}

	user := Context.GetUser(ctx)
	var entity model.SysJob
	entity.JobGroup = req.JobGroup
	entity.CreateTime = gtime.Now()
	entity.CreateBy = user.LoginName
	var editReq *define.JobApiUpdateReq
	if err := gconv.Struct(req, &editReq); err != nil {
		return 0, err
	}
	return s.save(&entity, editReq)
}

func (s *jobService) Delete(ids string) int64 {
	jobIdList := convert.ToInt64Array(ids, ",")
	result, err := dao.SysJob.Ctx(context.Background()).FindAll("job_id in (?)", jobIdList)
	if err != nil {
		g.Log().Error(err)
		return 0
	}
	var jobList []model.SysJob
	if err = result.Structs(jobList); err != nil {
		g.Log().Error(err)
		return 0
	}
	if jobList == nil {
		return 0
	}

	for _, job := range jobList {
		gcron.Remove(job.JobName)
	}

	delRst, err := dao.SysJob.Ctx(context.Background()).Delete("job_id in (?)", jobIdList)
	if err != nil {
		g.Log().Error(err)
		return 0
	}
	nums, _ := delRst.RowsAffected()
	return nums
}

func (s *jobService) Update(ctx context.Context, req *define.JobApiUpdateReq) (int64, error) {
	//获取库中的任务信息
	entity, err := s.GetJobInfoByJobId(req.JobId)
	if err != nil {
		return 0, nil
	}
	if entity == nil {
		return 0, gerror.New("数据不存在")
	}

	//判断任务是否已添加到cron
	tmp := gcron.Search(entity.JobName)
	if tmp != nil {
		gcron.Remove(entity.JobName)
	}

	//更新信息
	user := Context.GetUser(ctx)
	entity.UpdateTime = gtime.Now()
	entity.UpdateBy = user.LoginName
	return s.save(entity, req)
}

func (s *jobService) save(job *model.SysJob, req *define.JobApiUpdateReq) (int64, error) {
	var (
		rs  int64
		err error
	)
	job.JobName = req.JobName
	job.JobParams = req.JobParams
	job.InvokeTarget = req.InvokeTarget
	job.CronExpression = req.CronExpression
	job.MisfirePolicy = req.MisfirePolicy
	job.Concurrent = req.Concurrent
	job.Status = req.Status
	job.Remark = req.Remark

	//保存
	result, err := dao.SysJob.Ctx(context.Background()).Data(job).Save()
	if err != nil {
		return 0, err
	}
	if job.JobId == 0 {
		// 新增
		rs, err = result.LastInsertId()
	} else {
		// 更新
		rs, err = result.RowsAffected()
	}
	if err != nil {
		return 0, err
	}

	// 启动
	if req.Status == "0" {
		if err = s.Start(job); err != nil {
			return rs, gerror.New("启动任务失败，您可以手动启动")
		}
	}
	return rs, err
}

//Start 启动任务
func (s *jobService) Start(entity *model.SysJob) error {
	//判断task任务栈是否存在指定的任务接口
	f := task.GetByName(entity.InvokeTarget)
	if f == nil || f.Func == nil {
		return gerror.New("当前task目录没有绑定这个方法")
	}

	//判断任务接口是否有定时任务正在执行, 若果有停止之前的任务
	if jobRunning, err := s.GetJobInfoByInvokeTarget(entity.InvokeTarget); err == nil && jobRunning != nil {
		//停止之前的任务
		if err = s.Stop(jobRunning); err != nil {
			return gerror.New("存在定时任务正在执行当前任务方法,停止之前的任务失败")
		}
	}

	//更新参数
	paramArr := strings.Split(entity.JobParams, "|")
	f.EditParams(paramArr)

	rs := gcron.Search(entity.JobName)
	if rs == nil {
		var err error
		if entity.MisfirePolicy == "1" {
			// 判断是否并发执行
			if entity.Concurrent == "1" {
				// 允许并发
				_, err = gcron.Add(entity.CronExpression, f.Run, entity.JobName)
			} else {
				// 同时只允许有一个该任务运行
				_, err = gcron.AddSingleton(entity.CronExpression, f.Run, entity.JobName)
			}
			if err != nil {
				return err
			}
		} else {
			_, err = gcron.AddOnce(entity.CronExpression, f.Run, entity.JobName)
			if err != nil {
				return err
			}
		}
	}

	//启动
	gcron.Start(entity.JobName)
	if entity.MisfirePolicy == "1" {
		entity.Status = "0"
		if _, err := dao.SysJob.Ctx(context.Background()).Data(entity).Save(); err != nil {
			return err
		}
	}
	return nil
}

//Stop 停止任务
func (s *jobService) Stop(job *model.SysJob) error {
	f := task.GetByName(job.JobName)
	if f == nil {
		return gerror.New("当前task目录下没有绑定这个方法")
	}

	rs := gcron.Search(job.JobName)
	if rs != nil {
		gcron.Stop(job.JobName)
	}
	job.Status = "1"
	if _, err := dao.SysJob.Ctx(context.Background()).Data(job).Save(); err != nil {
		return err
	}
	return nil
}

//Restart 重启任务
func (s *jobService) Restart() bool {
	result, err := dao.SysJob.Ctx(context.Background()).Where(dao.SysJob.Columns.Status, "0").
		Where(dao.SysJob.Columns.MisfirePolicy, "1").All()
	if err != nil {
		g.Log().Error("重启定时任务失败")
		return false
	}
	var entityList []model.SysJob
	if err = result.Structs(&entityList); err != nil {
		g.Log().Error("重启定时任务失败")
		return false
	}
	if len(entityList) == 0 {
		return true
	}
	for _, entity := range entityList {
		if err = s.Start(&entity); err != nil {
			g.Log().Errorf("定时任务：[%s] 启动失败， 请手动启动", entity.JobName)
		}
	}
	return true
}
