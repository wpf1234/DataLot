package utils

import (
	"fmt"
	"strconv"
	"strings"
)

func Slice2Str(s interface{}) string {
	var str string
	if ps, ok := s.(*[]string); ok {
		for _, v := range *ps {
			str = str + v + ","
		}
	} else if ps, ok := s.(*[]int); ok {
		for _, v := range *ps {
			data := strconv.Itoa(v)
			str = str + data + ","
		}
	} else {
		fmt.Println("其它数据类型，暂不支持转换!")
		return ""
	}
	str = strings.TrimRight(str, ",")
	return str
}
