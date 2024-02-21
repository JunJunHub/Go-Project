package accessMCUAdapter

import "C"
import (
	"context"
	"fmt"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"gowebserver/app/common/utils/valueCheck"
	"time"
)

//McuAdapterCore MCU适配器核心业务基础结构
type McuAdapterCore struct {
	adapter      AccessMCUAdapter
	ctx          context.Context
	mcuKey       string
	mcuConnParam *MCUServiceConnParam
	mcuApiLevel  uint

	//用户登录缓存
	//key   = string        用户名
	//value = MCULoginCache 登录缓存数据
	mcuUserLoginCaches *gmap.StrAnyMap

	//订阅会议使用的MCU账户信息
	subByMCUUser *MCULoginCache

	bDestroy bool //实例销毁标识
}

//DebugPrintState 调试打印MCU对接状态
func (c *McuAdapterCore) DebugPrintState() {
	logger().Notice("MCUConnParam:", c.mcuConnParam)
	for _, user := range c.mcuConnParam.MCUUsers {
		loginCache := c.mcuUserLoginCaches.Get(user.Username)
		if loginCache != nil {
			logger().Notice("MCUUserLoginCache:", loginCache)
		} else {
			logger().Notice("MCU账户未登录:", user.Username, user.Password)
		}
	}
	logger().Notice("订阅会议使用账户：", c.subByMCUUser)
}

//GetMCUUserConnState 查MCU账户连接状态
func (c *McuAdapterCore) GetMCUUserConnState(username string) (connState EMCUConnStateCode) {
	loginCache := c.mcuUserLoginCaches.Get(username)
	if loginCache != nil {
		return loginCache.(*MCULoginCache).ConnState
	}
	return EMCUNetConnErr
}

//mcuConnErrCB mcu连接报错回传
func (c *McuAdapterCore) mcuConnErrCB(err MCUConnState) {
	msg := MCUNotifyMessage{
		Type:   NotifyMCUConnState,
		Mcu:    c.mcuConnParam,
		Detail: err,
	}
	chnMCUNotifyMsg <- msg
}

//connect
//@summary 连接MCU
//         1、用户自动登录MCU
//         2、创建订阅通知协程
//         3、创建登录保活协程
func (c *McuAdapterCore) connect() {
	//用户登录MCU
	if len(c.mcuConnParam.MCUUsers) == 0 {
		c.mcuConnErrCB(MCUConnState{LastErrCode: EMCUUserUnSet, McuUsername: ""})
	}
	for _, user := range c.mcuConnParam.MCUUsers {
		userLoginCache, err := c.getUserLoginCache(user.Username)
		if err != nil {
			logger().Error(err)
			if err.Error() == fmt.Sprintf("%d", EMCUNetConnErr) {
				break //MCU无法连接直接break
			}
			continue
		}
		logger().Notice(userLoginCache)
	}

	//启动用户登录保活协程
	go c.adapter.keepUserAlive()

	//启动订阅协程
	go c.adapter.subscribeMeeting()
}

//disconnect
//@summary 断开连接
//         1、登出MCU
//         2、停止保活协程
//         3、停止订阅协程
func (c *McuAdapterCore) disconnect() {
	//todo
	c.bDestroy = true
}

