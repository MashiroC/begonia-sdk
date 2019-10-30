// Time : 2019/10/30 19:51
// Author : MashiroC

// begonia
package begonia

import (
	"errors"
	"reflect"
)

// result.go something

func Result(res interface{}, err error) rResult {
	return rResult{
		Res: res,
		Err: err,
	}
}

type rResult struct {
	Res interface{}
	Err error
}

func (r rResult) Bind(target interface{}) (err error) {
	t:=reflect.TypeOf(target)
	v:=reflect.ValueOf(target)
	if t.Kind()!=reflect.Ptr{
		return errors.New("target not a ptr")
	}
	err = parseMapToStruct(t.Elem(), r.Res.(map[string]interface{}), v)
	return
}
