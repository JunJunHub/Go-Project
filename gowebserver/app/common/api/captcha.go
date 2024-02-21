// ==========================================================================
// 获取图片验证码
// ==========================================================================

package api

import (
	"gowebserver/app/common/global"
	"gowebserver/app/common/service"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var CaptchaApi = new(captchaApi)

type captchaApi struct {
	BaseController
}

func (a *captchaApi) Init(r *ghttp.Request) {
	a.Module = "Auth"
	r.SetCtxVar(global.Module, a.Module)
}

// CaptchaImage
// @Summary 获取验证码图片信息
// @Description 获取验证码图片信息
// @Tags Auth
// @Success 200 {object} response.Response "例:{"code": 200, "timestamp": 1654875463112, "msg":"操作成功", "data": {"idKey":"...", "base64Img": "..."}}"
// @Router /auth/captchaImage [get]
// @Security
func (a *captchaApi) CaptchaImage(r *ghttp.Request) {
	idKeyC, base64stringC := service.Captcha.GetVerifyImgString()
	a.RespJsonExit(r, nil, g.Map{
		"idKey":     idKeyC,
		"base64Img": base64stringC,
	})
}