//login
//@summary 登录MCU
func (c *McuAdapterCore) login(user *MCUUserInfo) (userLoginCache *MCULoginCache, err error) {
	userLoginCache = &MCULoginCache{
		User:         user,
		AccountToken: "",
		UserCookie:   "",
		ConnState:    EMCUNetConnErr,
	}
	//获取Token
	userLoginCache.AccountToken, err = mcuApiSystemGetToken(c.mcuApiLevel, *c.mcuConnParam)
	if err != nil {
		//解析报错原因
		if err.Error() == "10001" {
			userLoginCache.ConnState = EMCUSecretInvalid
			c.mcuConnErrCB(MCUConnState{LastErrCode: EMCUSecretInvalid, McuUsername: ""})
		} else {
			userLoginCache.ConnState = EMCUNetConnErr
			c.mcuConnErrCB(MCUConnState{LastErrCode: EMCUNetConnErr, McuUsername: ""})
			return nil, gerror.New(fmt.Sprintf("%d", EMCUNetConnErr))
		}
		return nil, err
	}
	c.mcuConnErrCB(MCUConnState{LastErrCode: EMCUConnOk, McuUsername: ""})
	//登录
	userLoginCache.UserCookie, err = mcuApiSystemLogin(c.mcuApiLevel, *c.mcuConnParam, mcuSystemLoginReq{
		AccountToken: userLoginCache.AccountToken,
		Username:     user.Username,
		Password:     user.Password,
	})
	if err != nil {
		userLoginCache.ConnState = EMCUUserInvalid
		c.mcuConnErrCB(MCUConnState{LastErrCode: EMCUUserInvalid, McuUsername: user.Username})
		return nil, err
	}
	//获取MCU版本信息
	c.mcuApiLevel, err = mcuApiSystemGetVersion(c.mcuApiLevel, *c.mcuConnParam, userLoginCache.AccountToken, userLoginCache.UserCookie)
	if err != nil {
		userLoginCache.ConnState = EMCUVerErr
		c.mcuConnErrCB(MCUConnState{LastErrCode: EMCUVerErr, McuUsername: ""})
		return nil, err
	}
	//获取域信息(只有MCU管理员账户才有此接口权限).	//TODO Tips:登录接口会返回该账户所属用户域ID
	var domains []mcuVMSDomains
	domains, err = mcuApiNMSGetDomains(mcuApiHead{
		mcuApiLevel:     c.mcuApiLevel,
		mcu:             c.mcuConnParam,
		mcuAccountToken: userLoginCache.AccountToken,
		mcuUserCookie:   userLoginCache.UserCookie,
	})
	if err != nil {
		logger().Warning(err)
		err = nil //登录流程此接口报错不认为登录失败,报错置空
	} else {
		for _, v := range domains {
			if v.Type == "user" {
				userLoginCache.User.UserDomain = v.Moid
			}
		}
	}
	userLoginCache.ConnState = EMCUConnOk
	c.mcuConnErrCB(MCUConnState{LastErrCode: EMCUConnOk, McuUsername: user.Username})
	return
}

//heartbeat
//@summary 心跳请求
//         保活token和cookie
func (c *McuAdapterCore) heartbeat(user *MCUUserInfo) error {
	userLoginCache, err := c.getUserLoginCache(user.Username)
	if err != nil {
		return err
	}
	if err = mcuApiSystemHeartbeat(c.mcuApiLevel, *c.mcuConnParam, userLoginCache.AccountToken, userLoginCache.UserCookie); err != nil {
		//用户登录缓存失效(cookie|token)
		c.mcuUserLoginCaches.Remove(user.Username)

		//尝试重新登录
		_, err = c.getUserLoginCache(user.Username)
	}
	return err
}

//keepUserAlive
//@summary 保活
func (c *McuAdapterCore) keepUserAlive() {
	for {
		//科达MCU V5.0 以上
		time.Sleep(20 * time.Second)
		if c.bDestroy {
			return
		}
		for _, value := range c.mcuConnParam.MCUUsers {
			user := value
			go func() { _ = c.adapter.heartbeat(&user) }()
		}
	}
}

//subscribeMeeting
//@summary 订阅所有会议更新通知
//         1、订阅链路保活
//         2、收到的MCU通知消息存放 -> chnMCUNotifyMsg -> 上层业务
//
//         MCU版本不同,订阅方式不同.各版本自己重写该订阅接口
//         科达 MCU5.0↑ 要根据用户域来订阅某个域下的会议更新通知，暂时仅支持添加的用户都属于同一个域
func (c *McuAdapterCore) subscribeMeeting() {
	c.adapter.subscribeMeeting()
}

//getUserInfo
//@summary 用户信息
func (c *McuAdapterCore) getUserInfo(username string) (userinfo *MCUUserInfo, err error) {
	for _, user := range c.mcuConnParam.MCUUsers {
		if user.Username == username {
			userinfo = &user
			return
		}
	}
	return nil, gerror.New("会议模板所属MCU用户已删除")
}

