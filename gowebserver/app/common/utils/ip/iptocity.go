package ip

import (
	"fmt"
	"github.com/gogf/gf/encoding/gcharset"
	"github.com/gogf/gf/encoding/gjson"
	"github.com/gogf/gf/frame/g"
	"time"
)

//GetCityByIp 根据IP地址查询对应城市
func GetCityByIp(ip string) string {
	if ip == "" {
		return ""
	}

	if ip == "[::1]" || ip == "127.0.0.1" {
		return fmt.Sprintf("IP[%s]", ip)
	}

	url := "http://whois.pconline.com.cn/ipJson.jsp?json=true&ip=" + ip
	bytes := g.Client().Timeout(200 * time.Millisecond).GetBytes(url)
	src := string(bytes)
	srcCharset := "GBK"
	tmp, _ := gcharset.ToUTF8(srcCharset, src)
	json, err := gjson.DecodeToJson(tmp)
	if err != nil {
		return fmt.Sprintf("IP[%s]", ip)
	}
	if json.GetInt("code") == 0 {
		addr := json.GetString("addr")
		return fmt.Sprintf("%s[%s]", addr, ip)
	} else {
		return fmt.Sprintf("IP[%s]", ip)
	}
}
