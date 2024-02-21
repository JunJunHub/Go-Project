package accessMCUAdapter

import (
	"context"
	"fmt"
	"github.com/eclipse/paho.mqtt.golang"
	"github.com/sigmavirus24/gobayeux"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"testing"
	"time"
)

var mcuHead mcuApiHead

//测试sdk登录4.7会议平台
func TestSdkConnMcu(t *testing.T) {
	//McsDllInit()
}

//测试登录会议平台
func TestConnMCU(t *testing.T) {
	mcu := MCUServiceConnParam{
		AdapterVer:          EMCUAdapterV6_0,
		Ip:                  "10.67.24.84",
		Port:                80,
		OAuthConsumerKey:    "snZymbs4kjUE",
		OAuthConsumerSecret: "kBQhLgPUzkha",
		SSLCertFile:         "",
		SSLKeyFile:          "",
		MCUUsers:            nil,
	}

	//获取Token
	mcuAccountToken, err := mcuApiSystemGetToken(1, mcu)
	if err != nil {
		t.Log(err)
		return
	}

	//登录
	userCookie, err := mcuApiSystemLogin(1, mcu, mcuSystemLoginReq{
		AccountToken: mcuAccountToken,
		Username:     "admin1",
		Password:     "keda8888",
	})
	if err != nil {
		t.Log(err)
		return
	}

	//	var mcuHead MCUApiHead
	mcuHead.mcuApiLevel = 1
	mcuHead.mcuAccountToken = mcuAccountToken
	mcuHead.mcuUserCookie = userCookie
	mcuHead.mcu = &mcu

	//获取个人模板列表
	query := mcuSelectPageReq{
		count: 10,
		start: 0,
		order: 0,
	}

	confList, err := mcuApiMCGetPersonalTemplates(mcuHead, query)
	if err != nil {
		t.Log(err)
	}

	for _, templatesInfo := range confList {
		logger().Debug("template id:", templatesInfo.TemplateId, "template name:", templatesInfo.Name, "e164:", templatesInfo.E164)
		//获取个人模板详细信息
		conf, err := mcuApiMCGetOnePersonalTemplates(mcuHead, templatesInfo.TemplateId)
		if err != nil {
			t.Log(err)
		}
		logger().Debug("template id:", conf.TemplateId, "template name:", conf.Name, "e164:", conf.E164)
		if conf.Name == "xuxiaowen" {
			//创建会议
			confMark, err := mcuApiMCCreateConf(mcuHead, 3, conf.TemplateId, "xuxiaowen的会议的名字改一改", 0)
			if err != nil {
				t.Log(err)
			}
			logger().Debug("meeting id:", confMark.MeetingId, "conf id:", confMark.ConfId, "moid:", confMark.MachineRoomMoid)

			if conf.Name == "xuxiaowen的会议的名字改一改" {
				//获取语音激励状态
				vadStat, err := mcuApiVCGetVadState(mcuHead, confMark.ConfId)
				if err != nil {
					t.Log(err)
				}
				logger().Debug(vadStat)
			}
		}
	}

	//获取正在召开的视频会议列表
	meetingList, err := mcuApiVCGetConfs(mcuHead, query)
	if err != nil {
		logger().Error(err)
	}
	logger().Debug("total:", len(meetingList))
	for _, meetingInfo := range meetingList {
		logger().Debug("name:", meetingInfo.Name, "confid:", meetingInfo.ConfId, "meetingid:", meetingInfo.MeetingId)
		//获取视频会议详情
		meeting, err := mcuApiVCGetOneConf(mcuHead, meetingInfo.ConfId)
		if err != nil {
			logger().Error(err)
		}
		logger().Debug("conf name:", meeting.Name, "conf id:", meeting.ConfId, "user domain name:", meeting.UserDomainName)
		if meeting.Name == "xuxiaowen的会议的名字改一改" {
			err = mcuApiMCReleaseConf(mcuHead, meeting.ConfId, "0")
			if err != nil {
				t.Log(err)
			}

		}
		if meeting.Name == "省厅会议" {
			//获取级联会议
			cascadesList, err := mcuApiVCGetConfCascades(mcuHead, meeting.ConfId)
			if err != nil {
				t.Log(err)
			}
			logger().Debug(cascadesList)
			//获取终端列表
			mtList, err := mcuApiVCGetCurConfMts(mcuHead, meeting.ConfId)
			if err != nil {
				t.Log(err)
			}
			logger().Debug(mtList)

			//发送短消息
			var message mcuVCSms
			message.Message = "jksjdljdkgjlsdjsldk"
			message.Type = 2
			message.RollNum = 1
			message.RollSpeed = 3
			for _, smsMt := range mtList {
				var strMtId mcuVCMtId
				strMtId.MtId = smsMt.MtId
				message.Mts = append(message.Mts, strMtId)
			}
			mcuApiVCSms(mcuHead, meeting.ConfId, message)

			//添加终端入会
			var mtsAdd []mcuVCMtSim
			mtsAdd = make([]mcuVCMtSim, 0, len(mtsAdd))
			var index = 1
			for _, reqInfo := range mtList {
				var reqRes mcuVCMtSim
				reqRes.AccountType = reqInfo.AccountType
				var alias string
				alias = fmt.Sprintf("10.67.24.%d", index)
				reqRes.Account = alias
				reqRes.Bitrate = reqInfo.Bitrate
				reqRes.Protocol = 0
				reqRes.ForcedCall = 0
				reqRes.CallMode = 0
				mtsAdd = append(mtsAdd, reqRes)
				index += 1
			}
			mtsAddOut, err := mcuApiVCAddCascadeMts(mcuHead, meeting.ConfId, "0", mtsAdd)

			//删除终端
			var mtIds []string
			for _, addOutMt := range mtsAddOut {
				var strMtId string
				strMtId = addOutMt.MtId
				mtIds = append(mtIds, strMtId)
			}
			mcuApiVCDelMts(mcuHead, meeting.ConfId, mtIds)

		}
		if meeting.Name == "会议终端" {
			//获取会议主席
			chairman, err := mcuApiVCGetChairman(mcuHead, meeting.ConfId)
			if err != nil {
				t.Log(err)
			}
			logger().Debug(chairman)
			//获取主席终端信息
			chairmanInfo, err := mcuApiVCGetMt(mcuHead, meeting.ConfId, chairman)
			logger().Debug(chairmanInfo)
			speaker, err := mcuApiVCGetSpeaker(mcuHead, meeting.ConfId)
			if err != nil {
				t.Log(err)
			}
			logger().Debug(speaker)
			//主席静音

			mcuApiVCMtSetSilence(mcuHead, meeting.ConfId, chairmanInfo.MtId, 1)

			//获取发言人终端信息
			speakerInfo, err := mcuApiVCGetMt(mcuHead, meeting.ConfId, speaker.MtId)
			logger().Debug(speakerInfo)

			//获取混音
			mtMix, err := mcuApiVCGetMixs(mcuHead, meeting.ConfId)
			if err != nil {
				t.Log(err)
			}
			logger().Debug(mtMix)

			//获取选看终端列表
			mtInpecList, err := mcuApiVCGetInspections(mcuHead, meeting.ConfId)
			if err != nil {
				t.Log(err)
			}
			logger().Debug(mtInpecList)

			for _, mtInpect := range mtInpecList {
				//取消选看
				mcuApiVCStopInspect(mcuHead, meeting.ConfId, mtInpect.Dst.MtId, mtInpect.Mode)
			}
		}
	}
}

