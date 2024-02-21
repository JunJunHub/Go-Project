// ======================
// grpc client 实例管理
// ======================

package mpuGrpcClient

import (
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"gowebserver/app/common/global"
)

var (
	//Instances RpcClient实例
	Instances = gmap.NewStrAnyMap(true)
)

// Instance
// 根据rpc服务地址[获取/创建]一个rpcClient实例
// 参数: rpc服务地址
func Instance(addr ...string) *RpcClient {
	var key string
	if len(addr) > 0 && addr[0] != "" {
		key = addr[0]
	}
	return Instances.GetOrSetFuncLock(key, func() interface{} {
		return new(key)
	}).(*RpcClient)
}

// RemoveInstance 删除RPC实例
func RemoveInstance(addr string) {
	Instance(addr).unInitInstance()
	Instances.Remove(addr)
}

func logger() *glog.Logger {
	return g.Log(global.LoggerMpuSrv)
}
