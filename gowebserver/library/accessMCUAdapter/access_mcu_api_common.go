//==================================
// 封装MCU各版本通用API
//==================================

package accessMCUAdapter

import (
	"fmt"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/util/gconv"
	"io"
	"strconv"
	"strings"
	"time"
)

// respDecodeToJson
// @summary 解析MCU响应数据,并判断返回状态代码
// @param1  respString string "http返回消息"
// @return1 *gjson.Json "返回解析后的json对象"
// @return2 error "解析错误"
func respDecodeToJson(respString string) (*gjson.Json, error) {
	jsonObj, err := gjson.DecodeToJson(respString)
	if err != nil {
		return nil, err
	}
	//返回错误码校验
	if jsonObj.GetInt("success") != 1 {
		logger().Error(respString)
		//TODO
		//有可能返回错误码为0,与自定义MCU连接状态错误码重复.导致错误码翻译错误
		MCUErrCode := jsonObj.GetInt("error_code")
		if 0 == MCUErrCode {
			return jsonObj, gerror.New(strconv.Itoa(int(EMCUNetConnErr)))
		}
		return jsonObj, gerror.New(strconv.Itoa(MCUErrCode))
	}
	return jsonObj, nil
}

// newMCUPostReq
// @summary 创建post类请求并携带必要的头部信息
func newMCUPostReq(mcuApiLevel uint, mcuUserCookie string, timeout ...time.Duration) *ghttp.Client {
	gClient := g.Client().SetHeader("Accept", "application/json").
		SetContentType("application/x-www-form-urlencoded").
		SetHeader("Api_level", gconv.String(mcuApiLevel)).
		SetCookie("SSO_COOKIE_KEY", mcuUserCookie)

	if len(timeout) > 0 {
		gClient = gClient.Timeout(timeout[0] * time.Second)
	} else {
		//默认超时时间5s
		gClient = gClient.Timeout(5 * time.Second)
	}
	return gClient
}

// newMCUPutReq
// @summary 创建put类请求并携带必要的头部信息
func newMCUPutReq(mcuApiLevel uint, mcuUserCookie string, timeout ...time.Duration) *ghttp.Client {
	gClient := g.Client().SetHeader("Accept", "application/json").
		SetContentType("application/x-www-form-urlencoded").
		SetHeader("Api_level", gconv.String(mcuApiLevel)).
		SetCookie("SSO_COOKIE_KEY", mcuUserCookie)

	if len(timeout) > 0 {
		gClient = gClient.Timeout(timeout[0] * time.Second)
	} else {
		//默认超时时间5s
		gClient = gClient.Timeout(5 * time.Second)
	}
	return gClient
}

// newMCUDelReq
// @summary 创建del类请求并携带必要的头部信息
func newMCUDelReq(mcuApiLevel uint, mcuUserCookie string, timeout ...time.Duration) *ghttp.Client {
	gClient := g.Client().SetHeader("Accept", "application/json").
		SetHeader("Api_level", gconv.String(mcuApiLevel)).
		SetCookie("SSO_COOKIE_KEY", mcuUserCookie)

	if len(timeout) > 0 {
		gClient = gClient.Timeout(timeout[0] * time.Second)
	} else {
		//默认超时时间5s
		gClient = gClient.Timeout(5 * time.Second)
	}
	return gClient
}

// newMCUGetReq
// @summary 创建get类请求并携带必要的头部信息
func newMCUGetReq(mcuApiLevel uint, mcuUserCookie string, timeout ...time.Duration) *ghttp.Client {
	gClient := g.Client().SetHeader("Accept", "application/json").
		SetHeader("Api_level", gconv.String(mcuApiLevel)).
		SetCookie("SSO_COOKIE_KEY", mcuUserCookie)

	if len(timeout) > 0 {
		gClient = gClient.Timeout(timeout[0] * time.Second)
	} else {
		//默认超时时间5s
		gClient = gClient.Timeout(5 * time.Second)
	}
	return gClient
}

// mcuApiSystemGetToken
// @summary 获取token
// @param1  mcuApiLevel uint "接口版本"
// @param2  mcu MCUServiceConnParam "对接MCU信息"
// @return1 accountToken string "返回token"
// @return2 err error "返回报错信息"
func mcuApiSystemGetToken(mcuApiLevel uint, mcu MCUServiceConnParam) (accountToken string, err error) {
	resp, err := newMCUPostReq(mcuApiLevel, "").Timeout(1*time.Second).
		Post(fmt.Sprintf("http://%s:%d/api/v1/system/token", mcu.Ip, mcu.Port), g.Map{
			"oauth_consumer_key":    mcu.OAuthConsumerKey,
			"oauth_consumer_secret": mcu.OAuthConsumerSecret,
		})
	if err != nil {
		return
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return "", gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
	}
	return jsonObj.GetString("account_token"), err
}

// mcuApiSystemLogin
// @summary 登录请求
// @param1  mcuApiLevel uint "接口版本"
// @param2  mcu MCUServiceConnParam "对接MCU信息"
// @param3  loginReq mcuSystemLoginReq "登录用户信息"
// @return1 accountToken string "返回用户cookie"
// @return2 err error "返回报错信息"
func mcuApiSystemLogin(mcuApiLevel uint, mcu MCUServiceConnParam, loginReq mcuSystemLoginReq) (cookie string, err error) {
	resp, err := newMCUPostReq(mcuApiLevel, "").Timeout(1*time.Second).
		Post(fmt.Sprintf("http://%s:%d/api/v1/system/login", mcu.Ip, mcu.Port), loginReq)
	if err != nil {
		logger().Error(err)
		return "", err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return "", gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	_, err = respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
		return "", err
	}
	cookie = resp.GetCookie("SSO_COOKIE_KEY")
	return cookie, nil
}

