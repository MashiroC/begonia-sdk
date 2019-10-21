// Time : 2019/10/12 19:00
// Author : MashiroC

// begonia
package begonia

import (
	"fmt"
	"sync"
)

// waitchan.go something

// waitChan.go 根据uuid获得响应回调的map

// WaitChan 根据uuid获得响应的回调
// 并发安全
type WaitChan struct {
	data map[string]WaitFun
	lock sync.Mutex
}

type WaitFun = func(interface{},int)

// 构造函数
func NewWaitChan(len uint) *WaitChan {
	return &WaitChan{
		data: make(map[string]WaitFun, len),
		lock: sync.Mutex{},
	}
}

// Get 取
func (w *WaitChan) Get(k string) (callback WaitFun, ok bool) {
	w.lock.Lock()
	defer w.lock.Unlock()
	callback, ok = w.data[k]
	if !ok {
		fmt.Println(callback)
	}
	return
}

// Set 加
func (w *WaitChan) Set(k string, callback WaitFun) {
	w.lock.Lock()
	defer w.lock.Unlock()
	w.data[k] = callback
}

// Remove 删
func (w *WaitChan) Remove(k string) {
	w.lock.Lock()
	defer w.lock.Unlock()
	delete(w.data, k)
}