//测试连接会议CometD服务
func TestConnMCUCometDServer(t *testing.T) {
	mcu := MCUServiceConnParam{
		AdapterVer:          EMCUAdapterV6_0,
		Ip:                  "10.67.24.84",
		Port:                80,
		OAuthConsumerKey:    "snZymbs4kjUE",
		OAuthConsumerSecret: "kBQhLgPUzkha",
		SSLCertFile:         "",
		SSLKeyFile:          "",
		MCUUsers:            nil,
	}
	//获取Token
	mcuAccountToken, err := mcuApiSystemGetToken(1, mcu)
	if err != nil {
		t.Log(err)
		return
	}

	//登录
	userCookie, err := mcuApiSystemLogin(1, mcu, mcuSystemLoginReq{
		AccountToken: mcuAccountToken,
		Username:     "admin1",
		Password:     "keda8888",
	})
	if err != nil {
		t.Log(err)
		return
	}

	//	var mcuHead MCUApiHead
	mcuHead.mcuApiLevel = 1
	mcuHead.mcuAccountToken = mcuAccountToken
	mcuHead.mcuUserCookie = userCookie
	mcuHead.mcu = &mcu

	//获取用户域
	userDomains, err := mcuApiNMSGetDomains(mcuHead)
	if err != nil {
		t.Log(err)
	}
	t.Log("mcuAccountToken:", mcuHead.mcuAccountToken, "userCookie:", mcuHead.mcuUserCookie)

	//订阅通知
	chMCUNoticeMsg := make(chan []gobayeux.Message, 1000)
	mcuBayeuxServerUrl := url.URL{Scheme: "http", Host: fmt.Sprintf("%s:%d", "10.67.24.84", 80), Path: "/api/v1/publish"}
	goBayeuxLogger := logrus.New()
	goBayeuxLogger.SetLevel(logrus.TraceLevel)

	//new gobayeux client
	httpClient := http.Client{}
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		t.Log(err)
		os.Exit(-1)
	}
	var cookies []*http.Cookie
	cookies = append(cookies, &http.Cookie{
		Name:  "SSO_COOKIE_KEY",
		Value: mcuHead.mcuUserCookie,
	})
	jar.SetCookies(&mcuBayeuxServerUrl, cookies)
	httpClient.Timeout = 60 * time.Second
	httpClient.Jar = jar

	goBayeuxClient, err := gobayeux.NewClient(mcuBayeuxServerUrl.String(), gobayeux.WithHTTPClient(&httpClient), gobayeux.WithLogger(goBayeuxLogger))
	if err != nil {
		fmt.Printf("error initializing client: %q\n", err)
		os.Exit(1)
	}
	t.Log("got gobayeux client")

	ctx := context.Background()
	//ctx, cancel := context.WithDeadline(ctx, time.Now().Add(60*time.Second))
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	//goBayeuxClient.Publish(ctx)
	chErr := goBayeuxClient.Start(ctx)

	for _, userDomain := range userDomains {
		if userDomain.Name == "默认用户域" {
			//订阅通道
			subChanel := "/userdomains/" + userDomain.Moid + "/confs/**"
			goBayeuxClient.Subscribe(gobayeux.Channel(subChanel), chMCUNoticeMsg)

			//接收消息
			for {
				select {
				case err := <-chErr:
					t.Log("==============>", err)

				case ms := <-chMCUNoticeMsg:
					for _, m := range ms {
						t.Log("==============>", m)

						goBayeuxLogger.WithFields(logrus.Fields{
							"channel": m.Channel,
							"data":    string(m.Data),
						}).Info()

						//构建MCU通知消息
						var msg MCUNotifyMessage
						//解析通知
						bRet, msg := mcuNotifyAssistant(string(m.Channel), string(m.Data))
						if bRet == false {
							logger().Error("error analysis notify message:", err)
						} else {
							logger().Error("notify message:", msg.Type, "notify confid:", msg.ConfId, "notify method:", msg.Method)
						}

					}
				// nolint:staticcheck
				default:
				}
			}
		}
	}
}

