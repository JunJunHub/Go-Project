// ===============================
// 通用文件、图片上传API
// ===============================

package api

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/common/utils/upload_adapter"
	"net/url"

	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
)

var UploadApi = new(upload)

type upload struct {
	BaseController
}

func (c *upload) Init(r *ghttp.Request) {
	c.Module = "文件上传"
	r.SetCtxVar(global.Module, c.Module)
}

// UpImg 上传单张图片
// @Summary 上传单张图片
// @Description 上传单张图片
// @Tags 公共
// @Param upFile body string  true "upFile"
// @Success 0 {object} response.Response "{"code": 200, "timestamp": 1654875463112, "msg":"上传成功", "data": null}"
// @Router /upload/img [post]
// @Security
func (c *upload) UpImg(r *ghttp.Request) {
	upFile := r.GetUploadFile("file")
	info, err := uploadAdapter.Upload.UpImg(upFile)
	if err != nil {
		c.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidOperation, err, "文件上传失败"))
	}
	res := g.Map{
		"fileInfo": info,
	}
	c.RespJsonExit(r, nil, res)
}

// UpImgs
// @Summary 批量上传图片
// @Description 批量上传图片
// @Tags 公共
// @Param file body string  true "file"
// @Success 0 {object} response.Response "{"code": 200, "timestamp": 1654875463112, "msg":"上传成功", "data": null}"
// @Router  /upload/upImgs [post]
// @Security
func (c *upload) UpImgs(r *ghttp.Request) {
	upFiles := r.GetUploadFiles("file")
	infos, err := uploadAdapter.Upload.UpImgs(upFiles)
	if err != nil {
		c.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidOperation, err, "文件上传失败"))
	}
	res := g.Map{
		"fileInfos": infos,
	}
	c.RespJsonExit(r, nil, res)
}

// UpFile
// @Summary 单文件上传
// @Description 单文件上传
// @Tags 公共
// @Param file body string  true "file"
// @Success 0 {object} response.Response "{"code": 200, "timestamp": 1654875463112, "msg":"上传成功", "data": null}"
// @Router  /upload/upFile [post]
// @Security
func (c *upload) UpFile(r *ghttp.Request) {
	upFile := r.GetUploadFile("file")
	info, err := uploadAdapter.Upload.UpFile(upFile)
	if err != nil {
		c.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidOperation, err, "文件上传失败"))
	}
	res := g.Map{
		"fileInfo": info,
	}
	c.RespJsonExit(r, nil, res)
}

// UpFiles
// @Summary 批量上传文件
// @Description 批量上传文件
// @Tags 公共
// @Param file body string  true "file"
// @Success 0 {object} response.Response "{"code": 200, "timestamp": 1654875463112, "msg":"上传成功", "data": null}"
// @Router  /upload/upFiles [post]
// @Security
func (c *upload) UpFiles(r *ghttp.Request) {
	upFiles := r.GetUploadFiles("file")
	infos, err := uploadAdapter.Upload.UpFiles(upFiles)
	if err != nil {
		c.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonOperationFailed, err))
	}
	res := g.Map{
		"fileInfos": infos,
	}
	c.RespJsonExit(r, nil, res)
}

// CkEditorUp
// @Summary CkEditor编辑器上传附件
// @Description CkEditor编辑器上传附件
// @Tags 公共
// @Param upFile body string  true "upFile"
// @Success 0 {object} response.Response "{"code": 200, "timestamp": 1654875463112, "msg":"上传成功", "data": null}"
// @Router  /upload/ckEditorUp [post]
// @Security
func (c *upload) CkEditorUp(r *ghttp.Request) {
	upFile := r.GetUploadFile("upload")
	fType := gstr.ToLower(r.GetString("type"))
	var info *uploadAdapter.FileInfo
	var err error
	if fType == "images" {
		info, err = uploadAdapter.Upload.UpImg(upFile)
	} else if fType == "files" {
		info, err = uploadAdapter.Upload.UpFile(upFile)
	}
	if err != nil {
		c.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonInvalidOperation, err, "文件上传失败"))
	} else {
		parseInfo, _ := url.Parse(r.GetUrl())
		var fileUrl = info.FileUrl
		if !gregex.IsMatchString("^http", info.FileUrl) {
			fileUrl = parseInfo.Scheme + "://" + parseInfo.Host + "/" + info.FileUrl
		}
		c.RespJsonExit(r, nil, g.Map{"fileName": info.FileName, "uploaded": 1, "url": fileUrl})
	}
}
