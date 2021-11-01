package iocsvc

import (
	"fmt"
	"reflect"

	underscore "github.com/ahl5esoft/golang-underscore"
)

const (
	instanceIsNotPtr       = "ioc: 注入实例必须是指针"
	invalidTypeFormat      = "ioc: 无效类型(Type = %v)"
	notImplementsFormat    = "ioc: %v没有实现%v"
	notInterfaceTypeFormat = "ioc: 非接口类型(%v)"
)

var instanceValues = make(map[reflect.Type]reflect.Value)

func Get(interfaceObj interface{}) interface{} {
	return getValue(interfaceObj).Interface()
}

func Inject(instance interface{}, selectorFunc func(reflect.Value) reflect.Value) {
	instanceValue := reflect.ValueOf(instance)
	if instanceValue.Kind() != reflect.Ptr {
		panic(instanceIsNotPtr)
	}

	inject(instanceValue, selectorFunc)
}

func Set(interfaceObj interface{}, instance interface{}) {
	interfaceType := getInterfaceType(interfaceObj)
	instanceType := reflect.TypeOf(instance)
	if !instanceType.Implements(interfaceType) {
		panic(
			fmt.Errorf(notImplementsFormat, instance, interfaceType),
		)
	}

	instanceValues[interfaceType] = reflect.ValueOf(instance)
}

func getInterfaceType(interfaceObj interface{}) reflect.Type {
	var interfaceType reflect.Type
	var ok bool
	if interfaceType, ok = interfaceObj.(reflect.Type); !ok {
		interfaceType = reflect.TypeOf(interfaceObj)
	}

	if interfaceType.Kind() == reflect.Ptr {
		interfaceType = interfaceType.Elem()
	}

	if interfaceType.Kind() != reflect.Interface {
		panic(
			fmt.Errorf(notInterfaceTypeFormat, interfaceType),
		)
	}

	return interfaceType
}

func getValue(interfaceObj interface{}) reflect.Value {
	interfaceType := getInterfaceType(interfaceObj)
	if v, ok := instanceValues[interfaceType]; ok {
		return v
	}

	panic(
		fmt.Errorf(invalidTypeFormat, interfaceType),
	)
}

func inject(instanceValue reflect.Value, selectorFunc func(reflect.Value) reflect.Value) {
	if instanceValue.Kind() == reflect.Ptr {
		instanceValue = instanceValue.Elem()
	}

	underscore.Range(
		0,
		instanceValue.Type().NumField(),
		1,
	).Each(func(r int, _ int) {
		field := instanceValue.Type().Field(r)
		fieldValue := instanceValue.FieldByIndex(field.Index)
		if field.Anonymous {
			if field.Type.Kind() == reflect.Struct {
				inject(fieldValue, selectorFunc)
			}
			return
		}

		if _, ok := field.Tag.Lookup("inject"); !ok {
			return
		}

		if fieldValue.Kind() == reflect.Ptr {
			value := reflect.New(
				field.Type.Elem(),
			)
			fieldValue.Set(value)
			fieldValue = fieldValue.Elem()
		}

		v := getValue(field.Type)
		if selectorFunc != nil {
			v = selectorFunc(v)
		}
		fieldValue.Set(v)
	})
}
