package utils

import "sort"

// 结构体排序必须重写数组 Len(),Swap(),Less()
type bodyWrapper struct {
	Body []interface{}
	by func(p,q *interface{}) bool
}

type SortBodyBy func(p,q* interface{}) bool

// 数组长度
func (b bodyWrapper) Len() int  {
	return len(b.Body)
}

// 元素交换
func (b bodyWrapper) Swap(i,j int){
	b.Body[i],b.Body[i]=b.Body[j],b.Body[i]
}

func (b bodyWrapper) Less(i,j int)bool{
	return b.by(&b.Body[i],&b.Body[j])
}

// 自定义排序字段
func SortBody(body []interface{},by SortBodyBy){
	sort.Sort(bodyWrapper{body,by})
}
