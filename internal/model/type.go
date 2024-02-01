package model

import "reflect"

// 获取指针/切片/数组内的真实类型
func realType(v reflect.Type, i int) (reflect.Type, int) {
	switch v.Kind() {
	case reflect.Slice, reflect.Array:
		i++
		fallthrough
	case reflect.Ptr:
		return realType(v.Elem(), i)
	default:
		return v, i
	}
}
