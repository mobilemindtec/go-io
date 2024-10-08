package util

import (
	"reflect"
)

type FuncInfo struct {
	Type        reflect.Type
	ArgsCount   int
	Value       reflect.Value
	ArgsTypes   []reflect.Type
	ReturnCount int
	ReturnTypes []reflect.Type
}

func (this *FuncInfo) ArgType(i int) reflect.Type {
	return this.ArgsTypes[i]
}

func (this *FuncInfo) Call(vals []reflect.Value) []reflect.Value {
	return this.Value.Call(vals)
}

func (this *FuncInfo) CallEmpty() []reflect.Value {
	return this.Call([]reflect.Value{})
}

func (this *FuncInfo) HasTypedArgs() bool {
	return len(this.ArgsTypes) > 0
}

func NewFuncInfo(f interface{}) *FuncInfo {
	fnType := reflect.TypeOf(f)
	fnArgsCount := fnType.NumIn()
	fnValue := reflect.ValueOf(f)
	argsTypes := []reflect.Type{}
	fnRetCount := fnType.NumOut()
	retTypes := []reflect.Type{}
	for i := 0; i < fnArgsCount; i++ {
		argsTypes = append(argsTypes, fnType.In(i))
	}

	for i := 0; i < fnRetCount; i++ {
		retTypes = append(retTypes, fnType.Out(i))
	}

	return &FuncInfo{
		Type:        fnType,
		ArgsCount:   fnArgsCount,
		Value:       fnValue,
		ArgsTypes:   argsTypes,
		ReturnCount: fnRetCount,
		ReturnTypes: retTypes,
	}
}

func IsFunc(f interface{}) bool {
	return reflect.TypeOf(f).Kind() == reflect.Func
}
