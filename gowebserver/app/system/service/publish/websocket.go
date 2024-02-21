package publish

import (
	"context"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"github.com/gorilla/websocket"
	"gowebserver/app/common/define"
	"gowebserver/app/common/global"
	"sync"
	"time"
)

var WsPublish = wsPublish{}

const publishMsgChnCacheMax = 500 //发布消息最大缓冲条数
const proPublishGoroutineMax = 4  //发布消息处理协程个数

//WsConnect Ws连接对象
type WsConnect struct {
	Token string
	Conn  *ghttp.WebSocket
}

//WsConnectCloseHandel mpuaps自定义Ws链路关闭处理函数
func (wc *WsConnect) WsConnectCloseHandel(code int, text string) error {
	g.Log().Noticef("订阅链路断开连接 %s %s\nwsCloseCode: %d, text: %s",
		wc.Conn.RemoteAddr().String(), wc.Token, code, text)

	message := websocket.FormatCloseMessage(code, "ws close")
	_ = wc.Conn.WriteControl(websocket.CloseMessage, message, time.Now().Add(time.Second))

	WsPublish.DelWsLink(wc.Token, wc.Conn.RemoteAddr().String(), false)
	return nil
}

//NewWsConnect 创建ws连接对象
func NewWsConnect(token string, gfWsConn *ghttp.WebSocket) WsConnect {
	mpuapsWsConn := WsConnect{
		Token: token,
		Conn:  gfWsConn,
	}
	gfWsConn.SetCloseHandler(mpuapsWsConn.WsConnectCloseHandel)

	//接收服务端保活响应,如果超时未收到客户端保活应答,删除客户端订阅链路
	go func() {
		lastRecvAlivePackTimeMs := time.Now().UnixNano() / 1e6
		bConnClose := false

		//判断保活是否超时
		go func() {
			for {
				if bConnClose {
					return
				}

				//todo 方便调试,可修改保活时间改为20次
				if time.Now().UnixNano()/1e6-lastRecvAlivePackTimeMs > 5*global.UserLoginKeepaliveIntervalTimeS*1000 {
					//(30s)超过保活时间未收到保活应答,主动关闭链路
					g.Log().Noticef("订阅链路超时退出 %s %s", gfWsConn.RemoteAddr().String(), token)
					WsPublish.DelWsLink(token, gfWsConn.RemoteAddr().String(), true)
					return
				}
				time.Sleep(5 * time.Second)
			}
		}()

		//接收客户端发送消息
		for {
			_, msg, err := gfWsConn.ReadMessage()
			if err != nil {
				g.Log().Errorf("读取链路消息错误,订阅链路已断开连接 %s %s\n%v", gfWsConn.RemoteAddr().String(), token, err)
				bConnClose = true
				WsPublish.DelWsLink(token, gfWsConn.RemoteAddr().String(), true)
				return
			}
			jsonObj, err := gjson.DecodeToJson(msg)
			if err != nil {
				g.Log().Error(err)
			}
			if jsonObj.GetString("topic") == define.ENotifyAlive.String() {
				g.Log().Debugf("wsSubLinkAlive[%s]\n%v", gfWsConn.RemoteAddr().String(), jsonObj.Export())
				lastRecvAlivePackTimeMs = time.Now().UnixNano() / 1e6
			}
		}
	}()
	return mpuapsWsConn
}

//WSClientInfo 订阅客户端,一个客户端可以有多条ws链路(对应web多个页面)
type WSClientInfo struct {
	Token           string               //连接token
	IsCascadeSubChn bool                 //标识是否是级联连接
	CascadeId       string               //级联ID(级联链路记录对端级联ID)
	WSSubLinkMap    map[string]WsConnect //客户端ws连接Map表 key=addr value=WsConnect
}

//wsPublish 显控发布消息管理业务模型,处理需要发送给客户端的消息
type wsPublish struct {
	bInit                   bool                               //业务初始化
	chPubMsg                chan PubCommMsg                    //发布消息缓冲通道
	clientList              map[string]WSClientInfo            //发起订阅的ws客户端列表 key=token value=WSClientInfo
	mutexClientList         sync.Mutex                         //客户端map表读写锁
	disConnectHandleFuncMap map[string]ClientOfflineHandelFunc //离线处理接口
}

