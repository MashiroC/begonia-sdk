// Time : 2019/10/12 8:55
// Author : MashiroC

// begonia
package begonia

import (
	"github.com/MashiroC/begonia-rpc/entity"
	"reflect"
	"sync"
)

// remotecall_service.go something

// service 注册的服务
type service struct {
	name string
	fun  []*remoteFun
	in   reflect.Value
}

// do 调用一个服务
// fun 函数名 in 输入参数
// 返回结果或者err
func (s *service) do(fun string, in []interface{}) (res interface{}, err error) {

	// 找到调用的func
	var f *remoteFun
	for _, v := range s.fun {
		if v.name == fun {
			f = v
			break
		}
	}

	if f != nil {
		return f.do(s.in, in)
	}
	return nil, entity.FunctionNotFoundErr
}

type resModel uint8

const (
	normal resModel = iota + 1
	normalWithErr
	empty
	emptyWithErr
)

// remoteFun 远程函数
// 这个是注册好的函数
type remoteFun struct {
	name     string         // 函数名
	in       []reflect.Type // 函数输入的类型
	resModel resModel       // 是否函数返回值最后一个参数是err
	fun      reflect.Value  // 函数的值
}

// do 调用一个函数
// 这个是服务中的函数
func (rf *remoteFun) do(value reflect.Value, obj []interface{}) (interface{}, error) {
	// 检查输入参数长度
	if len(rf.in) != len(obj) {
		err := entity.CallError{
			ErrCode:    "114514",
			ErrMessage: "params num failed",
		}
		return nil, err
	}

	// 这里是核心部分 使用反射调取服务
	// 首先检查请求参数 把float64转为int
	// 然后再把map[string]interface{}转回struct
	// in是检查和解析后的参数
	in, err := parseParam(rf.in, obj)
	if err != nil {
		return nil, err
	}
	// 调用就完事了
	res, err := rf.call(value, in)
	return res, err
}

// call 调函数
func (rf *remoteFun) call(in reflect.Value, params []reflect.Value) (interface{}, error) {
	//rf.fun.Call(values)
	values := make([]reflect.Value, 1, 1+len(params))
	values[0] = in
	values = append(values, params...)
	tmp := rf.fun.Call(values)
	res, err := parseResult(rf.resModel, tmp)
	return res, err
}

// serviceMap 存服务的实体 并发安全
// key是服务名 value是服务实体
type serviceMap struct {
	data map[string]*service
	lock sync.Mutex
}

// newServiceMap 构造函数
func newServiceMap(len uint) *serviceMap {
	m := make(map[string]*service, len)
	return &serviceMap{
		data: m,
		lock: sync.Mutex{},
	}
}

// Get 拿服务
func (s *serviceMap) Get(key string) (v *service, ok bool) {
	s.lock.Lock()
	defer s.lock.Unlock()
	v, ok = s.data[key]
	return
}

// Set 加服务
func (s *serviceMap) Set(k string, v *service) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.data[k] = v
}

// Remove 移除服务
func (s *serviceMap) Remove(k string) {
	s.lock.Lock()
	defer s.lock.Unlock()
	delete(s.data, k)
}
