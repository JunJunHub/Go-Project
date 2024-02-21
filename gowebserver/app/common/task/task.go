// ==========================================================================================
// 定时任务栈: 支持定时任务管理的函数接口列表. 只有添加到栈中的任务接口, 才能通过系统定时任务管理真正的执行任务
//           各模块可以通过任务栈管理接口, 将函数接口添加到定时任务, 然后通过web配置定时任务规则.
// ==========================================================================================

package task

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/gtime"
)

//TasksEntity 任务信息
type TasksEntity struct {
	FuncName  string             //执行函数名
	Func      func() interface{} //执行函数
	Param     []string           //执行函数参数
	StartTime int64              //任务开始时间
}

//任务列表
var taskList = make([]TasksEntity, 0)

//TasksRunRecordFunc 任务执行记录函数接口, 定时任务日志管理向该模块注册使用
type TasksRunRecordFunc func(startTime, endTime int64, result interface{})

var recordJobLogFunc TasksRunRecordFunc

//Init 初始注册一些定时任务测试接口
func init() {
	task1 := TasksEntity{
		FuncName: "TasksDemo1",
		Func:     TasksDemo1,
	}
	task2 := TasksEntity{
		FuncName: "TaskDemo2",
		Func:     TasksDemo2,
	}
	task3 := TasksEntity{
		FuncName: "TaskDemo3",
		Func:     TasksDemo3,
	}
	task4 := TasksEntity{
		FuncName: "TaskDemo4",
		Func:     TasksDemo4,
	}

	if err := Add(task1); err != nil {
		g.Log().Error(err)
		return
	}
	if err := Add(task2); err != nil {
		g.Log().Error(err)
		return
	}
	if err := Add(task3); err != nil {
		g.Log().Error(err)
		return
	}
	if err := Add(task4); err != nil {
		g.Log().Error(err)
		return
	}
}

//GetByName 根据任务函数名获取任务实体信息, 对应任务管理的InvokeTarget信息
func GetByName(funcName string) *TasksEntity {
	var rst TasksEntity
	for _, task := range taskList {
		if task.FuncName == funcName {
			rst = task
			break
		}
	}
	return &rst
}

//Add 新增定时任务接口
func Add(task TasksEntity) error {
	if task.FuncName == "" {
		return gerror.New("执行方法名能为空")
	}
	if task.Func == nil {
		return gerror.New("任务执行函数为空")
	}
	taskList = append(taskList, task)
	return nil
}

//SetRecordJobRunLogFunc 设置记录定时任务执行日志的函数接口
func SetRecordJobRunLogFunc(funcCb TasksRunRecordFunc) {
	recordJobLogFunc = funcCb
}

//EditParams 修改任务参数
func (t *TasksEntity) EditParams(params []string) {
	t.Param = params
}

//Run 执行任务
func (t *TasksEntity) Run() {
	t.before()
	result := t.Func()
	t.after(result)
}

//before 任务执行前, 记录任务开始执行时间
func (t *TasksEntity) before() {
	t.StartTime = gtime.TimestampMilli()
}

//after 任务执行后, 记录定时任务执行日志
func (t *TasksEntity) after(result interface{}) {
	endTime := gtime.TimestampMilli()

	if recordJobLogFunc != nil {
		recordJobLogFunc(t.StartTime, endTime, result)
	}
}
