// Time : 2019/10/11 17:19
// Author : MashiroC

// begonia
package begonia

import (
	"errors"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
	"reflect"
)

// remotecall.go something

// RemoteCallHandler 用来处理远程调用 接受的是center发来的请求
type RemoteCallHandler struct {
	service *serviceMap
}

// defaultRemoteCallHandler 默认的handler
func defaultRemoteCallHandler() *RemoteCallHandler {
	return &RemoteCallHandler{service: newServiceMap(5)}
}

// sign 注册服务
// name是服务名 in是服务结构体的一个实例
func (rc *RemoteCallHandler) sign(name string, in interface{}) (res []entity.FunEntity) {
	funs := make([]reflect.Method, 0)

	v := reflect.ValueOf(in)
	t := reflect.TypeOf(in)
	num := t.NumMethod()
	for i := 0; i < num; i++ {
		m := t.Method(i)

		funs = append(funs, m)
		en := entity.FunEntity{
			Name: m.Name,
			Size: m.Type.NumIn()-1,
		}
		res = append(res, en)
	}
	rc.addService(name, v, funs)
	return
}

// addService 添加一个服务
func (rc *RemoteCallHandler) addService(name string, in reflect.Value, funs []reflect.Method) {
	rfs := make([]*remoteFun, len(funs))
	for i, fun := range funs {
		// 大小减一 并且下面从1开始
		// 因为从type拿到的方法 第一个参是服务结构体指针
		in := make([]reflect.Type, fun.Type.NumIn()-1)
		t := fun.Type

		for i := 1; i <= len(in); i++ {
			funIn := t.In(i)
			in[i-1] = funIn
		}
		resModel := normal
		numOut := t.NumOut()
		if numOut == 0 {
			resModel = empty
		} else {
			lastOut := t.Out(numOut - 1)
			if lastOut.Name() == "error" {
				if numOut == 1 {
					resModel = emptyWithErr
				} else {
					resModel = normalWithErr
				}
			}
		}

		rf := &remoteFun{
			name:     fun.Name,
			in:       in,
			resModel: resModel,
			fun:      fun.Func,
		}
		rfs[i] = rf
	}
	s := &service{
		name: name,
		fun:  rfs,
		in:   in,
	}
	rc.service.Set(name, s)
}

func (rc *RemoteCallHandler) call(service, fun string, param []interface{}) (res interface{}, err error) {
	// 捕获panic出来的错误
	defer func() {
		if re := recover(); re != nil {
			if reErr, ok := re.(error); ok {
				log.Error("recover error:%s", reErr.Error())
				err = reErr
				return
			}
			log.Error("recover unknown:%s", re)
			err = errors.New("unknown error")
		}
	}()
	s, _ := rc.service.Get(service)
	return s.do(fun, param)
}

//func (rc *RemoteCallHandler) call(){
//
//}