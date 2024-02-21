// ==============================
// 显控平台GRPC客户端
// ==============================

package mpuGrpcClient

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"gowebserver/app/mpusrv/model"
	"io"
	"mpuaps/library/mpurpc"
	"sync"
	"time"
)

// RpcMsgCBFunc 其他模块注册接收RPC订阅消息函数类型
type RpcMsgCBFunc func(context.Context, RpcClientNotify)

type RpcClient struct {
	bInit         bool                     // 客户端初始化标志
	bConnect      bool                     // 连接状态: true 连接成功, false 连接失败
	bSubState     bool                     // 订阅状态: true 订阅成功, false 订阅失败
	mutexConnect  sync.Mutex               // 连接状态锁
	mutexSubState sync.Mutex               // 订阅状态锁
	rpcConn       *grpc.ClientConn         // RPC连接存根
	rpcInterface  mpurpc.MspRpcSrvClient   // rpc接口
	chSubMsg      chan RpcClientNotify     // 订阅消息缓冲通道, 用于订阅消息接收携程与订阅消息处理携程通讯
	msgCBFuncMap  map[string]RpcMsgCBFunc  // 订阅消息回调函数与对应的模块映射
	rpcServerAddr string                   // rpc服务端地址
	mpuPlatform   *model.MpuAccessPlatform // rpc服务对应显控信息(级联关系)
}

func (s *RpcClient) String() string {
	return fmt.Sprintf("\nrpcServerAddr: %s\nrpcConnectState: %v\nbSubState: %v\n",
		s.rpcServerAddr,
		s.bConnect,
		s.bSubState)
}

type RpcContextKeyType string

// RpcClientContextKey Rpc上下文变量存储键名
const RpcClientContextKey RpcContextKeyType = "RpcClientContextKey"

// RpcClientContext 自定义Rpc上下文结构
type RpcClientContext struct {
	RpcServerAddr string // rpc服务地址
}

// RpcClientNotify 通知消息
type RpcClientNotify struct {
	RpcServerAddr string
	AnyNotify     *mpurpc.AnyNotify
}

// new creates and returns a RpcClient.
func new(rpcServerAddr string) *RpcClient {
	rpcClt := &RpcClient{
		rpcServerAddr: rpcServerAddr,
	}
	return rpcClt.initInstance(rpcServerAddr)
}

// InitInstance RPC客户端初始化
func (s *RpcClient) initInstance(rpcServerAddr string) *RpcClient {
	if s.bInit {
		return nil
	}
	if rpcServerAddr == "" {
		logger().Warning("RPC Server Addr is nil! RpcClient.InitInstance false!")
		return nil
	}
	s.bInit = true
	logger().Notice("RpcClient.InitInstance! ServerAddr:", rpcServerAddr)

	//连接并订阅通知
	s.connect()
	s.subscribe()

	// Rpc连接保活
	go s.rpcClientKeepAlive()
	return s
}

func (s *RpcClient) unInitInstance() {
	//关闭RPC连接
	//结束保活协程
	//结束订阅消息接收协程
	//结束订阅消息处理协程
	s.bInit = false
	_ = s.rpcConn.Close()
	logger().Notice("RpcClient.unInitInstance! ServerAddr:", s.rpcServerAddr)
}

// SetRpcMsgCBFunc 设置接收RPC消息回调函数
func (s *RpcClient) SetRpcMsgCBFunc(module string, cbFunc RpcMsgCBFunc) *RpcClient {
	if module == "" || cbFunc == nil {
		logger().Warning("参数为nil")
	}

	// 初始化订阅消息回调函数与对应的模块的映射
	if s.msgCBFuncMap == nil {
		s.msgCBFuncMap = make(map[string]RpcMsgCBFunc)
	}

	s.msgCBFuncMap[module] = cbFunc
	return s
}

// SetMpuPlatformInfo 设置Rpc连接对应显控平台信息
func (s *RpcClient) SetMpuPlatformInfo(mpuPlatformInfo *model.MpuAccessPlatform) *RpcClient {
	s.mpuPlatform = mpuPlatformInfo
	return s
}

