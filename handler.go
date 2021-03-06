// Time : 2019/10/12 19:08
// Author : MashiroC

// begonia
package begonia

import (
	"encoding/json"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
)

// handler.go something

// handlerRequest 处理远程调用请求帧
// 客户端收到的请求帧只能是被rpc调用的帧
func (cli *Client) handlerRequest(data []byte) {

	// 先检查data的json对不对 json不对直接关了
	req := entity.Request{}
	if err := json.Unmarshal(data, &req); err != nil {
		log.Error("request frame json [%s] error [%s] req addr: [%s]", string(data), err.Error(), cli.conn.Addr())
		cli.closeWith(err)
		return
	}

	// 检查一个帧要调用的service和function是否存在
	if req.Service == "" || req.Fun == "" {
		cli.respError(req.UUID, entity.ServiceNotFoundErr)
		return
	}

	log.Info("received [%s] call %s.%s", cli.conn.Addr(), req.Service, req.Fun)

	res, err := cli.rc.call(req.Service, req.Fun, req.Data)
	if err != nil {
		// TODO:Err

		resp := entity.RespForm{
			Uuid: req.UUID,
			Type: entity.ErrorResponse,
			Data: entity.FromError(err),
		}

		_ = cli.conn.WriteResponse(resp)
		return
	}

	resp := entity.RespForm{
		Uuid: req.UUID,
		Type: entity.NormalResponse,
		Data: res,
	}
	_ = cli.conn.WriteResponse(resp)
}

// handlerResponse 处理远程调用响应帧
// 客户端收到的响应帧只能是自己调用的响应
func (cli *Client) handlerResponse(data []byte) {
	// 先检查data的json对不对 json不对直接关了
	form := entity.RespForm{}
	if err := json.Unmarshal(data, &form); err != nil {
		log.Error("resp frame json [%s] error [%s] form addr: [%s]", string(data), err.Error(), cli.conn.Addr())
		cli.closeWith(err)
		return
	}

	// 先找uuid有没有 uuid没有就是没有注册回调
	if form.Uuid == "" {
		// Uuid not found 这个应该直接返回给这条连接
		cli.respError("", entity.CallbackNotSignedErr)
		return
	}
	if form.Type == entity.NormalResponse {

		// 普通回调
		if err := cli.resp.callback(form.Uuid, form.Data); err != nil {
			cli.respError(form.Uuid, err)
			return
		}

	} else if form.Type == entity.ErrorResponse {

		// 错误回调
		data := form.Data.(map[string]interface{})
		cErr := entity.NewError(data["errCode"].(string), data["errMsg"].(string))
		if err := cli.resp.callbackErr(form.Uuid, cErr); err != nil {
			cli.respError(form.Uuid, err)
			return
		}

	}

}

// handlerError 处理错误帧
// 这个错误帧指的是收到的error frame 不是异常帧
// 客户端收到的错误帧只能是rpc响应的错误
func (cli *Client) handlerError(data []byte) {
	// 先检查data的json对不对 json不对直接关了
	form := entity.ErrForm{}
	if err := json.Unmarshal(data, &form); err != nil {
		log.Error("error frame json [%s] error [%s] form addr: [%s]", string(data), err.Error(), cli.conn.Addr())
		cli.closeWith(err)
		return
	}

	// 这里和响应的逻辑基本一样 只不过回调传的是error
	if form.Uuid == "" {
		//cli.respError(entity.CallbackNotSignedErr)
		log.Error("fuck Uuid")
		return
	}

	// 回调的error
	cErr := entity.CallError{
		ErrCode:    form.ErrCode,
		ErrMessage: form.ErrMsg,
	}

	_ = cli.resp.callbackErr(form.Uuid, cErr)
	log.Error("received some error %s", cErr.Error())
}
