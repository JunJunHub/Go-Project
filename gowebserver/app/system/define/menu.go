package define

//SysMenuInfo 菜单信息
type SysMenuInfo struct {
	Id            int64  `orm:"id,primary"     json:"id"            description:"菜单id"`
	Pid           int64  `orm:"pid"            json:"pid"           description:"父菜单id"`
	Title         string `orm:"title"          json:"title"         description:"菜单标题"`
	TitleTag      string `orm:"title_tag"      json:"titleTag"      description:"菜单标签"`
	Type          uint   `orm:"type"           json:"type"          description:"菜单类型（0目录 1菜单 2按钮）"`
	Component     string `orm:"component"      json:"component"     description:"前端组件路径"`
	Icon          string `orm:"icon"           json:"icon"          description:"菜单图标"`
	Router        string `orm:"router"         json:"router"        description:"后端路由规则"`
	RestInterface string `orm:"rest_interface" json:"restInterface" description:"请求接口"`
	Weigh         int    `orm:"weigh"          json:"weigh"         description:"菜单权重(菜单显示顺序)"`
	IsHide        uint   `orm:"is_hide"        json:"isHide"        description:"是否隐藏: 0不隐藏 1隐藏"`
	Remark        string `orm:"remark"         json:"remark"        description:"备注"`
}

//SysMenuTree 功能菜单树
type SysMenuTree struct {
	*SysMenuInfo
	ChildTree []*SysMenuTree `json:"childTree"`
}
