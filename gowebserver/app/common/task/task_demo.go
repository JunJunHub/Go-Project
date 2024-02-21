// ======================================
// 定时任务接口示例, 返回值空接口
// ======================================

package task

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
)

func TasksDemo1() interface{} {
	g.Log().Info("TasksDemo1 执行并返回空")
	return nil
}

func TasksDemo2() interface{} {
	g.Log().Info("TasksDemo2 执行并返回错误信息")
	return gerror.New("任务结束")
}

func TasksDemo3() interface{} {
	g.Log().Info("TasksDemo3 执行并返回BOOL值")
	return true
}

func TasksDemo4() interface{} {
	g.Log().Info("TasksDemo4 执行并返回字符串")
	return "ok"
}
