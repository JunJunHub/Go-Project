//========================
// v6.0版本API接口
//========================

package accessMCUAdapter

import (
	"context"
	"fmt"
	"github.com/sigmavirus24/gobayeux"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type MCUAdapterV6_0 struct {
	*McuAdapterCore
}

//subscribeMeeting
//@summary 订阅MCU通知
//         科达6.0版本MCU订阅服务未CometD
func (mcuApiV6_0 *MCUAdapterV6_0) subscribeMeeting() {
	logger().Debug("subscribeMeeting start:", mcuApiV6_0.mcuConnParam)
_reconnect_sub:
	if mcuApiV6_0.bDestroy {
		return
	}

	//找到一个在线用户
	var userCache *MCULoginCache
	for _, user := range mcuApiV6_0.mcuConnParam.MCUUsers {
		tmpUserCache, err := mcuApiV6_0.getUserLoginCache(user.Username)
		if err == nil && tmpUserCache != nil {
			userCache = tmpUserCache
			break
		}
	}
	if userCache == nil {
		mcuApiV6_0.subByMCUUser = nil
		logger().Error("MCU未登录:", mcuApiV6_0.mcuConnParam)
		time.Sleep(30 * time.Second)
		goto _reconnect_sub
	}
	mcuApiV6_0.subByMCUUser = userCache
	logger().Noticef("mcu[%s] user[%s] sub notify. userDomain[%s] mcuAccountToken[%s] userCookie[%s]",
		mcuApiV6_0.mcuKey, userCache.User.Username, userCache.User.UserDomain, userCache.AccountToken, userCache.UserCookie)

	//订阅通知
	chMCUMsg := make(chan []gobayeux.Message, 1000)
	mcuBayeuxServerUrl := url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", mcuApiV6_0.mcuConnParam.Ip, mcuApiV6_0.mcuConnParam.Port), Path: "/api/v1/publish"}
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
		if mcuApiV6_0.bDestroy {
			return
		}

		select {
		case e := <-chBayeuxErr:
			logger().Notice(e)
			time.Sleep(10 * time.Second)
			goto _reconnect_sub

		case messages := <-chMCUMsg:
			for _, tmpMsg := range messages {
				//logger().Debug(tmpMsg)
				//goBayeuxLogger.WithFields(logrus.Fields{
				//	"channel": tmpMsg.Channel,
				//	"data":    string(tmpMsg.Data),
				//}).Info()

				//构建MCU通知消息
				if bOk, msg := mcuNotifyAssistant(string(tmpMsg.Channel), string(tmpMsg.Data)); bOk {
					msg.Mcu = mcuApiV6_0.mcuConnParam
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
