//================================
// MCU配置管理
//=================================

package accessMCUAdapter

import (
	"sync"
)

//MCUConfigs
//@summary MCU配置信息管理
//         key   = MCU唯一标识(mcuKey)
//         value = MCU配置信息(MCUServiceConnParam)
type MCUConfigs map[string]MCUServiceConnParam

//internalMCUConfigs 内部MCU配置管理对象
var internalMCUConfigs struct {
	sync.RWMutex
	mcuConfigs MCUConfigs
}

//SetConfig
//@summary 设置MCU配置信息
//         此接口将会覆盖旧的配置
func SetConfig(config MCUConfigs) {
	defer instances.Clear()

	internalMCUConfigs.Lock()
	defer internalMCUConfigs.Unlock()
	internalMCUConfigs.mcuConfigs = config
}

//AddConfigNode
//@summary 添加MCU配置信息
func AddConfigNode(mcuKey string, mcu MCUServiceConnParam) {
	internalMCUConfigs.Lock()
	defer internalMCUConfigs.Unlock()
	if _, ok := internalMCUConfigs.mcuConfigs[mcuKey]; ok {
		logger().Errorf("mcuKey=%s repetition", mcuKey)
	} else {
		internalMCUConfigs.mcuConfigs[mcuKey] = mcu
	}
}

//DelConfigNode
//@summary 删除MCU配置信息
//         同时释放对应实例
func DelConfigNode(mcuKey string) {
	internalMCUConfigs.Lock()
	defer internalMCUConfigs.Unlock()
	if _, ok := internalMCUConfigs.mcuConfigs[mcuKey]; !ok {
		logger().Errorf("not find mcuConfigNode by mcuKey=%s ", mcuKey)
	} else {
		//删除并销毁对接实例
		InstanceDestroy(mcuKey)

		//删除MCU配置
		delete(internalMCUConfigs.mcuConfigs, mcuKey)
	}
}

//GetConfigNode
//@summary 检索并返回对应MCU配置信息
func GetConfigNode(mcuKey string) MCUServiceConnParam {
	internalMCUConfigs.Lock()
	defer internalMCUConfigs.Unlock()
	return internalMCUConfigs.mcuConfigs[mcuKey]
}

//IsConfigured
//@summary 校验是否已经配置MCU参数
func IsConfigured() bool {
	internalMCUConfigs.Lock()
	defer internalMCUConfigs.Unlock()
	return len(internalMCUConfigs.mcuConfigs) > 0
}
