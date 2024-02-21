package accessMCUAdapter

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"time"
)

type MCUAdapterV8_0 struct {
	*McuAdapterCore
}

//subscribeMeeting
//@summary 订阅MCU通知
//         科达8.0版本MCU订阅服务为MQTT
func (mcuApiV8_0 *MCUAdapterV8_0) subscribeMeeting() {
	logger().Debug("subscribeMeeting start:", mcuApiV8_0.mcuConnParam)
_reconnect_sub:
	if mcuApiV8_0.bDestroy {
		return
	}

	//找到一个在线用户
	var userCache *MCULoginCache
	for _, user := range mcuApiV8_0.mcuConnParam.MCUUsers {
		tmpUserCache, err := mcuApiV8_0.getUserLoginCache(user.Username)
		if err == nil && tmpUserCache != nil {
			userCache = tmpUserCache
			break
		}
	}
	if userCache == nil {
		mcuApiV8_0.subByMCUUser = nil
		logger().Error("MCU未登录:", mcuApiV8_0.mcuConnParam)
		time.Sleep(30 * time.Second)
		goto _reconnect_sub
	}

	mcuApiV8_0.subByMCUUser = userCache
	logger().Noticef("mcu user[%s] sub notify. userDomain[%s] mcuAccountToken[%s] userCookie[%s]",
		userCache.User.Username, userCache.User.UserDomain, userCache.AccountToken, userCache.UserCookie)

	//连接MQTT服务[不需要SetClientID(), MCU_MQTT服务器会为每个客户端连接自动生成一个ID]
	mqttClient := mqtt.NewClient(mqtt.NewClientOptions().
		SetUsername(userCache.UserCookie).
		SetPassword(userCache.AccountToken).
		SetCleanSession(true).
		AddBroker(fmt.Sprintf("mqtt://%s:%d", mcuApiV8_0.mcuConnParam.Ip, 1884)))
	if token := mqttClient.Connect(); token.Wait() && token.Error() != nil {
		logger().Error(token.Error())
		goto _reconnect_sub
	}

	logger().Notice("got mqtt client")

	//订阅MCU所有会议信息
	token := mqttClient.Subscribe(fmt.Sprintf("/userdomains/"+userCache.User.UserDomain+"/confs/#"), 0, mcuApiV8_0.mqttMessageCB)
	if token.Wait() && token.Error() != nil {
		logger().Error(token.Error())
		goto _reconnect_sub
	}

	//MQTT订阅连接保活
	for {
		time.Sleep(10 * time.Second)
		if mcuApiV8_0.bDestroy {
			return
		}
		if !mqttClient.IsConnected() {
			logger().Notice("mqtt sub close. mcu:", mcuApiV8_0.mcuConnParam)
			goto _reconnect_sub
		}
	}
}

//mqttMessageCB
func (mcuApiV8_0 *MCUAdapterV8_0) mqttMessageCB(client mqtt.Client, mqttMsg mqtt.Message) {
	//MQTT订阅消息
	logger().Debug("MCU Publish Msg:", mqttMsg.Topic(), mqttMsg.Payload())

	//构建MCU通知消息
	if bOk, msg := mcuNotifyAssistant(mqttMsg.Topic(), string(mqttMsg.Payload())); bOk {
		msg.Mcu = mcuApiV8_0.mcuConnParam
		msg.User = mcuApiV8_0.subByMCUUser
		chnMCUNotifyMsg <- msg
	} else {
		logger().Error("error analysis notify message:", mqttMsg)
	}
}
