// ================================================
// 文件上传适配器: 上传到本地,上传到云(腾讯云、阿里云...)
// ================================================

package uploadAdapter

import (
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/text/gregex"
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"
)

// FileInfo 上传的文件信息
type FileInfo struct {
	FileName string `json:"fileName"` //文件名称
	FileSize int64  `json:"fileSize"` //文件大小
	FileUrl  string `json:"fileUrl"`  //上传URL(保存路径)
	FileType string `json:"fileType"` //文件类型
}

//UploadParam 上传参数配置基本信息
type UploadParam struct {
	AdapterType     string //上传适配器类型(local、tencentCOS ...) 例: "local"
	FileSizeMax     string //支持上传文件最大大小                   例: "100M"
	FileTypeSupport string //支持上传的文件格式(后缀区分)            例: "doc,docx,xls,xlsx,,zip,rar,7z"
	ImgSizeMax      string //支持上传图片最大大小                   例: "10M"
	ImgTypeSupport  string //支持上传图片格式(后缀区分)              例: "jpg,jpeg,gif,png,npm"
}

type UploadAdapter interface {
	UpImg(file *ghttp.UploadFile) (fileInfo *FileInfo, err error)         //图片上传
	UpFile(file *ghttp.UploadFile) (fileInfo *FileInfo, err error)        //文件上传
	UpImgs(files []*ghttp.UploadFile) (fileInfos []*FileInfo, err error)  //批量上传图片
	UpFiles(files []*ghttp.UploadFile) (fileInfos []*FileInfo, err error) //批量上传文件
}

type upload struct {
	adapter UploadAdapter
}

var (
	Upload      *upload     //抽象的上传接口
	uploadParam UploadParam //上传参数
)

func init() {
	var adp UploadAdapter

	// 上传参数初始化(配置文件中读取)
	uploadParam = UploadParam{
		AdapterType:     g.Cfg().GetString("upload.type"),
		FileSizeMax:     g.Cfg().GetString("upload.fileSizeMax"),
		FileTypeSupport: g.Cfg().GetString("upload.fileTypeSupport"),
		ImgSizeMax:      g.Cfg().GetString("upload.imgSizeMax"),
		ImgTypeSupport:  g.Cfg().GetString("upload.imgTypeSupport"),
	}

	if uploadParam.AdapterType == "local" {
		adp = UploadLocalAdapter{
			UpPath: g.Cfg().GetString("upload.local.UpPath"),
		}
		Upload = &upload{
			adapter: adp,
		}
	} else if uploadParam.AdapterType == "tencentCOS" {
		// 使用腾讯云COS上传
		adp = UploadTencentCOSAdapter{
			UpPath:    g.Cfg().GetString("upload.tencentCOS.UpPath"),
			RawUrl:    g.Cfg().GetString("upload.tencentCOS.RawUrl"),
			SecretID:  g.Cfg().GetString("upload.tencentCOS.SecretID"),
			SecretKey: g.Cfg().GetString("upload.tencentCOS.SecretKey"),
		}
		Upload = &upload{
			adapter: adp,
		}
	} // todo 扩展上传适配器
}

func (u upload) UpImg(file *ghttp.UploadFile) (fileInfo *FileInfo, err error) {
	return u.adapter.UpImg(file)
}

func (u upload) UpFile(file *ghttp.UploadFile) (fileInfo *FileInfo, err error) {
	return u.adapter.UpFile(file)
}

func (u upload) UpImgs(files []*ghttp.UploadFile) (fileInfos []*FileInfo, err error) {
	return u.adapter.UpImgs(files)
}

func (u upload) UpFiles(files []*ghttp.UploadFile) (fileInfos []*FileInfo, err error) {
	return u.adapter.UpFiles(files)
}

//SetUploadParam 设置上传参数
func SetUploadParam(param UploadParam) {
	uploadParam = param
}

//判断上传文件类型是否合法
//传参示例: "dept.doc" "doc,docx,xls,xlsx,,zip,rar,7z"
func checkFileType(fileName, typeString string) bool {
	suffix := gstr.SubStrRune(fileName, gstr.PosRRune(fileName, ".")+1, gstr.LenRune(fileName)-1)
	imageType := gstr.Split(typeString, ",")
	rightType := false
	for _, v := range imageType {
		if gstr.Equal(suffix, v) {
			rightType = true
			break
		}
	}
	return rightType
}

//检查文件大小是否合法
//传参示例: "35M" 1024(单位字节)
func checkSize(configSize string, fileSize int64) (bool, error) {
	match, err := gregex.MatchString(`^([0-9]+)(?i:([a-z]*))$`, configSize)
	if err != nil {
		return false, err
	}
	if len(match) == 0 {
		err = gerror.New("上传文件大小未设置，请在后台配置，格式为（30M,30k,30MB）")
		return false, err
	}
	var cfSize int64
	switch gstr.ToUpper(match[2]) {
	case "MB", "M":
		cfSize = gconv.Int64(match[1]) * 1024 * 1024
	case "KB", "K":
		cfSize = gconv.Int64(match[1]) * 1024
	case "":
		cfSize = gconv.Int64(match[1])
	}
	if cfSize == 0 {
		err = gerror.New("上传文件大小未设置，请在后台配置，格式为（30M,30k,30MB），最大单位为MB")
		return false, err
	}
	return cfSize >= fileSize, nil
}
