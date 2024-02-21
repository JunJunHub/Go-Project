// =================================================================================
// 上下文信息数据结构定义
// =================================================================================

package define

import (
	"github.com/gogf/gf/frame/g"
	"gowebserver/app/system/model"
)

// Context 请求上下文结构
type Context struct {
	Uid            int64                // 用户ID
	Token          string               // 用户登录令牌
	AcceptLanguage string               // 客户端支持的语言: zh_CN、en
	User           *model.SysUserExtend // 用户详细信息
	Data           g.Map                // 自定KV变量，业务模块根据需要设置，不固定
}
