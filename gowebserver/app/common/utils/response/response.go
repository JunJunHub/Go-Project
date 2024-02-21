package response

import (
	"github.com/gogf/gf/encoding/gurl"
	"github.com/gogf/gf/errors/gcode"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/os/gview"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/utils"

	"context"
	"fmt"
	"net/http"
	"time"
)

type Response struct {
	//错误代码
	Code int `json:"code" example:"200"`
	//简单报错信息
	Msg string `json:"msg"`
	//响应时间戳(毫秒)
	TimeStamp int64 `json:"timestamp"`
	//数据集
	Data interface{} `json:"data"`
}

var response = new(Response)

//JsonExit 通用响应接口(exit)
func JsonExit(r *ghttp.Request, err error, data ...interface{}) {
	response.JsonExit(r, err, data...)
}

//RJson 通用响应接口(not exit)
func RJson(r *ghttp.Request, err error, data ...interface{}) {
	response.RJson(r, err, data...)
}

//UnAuthorizedJson 未认证返回JSON
func UnAuthorizedJson(r *ghttp.Request) {
	response.JsonExit(r, gerror.NewCode(errcode.ErrHttpUnauthorized))
}

//WriteTpl 模板
func WriteTpl(r *ghttp.Request, tpl string, view *gview.View, params ...gview.Params) error {
	return response.WriteTpl(r, tpl, view, params...)
}

//JsonExit 返回JSON数据并退出当前HTTP执行函数
func (res *Response) JsonExit(r *ghttp.Request, err error, data ...interface{}) {
	res.RJson(r, err, data...)
	r.Exit()
}

// RJson
// @summary Http接口返回数据,请求结果返回固定格式的JSON数据
// @param1  err error           "请求返回报错信息,包含错误码信息"
// @param2  data ...interface{} "返回数据内容"
// @return  nil
func (res *Response) RJson(r *ghttp.Request, err error, data ...interface{}) {
	responseData := interface{}(nil)
	if len(data) > 0 {
		responseData = data[0]
	}

	var (
		msg  string
		code = gerror.Code(err)
	)
	if err != nil {
		if code == gcode.CodeNil {
			code = errcode.ErrCommonUnknown
		}

		//1、错误码对象message统一格式为{#ErrCommonInvalidParameter}
		//例:
		//ErrCommonInvalidParameter = NewMpuapsErr(codeCommonInvalidParameter, "{#ErrCommonInvalidParameter}", nil)
		//
		//2、在配置文件中翻译成对应国家地区语言
		//例:
		//en.toml   : ErrCommonInvalidParameter = "Invalid Parameter"
		//zh_CN.toml: ErrCommonInvalidParameter = "无效的参数"
		//
		//依据Http请求头中 Accept-Language 返回对应国家区域的报错提示信息
		msg = err.Error()
		language := utils.ParseAcceptLanguage(r.GetHeader("Accept-Language"))
		g.I18n().SetLanguage(language)
		msg = g.I18n().T(context.Background(), fmt.Sprintf(`%s`, msg))

	} else if r.Response.Status > 0 && r.Response.Status != http.StatusOK {
		msg = http.StatusText(r.Response.Status)
		switch r.Response.Status {
		case http.StatusNotFound:
			code = errcode.ErrHttpNotFount
		case http.StatusForbidden:
			code = errcode.ErrHttpForbidden
		default:
			code = errcode.ErrCommonUnknown
		}
	} else {
		code = errcode.CodeOk
	}

	response = &Response{
		Code:      code.Code(),
		Msg:       msg,
		TimeStamp: time.Now().UnixNano() / 1e6,
		Data:      responseData,
	}
	r.SetParam("apiReturnRes", response)
	if err = r.Response.WriteJson(response); err != nil {
		g.Log().Error(err)
	}
}

//Redirect 重定向返回
func (res *Response) Redirect(r *ghttp.Request, location string, code ...int) {
	r.Response.RedirectTo(location, code...)
}

//WriteTpl 模板输出
func (res *Response) WriteTpl(r *ghttp.Request, tpl string, view *gview.View, params ...gview.Params) error {
	//绑定模板中需要用到的方法
	view.BindFuncMap(g.Map{
		// 根据长度i来切割字符串
		"subStr": func(str interface{}, i int) (s string) {
			s1 := gconv.String(str)
			if gstr.LenRune(s1) > i {
				s = gstr.SubStrRune(s1, 0, i) + "..."
				return s
			}
			return s1
		},
		// 格式化时间戳 年月日
		"timeFormatYear": func(time interface{}) string {
			return gtime.NewFromTimeStamp(gconv.Int64(time)).Format("Y-m-d")
		},
		// 格式化时间戳 年月日时分秒
		"timeFormatDateTime": func(time interface{}) string {
			return gtime.NewFromTimeStamp(gconv.Int64(time)).Format("Y-m-d H:i:s")
		},
		"add": func(a, b interface{}) int {
			return gconv.Int(a) + gconv.Int(b)
		},
	})
	//设置全局变量
	domain, _ := GetDomain(r)
	view.Assigns(g.Map{
		"domain": domain,
	})
	return r.Response.WriteTpl(tpl, params...)
}

// GetDomain 获取请求接口域名
func GetDomain(r *ghttp.Request) (string, error) {
	pathInfo, err := gurl.ParseURL(r.GetUrl(), -1)
	if err != nil {
		g.Log().Error(err)
		err = gerror.New("解析附件路径失败")
		return "", err
	}
	return fmt.Sprintf("%s://%s:%s/", pathInfo["scheme"], pathInfo["host"], pathInfo["port"]), nil
}
