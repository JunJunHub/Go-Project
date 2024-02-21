package utils

import (
	"fmt"
	"github.com/go-ping/ping"
	"github.com/gogf/gf/frame/g"
	"math"
	"net"
	"strings"
	"time"
)

type Rectangle struct {
	X int
	Y int
	W int
	H int
}

//IsRectangleOverlap
//@summary 判断两个矩形是否重叠
func IsRectangleOverlap(rect1, rect2 Rectangle) bool {
	//窗口脱出大屏区域时,横坐标值为负值,会导致判断是否重叠失效;此处特殊处理
	if rect1.X < 0 {
		rect1.X = 0
	}
	if rect2.X < 0 {
		rect2.X = 0
	}

	//相交区域对角线上的两个点坐标
	minX := math.Max(float64(rect1.X), float64(rect2.X))
	minY := math.Max(float64(rect1.Y), float64(rect2.Y))
	maxX := math.Min(float64(rect1.X+rect1.W), float64(rect2.X+rect2.W))
	maxY := math.Min(float64(rect1.Y+rect1.H), float64(rect2.Y+rect2.H))
	width := maxX - minX
	height := maxY - minY
	if width <= 0 || height <= 0 {
		return false
	}
	return true
}

//ParseAcceptLanguage
//@summary 解析客户端接收语言
//@param1  reqAcceptLanguage string "请求头中描述的accept language"
func ParseAcceptLanguage(reqAcceptLanguage string) string {
	//精准匹配
	switch reqAcceptLanguage {
	case "zh_CN", "zh-CN":
		return "zh_CN"
	case "zh_TW":
		return "zh_TW"
	case "en":
		return "en"
	case "ja":
		return "ja"
	case "ru":
		return "ru"
	}

	//优先返回中文、其次是英文
	if strings.Contains(reqAcceptLanguage, "zh_CN") || strings.Contains(reqAcceptLanguage, "zh-CN") {
		return "zh_CN"
	}
	if strings.Contains(reqAcceptLanguage, "en") {
		return "en"
	}
	if strings.Contains(reqAcceptLanguage, "zh_TW") {
		return "zh_TW"
	}

	//默认返回 zh_CN
	return "zh_CN"
}

//ServerPing
//@summary ping探测某个IP是否可ping通
func ServerPing(ip string) bool {
	pinger, err := ping.NewPinger(ip)
	if err != nil {
		g.Log().Error(err)
	}
	pinger.Count = 1
	pinger.Timeout = 100 * time.Millisecond
	pinger.SetPrivileged(true)
	_ = pinger.Run() // blocks until finished
	stats := pinger.Statistics()

	//有回包，就是说明IP是可用的
	return stats.PacketsRecv >= 1
}

//ProbeTCPPort
//@summary 探测 TCP端口能否连通连接
func ProbeTCPPort(host string, port int, timeout time.Duration) bool {
	address := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", address, timeout)
	if err != nil {
		g.Log().Error(err)
		return false
	}
	defer func() { _ = conn.Close() }()
	return true
}
