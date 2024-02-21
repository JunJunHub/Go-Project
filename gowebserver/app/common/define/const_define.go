//公共常量定义

package define

//收发盒用途
//ETRxFunc 收发盒用途枚举
type ETRxFunc int32

const (
	ETRxFuncForRelay    ETRxFunc = 0 //继电器
	ETRxFuncForPaiPai   ETRxFunc = 1 //拍一拍
	ETRxFuncForKeyBoard ETRxFunc = 2 //键盘
	ETRxFuncForKvmOnOff ETRxFunc = 3 //坐席开关机

)

//继电器状态枚举
const (
	ERelayOn           int = 0 //关
	ERelayOff          int = 1 //开
	ERelayStateUnknown int = 2 //未知
)