// GetMpuPlatformInfo 获取对应显控平台信息
func (s *RpcClient) GetMpuPlatformInfo() *model.MpuAccessPlatform {
	return s.mpuPlatform
}

// SetConnectState 设置连接状态方法
func (s *RpcClient) SetConnectState(bState bool) {
	s.mutexConnect.Lock()
	defer s.mutexConnect.Unlock()

	s.bConnect = bState
}

// CheckConnectState 检查连接状态方法: true已连接, false未连接
func (s *RpcClient) CheckConnectState() bool {
	s.mutexConnect.Lock()
	defer s.mutexConnect.Unlock()

	return s.bConnect
}

// SetSubState 设置订阅状态方法
func (s *RpcClient) SetSubState(bState bool) {
	s.mutexSubState.Lock()
	defer s.mutexSubState.Unlock()

	s.bSubState = bState
}

// CheckSubState 检查订阅状态方法: true已连接, false未连接
func (s *RpcClient) CheckSubState() bool {
	s.mutexSubState.Lock()
	defer s.mutexSubState.Unlock()

	return s.bSubState
}

// Interface 用于包外部直接访问调用rpc接口
func (s *RpcClient) Interface() mpurpc.MspRpcSrvClient {
	return s.rpcInterface
}

// rpc连接保活, 第一次判断连接成功和第一次判断断开连接会把连接状态回调给上层
func (s *RpcClient) rpcClientKeepAlive() {
	var bCBConnected = true  //标识是否需要回调过连接成功的通知
	var bCBDisconnect = true //标识是否需要回调过断开连接的通知
	for {
		if !s.bInit {
			logger().Notice("销毁RPC连接对象,RPC链路保活协程结束")
			return
		}
		if !s.CheckConnectState() {
			if bCBDisconnect {
				//回调给上层断开连接的通知
				bCBConnected = true
				bCBDisconnect = false
				s.callBackMsg(RpcClientNotify{
					AnyNotify: &mpurpc.AnyNotify{
						Topic:  mpurpc.ENotifyType_emNotifyDisconnect,
						Detail: nil,
					},
					RpcServerAddr: s.rpcServerAddr,
				})
				logger().Notice("显控平台：", s.rpcServerAddr, "断开连接")
			}

			//尝试重连
			s.connect()
		} else if bCBConnected {
			//初次连接成功,判断显控业务是否初始化完成,初始化完成之后,给上层发送通知消息,开始同步显控数据
			bCBDisconnect = true
			bCBConnected = false

			//异步等待显控资源准备OK
			go func() {
				for {
					if !s.bInit {
						logger().Notice("销毁RPC连接对象,RPC链路保活协程结束")
						return
					}
					if !s.CheckConnectState() {
						return
					}
					//通过查询通道信息接口返回结果来判断显控状态
					//要有个能准确感知显控状态的接口!!!
					chnStream, err := s.Interface().ChnnlListQuery(context.Background(), &mpurpc.QueryRequest{
						Sn:     "",
						Id:     1, //只查询显控本级的通道信息
						Subid:  uint32(mpurpc.EChannelType_emChnnlTypeVideoOutput),
						Offset: 0,
						Size:   0, //查询全部
					})
					if err != nil {
						logger().Error(err, "等待显控业务初始化...")
						time.Sleep(time.Second * 5)
						continue
					}
					_, err = chnStream.Recv()
					if err != nil && err != io.EOF {
						logger().Error(err, "等待显控业务初始化...")
						time.Sleep(time.Second * 5)
						continue
					}

					//给上层回调显控初始化完成通知
					logger().Notice("显控平台：", s.rpcServerAddr, "初始化完成,开始同步显控数据")
					s.callBackMsg(RpcClientNotify{
						AnyNotify: &mpurpc.AnyNotify{
							Topic:  mpurpc.ENotifyType_emNotifyAlive,
							Detail: nil,
						},
						RpcServerAddr: s.rpcServerAddr,
					})
					return
				}
			}()
		}
		if !s.CheckSubState() && s.CheckConnectState() {
			s.subscribe()
		}
		time.Sleep(5 * time.Second)
	}
}

