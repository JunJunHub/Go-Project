package define

type MPUApplianceCtrlSchemeGroup struct {
	GroupID   int64  `json:"groupId"`   // 分组ID，分组ID
	GroupName string `json:"groupName"` // 分组名称，分组名称
	ParentID  int64  `json:"parentId"`  // 父分组ID，父分组ID
}

type MPUApplianceCtrlSchemeGroupList struct {
	List  []MPUApplianceCtrlSchemeGroup `json:"list"`
	Page  int                           `json:"page"`
	Size  int                           `json:"size"`
	Total int                           `json:"total"`
}
