//========================
// v5.2版本API接口
//========================

package accessMCUAdapter

import (
	"context"
	"fmt"
	"github.com/gogf/gf/errors/gerror"
	"github.com/sigmavirus24/gobayeux"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type MCUAdapterV5_2 struct {
	*McuAdapterCore
}

//login
//@summary 登录MCU
func (mcuApiV5_2 *MCUAdapterV5_2) login(user *MCUUserInfo) (userLoginCache *MCULoginCache, err error) {
	userLoginCache = &MCULoginCache{
		User:         user,
		AccountToken: "",
		UserCookie:   "",
		ConnState:    EMCUNetConnErr,
	}
	//获取Token
	userLoginCache.AccountToken, err = mcuApiSystemGetToken(mcuApiV5_2.mcuApiLevel, *mcuApiV5_2.mcuConnParam)
	if err != nil {
		//解析报错原因
		if err.Error() == "10001" {
			userLoginCache.ConnState = EMCUSecretInvalid
			mcuApiV5_2.mcuConnErrCB(MCUConnState{LastErrCode: EMCUSecretInvalid, McuUsername: ""})
		} else {
			userLoginCache.ConnState = EMCUNetConnErr
			mcuApiV5_2.mcuConnErrCB(MCUConnState{LastErrCode: EMCUNetConnErr, McuUsername: ""})
			return nil, gerror.New(fmt.Sprintf("%d", EMCUNetConnErr))
		}
		return nil, err
	}
	mcuApiV5_2.mcuConnErrCB(MCUConnState{LastErrCode: EMCUConnOk, McuUsername: ""})
	//登录
	userLoginCache.UserCookie, err = mcuApiSystemLogin(mcuApiV5_2.mcuApiLevel, *mcuApiV5_2.mcuConnParam, mcuSystemLoginReq{
		AccountToken: userLoginCache.AccountToken,
		Username:     user.Username,
		Password:     user.Password,
	})
	if err != nil {
		userLoginCache.ConnState = EMCUUserInvalid
		mcuApiV5_2.mcuConnErrCB(MCUConnState{LastErrCode: EMCUUserInvalid, McuUsername: user.Username})
		return nil, err
	}

	//todo 实测5.2版本不支持获取Api版本接口,默认 ApiLevel=1
	mcuApiV5_2.mcuApiLevel = 1
	//获取MCU版本信息
	//mcuApiV5_2.mcuApiLevel, err = mcuApiSystemGetVersion(mcuApiV5_2.mcuApiLevel, *mcuApiV5_2.mcuConnParam, userLoginCache.AccountToken, userLoginCache.UserCookie)
	//if err != nil {
	//	logger().Error(err)
	//	return nil, err
	//}
	//获取域信息(只有MCU管理员账户才有此接口权限).	//TODO Tips:登录接口会返回该账户所属用户域ID
	var domains []mcuVMSDomains
	domains, err = mcuApiNMSGetDomains(mcuApiHead{
		mcuApiLevel:     mcuApiV5_2.mcuApiLevel,
		mcu:             mcuApiV5_2.mcuConnParam,
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
	mcuApiV5_2.mcuConnErrCB(MCUConnState{LastErrCode: EMCUConnOk, McuUsername: user.Username})
	return
}

//subscribeMeeting
//@summary 订阅MCU通知
//         科达5.2版本MCU订阅服务为CometD
func (mcuApiV5_2 *MCUAdapterV5_2) subscribeMeeting() {
	logger().Debug("subscribeMeeting start:", mcuApiV5_2.mcuConnParam)
_reconnect_sub:
	if mcuApiV5_2.bDestroy {
		return
	}

	//找到一个在线用户
	var userCache *MCULoginCache
	for _, user := range mcuApiV5_2.mcuConnParam.MCUUsers {
		tmpUserCache, err := mcuApiV5_2.getUserLoginCache(user.Username)
		if err == nil && tmpUserCache != nil {
			userCache = tmpUserCache
			break
		}
	}
	if userCache == nil {
		mcuApiV5_2.subByMCUUser = nil
		logger().Error("MCU未登录:", mcuApiV5_2.mcuConnParam)
		time.Sleep(30 * time.Second)
		goto _reconnect_sub
	}
	mcuApiV5_2.subByMCUUser = userCache
	logger().Noticef("mcu user[%s] sub notify. userDomain[%s] mcuAccountToken[%s] userCookie[%s]",
		userCache.User.Username, userCache.User.UserDomain, userCache.AccountToken, userCache.UserCookie)

	//订阅通知
	chMCUMsg := make(chan []gobayeux.Message, 1000)
	mcuBayeuxServerUrl := url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", mcuApiV5_2.mcuConnParam.Ip, mcuApiV5_2.mcuConnParam.Port), Path: "/api/v1/publish"}
	goBayeuxLogger := logrus.New()
	goBayeuxLogger.SetLevel(logrus.ErrorLevel)

	//new httpClient with cookie
	httpClient := http.Client{}
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		logger().Error(err)
		goto _reconnect_sub
	}
	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:  "SSO_COOKIE_KEY",
		Value: userCache.UserCookie,
	})
	jar.SetCookies(&mcuBayeuxServerUrl, cookies)
	httpClient.Timeout = 60 * time.Second
	httpClient.Jar = jar

	//new gobayeux client
	goBayeuxClient, err := gobayeux.NewClient(mcuBayeuxServerUrl.String(), gobayeux.WithHTTPClient(&httpClient), gobayeux.WithLogger(goBayeuxLogger))
	if err != nil {
		logger().Error("error initializing client:", err)
		goto _reconnect_sub
	}
	logger().Notice("got gobayeux client")

	//ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(60*time.Second))
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	//goBayeuxClient.Publish(ctx)
	chBayeuxErr := goBayeuxClient.Start(ctx)

	//订阅MCU所有会议信息
	subChanel := gobayeux.Channel("/userdomains/" + userCache.User.UserDomain + "/confs/**")
	goBayeuxClient.Subscribe(subChanel, chMCUMsg)

	//接收消息
	for {
		if mcuApiV5_2.bDestroy {
			return
		}

		select {
		case e := <-chBayeuxErr:
			logger().Notice(e)
			time.Sleep(10 * time.Second)
			goto _reconnect_sub

		case messages := <-chMCUMsg:
			for _, tmpMsg := range messages {
				//构建MCU通知消息
				if bOk, msg := mcuNotifyAssistant(string(tmpMsg.Channel), string(tmpMsg.Data)); bOk {
					msg.Mcu = mcuApiV5_2.mcuConnParam
					msg.User = userCache
					chnMCUNotifyMsg <- msg
				} else {
					logger().Error("error analysis notify message:", tmpMsg)
				}
			}

			//nolint:staticcheck
			//default:
		}
	}
}
