package tabler

import (
	"reflect"
	"sync"
)

const columnsStr = "columns"

var columnsCache sync.Map

type Tabler interface {
	Columns() []string
}

func Columns(in interface{}) []string {
	if t, ok := in.(Tabler); ok {
		return t.Columns()
	}
	return getColumnsByTag(in)
}

func getColumnsByTag(in interface{}) []string {
	tp := reflect.TypeOf(in)
	cache, ok := columnsCache.Load(tp)
	if ok {
		return cache.([]string)
	}
	result := make([]string, 0)
	for i := 0; i < tp.NumField(); i++ {
		field := tp.Field(i)
		//unexported field
		if field.PkgPath != "" {
			continue
		}
		tag := field.Tag.Get(columnsStr)
		if tag == "-" {
			continue
		}
		if tag == "" {
			result = append(result, unCapitalize(field.Name))
			continue
		}
		result = append(result, tag)
	}
	columnsCache.Store(tp, result)
	return result
}

func unCapitalize(tag string) string {
	ss := []rune(tag)
	if ss[0] >= 'A' && ss[0] <= 'Z' {
		ss[0] = ss[0] + 32
	}
	return string(ss)
}
