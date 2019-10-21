// Time : 2019/10/11 17:17
// Author : MashiroC

// begonia
package begonia

import (
	"github.com/MashiroC/begonia-rpc/conn"
	"github.com/MashiroC/begonia-rpc/entity"
)

// request.go something

// RequestHandler 处理本地请求
type RequestHandler struct {
	conn conn.Conn
}

type Request struct {
	Service string
	Fun     string
	Param   []interface{}
}

func (h *RequestHandler) call(uuid string, request Request) (err error) {
	req := entity.Request{
		UUID:    uuid,
		Service: request.Service,
		Fun:     request.Fun,
		Data:    request.Param,
	}

	err = h.conn.WriteRequest(req)
	return
}

func defaultRequestHandler() *RequestHandler {
	return &RequestHandler{}
}
