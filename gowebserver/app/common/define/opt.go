package define

import "fmt"

//EOptType 通用操作类型枚举
type EOptType int32

const (
	EOptNone      EOptType = 0  //无效
	EOptAdd       EOptType = 1  //添加(创建)
	EOptModify    EOptType = 2  //修改(更新)
	EOptDel       EOptType = 3  //删除
	EOptUp        EOptType = 4  //上移
	EOptDown      EOptType = 5  //下移
	EOptTop       EOptType = 6  //置顶
	EOptBottom    EOptType = 7  //置底
	EOptClean     EOptType = 8  //清理
	EOptReset     EOptType = 9  //重置
	EOptLoad      EOptType = 10 //加载
	EOptUnLoad    EOptType = 11 //卸载
	EOptReload    EOptType = 12 //重载
	EOptOn        EOptType = 13 //开启
	EOptOff       EOptType = 14 //关闭
	EOptOnline    EOptType = 15 //上线
	EOptOffline   EOptType = 16 //离线
	EOptPush      EOptType = 17 //推送
	EOptTakeBack  EOptType = 18 //撤销
	EOptStart     EOptType = 19 //开始
	EOptStop      EOptType = 20 //停止
	EOptEnable    EOptType = 21 //启用
	EOptDisable   EOptType = 22 //禁用
	EOptRegroup   EOptType = 23 //重新分组
	EOptSwitchOn  EOptType = 24 //开
	EOptSwitchOff EOptType = 25 //关
)

func (x EOptType) String() string {
	switch x {
	case EOptNone:
		return "emNone"
	case EOptAdd:
		return "emAdd"
	case EOptModify:
		return "emModify"
	case EOptDel:
		return "emDel"
	case EOptUp:
		return "emUp"
	case EOptDown:
		return "emDown"
	case EOptTop:
		return "emTop"
	case EOptBottom:
		return "emBottom"
	case EOptClean:
		return "emClean"
	case EOptReset:
		return "emReset"
	case EOptLoad:
		return "emLoad"
	case EOptUnLoad:
		return "emUnLoad"
	case EOptReload:
		return "emReload"
	case EOptOn:
		return "emOn"
	case EOptOff:
		return "emOff"
	case EOptOnline:
		return "emOnline"
	case EOptOffline:
		return "emOffline"
	case EOptPush:
		return "emPush"
	case EOptTakeBack:
		return "emTakeBack"
	case EOptStart:
		return "emStart"
	case EOptStop:
		return "emStop"
	case EOptEnable:
		return "emEnable"
	case EOptDisable:
		return "emDisable"
	case EOptRegroup:
		return "emRegroup"
	case EOptSwitchOn:
		return "emSwitchOn"
	case EOptSwitchOff:
		return "emSwitchOff"
	default:
		return fmt.Sprintf("EOptType(%d) undefined, Please check file \"app/common/define/opt.go\"", x)
	}
}
