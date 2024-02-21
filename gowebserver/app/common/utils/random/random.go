package random

import (
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gconv"
	"github.com/xiaojunjun-li/snowflake"

	"fmt"
	"math/rand"
	"strings"
	"time"
)

var sf *snowflake.SnowFlake

func init() {
	workerId := gconv.Uint32(g.Cfg().GetString("currentPlatform.PlatSn"))
	if workerId == 0 {
		workerId = 1
		cfgFilePath, err := g.Cfg().GetFilePath()
		if err != nil {
			g.Log().Line().Error(err)
		}
		g.Log().Warningf("%s未配置SnowFlake workerId, set default workerId=%d", cfgFilePath+g.Cfg().GetFileName(), workerId)
	}
	g.Log().Line().Infof("SnowFlake workerId=%d", workerId)

	var err error
	sf, err = snowflake.New(workerId)
	if err != nil {
		panic(err)
	}
}

//GenerateU64 生成全局唯一ID(uint64) -- snowflake算法
func GenerateU64() uint64 {
	id, err := sf.Generate()
	if err != nil {
		g.Log().Line().Error(err)
	}
	return id
}

//GenValidateCode 生成指定位数的数字
func GenValidateCode(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())

	var sb strings.Builder
	for i := 0; i < width; i++ {
		fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

//GenerateSubId 生成指定位数的字符
func GenerateSubId(width int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	rand.Seed(time.Now().UnixNano())

	b := make([]rune, width)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