// mcuApiSystemGetVersion
// @summary 获取MCU版本信息
// @param1  mcuApiLevel uint "接口版本"
// @param3  mcuAccountToken string "用户token"
// @param4  mcuUserCookie string "用户cookie"
// @return1 apiLevel uint "返回Api版本"
// @return2 err error "返回报错信息"
func mcuApiSystemGetVersion(mcuApiLevel uint, mcu MCUServiceConnParam, mcuAccountToken, mcuUserCookie string) (apiLevel uint, err error) {
	resp, err := newMCUGetReq(mcuApiLevel, mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/system/version", mcu.Ip, mcu.Port), g.Map{
			"account_token": mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return 0, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
		return
	}
	apiLevel = jsonObj.GetUint("api_level")
	return apiLevel, nil
}

// mcuApiVCGetVersion
// @summary 获取会控API版本信息
// @param1  mcuApiLevel uint "接口版本"
// @param2  mcu MCUServiceConnParam "对接MCU信息"
// @param3  mcuAccountToken string "用户token"
// @param4  mcuUserCookie string "用户cookie"
// @return1 apiLevel uint "返回Api版本"
// @return2 err error "返回报错信息"
func mcuApiVCGetVersion(mcuApiLevel uint, mcu MCUServiceConnParam, mcuAccountToken, mcuUserCookie string) (apiLevel uint, err error) {
	resp, err := newMCUGetReq(mcuApiLevel, mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/vc/version", mcu.Ip, mcu.Port), g.Map{
			"account_token": mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return 0, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
		return
	}
	apiLevel = jsonObj.GetUint("api_level")
	return apiLevel, nil
}

// mcuApiMCGetVersion
// @summary 获取会管API版本信息
// @param1  mcuApiLevel uint "接口版本"
// @param2  mcu MCUServiceConnParam "对接MCU信息"
// @param3  mcuAccountToken string "用户token"
// @param4  mcuUserCookie string "用户cookie"
// @return1 apiLevel uint "返回Api版本"
// @return2 err error "返回报错信息"
func mcuApiMCGetVersion(mcuApiLevel uint, mcu MCUServiceConnParam, mcuAccountToken, mcuUserCookie string) (apiLevel uint, err error) {
	resp, err := newMCUGetReq(mcuApiLevel, mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/mc/version", mcu.Ip, mcu.Port), g.Map{
			"account_token": mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return 0, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
		return
	}
	apiLevel = jsonObj.GetUint("api_level")
	return apiLevel, nil
}

// mcuApiSystemHeartbeat
// @summary 心跳请求,保活token和cookie
// @param1  mcuApiLevel uint "接口版本"
// @param2  mcu MCUServiceConnParam "对接MCU信息"
// @param3  mcuAccountToken string "用户token"
// @param4  mcuUserCookie string "用户cookie"
// @return1 err error "返回报错信息"
func mcuApiSystemHeartbeat(mcuApiLevel uint, mcu MCUServiceConnParam, mcuAccountToken, mcuUserCookie string) error {
	resp, err := newMCUPostReq(mcuApiLevel, mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/system/heartbeat", mcu.Ip, mcu.Port), g.Map{
			"account_token": mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	if _, err = respDecodeToJson(resp.ReadAllString()); err != nil {
		logger().Error(err)
		return err
	}
	return nil
}

// mcuApiAMSGetAccount
// @doc     http://10.8.0.240:808/docs/apiCore/api5_amsapi_accounts_restful
// @summary 根据[moid/e164/jid/account/email/e164/mobile]查询账户信息
//
// @param1  mcuHead mcuApiHead "协议头"
// @param2  query mcuSelectPageReq "分页条件查询"
// @return1 confList []mcuVCConf "返回正在召开的会议列表"
// @return2 err error "返回报错信息"
func mcuApiAMSGetAccount(mcuHead mcuApiHead, accountId string) (account mcuAccount, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api5/ams/accounts", mcuHead.mcu.Ip, mcuHead.mcu.Port), g.Map{
			"account_token": mcuHead.mcuAccountToken,
			"username":      accountId,
		})
	if err != nil {
		logger().Error(err)
		return
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return account, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
		return account, err
	}

	err = jsonObj.GetJson("user").Structs(&account)
	return account, err
}

// mcuApiAMCGetPublicGroups
// @summary 获取公共群组列表
// @return1 err error "返回报错信息"
func mcuApiAMCGetPublicGroups(mcuHead mcuApiHead) (err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/amc/public_groups", mcuHead.mcu.Ip, mcuHead.mcu.Port), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	_, err = respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
		return err
	}
	return
}

// mcuApiNMSGetDomains
// @doc     nms.1
// @summary 获取用户所在域信息
//          kernel(核心域), service(服务域), platform(平台域), user(用户域), machine_room(机房)
// @param1  mcuApiLevel uint "接口版本"
// @param2  mcu MCUServiceConnParam "对接MCU信息"
// @param3  mcuAccountToken string "用户token"
// @param4  mcuUserCookie string "用户cookie"
// @return1 domains mcuVMSDomains "返回用户所属域信息"
// @return2 err error "返回报错信息"
func mcuApiNMSGetDomains(mcuHead mcuApiHead) (domains []mcuVMSDomains, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/nms/domains", mcuHead.mcu.Ip, mcuHead.mcu.Port), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return nil, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	err = jsonObj.GetJson("domains").Struct(&domains)
	return
}

// mcuApiNMSGetDomains
// @doc     nms.4
// @summary 获取用户域下所有会议终端设备
// @param1  mcuHead mcuApiHead "协议头"
// @param2  userDomainMoid string "用户域moid"
// @return1 terminals []mcuTerminalBaseInfo "返回终端设备列表"
// @return2 err error "返回报错信息"
func mcuApiNMSGetTerminalsByUserDomain(mcuHead mcuApiHead, userDomainMoid string) (terminals []mcuTerminalBaseInfo, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/nms/user_domains/%s/terminals", mcuHead.mcu.Ip, mcuHead.mcu.Port, userDomainMoid), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return nil, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	err = jsonObj.GetJson("terminals").Structs(&terminals)
	return
}

