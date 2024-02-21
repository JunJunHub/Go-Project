package define

import "gowebserver/app/common/define"

type SysBlackListInfo struct {
	Id     int    `json:"id"`     // 配置ID
	Type   uint   `json:"type"`   // 1=黑名单 2=白名单
	IpAddr string `json:"ipAddr"` // 访问端IP地址
}

type SysBlackListAddReq struct {
	Type   uint   `json:"type"`   // 1=黑名单 2=白名单
	IpAddr string `json:"ipAddr"` // 访问端IP地址
}

type SysBlackListSelectReq struct {
	Type    uint   `p:"type"`    // 1=黑名单 2=白名单
	Keyword string `p:"keyword"` // 搜索查询关键字
	define.SelectPageReq
}

type SysBlackListInfoRes struct {
	List  []SysBlackListInfo `json:"list"`  // 系统访问黑名单列表
	Page  int                `json:"page"`  // 页码
	Size  int                `json:"size"`  // 每页条数
	Total int                `json:"total"` // 总数
}
