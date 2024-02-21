package accessMCUAdapter

import (
	"os"
	"testing"
)

//测试登录会议平台
func TestConnMCU_V6_0(t *testing.T) {
	//MCU配置信息
	mcuConfigs := make(MCUConfigs)
	mcuConnParam := MCUServiceConnParam{
		AdapterVer:          EMCUAdapterV6_0,
		Ip:                  "10.67.24.84",
		Port:                80,
		OAuthConsumerKey:    "snZymbs4kjUE",
		OAuthConsumerSecret: "kBQhLgPUzkha",
		SSLCertFile:         "",
		SSLKeyFile:          "",
		MCUUsers:            nil,
	}
	mcuConnParam.MCUUsers = append(mcuConnParam.MCUUsers, MCUUserInfo{
		Username:   "admin1",
		Password:   "keda8888",
		UserDomain: "",
	})
	mcuKey := mcuConnParam.Ip
	mcuConfigs[mcuKey] = mcuConnParam

	//设置MCU配置
	SetConfig(mcuConfigs)

	//构建MCU适配器实例
	mcuApiAdapter, err := newAdapter(mcuKey)
	if err != nil {
		t.Log(err)
		os.Exit(-1)
	}

	//查询MCU版本信息
	if ver, e := mcuApiAdapter.MCUApiSystemGetVersion(); e != nil {
		os.Exit(-1)
	} else {
		t.Log("MCUVersion:", ver)
	}

}