// mc1.mcuApiMCGetPersonalTemplates
// @summary 获取个人模板列表
// @param1  mcuHead mcuApiHead "协议头"
// @param2  query mcuSelectPageReq "条件查询"
// @return1 confList []mcuMCPersonalTemplatesSimple "返回个人模板列表"
// @return2 err error "返回报错信息"
func mcuApiMCGetPersonalTemplates(mcuHead mcuApiHead, query mcuSelectPageReq) (confList []mcuMCPersonalTemplatesSimple, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/mc/personal_templates", mcuHead.mcu.Ip, mcuHead.mcu.Port), g.Map{
			"account_token": mcuHead.mcuAccountToken,
			"count":         query.count,
			"order":         query.order,
			"start":         query.start,
		})
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return nil, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	err = jsonObj.GetJson("personal_templates").Structs(&confList)
	return
}

// mc2.mcuApiMCGetOnePersonalTemplates
// @summary 获取个人模板详情
// @param1  mcuHead mcuApiHead "协议头"
// @param2  templateId string "个人会议默认Id"
// @return1 confInfo mcuMCPersonalTemplatesDetail "返回个人模板信息"
// @return2 err error "返回报错信息"
func mcuApiMCGetOnePersonalTemplates(mcuHead mcuApiHead, templateId string) (confInfo mcuMCPersonalTemplatesDetail, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/mc/personal_templates/%s", mcuHead.mcu.Ip, mcuHead.mcu.Port, templateId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return confInfo, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return confInfo, gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return confInfo, err
	}

	err = jsonObj.Structs(&confInfo)

	logger().Notice(err, confInfo)

	return confInfo, err
}

// mc11.mcuApiMCCreateConf
// @summary 创建会议
// @param1  mcuHead mcuApiHead "协议头"
// @param2  createType int "创建会议类型" 1=及时会议，2=公共模板，3=个人模板，4=根据虚拟会议室创建，5=预约会议提前召开，当前只支持根据个人模板召开会议
// @param3  templateId string "会议模板ID"
// @param4  name string "会议名称"
// @param5  duration uint "会议时长，0=永久会议（即手动结束），否则按时长自动结束，单位分钟"
// @return1 confMark mcuMCConfMark "返回会议标识信息"
// @return2 err error "返回报错信息"
func mcuApiMCCreateConf(mcuHead mcuApiHead, createType int, templateId string, name string, duration uint) (confMark mcuMCConfMark, err error) {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"create_type": createType,
		"template_id": templateId,
		"name":        name,
		"duration":    duration,
	}
	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/mc/confs", mcuHead.mcu.Ip, mcuHead.mcu.Port), &req)
	if err != nil {
		logger().Error(err)
		return
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return confMark, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		return
	}
	err = jsonObj.Structs(&confMark)
	return confMark, err
}

// mc12.mcuApiMCReleaseConf
// @summary 结束会议
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议Id"
// @param3  releaseSubConf string "是否结束级联会议:0=否 1=是"
// @return1 err error "返回报错信息"
func mcuApiMCReleaseConf(mcuHead mcuApiHead, confId string, releaseSubConf string) error {
	resp, err := newMCUDelReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Delete(fmt.Sprintf("http://%s:%d/api/v1/mc/confs/%s", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
			"end_subconf":   releaseSubConf,
		})
	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	if _, err = respDecodeToJson(strRespMsg); err != nil && err != io.EOF {
		logger().Error(err)
		return err
	}
	return nil
}

// vc1.mcuApiVCGetConfs
// @summary 获取当前正在召开的视频会议列表
// @param1  mcuHead mcuApiHead "协议头"
// @param2  query mcuSelectPageReq "分页条件查询"
// @return1 confList []mcuVCConf "返回正在召开的会议列表"
// @return2 err error "返回报错信息"
func mcuApiVCGetConfs(mcuHead mcuApiHead, query mcuSelectPageReq) (confList []mcuVCConf, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs", mcuHead.mcu.Ip, mcuHead.mcu.Port), g.Map{
			"account_token": mcuHead.mcuAccountToken,
			"own_conf":      0, //是否仅获取与自己相关的会议
			"start":         query.start,
			"count":         query.count,
		})
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return nil, gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return nil, err
	}

	err = jsonObj.GetJson("confs").Structs(&confList)
	logger().Notice(err, confList)
	for _, temp := range confList {
		logger().Debug("conf id:", temp.ConfId, "conf name:", temp.Name, "meeting id:", temp.MeetingId)
	}
	return confList, err
}

// vc2.mcuApiVCGetOneConf
// @summary 获取视频会议详情
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议ID"
// @return1 confInfo mcuVCConfDetail "返回会议信息"
// @return2 err error "返回报错信息"
func mcuApiVCGetOneConf(mcuHead mcuApiHead, confId string) (confInfo mcuVCConfDetail, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return confInfo, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {

		return
	}
	err = jsonObj.Structs(&confInfo)
	return
}

// vc3.mcuApiVCGetConfCascades
// @summary 获取会议级联信息
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议Id"
// @return1 confList []mcuVCCascadeConf "返回级联会议列表"
// @return2 err error "返回报错信息"
func mcuApiVCGetConfCascades(mcuHead mcuApiHead, confId string) (confList []mcuVCCascadeConf, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/cascades", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return nil, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		return nil, err
	}
	err = jsonObj.Structs(&confList)
	return confList, err
}

// vc4.mcuApiVCGetCascadesMts
// @summary 获取级联会议终端列表
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议Id"
// @param3  cascadeId string "级联会议Id" 0表示本机会议
// @return1 mts []mcuVCMt "返回级联会议终端列表"
// @return2 err error "返回报错信息"
func mcuApiVCGetCascadesMts(mcuHead mcuApiHead, confId string, cascadeId string) (mts []mcuVCMt, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/cascades/%s/mts", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, cascadeId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return nil, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		return nil, err
	}
	err = jsonObj.GetJson("mts").Structs(&mts)
	return mts, err
}

