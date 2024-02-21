package api

import (
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/net/ghttp"
	"gowebserver/app/common/api"
	"gowebserver/app/common/errcode"
	"gowebserver/app/common/global"
	"gowebserver/app/customization/define"
)

type mpuApplianceSchemeApi struct {
	api.BaseController
}

var MpuApplianceSchemeApi = new(mpuApplianceSchemeApi)

func (a *mpuApplianceSchemeApi) Init(r *ghttp.Request) {
	a.Module = "电器设备调度预案分组"
	r.SetCtxVar(global.Module, a.Module)
}

func (a *mpuApplianceSchemeApi) GetList(r *ghttp.Request) {
	var req *define.MPUApplianceCtrlSchemeSelectPageReq
	if err := r.Parse(&req); err != nil {
		a.RespJsonExit(r, gerror.WrapCode(errcode.ErrCommonParsingParameter, err))
	}

	var schemeList []define.MPUApplianceCtrlScheme
	for schemeId := 1; schemeId < 10; schemeId++ {
		schemeList = append(schemeList, define.MPUApplianceCtrlScheme{
			GroupID:    req.GroupId,
			LoadState:  0,
			LoadTimeMS: 0,
			SchemeID:   int64(schemeId),
			SchemeName: fmt.Sprintf("%d楼会议室设备调度", schemeId),
		})
	}

	a.RespJsonExit(r, nil, define.MPUApplianceCtrlSchemeList{
		List:  schemeList,
		Page:  1,
		Size:  len(schemeList),
		Total: len(schemeList),
	})
}

func BuildLampTestData(schemeId int64) []define.MPUApplianceCtrlParam {
	var applianceCtrlParams []define.MPUApplianceCtrlParam
	for devId := 1; devId < 10; devId++ {
		applianceCtrlParams = append(applianceCtrlParams, define.MPUApplianceCtrlParam{
			SchemeID: schemeId,
			DevID:    fmt.Sprintf("%d", devId),
			DevType:  1,
			CtrlCmd:  "ON=1;Brilliance=0.60;",
			DevCapabilities: define.MPULampCapabilities{
				BSupportBrillianceCtrl: 1,
				BSupportOnOff:          1,
			},
		})
	}
	return applianceCtrlParams
}

func BuildAirConditionerTestData(schemeId int64) []define.MPUApplianceCtrlParam {
	var applianceCtrlParams []define.MPUApplianceCtrlParam
	for devId := 11; devId < 20; devId++ {
		applianceCtrlParams = append(applianceCtrlParams, define.MPUApplianceCtrlParam{
			SchemeID: schemeId,
			DevID:    fmt.Sprintf("%d", devId),
			DevType:  2,
			CtrlCmd:  "ON=1;Mode=0;Temperature=26;",
			DevCapabilities: define.MPUAirConditionerCapabilities{
				BSupportModes: define.BSupportModes{
					BSupportAutomaticMode:       1,
					BSupportCoolMode:            1,
					BSupportHeatMode:            1,
					BSupportAirDistributionMode: 1,
				},
				BSupportOnOff:                 1,
				BSupportTemperatureAdjustment: 1,
				BSupportWindSpeedAdjustment:   1,
			},
		})
	}
	return applianceCtrlParams
}

func BuildCurtainTestData(schemeId int64) []define.MPUApplianceCtrlParam {
	var applianceCtrlParams []define.MPUApplianceCtrlParam
	for devId := 21; devId < 30; devId++ {
		applianceCtrlParams = append(applianceCtrlParams, define.MPUApplianceCtrlParam{
			SchemeID: schemeId,
			DevID:    fmt.Sprintf("%d", devId),
			DevType:  3,
			CtrlCmd:  "ON=1;",
			DevCapabilities: define.MPUCurtainCapabilities{
				BSupportOnOff: 1,
			},
		})
	}
	return applianceCtrlParams
}

func (a *mpuApplianceSchemeApi) GetDetailInfo(r *ghttp.Request) {
	schemeId := r.GetInt64("schemeId")
	if schemeId == 0 {
		a.RespJsonExit(r, gerror.New("设备调度预案不存在"))
	}

	var applianceCtrlParams []define.MPUApplianceCtrlParam
	devType := r.GetInt64("devType")
	switch devType {
	case 1:
		applianceCtrlParams = BuildLampTestData(schemeId)
	case 2:
		applianceCtrlParams = BuildAirConditionerTestData(schemeId)
	case 3:
		applianceCtrlParams = BuildCurtainTestData(schemeId)
	}

	a.RespJsonExit(r, nil, define.MPUApplianceCtrlParamList{
		List:  applianceCtrlParams,
		Page:  1,
		Size:  len(applianceCtrlParams),
		Total: len(applianceCtrlParams),
	})
}

func (a *mpuApplianceSchemeApi) Load(r *ghttp.Request) {
	a.RespJsonExit(r, nil)
}

func (a *mpuApplianceSchemeApi) UnLoad(r *ghttp.Request) {
	a.RespJsonExit(r, nil)
}
func (a *mpuApplianceSchemeApi) Control(r *ghttp.Request) {
	a.RespJsonExit(r, nil)
}
