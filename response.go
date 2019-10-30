// Time : 2019/10/11 17:18
// Author : MashiroC

// begonia
package begonia

import (
	"github.com/MashiroC/begonia-rpc/entity"
)

// response.go something

type ResponseHandler struct {
	cbMap *WaitChan
}

func (h *ResponseHandler) signCallback(uuid string) (CallbackChan, error) {
	ch := make(CallbackChan, 1)
	h.cbMap.Set(uuid, func(resp interface{}, model int) {
		cbEntity := CallbackEntity{
			model: model,
			data:  resp,
		}
		ch <- cbEntity
	})
	return ch, nil
}

func (h *ResponseHandler) callback(uuid string, params interface{}) (err error) {
	f, ok := h.cbMap.Get(uuid)
	h.cbMap.Remove(uuid)

	if !ok {
		err = entity.CallbackNotSignedErr
		return
	}

	f(params, entity.NormalResponse)

	return
}

func (h *ResponseHandler) callbackErr(uuid string, cErr entity.CallError) (err error) {
	f, ok := h.cbMap.Get(uuid)
	if !ok {
		return entity.CallbackNotSignedErr
	}
	f(cErr, entity.ErrorResponse)

	return
}

func defaultResponseHandler() *ResponseHandler {
	return &ResponseHandler{cbMap: NewWaitChan(5)}
}
