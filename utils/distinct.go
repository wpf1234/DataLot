package utils

import "reflect"

//去重函数
func Distinct(tempItem []string) (newArr []string) {
	newArr = make([]string, 0)
	for i := 0; i < len(tempItem); i++ {
		repect := false
		for j := i + 1; j < len(tempItem); j++ {
			if tempItem[i] == tempItem[j] {
				repect = true
				break
			}
		}
		if !repect {
			newArr = append(newArr, tempItem[i])
		}
	}
	return
}

func Duplicate(a interface{}) (ret []interface{}) {
	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}
	return ret
}