// Time : 2019/10/12 10:21
// Author : MashiroC

// begonia
package begonia

import (
	"errors"
	"github.com/MashiroC/begonia-rpc/util/log"
	"reflect"
)

// remotecall_reflect.go something

func parseNum(num float64, kind string) (res reflect.Value) {
	var tmp interface{}
	switch kind {
	case "int":
		tmp = int(num)
	case "uint":
		tmp = uint(num)
	case "int8":
		tmp = int8(num)
	case "uint8":
		tmp = uint8(num)
	case "int16":
		tmp = int16(num)
	case "uint16":
		tmp = uint16(num)
	case "int32":
		tmp = int32(num)
	case "uint32":
		tmp = uint32(num)
	case "int64":
		tmp = int64(num)
	case "uint64":
		tmp = uint64(num)
	case "float32":
		tmp = float32(num)
	case "float64":
		tmp = num
	}
	res=reflect.ValueOf(tmp)
	return
}

func parseParam(inType []reflect.Type, obj []interface{}) (res []reflect.Value, err error) {
	res = make([]reflect.Value, len(obj))

	// 检查每个参数 针对各种类型做处理
	for i, param := range obj {
		in := inType[i]
		kind := in.Kind().String()
		if num, ok := param.(float64); ok {
			res[i] = parseNum(num, kind)
			continue
		}

		if kind == "struct" {
			// struct会变成map[string]interface{} 这里变回结构体
			value,er:=parseStruct(param,in)
			if er!=nil{
				err =er
				return
			}
			res[i]=value
			continue
		}

		res[i]=reflect.ValueOf(param)


	}

	return
}

func parseStruct(param interface{}, inType reflect.Type) (value reflect.Value, err error) {
	v := reflect.ValueOf(param)
	if v.Kind().String() != "map" {
		err = errors.New("param not struct")
		return
	}

	obj := reflect.New(inType)

	for i := 0; i < inType.NumField(); i++ {
		var fieldValue reflect.Value
		f := inType.Field(i)
		name := f.Type.Kind().String()
		if name == "string" {
			fieldValue = parseStringField(param.(map[string]interface{}), f)
		} else if name == "struct" {
			tmp := param.(map[string]interface{})
			key := f.Name
			if tag := f.Tag.Get("json"); tag != "" {
				key = tag
			}
			fieldValue, err = parseStruct(tmp[key], f.Type)
			if err != nil {
				return
			}
		} else {
			tmp := param.(map[string]interface{})
			key := f.Name
			if tag := f.Tag.Get("json"); tag != "" {
				key = tag
			}
			fieldValue = parseNum(tmp[key].(float64), f.Type.Name())
		}

		obj.Elem().Field(i).Set(fieldValue)
	}

	return obj.Elem(), nil
}

func parseToInterface(value []reflect.Value) (res []interface{}) {
	res = make([]interface{}, len(value))

	for i, v := range value {
		res[i] = v.Interface()
	}
	return
}

// parseResult 根据返回值的model来处理
// 如果调用的函数除了error外有多个返回值
// 那么将返回值解析为数组
// 当调用的函数最后一个值为error的情况下，error则作为这个函数的error返回
func parseResult(model resModel, value []reflect.Value) (interface{}, error) {
	switch model {
	// 普通 就是有返回值 但最后一个不为error
	case normal:
		if len(value) == 1 {
			return value[0].Interface(), nil
		} else {
			return parseToInterface(value), nil
		}
	case normalWithErr:
		err := value[len(value)-1]
		if err.IsNil() {
			return parseResult(normal, value[:len(value)-1])
		} else {
			return nil, err.Interface().(error)
		}
	case empty:
		return true, nil
	case emptyWithErr:
		err := value[0]
		if err.IsNil() {
			return err.Interface(), nil
		} else {
			return nil, err.Interface().(error)
		}
	default:
		log.Fatal("some error:", model)
		return nil, nil
	}
}

func parseStringField(params map[string]interface{}, f reflect.StructField) reflect.Value {
	key := f.Name
	if tag := f.Tag.Get("json"); tag != "" {
		key = tag
	}
	in := params[key]
	return reflect.ValueOf(in)
}

func parseAnonymousField(params map[string]interface{}, v reflect.Value) {
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if f.Anonymous {
			parseAnonymousField(params, v.Field(i))
			continue
		}
		//value := parseStringField(params, f)
	}
}
