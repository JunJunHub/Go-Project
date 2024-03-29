package convert

import (
	"github.com/gogf/gf/text/gstr"
	"github.com/gogf/gf/util/gconv"

	"fmt"
	"strings"
)

//DynamicParamParse 自定义动态参数解析
func DynamicParamParse(include string) string {
	var params string
	paramList := strings.Split(include, ",")
	for idx, param := range paramList {
		if idx == 0 {
			params += fmt.Sprintf("t.%s", gstr.CaseSnake(param))
		} else {
			params += fmt.Sprintf(",t.%s", gstr.CaseSnake(param))
		}
	}
	return params
}

//ToInt64Array 将带分隔符的字符串切成int64数组
func ToInt64Array(str, split string) []int64 {
	result := make([]int64, 0)
	if str == "" {
		return result
	}

	arr := strings.Split(str, split)
	if len(arr) > 0 {
		for i := range arr {
			if arr[i] != "" {
				result = append(result, gconv.Int64(arr[i]))
			}
		}
	}

	return result
}

//ReplaceHeadAndEndStr 过滤首尾有分隔符的数字字符串
func ReplaceHeadAndEndStr(str, split string) string {
	result := ""
	arr := strings.Split(str, split)
	if len(arr) <= 0 {
		return result
	}

	for i := range arr {
		if arr[i] != "" {
			if i == 0 {
				result = arr[i]
			} else {
				result += "," + arr[i]
			}
		}
	}

	return result
}

// Capitalize 字符首字母大写
func Capitalize(str string) string {
	var upperStr string
	vv := []rune(str) // 后文有介绍
	for i := 0; i < len(vv); i++ {
		if i == 0 {
			if vv[i] >= 97 && vv[i] <= 122 { // 后文有介绍
				vv[i] -= 32 // string的码表相差32位
				upperStr += string(vv[i])
			} else {
				fmt.Println("Not begins with lowercase letter,")
				return str
			}
		} else {
			upperStr += string(vv[i])
		}
	}
	return upperStr
}