// vc5.mcuApiVCGetCurConfMts
// @summary 获取本级会议终端列表
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议Id"
// @return1 mts []mcuVCMt "返回本级会议终端列表"
// @return2 err error "返回报错信息"
func mcuApiVCGetCurConfMts(mcuHead mcuApiHead, confId string) (mts []mcuVCMt, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mts", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return nil, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		return nil, err
	}
	err = jsonObj.GetJson("mts").Structs(&mts)
	return mts, err
}

// vc6.mcuApiVCGetMt
// @summary 获取与会终端详情
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议id"
// @param3  mtId string "终端id"
// @return1 mtInfo mcuVCMt "返回与会终端信息"
// @return2 err error "返回报错信息"
func mcuApiVCGetMt(mcuHead mcuApiHead, confId string, mtId string) (mtInfo mcuVCMt, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mts/%s", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mtId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return mtInfo, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return mtInfo, gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return mtInfo, err
	}
	err = jsonObj.Structs(&mtInfo)
	logger().Notice(err, mtInfo)
	return mtInfo, err
}

// vc7.mcuApiVCAddCascadeMts
// @summary 批量添加级联终端，添加本级或下级终端，异步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议id"
// @param3  cascadeId string "级联会议id" 0表示本机会议
// @param4  mtList []mcuVCMtSim "批量添加级联会议终端"
// @return1 mts []mcuVCMtSim "返回终端信息"
// @return2 err error "返回报错信息"
func mcuApiVCAddCascadeMts(mcuHead mcuApiHead, confId string, cascadeId string, mtList []mcuVCMtSim) (mts []mcuVCMtSim, err error) {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"mts": mtList,
	}

	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/cascades/%s/mts", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, cascadeId), &req)
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return nil, gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return nil, err
	}

	err = jsonObj.GetJson("mts").Structs(&mts)
	logger().Notice(err, mts)
	return mts, err
}

// vc8.mcuApiVCAddMts
// @summary 批量添加本级终端，添加终端为同步操作，等待终端上线为异步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  mtList []mcuVCMtSim "批量添加本机会议终端"
// @return1 mts []mcuVCMtSim "返回终端信息"
// @return2 err error "返回报错信息"
func mcuApiVCAddMts(mcuHead mcuApiHead, confId string, mtList []mcuVCMtSim) (mts []mcuVCMtSim, err error) {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"mts": mtList,
	}
	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mts", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return nil, gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return nil, err
	}

	err = jsonObj.GetJson("mts").Structs(&mts)
	logger().Notice(err, mts)
	return mts, err
}

// vc9.mcuApiVCDelMts
// @summary 批量添加本级终端，添加终端为同步操作，等待终端上线为异步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议id"
// @param3  mtIds string "批量删除会议终端id列表"
// @return1 err error "返回报错信息"
func mcuApiVCDelMts(mcuHead mcuApiHead, confId string, mtIds []string) error {
	var req mcuPostReq
	var reqBody []mcuVCMtId
	for _, mtIdx := range mtIds {
		var strMtId mcuVCMtId
		strMtId.MtId = mtIdx
		reqBody = append(reqBody, strMtId)
	}
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"mts": reqBody,
	}
	resp, err := newMCUDelReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Delete(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mts", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)
	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}

	return err
}

// vc15.mcuApiVCCallMts
// @summary 批量呼叫终端，异步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议id"
// @param3  mts []mcuVCMtId "批量呼叫会议终端"
// @return1 err error "返回报错信息"
func mcuApiVCCallMts(mcuHead mcuApiHead, confId string, mts []mcuVCMtId) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"mts": mts,
	}
	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/online_mts", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)
	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	_, err = respDecodeToJson(strRespMsg)
	return err
}

// vc16.mcuApiVCDropMts
// @summary 批量挂断终端，同步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议id"
// @param3  mtIds [] string "批量呼叫会议终端Id"
// @return1 err error "返回报错信息"
func mcuApiVCDropMts(mcuHead mcuApiHead, confId string, mtIds []string) error {
	var req mcuPostReq
	var reqBody []mcuVCMtId
	for _, mtIdx := range mtIds {
		var strMtId mcuVCMtId
		strMtId.MtId = mtIdx
		reqBody = append(reqBody, strMtId)
	}
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"mts": reqBody,
	}

	resp, err := newMCUDelReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Delete(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/online_mts", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)
	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc17.mcuApiVCGetChairman
// @summary 获取会议主席
// @param1  mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @return1 mtId string "终端id"
// @return2 err error "返回报错信息"
func mcuApiVCGetChairman(mcuHead mcuApiHead, confId string) (mtId string, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/chairman", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return "", err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return "", gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
		return "", err
	}
	mtId = jsonObj.GetString("mt_id")
	return mtId, err
}

// vc18.mcuApiVCSetChairman
// @summary 指定会议主席，同步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  mtId string "指定终端为主席参数"
// @return1 err error "返回报错信息"
func mcuApiVCSetChairman(mcuHead mcuApiHead, confId string, mtId string) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{"mt_id": mtId}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/chairman", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc19.mcuApiVCGetSpeaker
// @summary 获取会议发言人
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议id"
// @return1 mtId mcuVCMtSet "终端id"
// @return2 err error "返回报错信息"
func mcuApiVCGetSpeaker(mcuHead mcuApiHead, confId string) (mtId mcuVCMtSet, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/speaker", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return mtId, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return mtId, gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return mtId, err
	}
	err = jsonObj.Structs(&mtId)
	logger().Notice(err, mtId)
	return mtId, err
}

// vc20.mcuApiVCSetSpeaker
// @summary 指定会议发言人，同步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  mtId string "指定终端为发言人"
// @param5  broadcast int "是否设置发言人强制广播" 0=否，1=是
// @return1 err error "返回报错信息"
func mcuApiVCSetSpeaker(mcuHead mcuApiHead, confId string, mtId string, broadcast int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"mt_id":           mtId,
		"force_broadcast": broadcast,
	}

	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/speaker", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc21.mcuApiVCGetDualStream