//InitInstance websocket通知业务逻辑初始化
func (s *wsPublish) InitInstance() {
	if s.bInit {
		return
	}

	//初始化客户端列表
	if s.clientList == nil {
		s.clientList = make(map[string]WSClientInfo)
	}

	//初始化发布消息通道
	if s.chPubMsg == nil {
		s.chPubMsg = make(chan PubCommMsg, publishMsgChnCacheMax)
	}

	//创建通知消息处理协程
	for i := 0; i < proPublishGoroutineMax; i++ {
		go s.proPubMsg(s.chPubMsg)
	}

	//创建订阅通道保活探测协程
	go s.keepAlive()

	s.bInit = true
}

// ClientOfflineHandelFunc
// @summary 其他模块注册的离线处理接口函数类型
//          注意：各模块注册的离线处理接口不能阻塞
// @param1  context.Context
// @param2  WSClientInfo
type ClientOfflineHandelFunc func(context.Context, WSClientInfo)

// SetWsClientDisConnectHandelFunc
// @summary 注册客户端离线处理接口
// @param1  module string "模块名"
// @param2  handleFunc ClientOfflineHandelFunc "离线处理接口"
// @return  nil
func (s *wsPublish) SetWsClientDisConnectHandelFunc(moduleName string, handleFunc ClientOfflineHandelFunc) {
	if moduleName == "" || handleFunc == nil {
		g.Log().Warning("参数为nil")
	}

	// 初始化订阅消息回调函数与对应的模块的映射
	if s.disConnectHandleFuncMap == nil {
		s.disConnectHandleFuncMap = make(map[string]ClientOfflineHandelFunc)
	}
	s.disConnectHandleFuncMap[moduleName] = handleFunc
}

//DebugPrintWsSubLinkInfo 打印所有Ws订阅客户端信息
func (s *wsPublish) DebugPrintWsSubLinkInfo() {
	if len(s.clientList) == 0 {
		g.Log().Info("ws订阅客户端信息为空")
		return
	}

	s.mutexClientList.Lock()
	g.Log().Info("===================================")
	for _, wsClient := range s.clientList {
		g.Log().Info("clientToken:", wsClient.Token)
		g.Log().Info("isCascadeLink:", wsClient.IsCascadeSubChn)
		g.Log().Info("cascadeId:", wsClient.CascadeId)
		g.Log().Info("wsLinkNum:", len(wsClient.WSSubLinkMap))
		for idx, wsConnect := range wsClient.WSSubLinkMap {
			g.Log().Infof("wsLinkAddr[%s] wsHandel: %v", idx, wsConnect)
		}
		g.Log().Info("===================================")
	}
	s.mutexClientList.Unlock()
}

//proPubMsg 通知所有已订阅的客户端
func (s *wsPublish) proPubMsg(chPublishMsg chan PubCommMsg) {
	defer g.Log().Warning("发送订阅通知协程退出")

	for {
		pubMsg := <-chPublishMsg

		switch pubMsg.PubMode {
		case PubModeToSpecifiedClient:
			s.mutexClientList.Lock()
			for tokenKey := range s.clientList {
				if s.clientList[tokenKey].Token == pubMsg.SpecialClientToken {
					//发送消息给指定客户端所有的订阅连接
					for linkAddr := range s.clientList[tokenKey].WSSubLinkMap {
						if err := s.clientList[tokenKey].WSSubLinkMap[linkAddr].Conn.WriteJSON(pubMsg); err != nil {
							g.Log().Errorf("%v\ntoken=\nsendMsg=%s\n", err, tokenKey, pubMsg)
							s.DelWsLink(tokenKey, linkAddr, false)
						} else {
							g.Log(global.LoggerWebsocket).Debug("notify client:", linkAddr, pubMsg)
						}
					}
					break
				}
			}
			s.mutexClientList.Unlock()

		case PubModeOutOfSpecifiedClient:
			s.mutexClientList.Lock()
			for tokenKey := range s.clientList {
				if s.clientList[tokenKey].Token != pubMsg.SpecialClientToken {
					//发送消息给客户端所有的订阅连接,排除某一特定客户端
					for linkAddr := range s.clientList[tokenKey].WSSubLinkMap {
						if err := s.clientList[tokenKey].WSSubLinkMap[linkAddr].Conn.WriteJSON(pubMsg); err != nil {
							g.Log().Errorf("%v\ntoken=\nsendMsg=%s\n", err, tokenKey, pubMsg)
							s.DelWsLink(tokenKey, linkAddr, false)
						} else {
							g.Log(global.LoggerWebsocket).Debug("notify client:", linkAddr, pubMsg)
						}
					}
				}
			}
			s.mutexClientList.Unlock()

		default:
			//消息发送给所有订阅的客户端
			s.mutexClientList.Lock()
			for tokenKey := range s.clientList {
				//发送消息给所有客户端所有的订阅连接
				for linkAddr := range s.clientList[tokenKey].WSSubLinkMap {
					if err := s.clientList[tokenKey].WSSubLinkMap[linkAddr].Conn.WriteJSON(pubMsg); err != nil {
						g.Log().Errorf("%v\ntoken=\nsendMsg=%s\n", err, tokenKey, pubMsg)
						s.DelWsLink(tokenKey, linkAddr, false)
					} else {
						g.Log(global.LoggerWebsocket).Debug("notify client:", linkAddr, pubMsg)
					}
				}
			}
			s.mutexClientList.Unlock()
		}
	}
}

