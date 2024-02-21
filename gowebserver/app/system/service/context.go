package service

import (
	"gowebserver/app/common/global"
	"gowebserver/app/system/define"
	"gowebserver/app/system/model"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"

	"context"
)

// Context 上下文管理服务
var Context = contextShared{}

type contextShared struct{}

// Init 初始化上下文对象指针到上下文对象中，以便后续的请求流程中可以修改。
func (s *contextShared) Init(r *ghttp.Request, customCtx *define.Context) {
	r.SetCtxVar(global.ContextKey, customCtx)
}

// GetCtx 获取上下文变量,如果没有设置,那么返回nil
func (s *contextShared) GetCtx(ctx context.Context) *define.Context {
	value := ctx.Value(global.ContextKey)
	if value == nil {
		return nil
	}
	if localCtx, ok := value.(*define.Context); ok {
		return localCtx
	}
	return nil
}

// SetUser 将用户详细信息设置到请求上下文中
func (s *contextShared) SetUser(ctx context.Context, ctxUser *model.SysUserExtend) {
	s.GetCtx(ctx).Uid = ctxUser.UserId
	s.GetCtx(ctx).User = ctxUser
}

// SetToken 将用户登录令牌设置到请求上下文
func (s *contextShared) SetToken(ctx context.Context, token string) {
	s.GetCtx(ctx).Token = token
}

// SetData 请求上下文中记录自定义扩展信息
func (s *contextShared) SetData(ctx context.Context, data g.Map) {
	s.GetCtx(ctx).Data = data
}

// GetUser 获取登录用户
func (s *contextShared) GetUser(ctx context.Context) (user *model.SysUserExtend) {
	ctxInfo := s.GetCtx(ctx)
	if ctxInfo == nil || ctxInfo.User == nil {
		//g.Log().Debugf("获取上下文用户信息为空: %v\n%v", ctx, g.Log().GetStack())
		return nil
	}
	return ctxInfo.User
}

// GetUserId 获取登录用户ID
func (s *contextShared) GetUserId(ctx context.Context) (userId int64) {
	ctxInfo := s.GetCtx(ctx)
	if ctxInfo == nil || ctxInfo.User == nil {
		//g.Log().Debugf("获取上下文用户信息为空: %v\n", ctx, g.Log().GetStack())
		return
	}
	return ctxInfo.User.UserId
}

// GetUserName 获取用户登录名
func (s *contextShared) GetUserName(ctx context.Context) (loginName string) {
	ctxInfo := s.GetCtx(ctx)
	if ctxInfo == nil || ctxInfo.User == nil {
		//g.Log().Debugf("获取上下文用户信息为空: %v\n%v", ctx, g.Log().GetStack())
		return ""
	}
	return ctxInfo.User.LoginName
}

// GetUserToken 获取用户令牌
func (s *contextShared) GetUserToken(ctx context.Context) (loginToken string) {
	ctxInfo := s.GetCtx(ctx)
	if ctxInfo == nil {
		//g.Log().Debugf("获取上下文用户信息为空: %v\n", ctx, g.Log().GetStack())
		return ""
	}
	return ctxInfo.Token
}

// SetAcceptLanguage 将客户端支持的语言类型设置到请求上下文
func (s *contextShared) SetAcceptLanguage(ctx context.Context, language string) {
	s.GetCtx(ctx).AcceptLanguage = language
}

// GetAcceptLanguage 获取客户端支持的语言类型
func (s *contextShared) GetAcceptLanguage(ctx context.Context) (language string) {
	ctxInfo := s.GetCtx(ctx)
	if ctxInfo == nil {
		//g.Log().Debugf("获取上下文用户信息为空: %v\n", ctx,g.Log().GetStack())
		return ""
	}
	return ctxInfo.AcceptLanguage
}
