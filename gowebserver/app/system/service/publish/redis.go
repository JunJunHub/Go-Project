// ==========================================================================
// 消息发布到redis服务器.
// ==========================================================================

package publish

import (
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/frame/g"

	"sync"
)

var RedisPublish = redisPublish{}

type redisPublish struct {
	bInit        bool            // 模块初始化
	bConnect     bool            // 连接状态: true 连接成功, false 连接失败
	mutexConnect sync.Mutex      // 连接状态锁
	redis        *gredis.Redis   // 连接对象
	chPubMsg     chan PubCommMsg // 发布消息缓冲通道
}

func (s *redisPublish) InitInstance() {
	if s.bInit {
		return
	}
	if s.chPubMsg == nil {
		s.chPubMsg = make(chan PubCommMsg, 20)
	}
	go s.proPubMsg()
	s.bInit = true
}

// Connect 连接Redis服务器
func (s *redisPublish) Connect() {
	var (
		group  = "mpuaps"
		config = gredis.Config{
			Host:            "127.0.0.1",
			Port:            6379,
			Db:              0,
			Pass:            "123456",
			MaxIdle:         10,
			MaxActive:       100,
			IdleTimeout:     100,
			MaxConnLifetime: 30,
			ConnectTimeout:  30,
			TLS:             false,
			TLSSkipVerify:   false,
		}
	)
	gredis.SetConfig(&config, group)
	if s.redis = gredis.Instance(group); s.redis == nil {
		g.Log().Line().Warning("Redis Publish Instance err! config:", group, config)
		return
	}

	for {
		if s.CheckConnected() {
			continue
		}
		gredis.Instance(group).Stats()

		s.SetConnectState(true)
		g.Log().Line().Info("Redis Publish Instance, config:", group, config)
	}
}

// SetConnectState 设置连接状态
func (s *redisPublish) SetConnectState(bState bool) {
	s.mutexConnect.Lock()
	defer s.mutexConnect.Unlock()

	s.bConnect = bState
}

// CheckConnected 检查连接状态: true已连接, false未连接
func (s *redisPublish) CheckConnected() bool {
	s.mutexConnect.Lock()
	defer s.mutexConnect.Unlock()

	return s.bConnect
}

//PublishMsg 发布消息
func (s *redisPublish) PublishMsg(pubMsg PubCommMsg) {
	s.chPubMsg <- pubMsg
}

//proPubMsg 通知所有已订阅的客户端
func (s *redisPublish) proPubMsg() {
	pubMsg := <-s.chPubMsg
	g.Log().Info(pubMsg)
}