//PublishMsgToAllClient 发布消息并通知所有客户端
func (s *wsPublish) PublishMsgToAllClient(topic string, data ...interface{}) {
	notifyData := interface{}(nil)
	if len(data) > 0 {
		notifyData = data[0]
	}
	pubMsg := PubCommMsg{
		PubMode:   PubModeToAll,
		Topic:     topic,
		TimeStamp: time.Now().UnixNano() / 1e6,
		Data:      notifyData,
	}
	s.pushMsgToChnCache(pubMsg)
}

//PublishMsgToSpecifiedClient 发布消息并通知给特定客户端
func (s *wsPublish) PublishMsgToSpecifiedClient(clientToken, topic string, data ...interface{}) {
	notifyData := interface{}(nil)
	if len(data) > 0 {
		notifyData = data[0]
	}
	pubMsg := PubCommMsg{
		PubMode:            PubModeToSpecifiedClient,
		SpecialClientToken: clientToken,
		Topic:              topic,
		TimeStamp:          time.Now().UnixNano() / 1e6,
		Data:               notifyData,
	}
	s.pushMsgToChnCache(pubMsg)
}

//PublishMsgOutOfSpecifiedClient 发布消息并排除特定客户端
func (s *wsPublish) PublishMsgOutOfSpecifiedClient(clientToken, topic string, data ...interface{}) {
	notifyData := interface{}(nil)
	if len(data) > 0 {
		notifyData = data[0]
	}
	pubMsg := PubCommMsg{
		PubMode:            PubModeOutOfSpecifiedClient,
		SpecialClientToken: clientToken,
		Topic:              topic,
		TimeStamp:          time.Now().UnixNano() / 1e6,
		Data:               notifyData,
	}
	s.pushMsgToChnCache(pubMsg)
}

//消息写入待发送消息缓冲通道
func (s *wsPublish) pushMsgToChnCache(msg PubCommMsg) {
	if len(s.chPubMsg) == publishMsgChnCacheMax {
		g.Log().Warningf("websocket通知消息缓冲通道已满[%d]", publishMsgChnCacheMax)
		return
	} else if len(s.chPubMsg) > 20 {
		g.Log().Warningf("websocket通知消息缓冲通道有[%d]条数据待发送", len(s.chPubMsg))
		if len(s.chPubMsg) > 100 && (msg.Topic == define.ENotifyChnUpdate.String() || msg.Topic == define.ENotifyChnGroupMemUpdate.String()) {
			g.Log().Warningf("不向客户端发送通道更新通知") //todo
			return
		}
	}
	s.chPubMsg <- msg
}

//keepAlive
//@summary 定时向客户端发送保活探测
//         客户端收到保活探测包,会回复保活应答
//@return1 nil
func (s *wsPublish) keepAlive() {
	for {
		aliveMsg := PubCommMsg{
			Topic:     define.ENotifyAlive.String(),
			TimeStamp: time.Now().UnixNano() / 1e6,
			Data:      nil,
		}

		//保活探测消息发送给所有的客户端
		s.mutexClientList.Lock()
		for tokenKey := range s.clientList {
			//发送消息给客户端所有的订阅连接
			for linkAddr := range s.clientList[tokenKey].WSSubLinkMap {
				if err := s.clientList[tokenKey].WSSubLinkMap[linkAddr].Conn.WriteJSON(aliveMsg); err != nil {
					g.Log().Errorf("%v\ntoken=%s\nsendMsg=%v\n", err, tokenKey, aliveMsg)
					s.DelWsLink(tokenKey, linkAddr, false)
				} else {
					g.Log().Debug("notify client:", linkAddr, aliveMsg)
				}
			}
		}
		s.mutexClientList.Unlock()
		time.Sleep(time.Second * global.UserLoginKeepaliveIntervalTimeS)
	}
}