//测试连接会议MQTT服务
func TestConnMCUMQTTServer(t *testing.T) {
	mcu := MCUServiceConnParam{
		AdapterVer:          EMCUAdapterV6_0,
		Ip:                  "10.67.24.84",
		Port:                80,
		OAuthConsumerKey:    "snZymbs4kjUE",
		OAuthConsumerSecret: "kBQhLgPUzkha",
		SSLCertFile:         "",
		SSLKeyFile:          "",
		MCUUsers:            nil,
	}

	//获取Token
	mcuAccountToken, err := mcuApiSystemGetToken(1, mcu)
	if err != nil {
		t.Log(err)
		return
	}

	//登录
	userCookie, err := mcuApiSystemLogin(1, mcu, mcuSystemLoginReq{
		AccountToken: mcuAccountToken,
		Username:     "admin1",
		Password:     "keda8888",
	})
	if err != nil {
		t.Log(err)
		return
	}
	t.Log("mcuAccountToken:", mcuAccountToken, "userCookie:", userCookie)

	//获取版本信息
	_, err = mcuApiSystemGetVersion(1, mcu, mcuAccountToken, userCookie)
	if err != nil {
		t.Log(err)
	}

	//连接MQTT服务
	mqttClient := mqtt.NewClient(mqtt.NewClientOptions().
		SetClientID("fooTest").
		SetUsername(userCookie).
		SetPassword(mcuAccountToken).
		AddBroker("mqtt://10.67.24.84:1884"))
	if mqttClient == nil {
		logger().Error("mqtt new client failed!")
		return
	}
	mqttClient.Connect()
}

//测试获取公共群组列表(MCU地址簿)
func TestGetPublicGroups(t *testing.T) {
	mcu := MCUServiceConnParam{
		AdapterVer:          EMCUAdapterV6_0,
		Ip:                  "10.67.24.84",
		Port:                80,
		OAuthConsumerKey:    "snZymbs4kjUE",
		OAuthConsumerSecret: "kBQhLgPUzkha",
		SSLCertFile:         "",
		SSLKeyFile:          "",
		MCUUsers:            nil,
	}

	//获取Token
	mcuAccountToken, err := mcuApiSystemGetToken(1, mcu)
	if err != nil {
		t.Log(err)
		return
	}

	//登录
	userCookie, err := mcuApiSystemLogin(1, mcu, mcuSystemLoginReq{
		AccountToken: mcuAccountToken,
		Username:     "admin1",
		Password:     "keda8888",
	})
	if err != nil {
		t.Log(err)
		return
	}

	//MCU连接句柄
	mcuHead.mcuApiLevel = 1
	mcuHead.mcuAccountToken = mcuAccountToken
	mcuHead.mcuUserCookie = userCookie
	mcuHead.mcu = &mcu

	_ = mcuApiAMCGetPublicGroups(mcuHead)
}