// connect gRPC客户端连接
func (s *RpcClient) connect() {
	// 检查连接状态
	if s.CheckConnectState() {
		logger().Line().Info("gRPC客户端已连接")
		return
	}

	// 检查rpc服务地址
	if s.rpcServerAddr == "" {
		logger().Panic("gRPC服务地址配置不能为空")
		return
	}

	var err error
	if s.rpcConn == nil {
		// 新建虚拟连接(设置grpc.WithInsecure, 否则没有证书会报错)
		s.rpcConn, err = grpc.Dial(s.rpcServerAddr, grpc.WithInsecure())
		if err != nil {
			logger().Line().Panic("grpc.Dial() FailJson:", err)
			return
		}
	}
	//s.rpcConn.Connect()
	if s.rpcInterface == nil {
		s.rpcInterface = mpurpc.NewMspRpcSrvClient(s.rpcConn)
		if s.rpcInterface == nil {
			logger().Line().Warning("Rpc connect err! Server Addr:", s.rpcServerAddr)
			return
		}
	}

	//测试rpc连接
	if _, err := s.rpcInterface.SayHello(context.Background(), &emptypb.Empty{}); err != nil {
		logger().Line().Warning("Rpc connect err! Server Addr:", s.rpcServerAddr)
		return
	}

	s.SetConnectState(true)
	logger().Notice("Rpc Connected! ServerAddr: ", s.rpcServerAddr)
}

// subscribe gRPC客户端订阅主业务消息
func (s *RpcClient) subscribe() {
	// 检查连接状态
	if !s.CheckConnectState() {
		logger().Line().Warning("gRPC客户端未连接")
		return
	}
	// 检查订阅状态
	if s.CheckSubState() {
		logger().Line().Info("已订阅服务端消息")
		return
	}
	// 订阅服务端所有消息
	resp, err := s.rpcInterface.SubScribe(context.Background(), &mpurpc.SubRequest{
		IsAll:  true,
		Topics: nil,
	})
	if err != nil {
		logger().Warning("SubScribe failed! ", err)
		return
	}
	s.SetSubState(true)

	// 初始化订阅通道
	if s.chSubMsg == nil {
		s.chSubMsg = make(chan RpcClientNotify, 20)
	}

	// 创建携程循环读取订阅消息
	go func() {
		for {
			if !s.bInit {
				logger().Notice("销毁RPC连接对象,订阅协程结束")
				return
			}
			msg, err := resp.Recv()
			if err != nil {
				logger().Info("RPC连接已断开,订阅协程结束. gRPCServerAddr:", s.rpcServerAddr, err)
				s.SetConnectState(false)
				s.SetSubState(false)
				break
			}
			notify := RpcClientNotify{
				AnyNotify:     msg,
				RpcServerAddr: s.rpcServerAddr,
			}
			s.chSubMsg <- notify
		}
	}()

	// 创建订阅消息处理协程
	go s.proSubScribeMsg(s.chSubMsg)

	logger().Notice("gRPC SubScribe Start")
}

// callBackMsg 通知订阅回调
func (s *RpcClient) callBackMsg(msg RpcClientNotify) {
	for module, callBackFunc := range s.msgCBFuncMap {
		logger().Debug("module:", module, "callBackFunc:", callBackFunc)
		if callBackFunc != nil {
			callBackFunc(s.rpcClientCtxInit(), msg)
		}
	}
}

// rpc上下文信息
func (s *RpcClient) rpcClientCtxInit() context.Context {
	return context.WithValue(context.Background(), RpcClientContextKey, RpcClientContext{
		RpcServerAddr: s.rpcServerAddr,
	})
}

// proSubScribeMsg 订阅消息处理接口
func (s *RpcClient) proSubScribeMsg(chSubMsg chan RpcClientNotify) {
	for {
		if !s.bInit {
			return
		}
		if !s.CheckConnectState() || !s.CheckSubState() {
			break
		}

		msg := <-chSubMsg
		switch msg.AnyNotify.GetTopic() {
		case mpurpc.ENotifyType_emNotifyAlive:
			logger().Debug("RPCServer: ", s.rpcServerAddr, msg.AnyNotify.Topic.String())
			continue
		default:
		}

		// 订阅消息回调给其他业务模块处理
		s.callBackMsg(msg)
	}
}