// @summary 获取会议双流源
// @param1  mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @return1 mtId string "终端id"
// @return2 err error "返回报错信息"
func mcuApiVCGetDualStream(mcuHead mcuApiHead, confId string) (mtId string, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/dualstream", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return "", err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return "", gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return "", err
	}
	mtId = jsonObj.GetString("success")
	logger().Notice(err, mtId)
	return mtId, err
}

// vc22.mcuApiVCSetDualStream
// @summary 指定会议双流源，异步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  mtId string "指定会议双流源"
// @return1 err error "返回报错信息"
func mcuApiVCSetDualStream(mcuHead mcuApiHead, confId string, mtId string) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"mt_id": mtId,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/dualstream", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc23.mcuApiVCDelDualStream
// @summary 取消会议双流源，异步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  mtId string "取消会议双流源"
// @return1 err error "返回报错信息"
func mcuApiVCDelDualStream(mcuHead mcuApiHead, confId string, mtId string) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"mt_id": mtId,
	}
	resp, err := newMCUDelReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Delete(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/dualstream", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc24.mcuApiVCSetDelay
// @summary 延长会议时间，同步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  timeDelay uint "延长的时间"
// @return1 err error "返回报错信息"
func mcuApiVCSetDelay(mcuHead mcuApiHead, confId string, timeDelay uint) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"delay_time": timeDelay,
	}
	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/delay", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc25.mcuApiVCSetSilence
// @summary 会场静音操作，同步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  value int "静音状态" 0=停止静音，1=静音
// @return1 err error "返回报错信息"
func mcuApiVCSetSilence(mcuHead mcuApiHead, confId string, value int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"value": value,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/silence", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc26.mcuApiVCSetMute
// @summary 会场哑音操作，同步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议id"
// @param3  state int "哑音状态" 0=停止哑音，1=哑音
// @param4  forceMute int "全场哑音下是否禁止终端取消自身哑音"
// @return1 err error "返回报错信息"
func mcuApiVCSetMute(mcuHead mcuApiHead, confId string, state int, forceMute int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"value":      state,
		"force_mute": forceMute,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mute", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc28.mcuApiVCSetName
// @summary 修改会议名称，同步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  name string "会议名称" 最大128字符
// @return1 err error "返回报错信息"
func mcuApiVCSetName(mcuHead mcuApiHead, confId string, name string) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"value": name,
	}
	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/name", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc29.mcuApiVCSms
// @summary 发送短消息，SFU会议不支持
// @param1  mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  sms mcuVCSms "短消息参数"
// @return1 err error "返回报错信息"
func mcuApiVCSms(mcuHead mcuApiHead, confId string, sms mcuVCSms) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = sms
	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/sms", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc30.mcuApiVCSetForceBroadcast
// @summary 设置是否强制广播
// @param1  mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  mode int "是否强制广播" 0=不强制广播，1=强制广播
// @return1 err error "返回报错信息"
func mcuApiVCSetForceBroadcast(mcuHead mcuApiHead, confId string, mode int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"force_broadcast": mode,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/forcebroadcast", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc31.mcuApiVCSetVadState
// @summary 开启/关闭语言激励
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议id"
// @param3  state int "语音激励" 0=关闭，1=开启
// @param4  interval uint "语音激励敏感度" 单位s，最小值3s
// @return1 err error "返回报错信息"
func mcuApiVCSetVadState(mcuHead mcuApiHead, confId string, state int, interval uint) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"state":       state,
		"vacinterval": interval,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/vad", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)
	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}

	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		return err
	}
	return err
}

// vc32.mcuApiVCGetVadState
// @summary 开启/关闭语言激励
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议id"
// @return1 state mcuVCVac "语音激励状态"
// @return2 err error "返回报错信息"
func mcuApiVCGetVadState(mcuHead mcuApiHead, confId string) (state mcuVCVac, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/vad", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return state, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return state, gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		return state, err
	}
	state.VoiceActivityDetection = jsonObj.GetInt("voice_activity_detection")
	state.VacInterval = jsonObj.GetInt("vacinterval")
	return state, err
}

// vc33.mcuApiVCGetMixs
// @summary  获取会议混音信息
// @param1   mcuHead mcuApiHead "协议头"
// @param2   confId string "会议id"
// @return1  mix mcuMCMix "混音参数"
// @return2  err error "返回报错信息"
func mcuApiVCGetMixs(mcuHead mcuApiHead, confId string) (mix mcuMCMix, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mixs/1", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return mix, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return mix, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	//解析消息
	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		return mix, err
	}
	err = jsonObj.Structs(&mix)
	return mix, err
}

