package api

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"gowebserver/app/common/errcode"
	"os"
)

var DownloadApi = new(downloadApi)

type downloadApi struct {
	BaseController
}

// Download
// @Summary 下载文件
// @Tags    公共
// @Param   fileName query string true "下载文件名"
// @Param   delete   query bool   true "下载完之后,是否删除服务器上的文件:false不删除 true删除"
// @Success 200 {object} response.Response "下载失败返回"
// @Router  /download [GET]
func (a *downloadApi) Download(r *ghttp.Request) {
	fileName := r.GetQueryString("fileName")
	del := r.GetQueryBool("delete")
	if fileName == "" {
		a.RespJsonExit(r, gerror.NewCode(errcode.ErrCommonInvalidParameter))
		return
	}

	//确定文件是否存在
	curDir, err := os.Getwd()
	if err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInternalError, err))
		return
	}
	filepath := curDir + g.Cfg().GetString("download.downPath") + "/" + fileName
	_, err = os.Stat(filepath)
	if err != nil || os.IsNotExist(err) {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
		return
	}

	//读取并返回要下载的文件
	r.Response.ServeFileDownload(filepath, fileName)
	if del {
		if err = os.Remove(filepath); err != nil {
			g.Log().Error(err)
		}
	}
}
