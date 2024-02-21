// ============================================
// 应用状态管理：http服务状态、数据库状态、资源同步状态
// ============================================

package appstate

import (
	"github.com/gogf/gf/frame/g"
	"sync"
	"time"
)

type MpuApsStatusType uint

const (
	EMpuApsStatusHttpServer  MpuApsStatusType = 0x00000001 //Http服务状态
	EMpuApsStatusDBConnect   MpuApsStatusType = 0x00000002 //数据库连接状态
	EMpuApsStatusDBInit      MpuApsStatusType = 0x00000004 //数据表初始化状态
	EMpuApsStatusSyncMPURes  MpuApsStatusType = 0x00000008 //与显控平台数据同步状态(本级)
	EMpuApsStatusSyncReduMPU MpuApsStatusType = 0x00000010 //正在与冗余主控同步数据
)

var mpuApsStatus MpuApsStatusType //服务状态
var mutexMpuApsStatus sync.Mutex  //状态锁
var syncUmtChnWarnState uint      //同步Umt点位异常告警状态: 0无告警、1kafka堆积通知过多,数据更新不及时 ...

//初始化
func init() {
	SetMpuApsStatusOk(EMpuApsStatusHttpServer)
	syncUmtChnWarnState = 0
}

//GetMpuApsStatus 获取当前状态
func GetMpuApsStatus() MpuApsStatusType {
	mutexMpuApsStatus.Lock()
	defer mutexMpuApsStatus.Unlock()

	g.Log().Noticef("\n"+
		"程序状态码含义:\n"+
		"EMpuApsStatusHttpServer: 0x00000001, EMpuApsStatusDBConnect: 0x00000002, EMpuApsStatusDBInit: 0x00000004, EMpuApsStatusSyncMPURes: 0x00000008, EMpuApsStatusSyncReduMPU: 0x00000010\n"+
		"当前状态值: %d", mpuApsStatus)
	return mpuApsStatus
}

//SetMpuApsStatusOk 设置服务对应功能状态ok
func SetMpuApsStatusOk(statusType MpuApsStatusType) {
	mutexMpuApsStatus.Lock()
	defer mutexMpuApsStatus.Unlock()

	mpuApsStatus |= statusType
	g.Log().Noticef("\n"+
		"状态码：httpServerStatus: 1, DBConnectStatus: 2, DBInitStatus: 4, SyncMpuResource: 8, EMpuApsStatusSyncReduMPU: 16\n"+
		"%d ok, CurrentMpuapsStatus: %d", statusType, mpuApsStatus)
}

//SetMpuApsStatusErr 设置服务对应功能状态Err
func SetMpuApsStatusErr(statusType MpuApsStatusType) {
	mutexMpuApsStatus.Lock()
	defer mutexMpuApsStatus.Unlock()

	/**
	go 不支持取反符号~
	c.Flag &= ~status
	取反 ^status
	*/
	mpuApsStatus &= ^statusType

	g.Log().Noticef("\n"+
		"状态码：httpServerStatus: 1, DBConnectStatus: 2, DBInitStatus: 4, SyncMpuResource: 8, EMpuApsStatusSyncRedu: 16\n"+
		"%d err, CurrentMpuapsStatus: %d", statusType, mpuApsStatus)
}

//SetSyncUmtChnWarnState 更新网络点位同步异常告警状态
func SetSyncUmtChnWarnState(state uint) {
	syncUmtChnWarnState = state
}
func GetSyncUmtChnWarnState() (state uint) {
	return syncUmtChnWarnState
}

//CheckHttpServerIsOk 检查Http服务状态
func CheckHttpServerIsOk() bool {
	mutexMpuApsStatus.Lock()
	defer mutexMpuApsStatus.Unlock()

	return (mpuApsStatus & EMpuApsStatusHttpServer) == EMpuApsStatusHttpServer
}

//CheckDBConnectIsOk 检查数据库连接状态
func CheckDBConnectIsOk() bool {
	mutexMpuApsStatus.Lock()
	defer mutexMpuApsStatus.Unlock()

	return (mpuApsStatus & EMpuApsStatusDBConnect) == EMpuApsStatusDBConnect
}

//CheckDBInitIsOk 检查数据库初始化状态
func CheckDBInitIsOk() bool {
	mutexMpuApsStatus.Lock()
	defer mutexMpuApsStatus.Unlock()

	return (mpuApsStatus & EMpuApsStatusDBInit) == EMpuApsStatusDBInit
}

//CheckSyncMPUResourcesIsOk 检查资源是否准备就绪
func CheckSyncMPUResourcesIsOk() bool {
	mutexMpuApsStatus.Lock()
	defer mutexMpuApsStatus.Unlock()

	return (mpuApsStatus & EMpuApsStatusSyncMPURes) == EMpuApsStatusSyncMPURes
}

//CheckSyncReduMPUResIsOk 检查与冗余主控同步状态
func CheckSyncReduMPUResIsOk() bool {
	mutexMpuApsStatus.Lock()
	defer mutexMpuApsStatus.Unlock()

	return (mpuApsStatus & EMpuApsStatusSyncReduMPU) == EMpuApsStatusSyncReduMPU
}

//WaitDBConnect 等待数据库连接并初始化
func WaitDBConnect() {
	for {
		if CheckDBConnectIsOk() {
			break
		}
		g.Log().Info("wait db connect!")
		time.Sleep(1 * time.Second)
	}
}

//WaitDBInit 等待数据库初始化
func WaitDBInit() {
	for {
		if CheckDBInitIsOk() {
			break
		}
		g.Log().Info("wait db init!")
		time.Sleep(10 * time.Second)
	}
}

//WaitSyncMPUResourcesOk 等待应用代理同步显控平台资源
func WaitSyncMPUResourcesOk() {
	for {
		if CheckSyncMPUResourcesIsOk() {
			break
		}
		time.Sleep(1 * time.Second)
	}
}
