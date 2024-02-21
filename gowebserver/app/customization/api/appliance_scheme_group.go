package api

import (
	"fmt"
	"github.com/gogf/gf/net/ghttp"
	"gowebserver/app/common/api"
	"gowebserver/app/common/global"
	"gowebserver/app/customization/define"
)

type mpuApplianceSchemeGroupApi struct {
	api.BaseController
}

var MpuApplianceSchemeGroupApi = new(mpuApplianceSchemeGroupApi)

func (a *mpuApplianceSchemeGroupApi) Init(r *ghttp.Request) {
	a.Module = "电器设备调度预案分组"
	r.SetCtxVar(global.Module, a.Module)
}

func (a *mpuApplianceSchemeGroupApi) GetList(r *ghttp.Request) {
	var groupList []define.MPUApplianceCtrlSchemeGroup
	for groupId := 1; groupId < 5; groupId++ {
		groupList = append(groupList, define.MPUApplianceCtrlSchemeGroup{
			GroupID:   int64(groupId),
			GroupName: fmt.Sprintf("设备调度预案组%d", groupId),
			ParentID:  0,
		})
	}

	a.RespJsonExit(r, nil, define.MPUApplianceCtrlSchemeGroupList{
		List:  groupList,
		Page:  1,
		Size:  len(groupList),
		Total: len(groupList),
	})
}

func (a *mpuApplianceSchemeGroupApi) Add(r *ghttp.Request) {

}

func (a *mpuApplianceSchemeGroupApi) Del(r *ghttp.Request) {

}

func (a *mpuApplianceSchemeGroupApi) Put(r *ghttp.Request) {

}
