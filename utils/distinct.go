package utils

import (
	"fmt"
	"reflect"
)

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


/* 在slice中去除重复的元素，其中a必须是已经排序的序列。
 * params:
 *   a: slice对象，如[]string, []int, []float64, ...
 * return:
 *   []interface{}: 已经去除重复元素的新的slice对象
 */
func SliceRemoveDuplicate(a interface{}) (ret []interface{}) {
	if reflect.TypeOf(a).Kind() != reflect.Slice {
		fmt.Printf("<SliceRemoveDuplicate> <a> is not slice but %T\n", a)
		return ret
	}

	va := reflect.ValueOf(a)
	for i := 0; i < va.Len(); i++ {
		if i > 0 && reflect.DeepEqual(va.Index(i-1).Interface(), va.Index(i).Interface()) {
			continue
		}
		ret = append(ret, va.Index(i).Interface())
	}

	return ret
}