// vc34.mcuApiVCStartMixs
// @summary 开启会议混音信息，异步操作
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  mode int "混音模式" 1=智能混音，2=定制混音
// @param6  mtIds []string "混音成员列表"
// @return1 err error "返回报错信息"
func mcuApiVCStartMixs(mcuHead mcuApiHead, confId string, mode int, mtIds []string) error {
	var req mcuPostReq
	var mixs mcuMCMix
	mixs.Mode = mode
	for _, mtIdx := range mtIds {
		var tMtId mcuMCMt
		tMtId.MtId = mtIdx
		mixs.Members = append(mixs.Members, tMtId)
	}
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = mixs
	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mixs", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc35.mcuApiVCStopMixs
// @summary 停止会议混音器，异步操作
// @param1  mcuHead mcuApiHead "协议头"
// @param2  confId string "会议id"
// @param3  mixId string "混音器ID，默认mix_id为1"
// @return1 err error "返回报错信息"
func mcuApiVCStopMixs(mcuHead mcuApiHead, confId, mixId string) error {
	resp, err := newMCUDelReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Delete(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mixs/%s", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mixId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc36.mcuApiVCAddMixMts
// @summary 添加混音成员
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  mtIds []string "混音成员"
// @return1 err error "返回报错信息"
func mcuApiVCAddMixMts(mcuHead mcuApiHead, confId, mixId string, mtIds []string) error {
	var req mcuPostReq
	var members mcuVCMtMember
	for _, mtIdx := range mtIds {
		var tMtId mcuVCMtId
		tMtId.MtId = mtIdx
		members.Members = append(members.Members, tMtId)
	}
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = members

	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mixs/%s/members", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mixId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc37.mcuApiVCDelMixMts
// @summary 删除混音成员
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  mtIds []string "混音成员"
// @return1 err error "返回报错信息"
func mcuApiVCDelMixMts(mcuHead mcuApiHead, confId, mixId string, mtIds []string) error {
	var req mcuPostReq
	var members mcuVCMtMember
	for _, mtIdx := range mtIds {
		var tMtId mcuVCMtId
		tMtId.MtId = mtIdx
		members.Members = append(members.Members, tMtId)
	}
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = members

	resp, err := newMCUDelReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Delete(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mixs/%s/members", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mixId), &req)
	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc38.mcuApiVCGetVmps
// @summary 获取会议画面合成信息，无画面合成时返回错误
// @param1  mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @return2 vmp mcuMCVmp "混音参数"
// @return1 err error "返回报错信息"
func mcuApiVCGetVmps(mcuHead mcuApiHead, confId string) (vmp mcuMCVmp, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/vmps/1", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return vmp, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return vmp, gerror.New(resp.Status)
	}
	strRespMsg := resp.ReadAllString()
	logger().Debug(strRespMsg)

	jsonObj, err := respDecodeToJson(strRespMsg)
	if err != nil {
		logger().Error(err)
		return vmp, err
	}
	_ = jsonObj.Structs(&vmp)
	return vmp, err
}

// vc39.mcuApiVCStartVmps
// @summary 开启会议画面合成，异步操作
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  vmp mcuMCVmp "画面合成参数"
// @return1 err error "返回报错信息"
func mcuApiVCStartVmps(mcuHead mcuApiHead, confId string, vmp mcuMCVmp) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = vmp
	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/vmps", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc40.mcuApiVCModVmps
// @summary 修改会议画面合成，支持模式、开关、风格、广播、台标、成员、合成轮询等参数修改
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  vmp mcuMCVmp "画面合成参数"
// @return1 err error "返回报错信息"
func mcuApiVCModVmps(mcuHead mcuApiHead, confId string, vmp mcuMCVmp) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = vmp

	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/vmps/1", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc41.mcuApiVCStopVmps
// @summary 停止会议主混音器，默认mix_id为1，异步操作
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @return1 err error "返回报错信息"
func mcuApiVCStopVmps(mcuHead mcuApiHead, confId string) error {
	resp, err := newMCUDelReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Delete(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/vmps/1", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}

	return err
}

// vc42.mcuApiVCGetChairmanPoll
// @summary 获取主席轮询信息
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @return2  poll mcuMCPoll "轮询参数"
// @return1 err error "返回报错信息"
func mcuApiVCGetChairmanPoll(mcuHead mcuApiHead, confId string) (poll mcuMCPoll, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/chairman_poll", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})

	if err != nil {
		logger().Error(err)
		return poll, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return poll, gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return poll, err
	}
	jsonObj.Structs(&poll)
	return poll, err
}

// vc43.mcuApiVCStartChairmanPoll
// @summary 开启主席轮询，异步操作，SFU会议不支持该功能
// @param1  mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  poll mcuMCPoll "主席轮询参数"
// @return1 err error "返回报错信息"
func mcuApiVCStartChairmanPoll(mcuHead mcuApiHead, confId string, poll mcuMCPoll) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = poll
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/chairman_poll", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc44.mcuApiVCSetChairmanPoll
// @summary 修改主席轮询状态和模式
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  state int "轮询状态" 详见轮询状态枚举
// @param7  mode int "轮询模式" 详见轮询模式枚举，支持主席视频轮询和主席音视频轮询模式切换
// @return1 err error "返回报错信息"
func mcuApiVCSetChairmanPoll(mcuHead mcuApiHead, confId string, state int, mode int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"value": state,
		"mode":  mode,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/chairman_poll/state", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc45.mcuApiVCGetConfPoll
// @summary 获取会议轮询信息
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @return2  poll mcuMCPoll "轮询参数"
// @return1 err error "返回报错信息"
func mcuApiVCGetConfPoll(mcuHead mcuApiHead, confId string) (poll mcuMCPoll, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/poll", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})

	if err != nil {
		logger().Error(err)
		return poll, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return poll, gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return poll, err
	}
	jsonObj.Structs(&poll)
	return poll, err
}

// vc46.mcuApiVCStartConfPoll
// @summary 开启会议轮询，异步操作，SFU会议不支持该功能
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  poll mcuMCPoll "会议轮询参数"
// @return1 err error "返回报错信息"
func mcuApiVCStartConfPoll(mcuHead mcuApiHead, confId string, poll mcuMCPoll) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = poll

	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/poll", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc47.mcuApiVCSetConfPoll
// @summary 修改主席轮询状态和模式
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  state int "轮询状态" 详见轮询状态枚举
// @param7  mode int "轮询模式" 详见轮询模式枚举，支持视频轮询和音视频轮询模式切换
// @return1 err error "返回报错信息"
func mcuApiVCSetConfPoll(mcuHead mcuApiHead, confId string, state int, mode int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"value": state,
		"mode":  mode,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/poll/state", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc81.mcuApiVCGetInspections
// @summary 获取终端选看列表，终端选看终端或非广播的画面合成
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @return2 mts []mcuVCInspec "选看参数"
// @return1 err error "返回报错信息"
func mcuApiVCGetInspections(mcuHead mcuApiHead, confId string) (mts []mcuVCInspec, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/inspections", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})

	if err != nil {
		logger().Error(err)
		return nil, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return nil, gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return nil, err
	}
	jsonObj.GetJson("inspections").Structs(&mts)

	logger().Notice(err, mts)
	for _, temp := range mts {
		logger().Debug("mode:", temp.Mode, "srd:", temp.Src.MtId, "dst:", temp.Dst.MtId)
	}
	return mts, err
}

// vc82.mcuApiVCStartInspect
// @summary 开启选看，异步操作，传统会议，终端选看终端或画面合成，SFU会议不支持
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  inspect mcuVCInspec "选看参数"
// @return1 err error "返回报错信息"
func mcuApiVCStartInspect(mcuHead mcuApiHead, confId string, inspect mcuVCInspec) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = inspect
	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/inspections", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc83.mcuApiVCStopInspect
// @summary 取消终端选看，异步操作
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param5  mtId string "被选看终端id"
// @param5  mode int "选看模式" 1=视频 2=音频
// @return1 err error "返回报错信息"
func mcuApiVCStopInspect(mcuHead mcuApiHead, confId string, mtId string, mode int) error {
	resp, err := newMCUDelReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Delete(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/inspections/%s/%d", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mtId, mode), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc87.mcuApiVCMtSetCallMode
// @summary 修改终端呼叫模式
// @param1   mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  mtId string "终端id"
// @param6  mode int "呼叫模式" 0=手动，2=自动，3=追呼
// @return1 err error "返回报错信息"
func mcuApiVCMtSetCallMode(mcuHead mcuApiHead, confId string, mtId string, mode int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"value": mode,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mts/%s/call_mode", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mtId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc88.mcuApiVCMtSetSilence
// @summary 终端静音操作，同步操作
// @param1   mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  mtId string "终端id"
// @param6  silence int "静音状态" 0=停止静音，1=静音
// @return1 err error "返回报错信息"
func mcuApiVCMtSetSilence(mcuHead mcuApiHead, confId string, mtId string, silence int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"value": silence,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mts/%s/silence", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mtId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc89.mcuApiVCMtSetMute
// @summary 终端哑音操作，同步操作
// @param1   mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  mtId string "终端id"
// @param6  mute int "哑音状态" 0=停止哑音，1=哑音
// @return1 err error "返回报错信息"
func mcuApiVCMtSetMute(mcuHead mcuApiHead, confId string, mtId string, mute int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"value": mute,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mts/%s/mute", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mtId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc90.mcuApiVCMtSetVolume
// @summary 修改终端音量，同步操作
// @param1   mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  mtId string "终端id"
// @param6  mode int "终端音量设备" 1=扬声器，2=麦克风
// @param6  volume int "音量" 0-35
// @return1 err error "返回报错信息"
func mcuApiVCMtSetVolume(mcuHead mcuApiHead, confId string, mtId string, mode int, volume int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"vol_mode":  mode,
		"vol_value": volume,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mts/%s/volume", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mtId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc91.mcuApiVCMtSetPtz
// @summary 控制终端摄像头位移、焦距、视野、亮度
// @param1   mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  mtId string "终端id"
// @param6  state int "摄像头状态"
// @param6  ptz int "控制类型" 详见云台控制类型枚举
// @return1 err error "返回报错信息"
func mcuApiVCMtSetPtz(mcuHead mcuApiHead, confId string, mtId string, state int, ptz int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"state": state,
		"type":  ptz,
	}

	resp, err := newMCUPostReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Post(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mts/%s/camera", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mtId), &req)

	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}
	return err
}

// vc92.mcuApiVCMtGetVideos
// @summary 获取终端视频源信息，rtc协议注册的终端不支持
// @param1   mcuHead mcuApiHead "协议头"
// @param5  confId string "会议id"
// @param6  mtId string "终端id"
// @return1 mtVideos mcuVCMtVideos "返回终端当前通道视频源信息"
// @return1 err error "返回报错信息"
func mcuApiVCMtGetVideos(mcuHead mcuApiHead, confId string, mtId string) (mtVideos mcuVCMtVideos, err error) {
	resp, err := newMCUGetReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Get(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mts/%s/videos", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mtId), g.Map{
			"account_token": mcuHead.mcuAccountToken,
		})
	if err != nil {
		logger().Error(err)
		return mtVideos, err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return mtVideos, gerror.New(resp.Status)
	}
	jsonObj, err := respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return mtVideos, err
	}
	jsonObj.Structs(&mtVideos)
	logger().Notice(err, mtId)
	return mtVideos, err
}

// vc93.mcuApiVCMtSetVideos
// @summary 设置终端当前视频源，SFU会议不支持此功能
// @param1   mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  mtId string "终端id"
// @param6 videoIdx int "终端接收视频源的通道号"	只填通道号
// @return1 err error "返回报错信息"
func mcuApiVCMtSetVideos(mcuHead mcuApiHead, confId string, mtId string, videoIdx int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{
		"video_idex": videoIdx,
	}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/mts/%s/videos", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId, mtId), &req)
	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}

	logger().Notice(err)
	return err
}

// vc98.mcuApiVCMtSetDualMode
// @summary 设置双流模式，设置双流共享模式为任意会场或发言会场
// @param1   mcuHead mcuApiHead "协议头"
// @param4  confId string "会议id"
// @param5  mode int "双流开启模式" 0=任意会场，1=任意会场
// @return1 err error "返回报错信息"
func mcuApiVCMtSetDualMode(mcuHead mcuApiHead, confId string, mode int) error {
	var req mcuPostReq
	req.AccountToken = mcuHead.mcuAccountToken
	req.Params = g.Map{"value": mode}
	resp, err := newMCUPutReq(mcuHead.mcuApiLevel, mcuHead.mcuUserCookie).
		Put(fmt.Sprintf("http://%s:%d/api/v1/vc/confs/%s/dualmode", mcuHead.mcu.Ip, mcuHead.mcu.Port, confId), &req)
	if err != nil {
		logger().Error(err)
		return err
	}
	defer func() { _ = resp.Close() }()
	if resp.StatusCode != 200 {
		logger().Error(resp.Status)
		return gerror.New(resp.Status)
	}
	_, err = respDecodeToJson(resp.ReadAllString())
	if err != nil {
		logger().Error(err)
		return err
	}

	logger().Notice(err)
	return err
}

//mcuNotifyAssistant
//@summary 解析MCU通知消息
func mcuNotifyAssistant(url string, data string) (bOk bool, msg MCUNotifyMessage) {
	//MCU推送URL格式: /userdomains/{domain_id}/confs/{conf_id}/**
	strList := strings.Split(url, "/")
	if len(strList) == 0 || len(strList) < 4 {
		return false, msg
	}
	//解析会议ID
	for index, eleStr := range strList {
		if eleStr == "confs" && index+1 < len(strList) {
			msg.ConfId = strList[index+1] //会议ID
			break
		}
	}

	//解析消息体
	jsonObj, err := gjson.DecodeToJson(data)
	if err != nil {
		logger().Error(err)
		return
	}
	msg.Method = jsonObj.GetString("method")

	//URL后三个子串
	strUrlSuffix := strList[len(strList)-1]
	strUrlPenult := strList[len(strList)-2]
	strUrlAntepenult := strList[len(strList)-3]

	//判断URL最后一个子串
	if parseMCUNotifyByTopic(strUrlSuffix, strList, jsonObj, &msg) {
		return true, msg
	}
	//判断URL倒数第二个子串
	if parseMCUNotifyByTopic(strUrlPenult, strList, jsonObj, &msg) {
		return true, msg
	}
	//判断URL倒数第三个子串
	if parseMCUNotifyByTopic(strUrlAntepenult, strList, jsonObj, &msg) {
		return true, msg
	}
	return false, msg
}

//解析MCU会议更新通知消息
func parseMCUNotifyByTopic(strUrlSuffix string, strUrlParamList []string, jsonObj *gjson.Json, msg *MCUNotifyMessage) bool {
	switch EMeetingUpdateTopic(strUrlSuffix) {
	case NotifyConfs:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		return true

	case NotifyCascades:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		return true

	case NotifySpeaker:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		return true

	case NotifyChairman:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		return true

	case NotifyDualStream:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		return true

	case NotifyConfPoll:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		return true

	case NotifyChairmanPoll:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		return true

	case NotifyVad:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		return true

	case NotifyMixs:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)
		if msg.Method == "notify" {
			//混音失败通知
			var mixInfo MCUNotifyMixs
			mixInfo.MixId = strUrlParamList[len(strUrlParamList)-1]
			mixInfo.ErrorCode = jsonObj.GetInt("error_code")
			mixInfo.Mode = jsonObj.GetInt("mode")

			msg.Detail = mixInfo
		} else {
			//开启混音通知
			var mixId string
			mixId = strUrlParamList[len(strUrlParamList)-1]

			msg.Detail = mixId
		}
		return true

	case NotifyVmps:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		if msg.Method == "notify" {
			//画面合成失败通知
			var vmpInfo MCUNotifyVmps
			vmpInfo.VmpId = strUrlParamList[len(strUrlParamList)-1]
			vmpInfo.ErrorCode = jsonObj.GetInt("error_code")

			msg.Detail = vmpInfo
		} else {
			//开启画面合成通知
			var vmpId string
			vmpId = strUrlParamList[len(strUrlParamList)-1]

			msg.Detail = vmpId
		}

		return true

	case NotifyMts:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)
		if msg.Method == "notify" {
			// /userdomains/{domain_id}/confs/{conf_id}/mts

			//添加终端失败通知
			var addMtFail MCUNotifyAddMts
			addMtFail.ErrorCode = jsonObj.GetInt("error_code")
			addMtFail.OccupyConfName = jsonObj.GetString("occupy_confname")
			addMtFail.MtAccount = jsonObj.GetJson("mt").GetString("account")
			addMtFail.MtAccountType = jsonObj.GetJson("mt").GetString("account_type")

			msg.Detail = addMtFail
		} else {
			// /userdomains/{domain_id}/confs/{conf_id}/cascades/{cascade_id}/mts/{mt_id}
			//[添加|删除]会议终端
			var mtInfo MCUUpdateMts
			mtInfo.MtId = strUrlParamList[len(strUrlParamList)-1]
			if strUrlParamList[len(strUrlParamList)-4] == "cascades" {
				mtInfo.CascadeId = strUrlParamList[len(strUrlParamList)-3]
			}
			msg.Detail = mtInfo
		}
		return true

	case NotifyMtVmps:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		var vmpId string
		vmpId = strUrlParamList[len(strUrlParamList)-1]

		msg.Detail = vmpId

		return true

	case NotifyMtChnVmps: //终端单通道自主画面合成
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		var mtInfo MCUUpdateMtChnVmps
		mtInfo.MtId = strUrlParamList[len(strUrlParamList)-2]
		mtInfo.MtChnIdx = strUrlParamList[len(strUrlParamList)-1]

		msg.Detail = mtInfo

		return true

	case NotifyMtInspections:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)
		if msg.Method == "notify" {
			var specInfo MCUNotifyInspect
			specInfo.Dst.MtId = strUrlParamList[len(strUrlParamList)-2]
			specInfo.Mode, _ = strconv.Atoi(strUrlParamList[len(strUrlParamList)-1])
			specInfo.ErrorCode = jsonObj.GetInt("error_code")
			specInfo.Src.MtId = jsonObj.GetJson("src").GetString("mt_id")
			specInfo.Src.Type = jsonObj.GetJson("src").GetInt("type")
			msg.Detail = specInfo
		} else {
			var mtInfo MCUUpdateMtInspect
			mtInfo.MtId = strUrlParamList[len(strUrlParamList)-2]
			mtInfo.Mode, _ = strconv.Atoi(strUrlParamList[len(strUrlParamList)-1])
			msg.Detail = mtInfo
		}

		return true

	case NotifyMtMix:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		var mtId string
		mtId = jsonObj.GetString("mt_id")
		msg.Detail = mtId

		return true

	case NotifyMtSpeaker:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		var mtId string
		mtId = jsonObj.GetString("mt_id")
		msg.Detail = mtId

		return true

	case NotifyMtChairman:
		msg.Type = EMeetingUpdateTopic(strUrlSuffix)

		var mtId string
		mtId = jsonObj.GetString("mt_id")
		msg.Detail = mtId
		return true
	}
	return false
}
