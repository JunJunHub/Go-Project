// ===========================================
// 上传文件到腾讯云
// ============================================

package uploadAdapter

import (
	"context"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gfile"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/util/grand"
	"github.com/tencentyun/cos-go-sdk-v5"
	"github.com/tencentyun/cos-go-sdk-v5/debug"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type UploadTencentCOSAdapter struct {
	UpPath    string
	SecretID  string
	SecretKey string
	RawUrl    string
}

func (u UploadTencentCOSAdapter) UpImg(file *ghttp.UploadFile) (fileInfo *FileInfo, err error) {
	return u.upByType(file, "img")
}

func (u UploadTencentCOSAdapter) UpFile(file *ghttp.UploadFile) (fileInfo *FileInfo, err error) {
	return u.upByType(file, "file")
}

func (u UploadTencentCOSAdapter) UpImgs(files []*ghttp.UploadFile) (fileInfos []*FileInfo, err error) {
	return u.upBathByType(files, "img")
}

func (u UploadTencentCOSAdapter) UpFiles(files []*ghttp.UploadFile) (fileInfos []*FileInfo, err error) {
	return u.upBathByType(files, "file")
}

//文件上传 img|file
func (u UploadTencentCOSAdapter) upByType(file *ghttp.UploadFile, fType string) (fileInfo *FileInfo, err error) {
	if file == nil {
		err = gerror.New("未上传任何文件")
		return
	}

	var fileTypeList string
	var fileSizeMax string

	if fType == "img" {
		fileTypeList = uploadParam.ImgTypeSupport
		fileSizeMax = uploadParam.ImgSizeMax
	} else if fType == "file" {
		fileTypeList = uploadParam.FileTypeSupport
		fileSizeMax = uploadParam.FileSizeMax
	}

	//检测文件类型
	rightType := checkFileType(file.Filename, fileTypeList)
	if !rightType {
		err = gerror.New("上传文件类型错误，只能包含后缀为：" + fileTypeList + "的文件。")
		return
	}

	rightSize, err := checkSize(fileSizeMax, file.Size)
	if err != nil {
		return
	}
	if !rightSize {
		err = gerror.New("上传文件超过最大尺寸：" + fileSizeMax)
		return
	}
	var path string
	path, err = u.upAction(file)
	if err != nil {
		return
	}
	fileInfo = &FileInfo{
		FileName: file.Filename,
		FileSize: file.Size,
		FileUrl:  u.getUrl(path),
		FileType: file.Header.Get("Content-type"),
	}
	return
}

//批量上传 img|file
func (u UploadTencentCOSAdapter) upBathByType(files []*ghttp.UploadFile, fType string) (fileInfos []*FileInfo, err error) {
	if len(files) == 0 {
		err = gerror.New("未上传任何文件")
		return
	}

	var fileTypeList string
	var fileSizeMax string

	if fType == "img" {
		fileTypeList = uploadParam.ImgTypeSupport
		fileSizeMax = uploadParam.ImgSizeMax
	} else if fType == "file" {
		fileTypeList = uploadParam.FileTypeSupport
		fileSizeMax = uploadParam.FileSizeMax
	}

	for _, file := range files {
		//检测文件类型
		rightType := checkFileType(file.Filename, fileTypeList)
		if !rightType {
			err = gerror.New("上传文件类型错误，只能包含后缀为：" + fileTypeList + "的文件。")
			return
		}
		var rightSize bool
		rightSize, err = checkSize(fileSizeMax, file.Size)
		if err != nil {
			return
		}
		if !rightSize {
			err = gerror.New("上传文件超过最大尺寸：" + fileSizeMax)
			return
		}
	}
	for _, file := range files {
		var path string
		path, err = u.upAction(file)
		if err != nil {
			return
		}
		fileInfo := &FileInfo{
			FileName: file.Filename,
			FileSize: file.Size,
			FileUrl:  u.getUrl(path),
			FileType: file.Header.Get("Content-type"),
		}
		fileInfos = append(fileInfos, fileInfo)
	}
	return
}

// 上传到腾讯cos操作
func (u UploadTencentCOSAdapter) upAction(file *ghttp.UploadFile) (path string, err error) {
	name := gfile.Name(file.Filename) + "_" + strings.ToLower(strconv.FormatInt(gtime.TimestampNano(), 36)+grand.S(6))
	name = name + gfile.Ext(file.Filename)

	path = u.getUpPath() + name
	url, _ := url.Parse(u.RawUrl)
	b := &cos.BaseURL{BucketURL: url}
	c := cos.NewClient(b, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  u.SecretID,
			SecretKey: u.SecretKey,
			Transport: &debug.DebugRequestTransport{
				RequestHeader:  false,
				RequestBody:    false,
				ResponseHeader: false,
				ResponseBody:   false,
			},
		},
	})
	opt := &cos.ObjectPutOptions{
		ObjectPutHeaderOptions: &cos.ObjectPutHeaderOptions{
			ContentLength: int64(file.Size),
		},
	}
	var f io.ReadCloser
	f, err = file.Open()
	if err != nil {
		return
	}
	defer f.Close()
	_, err = c.Object.Put(context.Background(), path, f, opt)
	return
}

func (u UploadTencentCOSAdapter) getUpPath() (upPath string) {
	upPath = u.UpPath + gtime.Date() + "/"
	return
}

func (u UploadTencentCOSAdapter) getUrl(path string) string {
	url := u.RawUrl + path
	return url
}
