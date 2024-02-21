// ==========================
// 封装值校验
// ==========================

package valueCheck

import (
	"context"
	"github.com/gogf/gf/frame/g"
)

//ErrIsNil
//@summary 校验返回错误对象是否为空
//         tips: 基于 defer recover 机制
//         1、需结合 g.Try 使用
//         2、error不为空,抛出错误信息;
//@param1  ctx context.Context "上下文"
//@param2  err error "错误对象"
//@param3  msg ...string "抛出报错信息"
//@return1 nil
func ErrIsNil(ctx context.Context, err error, msg ...string) {
	if !g.IsNil(err) {
		g.Log().Error(ctx, err)
		if len(msg) > 0 {
			panic(msg[0])
		} else {
			panic(err.Error())
		}
	}
}

//ValueIsNil
//@summary 校验值是否为空
//         tips: 基于 defer recover 机制
//         1、需结合 g.Try 使用
//         2、value为空,抛出错误信息;
//@param1  value interface{} "值"
//@param2  msg string "抛出报错信息"
func ValueIsNil(value interface{}, msg string) {
	if g.IsNil(value) {
		g.Log().Warning(msg, g.Log().GetStack())
		panic(msg)
	}
}