//getUserLoginCache
//@summary 获取用户登录缓存信息
func (c *McuAdapterCore) getUserLoginCache(username string) (userLoginCache *MCULoginCache, err error) {
	v := c.mcuUserLoginCaches.GetOrSetFuncLock(username, func() interface{} {
		//校验用户信息
		var userinfo *MCUUserInfo
		userinfo, err = c.getUserInfo(username)
		if userinfo == nil {
			return nil
		}
		//登录MCU
		userLoginCache, err = c.adapter.login(userinfo)
		if err != nil {
			return nil
		}
		return userLoginCache
	})
	if v != nil {
		return v.(*MCULoginCache), err
	}
	return
}

//delUserLoginCache
//@summary 删除用户登录缓存信息
func (c *McuAdapterCore) delUserLoginCache(username string) {
	delValue := c.mcuUserLoginCaches.Remove(username)
	if delValue != nil {
		logger().Noticef("MCU账户[%s@%s]登录失效,删除登录缓存信息:%v", c.mcuKey, username, delValue)
	}
}

//checkNeedDelLoginCacheByMCUErr
//@summary 校验MCU接口返回错误码,判断是否需要清除登录缓存
func (c *McuAdapterCore) checkNeedDelLoginCacheByMCUErr(err error) bool {
	if err.Error() == "10001" || err.Error() == "10002" || err.Error() == "10102" {
		return true
	}
	return false
}

//MCUApiSystemGetVersion
//@summary 获取MCU版本信息
//@return1 apiLevel uint "返回Api版本"
//@return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiSystemGetVersion() (apiLevel uint, err error) {
	//找到一个在线用户
	var userCache *MCULoginCache
	for _, user := range c.mcuConnParam.MCUUsers {
		userCache, err = c.getUserLoginCache(user.Username)
		if err != nil && userCache != nil {
			break
		}
	}
	if userCache == nil {
		return 0, gerror.New("login failed!")
	}

	//获取MCU版本信息
	c.mcuApiLevel, err = mcuApiSystemGetVersion(c.mcuApiLevel, *c.mcuConnParam, userCache.AccountToken, userCache.UserCookie)
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(userCache.User.Username)
	}
	return c.mcuApiLevel, err
}

