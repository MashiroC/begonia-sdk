package begonia

import (
	"encoding/json"
	"github.com/MashiroC/begonia-rpc/conn"
	"github.com/MashiroC/begonia-rpc/entity"
	"github.com/MashiroC/begonia-rpc/util/log"
	uuid "github.com/satori/go.uuid"
	"time"
)

// handler有三种 分别是
// 本地cli请求远程
// 本地cli请求响应
// 处理远程请求

type Client struct {
	addr      string
	signCache []entity.SignEntity
	req       *RequestHandler
	resp      *ResponseHandler
	rc        *RemoteCallHandler
	conn      conn.Conn
	wait      chan bool
}

func (cli *Client) TestCall(service, fun string, param ...interface{}) (res interface{}, err error) {
	s, _ := cli.rc.service.Get(service)
	return s.do(fun, param)
}

type Callback = func(interface{}, error)

type CallbackEntity struct {
	model int
	data  interface{}
}

type CallbackChan = chan CallbackEntity

// New 创建客户端并监听端口
func Default(addr string) *Client {
	cli := &Client{
		addr:      addr,
		signCache: make([]entity.SignEntity, 0),
		req:       defaultRequestHandler(),
		resp:      defaultResponseHandler(),
		rc:        defaultRemoteCallHandler(),
		wait:      make(chan bool, 2),
	}
	cli.connectAndListen()
	return cli
}

func (cli *Client) connectAndListen() {
	c := cli.connection(cli.addr)
	cli.conn = c
	cli.req.conn = c
	go cli.listen()
}

// Sign 注册服务
func (cli *Client) Sign(name string, in interface{}) {
	fun := cli.rc.sign(name, in)

	e := entity.SignEntity{
		Name:   name,
		Fun:    fun,
		IsMore: false,
	}
	cli.signCache = append(cli.signCache, e)
	form := entity.SignForm{Sign: []entity.SignEntity{e}}
	_ = cli.conn.WriteSign(form)
}

func (cli *Client) Wait() {
	<-cli.wait
}

func (cli *Client) KeepConnect() {

	for {
		<-cli.wait

		time.Sleep(3 * time.Second)
		ok := Must(cli.connectAndListen, cli.wait)
		if ok {
			Must(cli.reSign, cli.wait)
		}

	}

}

// reSign 断开连接后重新注册服务
func (cli *Client) reSign() {
	if len(cli.signCache) != 0 {
		form := entity.SignForm{Sign: cli.signCache}

		_ = cli.conn.WriteSign(form)
	}
}

// RemoteService 远程服务
type RemoteService struct {
	cli     *Client
	Service string
}

// 远程函数
type RemoteFun = func(param ...interface{}) (interface{}, error)

// Service 构造一个服务
func (cli *Client) Service(s string) RemoteService {
	return RemoteService{
		cli:     cli,
		Service: s,
	}
}

func (s RemoteService) FunTest(f string) RemoteFun {
	return func(param ...interface{}) (interface{}, error) {
		b, _ := json.Marshal(param)
		var ins []interface{}
		json.Unmarshal(b, &ins)
		return s.cli.TestCall(s.Service, f, ins...)
	}
}

// FunAsync 返回一个异步函数
func (s RemoteService) FunAsync(f string) func(Callback, ...interface{}) {
	return func(callback Callback, i ...interface{}) {
		req := Request{
			Service: s.Service,
			Fun:     f,
			Param:   i,
		}
		s.cli.callAsync(req, callback)
	}
}

// FunSync 返回一个同步函数
func (s RemoteService) FunSync(f string) RemoteFun {
	return func(param ...interface{}) (interface{}, error) {
		req := Request{
			Service: s.Service,
			Fun:     f,
			Param:   param,
		}
		return s.cli.callSync(req)
	}
}

// call 同步调用
// 同步调用其实是异步调用的包装
// 使用异步调用 在回调里向管道写数据 然后在另一边从管道拿数据
func (cli *Client) callSync(req Request) (res interface{}, err error) {
	ch := make(chan interface{}, 1)

	// 使用异步调用 回调是向这个管道写数据 然后调用协程收到数据 返回
	cli.callAsync(req, func(res interface{}, err error) {
		if err != nil {
			ch <- err
		} else {
			ch <- res
		}
	})

	tmp := <-ch

	if err, ok := tmp.(error); ok {
		return nil, err
	}

	return tmp, nil
}

// CallAsync 异步调用
func (cli *Client) callAsync(req Request, cb Callback) {

	// 弄个uuid
	uuid := uuid.NewV4().String()

	// 先注册回调
	// 如果后注册回调 高并发情况下协程调度 可能先返回再注册回调成功
	cbCh, err := cli.resp.signCallback(uuid)
	if err != nil {
		cb(nil, err)
	}

	if err := cli.req.call(uuid, req); err != nil {
		//TODO: handler Err & 解绑回调
		cb(nil, err)
	}

	// 等待回调
	go func(cbCh CallbackChan) {
		resp := <-cbCh

		if resp.model == entity.ErrorResponse {
			cb(nil, resp.data.(error))
		} else if resp.model == entity.NormalResponse {
			cb(resp.data, nil)
		} else {
			log.Fatal("what`s your problem")
		}
	}(cbCh)

}

func Must(fun func(), ch chan bool) (res bool) {
	defer func() {
		if re := recover(); re != nil {
			log.Warn("recover something : %s", re)
			ch <- true
			res = false
		}
	}()
	fun()
	res = true
	return
}
