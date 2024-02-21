// =================================
// 上传文件到本地
// ================================

package uploadAdapter

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/text/gstr"
	"os"
)

type UploadLocalAdapter struct {
	UpPath string
}

// UpImg 上传图片
func (up UploadLocalAdapter) UpImg(file *ghttp.UploadFile) (fileInfo *FileInfo, err error) {
	return up.upByType(file, "img")
}

// UpFile 上传文件
func (up UploadLocalAdapter) UpFile(file *ghttp.UploadFile) (fileInfo *FileInfo, err error) {
	return up.upByType(file, "file")
}

// UpImgs 批量上传图片
func (up UploadLocalAdapter) UpImgs(files []*ghttp.UploadFile) (fileInfos []*FileInfo, err error) {
	return up.upBathByType(files, "img")
}

// UpFiles 批量上传文件
func (up UploadLocalAdapter) UpFiles(files []*ghttp.UploadFile) (fileInfos []*FileInfo, err error) {
	return up.upBathByType(files, "file")
}

//文件上传 img|file
func (up UploadLocalAdapter) upByType(file *ghttp.UploadFile, fType string) (fileInfo *FileInfo, err error) {
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
	//检查文件大小
	rightSize, err := checkSize(fileSizeMax, file.Size)
	if err != nil {
		return
	}
	if !rightSize {
		err = gerror.New("上传文件超过最大尺寸：" + fileSizeMax)
		return
	}

	path := up.getUpPath()
	fileName, err := file.Save(path)
	if err != nil {
		return
	}
	fileInfo = &FileInfo{
		FileName: file.Filename,
		FileSize: file.Size,
		FileUrl:  up.getUrl(path, fileName),
		FileType: file.Header.Get("Content-type"),
	}
	return
}

//批量上传 img|file
func (up UploadLocalAdapter) upBathByType(files []*ghttp.UploadFile, fType string) (fileInfos []*FileInfo, err error) {
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
			err = gerror.New("上传文件类型错误，只能包含后缀为：" + fileTypeList + "的文件")
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
	path := up.getUpPath()
	for _, file := range files {
		var fileName string
		fileName, err = file.Save(path)
		if err != nil {
			return
		}
		fileInfo := &FileInfo{
			FileName: file.Filename,
			FileSize: file.Size,
			FileUrl:  up.getUrl(path, fileName),
			FileType: file.Header.Get("Content-type"),
		}
		fileInfos = append(fileInfos, fileInfo)
	}
	return
}

func (up UploadLocalAdapter) getUpPath() (upPath string) {
	curDir, err := os.Getwd()
	if err != nil {
		curDir = "/msp/mpuaps/"
	}
	upPath = curDir + up.UpPath + gtime.Date() + "/"
	return
}

func (up UploadLocalAdapter) getUrl(path, fileName string) string {
	url := gstr.SubStr(path, gstr.Pos(path, up.UpPath)+1) + fileName
	return url
}
