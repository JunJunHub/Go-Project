package define

import (
	"gowebserver/app/common/define"
	"gowebserver/app/system/model"
)

//LoginHistoryApiSelectPageReq 查询列表请求参数
type LoginHistoryApiSelectPageReq struct {
	LoginName string `p:"loginName"` //登陆名
	Status    string `p:"status"`    //状态
	Ipaddr    string `p:"ipaddr"`    //登录地址
	define.SelectPageReq
}

//LoginHistoryApiDeleteReq API执行删除内容
type LoginHistoryApiDeleteReq struct {
	Ids string `p:"ids"  v:"required#请选择要删除的数据记录"`
}

//LoginHistoryServiceList 查询列表返回值
type LoginHistoryServiceList struct {
	List  []model.SysLoginHistory `json:"list"`
	Page  int                     `json:"page"`
	Size  int                     `json:"size"`
	Total int                     `json:"total"`
}
