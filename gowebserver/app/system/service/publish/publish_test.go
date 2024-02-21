package publish

import (
	"github.com/gogf/gf/frame/g"
	"testing"
)

func TestWsPublish_PublishMsg(t *testing.T) {
	WsPublish.PublishMsgToAllClient("test", g.Map{
		"wsPublishMsg": "hello ws client",
	})
}
