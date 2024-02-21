package publish

import _ "gowebserver/app/common/define"

type PubMode int

const (
	PubModeToAll                PubMode = 1 //通知所有订阅客户端
	PubModeToSpecifiedClient    PubMode = 2 //通知特定订阅客户端
	PubModeOutOfSpecifiedClient PubMode = 3 //排除特定订阅客户端
)

// PubCommMsg 通用发布消息结构
type PubCommMsg struct {
	PubMode            PubMode     `json:"-"`         //发布模式 PubMode
	SpecialClientToken string      `json:"-"`         //根据发布模式需特别处理的客户端令牌
	TimeStamp          int64       `json:"timestamp"` //时间戳
	Topic              string      `json:"topic"`     //订阅主题 define.ENotifyTopic
	Data               interface{} `json:"data"`      //数据内容
}

func init() {
	WsPublish.InitInstance()
	RedisPublish.InitInstance()
}