//AddWsLink
//@summary 客户端添加订阅链路,支持一个客户端创建多个ws订阅连接(web多个标签页)
//@param1  token string "客户端令牌"
//@param2  isCascadeSubChn bool "表示是否是级联连接订阅链路"
//@param3  cascadeId string "级联平台级联ID"
//@return  wsConn *ghttp.WebSocket "ws连接对象"
func (s *wsPublish) AddWsLink(token string, isCascadeSubChn bool, cascadeId string, wsConn *ghttp.WebSocket) {
	s.mutexClientList.Lock()
	defer s.mutexClientList.Unlock()

	if _, ok := s.clientList[token]; ok {
		s.clientList[token].WSSubLinkMap[wsConn.RemoteAddr().String()] = NewWsConnect(token, wsConn)
	} else {
		clientInfo := WSClientInfo{
			Token:           token,
			IsCascadeSubChn: isCascadeSubChn,
			CascadeId:       cascadeId,
			WSSubLinkMap:    make(map[string]WsConnect),
		}
		clientInfo.WSSubLinkMap[wsConn.RemoteAddr().String()] = NewWsConnect(token, wsConn)
		s.clientList[token] = clientInfo
	}

	if s.clientList[token].IsCascadeSubChn {
		g.Log().Noticef("检测到上级平台连接 cascadeId:%s, remoteAddr:%s, isCascadeSubChn: %t",
			s.clientList[token].CascadeId, wsConn.RemoteAddr(), s.clientList[token].IsCascadeSubChn)
	}
}

//DelWsLink 删除订阅客户端链路
func (s *wsPublish) DelWsLink(token string, linkAddr string, bLockClientList bool) {
	if bLockClientList {
		s.mutexClientList.Lock()
		defer s.mutexClientList.Unlock()
	}

	if client, clientOk := s.clientList[token]; clientOk {
		//各业务断链处理
		for modelName, handleFunc := range s.disConnectHandleFuncMap {
			if handleFunc != nil {
				g.Log().Noticef("[%s]处理ws订阅链路断开: %s %s", modelName, linkAddr, token)
				handleFunc(context.Background(), client)
			}
		}

		//删除对应链路信息
		if _, linkOk := s.clientList[token].WSSubLinkMap[linkAddr]; linkOk {
			_ = s.clientList[token].WSSubLinkMap[linkAddr].Conn.Close()
			delete(s.clientList[token].WSSubLinkMap, linkAddr)
			if len(s.clientList[token].WSSubLinkMap) == 0 {
				delete(s.clientList, token)
			}
		}
	}
}

//GetClient 获取订阅客户端
func (s *wsPublish) GetClient(token string) WSClientInfo {
	s.mutexClientList.Lock()
	defer s.mutexClientList.Unlock()

	return s.clientList[token]
}

//DelClient 删除订阅客户端
func (s *wsPublish) DelClient(token string, closeMsg ...string) {
	s.mutexClientList.Lock()
	defer s.mutexClientList.Unlock()

	g.Log().Noticef("断开客户端【%s】所有订阅连接", token)
	if client, clientOk := s.clientList[token]; clientOk {
		for _, link := range client.WSSubLinkMap {
			if len(closeMsg) > 0 {
				//关闭前发送关闭帧给客户端
				_ = link.Conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseGoingAway, closeMsg[0]), time.Now().Add(time.Second))
			}
			_ = link.Conn.Close()
		}
		//删除客户端订阅链路信息
		delete(s.clientList, token)
	}
}

//CheckClient 检查登录用户订阅客户端是否存在
func (s *wsPublish) CheckClient(token string) bool {
	s.mutexClientList.Lock()
	defer s.mutexClientList.Unlock()

	if _, ok := s.clientList[token]; ok {
		return true
	}
	g.Log().Noticef("=========== clientSubLink[%s] exit", token)
	return false
}
