package model

import (
	"fmt"
	"reflect"
)

func Print(foo interface{}) string {
	var ret string = ""
	v := reflect.ValueOf(foo)
	switch v.Kind() {
	case reflect.Struct:
		ret += "{"
		for i := 0; i < v.NumField(); i++ {
			if reflect.Indirect(v.Field(i)).IsValid() {
				val := Print(reflect.Indirect(v.Field(i)).Interface())
				if val != "[]" && val != "" && val != "0" {
					ret += fmt.Sprintf("%v:%v, ", v.Type().Field(i).Name, val)
				}
			}
		}
		if len(ret) < 2 {
			return ""
		}
		ret = fmt.Sprintf("%v}", ret[:len(ret)-2])
		return ret
	case reflect.Slice:
		ret += "["
		for i := 0; i < v.Len(); i++ {
			if v.Index(i).IsValid() {
				val := Print(v.Index(i).Interface())
				ret += fmt.Sprintf("%v, ", val)
			}
		}
		if len(ret) < 2 {
			return ""
		}
		ret = fmt.Sprintf("%v]", ret[:len(ret)-2])
		return ret
	case reflect.Int:
		return fmt.Sprintf("%v", v.Int())
	// case reflect.Ptr:
	// 	ret += "*"
	// 	fallthrough
	default:
		ret += reflect.Indirect(v).String()
		return ret
	}
}

func (f GameFilter) String() string {
	return Print(f)
}

func (f PlayerFilter) String() string {
	return Print(f)
}

func (f TeamFilter) String() string {
	return Print(f)
}