//MCUApiAMSGetAccount
//@summary 获取MCU账户信息
//         科达MCU内部接口：http://10.8.0.240:808/docs/apiCore/api5_amsapi_accounts_restful
//@param1  username string "MCU登录账户"
//@param2  accountId string "要查询的账户ID，可以是moid/jid/account/email/e164/mobile"
func (c *McuAdapterCore) MCUApiAMSGetAccount(username string, accountId string) (account MCUAccount, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpAccount mcuAccount
		tmpAccount, err = mcuApiAMSGetAccount(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, accountId)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		account.mcuAccount = tmpAccount
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiNMSGetUserDomains
//@summary 获取用户所属域信息
//         kernel(核心域), service(服务域), platform(平台域), user(用户域), machine_room(机房)
func (c *McuAdapterCore) MCUApiNMSGetUserDomains(username string) (domains []MCUVMSDomains, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpDomains []mcuVMSDomains
		tmpDomains, err = mcuApiNMSGetDomains(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		})
		valueCheck.ErrIsNil(context.TODO(), err)

		for _, v := range tmpDomains {
			domains = append(domains, MCUVMSDomains{
				mcuVMSDomains: v,
			})
		}
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiNMSGetAllTerminalsByUserDomain
//@summary 获取指定用户域下的所有会议终端
func (c *McuAdapterCore) MCUApiNMSGetAllTerminalsByUserDomain(username, userDomainMoid string) (meetTerminals []MCUMeetTerminal, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		//查询指定用户域下所有会议终端,分页查询查完为止
		selectPageReq := mcuSelectPageReq{
			count: 100,
			start: 0,
			order: 0,
		}
		for {
			var tmpTerminals []mcuTerminalBaseInfo
			tmpTerminals, err = mcuApiNMSGetTerminalsByUserDomain(mcuApiHead{
				mcuApiLevel:     c.mcuApiLevel,
				mcu:             c.mcuConnParam,
				mcuAccountToken: userLoginCache.AccountToken,
				mcuUserCookie:   userLoginCache.UserCookie,
			}, userDomainMoid)
			valueCheck.ErrIsNil(context.TODO(), err)

			//todo 数据转换

			if len(tmpTerminals) < selectPageReq.count {
				break
			}
			selectPageReq.start += selectPageReq.count
		}
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiMCGetPersonalTemplates
//@summary 获取用户所有会议模板
//@param1  username string
//@return1
//@return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiMCGetPersonalTemplates(username string) (templates []MCUMCPersonalTemplateBaseInfo, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		//查询所有个人会议模板,分页查询查完为止
		selectPageReq := mcuSelectPageReq{
			count: 100,
			start: 0,
			order: 0,
		}
		for {
			var tmpTemplates []mcuMCPersonalTemplatesSimple
			tmpTemplates, err = mcuApiMCGetPersonalTemplates(mcuApiHead{
				mcuApiLevel:     c.mcuApiLevel,
				mcu:             c.mcuConnParam,
				mcuAccountToken: userLoginCache.AccountToken,
				mcuUserCookie:   userLoginCache.UserCookie,
			}, selectPageReq)
			valueCheck.ErrIsNil(context.TODO(), err)

			//数据转换
			for _, v := range tmpTemplates {
				templates = append(templates, MCUMCPersonalTemplateBaseInfo{
					mcuMCPersonalTemplatesSimple: v,
				})
			}

			if len(tmpTemplates) < selectPageReq.count {
				break
			}
			selectPageReq.start += selectPageReq.count
		}
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiMCGetOnePersonalTemplateDetail
//@summary 获取会议模板详情
//@param1  username string
//@return1
//@return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiMCGetOnePersonalTemplateDetail(username, templateId string) (templateDetail MCUMCPersonalTemplateDetailInfo, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpTemplateDetail mcuMCPersonalTemplatesDetail
		tmpTemplateDetail, err = mcuApiMCGetOnePersonalTemplates(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, templateId)
		valueCheck.ErrIsNil(context.TODO(), err)

		templateDetail.mcuMCPersonalTemplatesDetail = tmpTemplateDetail
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiMCCreateConf
//@summary 创建会议
//@param1  username string
//@param2  createReq MCUCreateConfReq "创建会议请求参数"
//@return1 meetingMark MCUMeetingMark "会议标识"
//@return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiMCCreateConf(username string, createReq MCUCreateConfReq) (meetingMark MCUMeetingMark, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var confMark mcuMCConfMark
		confMark, err = mcuApiMCCreateConf(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, createReq.CreateType, createReq.TemplateId, createReq.Name, createReq.Duration)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		err = gconv.Struct(confMark, &meetingMark)
		valueCheck.ErrIsNil(context.TODO(), err)

	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiMCReleaseConf
//@summary 结束会议
//@param1  username string
//@param2  confId string "会议Id"
//@param3  bReleaseSubConf string "是否结束级联会议:0=否 1=是"
//@return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiMCReleaseConf(username string, confId string, bReleaseSubConf bool) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		releaseSubConf := "0"
		if bReleaseSubConf {
			releaseSubConf = "1"
		}
		err = mcuApiMCReleaseConf(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, releaseSubConf)
		valueCheck.ErrIsNil(context.TODO(), err)

	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiVCGetConfDetail
//@summary 获取视频会议详情
//@param1  username string "用户"
//@param2  confId string "会议id"
//@return1 confDetail MCUMeetingDetail "返回会议信息"
//@return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCGetConfDetail(username string, confId string) (confDetail MCUMeetingDetail, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpConfInfo mcuVCConfDetail
		tmpConfInfo, err = mcuApiVCGetOneConf(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		err = gconv.Struct(tmpConfInfo, &confDetail)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiVCGetConfCascades
//@summary 获取会议级联详情
//@param1  username string "用户"
//@param2  confId string "会议id"
//@return1 confCascadeInfo []MCUMeetingCascade "返回会议信息"
//@return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCGetConfCascades(username string, confId string) (confCascadeInfo []MCUMeetingCascade, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpConfCascadeInfo []mcuVCCascadeConf
		tmpConfCascadeInfo, err = mcuApiVCGetConfCascades(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		err = gconv.Struct(tmpConfCascadeInfo, &confCascadeInfo)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiVCSetSilence
//@summary 会场静音操作
//@param1  username string "用户"
//@param2  confId string "会议id"
//@param3  value int "静音状态" 0=停止静音，1=静音
//@return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCSetSilence(username string, confId string, value int) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCSetSilence(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, value)
		valueCheck.ErrIsNil(context.TODO(), err)

	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiVCSetMute
//@summary 会场哑音操作
//@param1  username string "用户"
//@param2  confId string "会议id"
//@param3  state int "静音状态" 0=停止静音，1=静音
//@param4  forceMute int "全场哑音下是否禁止终端取消自身哑音"
//@return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCSetMute(username string, confId string, state int, forceMute int) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCSetMute(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, state, forceMute)
		valueCheck.ErrIsNil(context.TODO(), err)

	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCGetChairman
// @summary 获取会议主席
// @param1  username string "用户"
// @param2  confId string "会议id"
// @return1 mtId string "终端id"
// @return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCGetChairman(username string, confId string) (mtId string, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		mtId, err = mcuApiVCGetChairman(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCSetChairman
// @summary 指定会议主席，同步操作
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mtId string "指定终端为主席参数"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCSetChairman(username string, confId string, mtId string) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCSetChairman(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mtId)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCGetSpeaker
// @summary 获取会议发言人
// @param1  username string "用户"
// @param2  confId string "会议id"
// @return1 mtId string "终端id"
// @return2 forceBroadCase int "是否设置了发言人强制广播，仅SFU会议有效"
// @return3 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCGetSpeaker(username, confId string) (mtId string, forceBroadCase int, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpResult mcuVCMtSet
		tmpResult, err = mcuApiVCGetSpeaker(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId)
		valueCheck.ErrIsNil(context.TODO(), err)

		mtId = tmpResult.MtId
		forceBroadCase = tmpResult.ForceBroadcase
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCSetSpeaker
// @summary 指定会议发言人，同步操作
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mtId string "指定终端为发言人"
// @param4  broadcast int "是否设置发言人强制广播,0=否，1=是"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCSetSpeaker(username, confId string, mtId string, broadcast int) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCSetSpeaker(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mtId, broadcast)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCSetForceBroadcast
// @summary 设置会议是否强制广播
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mode int "状态"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCSetForceBroadcast(username, confId string, mode int) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCSetForceBroadcast(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mode)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCGetVadState
// @summary 获取会议语音激励状态
// @param1  username string "用户"
// @param2  confId string "会议id"
// @return1 vadState MCUMeetingVadState "语音激励状态"
// @return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCGetVadState(username, confId string) (vadState MCUMeetingVadState, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpVadState mcuVCVac
		tmpVadState, err = mcuApiVCGetVadState(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		_ = gconv.Structs(tmpVadState, &vadState)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCSetVadState
// @summary 设置会议语音激励状态
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  state int "语音激励" 0=关闭，1=开启
// @param4  interval uint "语音激励敏感度" 单位s，最小值3s
// @return1 vadState MCUMeetingVadState "语音激励状态"
// @return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCSetVadState(username, confId string, state int, interval uint) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCSetVadState(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, state, interval)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCGetMixs
// @summary 查询会议混音状态
// @param1  username string "用户"
// @param2  confId string "会议id"
// @return1 meetingMixState MCUMeetingMixState "语音激励状态"
// @return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCGetMixs(username, confId string) (meetingMixState MCUMeetingMixState, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpMixState mcuMCMix
		tmpMixState, err = mcuApiVCGetMixs(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		_ = gconv.Structs(tmpMixState, &meetingMixState)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCStartMixs
// @summary 开启会议混音状态
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mode int "混音模式" 1=智能混音，2=定制混音
// @param4  mtIds []string "混音成员列表"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCStartMixs(username, confId string, mode int, mtIds []string) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCStartMixs(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mode, mtIds)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCStopMixs
// @summary 开启会议混音状态
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mixId string "混音器ID，默认mix_id为1"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCStopMixs(username, confId string, mixId string) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCStopMixs(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mixId)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCAddMixMts
// @summary 添加混音成员
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mixId string "混音器ID，默认mix_id为1"
// @param4  mtIds []string "混音成员列表"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCAddMixMts(username, confId, mixId string, mtIds []string) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCAddMixMts(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mixId, mtIds)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCDelMixMts
// @summary 删除混音成员
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mixId string "混音器ID，默认mix_id为1"
// @param4  mtIds []string "混音成员列表"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCDelMixMts(username, confId, mixId string, mtIds []string) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCDelMixMts(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mixId, mtIds)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCGetVmps
// @summary 获取会议画面合成信息
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  vmpId string "画面合成ID，默认vmpId为1"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCGetVmps(username, confId, vmpId string) (vmpState MCUMeetingVmpParam, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpVmpState mcuMCVmp
		tmpVmpState, err = mcuApiVCGetVmps(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		//TODO Tips: 内部有重名的字段时 gconv.Struct 数据转换有bug
		vmpState = MCUMeetingVmpParam{
			Mode:       tmpVmpState.Mode,
			Layout:     tmpVmpState.Layout,
			ExceptSelf: tmpVmpState.ExceptSelf,
			VoiceHint:  tmpVmpState.VoiceHint,
			Broadcast:  tmpVmpState.Broadcast,
			ShowMtName: tmpVmpState.ShowMtName,
			MtNameStyle: MCUVmpMTNameStyle{
				FontSize:  tmpVmpState.MtNameStyle.FontSize,
				FontColor: tmpVmpState.MtNameStyle.FontColor,
				Position:  tmpVmpState.MtNameStyle.Position,
			},
		}
		for v := range tmpVmpState.Members {
			var member MCUVmpMember
			_ = gconv.Struct(v, &member)
			vmpState.Members = append(vmpState.Members, member)
		}
		_ = gconv.Struct(tmpVmpState.Poll, &vmpState.Poll)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCStartVmps
// @summary 开启会议画面合成
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  vmp mcuMCVmp "画面合成参数"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCStartVmps(username, confId string, vmpParam MCUMeetingVmpParam) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		data := mcuMCVmp{
			Mode:       vmpParam.Mode,
			Layout:     vmpParam.Layout,
			ExceptSelf: vmpParam.ExceptSelf,
			VoiceHint:  vmpParam.VoiceHint,
			Broadcast:  vmpParam.Broadcast,
			ShowMtName: vmpParam.ShowMtName,
			MtNameStyle: mcuMCVmpMtNameStyle{
				FontSize:  vmpParam.MtNameStyle.FontSize,
				FontColor: vmpParam.MtNameStyle.FontColor,
				Position:  vmpParam.MtNameStyle.Position,
			},
		}
		for v := range vmpParam.Members {
			var member mcuMCVmpMember
			_ = gconv.Struct(v, &member)
			data.Members = append(data.Members, member)
		}
		_ = gconv.Struct(vmpParam.Poll, &data.Poll)

		err = mcuApiVCStartVmps(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, data)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCUpdateVmps
// @summary 更新会议画面合成参数
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  vmp mcuMCVmp "画面合成参数"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCUpdateVmps(username, confId string, vmpParam MCUMeetingVmpParam) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		data := mcuMCVmp{
			Mode:       vmpParam.Mode,
			Layout:     vmpParam.Layout,
			ExceptSelf: vmpParam.ExceptSelf,
			VoiceHint:  vmpParam.VoiceHint,
			Broadcast:  vmpParam.Broadcast,
			ShowMtName: vmpParam.ShowMtName,
			MtNameStyle: mcuMCVmpMtNameStyle{
				FontSize:  vmpParam.MtNameStyle.FontSize,
				FontColor: vmpParam.MtNameStyle.FontColor,
				Position:  vmpParam.MtNameStyle.Position,
			},
		}
		for v := range vmpParam.Members {
			var member mcuMCVmpMember
			_ = gconv.Struct(v, &member)
			data.Members = append(data.Members, member)
		}
		_ = gconv.Struct(vmpParam.Poll, &data.Poll)

		err = mcuApiVCModVmps(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, data)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCStopVmps
// @summary 更新会议画面合成参数
// @param1  username string "用户"
// @param2  confId string "会议id"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCStopVmps(username, confId string) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCStopVmps(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiVCGetCurConfMts
//@summary 获取本级会议终端列表
//@param1  username string "用户"
//@param2  confId string "会议id"
//@return1 terminals []MCUMeetTerminal "返回会议终端列表"
//@return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCGetCurConfMts(username string, confId string) (terminals []MCUMeetTerminal, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpMts []mcuVCMt
		tmpMts, err = mcuApiVCGetCurConfMts(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		_ = gconv.Structs(tmpMts, &terminals)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiVCGetCascadesMts
//@summary 获取级联会议终端列表
//@param1  username string "用户"
//@param2  confId string "会议id"
//@param2  cascadeId string "级联会议Id, 0表示本级会议"
//@return1 terminals []MCUMeetTerminal "返回会议终端列表"
//@return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCGetCascadesMts(username, confId, cascadeId string) (terminals []MCUMeetTerminal, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpMts []mcuVCMt
		tmpMts, err = mcuApiVCGetCascadesMts(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, cascadeId)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		_ = gconv.Structs(tmpMts, &terminals)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

//MCUApiVCGetMt
//@summary 获取会议终端详细信息
//@param1  username string "用户"
//@param2  confId string "会议id"
//@param2  mtId string "会议终端ID"
//@return1 terminal MCUMeetTerminal "返回会议终端信息"
//@return2 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCGetMt(username, confId, mtId string) (terminal MCUMeetTerminal, err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var tmpMt mcuVCMt
		tmpMt, err = mcuApiVCGetMt(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mtId)
		valueCheck.ErrIsNil(context.TODO(), err)

		//数据转换
		_ = gconv.Struct(tmpMt, &terminal)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCCallMts
// @summary 批量呼叫终端，异步操作
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mts []mcuVCMtId "批量呼叫会议终端"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCCallMts(username, confId string, mts []MCUCallMtReq) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		var req []mcuVCMtId
		for _, v := range mts {
			req = append(req, mcuVCMtId{
				MtId:       v.MtId,
				ForcedCall: v.ForcedCall,
			})
		}
		err = mcuApiVCCallMts(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, req)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCDropMts
// @summary 批量挂断终端，同步操作
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mtIds [] string "批量挂断会议终端Id"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCDropMts(username, confId string, mtIds []string) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCDropMts(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mtIds)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCMtSetSilence
// @summary 设置与会终端静音
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mtId string "会议终端Id"
// @param4  silence int "是否静音: 0=否 1=是"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCMtSetSilence(username, confId string, mtId string, silence int) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCMtSetSilence(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mtId, silence)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCMtSetMute
// @summary 设置与会终端哑音
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mtId string "会议终端Id"
// @param4  silence int "是否哑音: 0=否 1=是"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCMtSetMute(username, confId string, mtId string, mute int) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCMtSetMute(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mtId, mute)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}

// MCUApiVCMtSetVolume
// @summary 设置与会终端音量
// @param1  username string "用户"
// @param2  confId string "会议id"
// @param3  mtId string "会议终端Id"
// @param4  mode int "扬声器or麦克风: 1=扬声器 2=麦克风"
// @param5  silence int "声音大小: 0-35"
// @return1 err error "返回报错信息"
func (c *McuAdapterCore) MCUApiVCMtSetVolume(username, confId string, mtId string, mode int, volume int) (err error) {
	err = g.Try(func() {
		var userLoginCache *MCULoginCache
		userLoginCache, err = c.getUserLoginCache(username)
		valueCheck.ErrIsNil(context.TODO(), err)

		err = mcuApiVCMtSetVolume(mcuApiHead{
			mcuApiLevel:     c.mcuApiLevel,
			mcu:             c.mcuConnParam,
			mcuAccountToken: userLoginCache.AccountToken,
			mcuUserCookie:   userLoginCache.UserCookie,
		}, confId, mtId, mode, volume)
		valueCheck.ErrIsNil(context.TODO(), err)
	})
	if err != nil && c.checkNeedDelLoginCacheByMCUErr(err) {
		c.mcuUserLoginCaches.Remove(username)
	}
	return
}
