package cmn

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

// Fasthttp服务器结构体
type FasthttpServer struct {
	router       *fasthttprouter.Router
	server       *fasthttp.Server
	port         string
	beforeHandle GlobalBeforeRequestHandler
}

// 全局前置拦截器
type GlobalBeforeRequestHandler func(ctx *fasthttp.RequestCtx) bool

// 创建Fasthttp服务器对象
func NewFasthttpServer() *FasthttpServer {
	return &FasthttpServer{
		router: fasthttprouter.New(),
	}
}

// 注册全局前置拦截器（前置拦截器返回true时才会继续正常处理后续请求）
func (f *FasthttpServer) BeforeRequestHandle(beforeHandle GlobalBeforeRequestHandler) *FasthttpServer {
	f.beforeHandle = beforeHandle
	return f
}

// 注册POST方法的请求控制器
func (f *FasthttpServer) HandlePost(path string, handle fasthttp.RequestHandler) *FasthttpServer {
	f.router.POST(path, func(ctx *fasthttp.RequestCtx) {
		if f.beforeHandle == nil || f.beforeHandle(ctx) {
			handle(ctx)
		}
	})
	return f
}

// 注册GET方法的请求控制器
func (f *FasthttpServer) HandleGet(path string, handle fasthttp.RequestHandler) *FasthttpServer {
	f.router.GET(path, func(ctx *fasthttp.RequestCtx) {
		if f.beforeHandle == nil || f.beforeHandle(ctx) {
			handle(ctx)
		}
	})
	return f
}

// 注册指定方法的请求控制器
func (f *FasthttpServer) Handle(method string, path string, handle fasthttp.RequestHandler) *FasthttpServer {
	f.router.Handle(method, path, func(ctx *fasthttp.RequestCtx) {
		if f.beforeHandle == nil || f.beforeHandle(ctx) {
			handle(ctx)
		}
	})
	return f
}

// 设定服务端口
func (f *FasthttpServer) SetPort(port string) *FasthttpServer {
	f.port = port
	return f
}

// 设定服务配置项（参数中的Handler配置项将被忽略）
func (f *FasthttpServer) SetServer(server *fasthttp.Server) *FasthttpServer {
	f.server = server
	return f
}

// 启动服务
func (f *FasthttpServer) Start() error {

	// 端口未设定时默认为8080
	addr := ":8080"
	if !IsBlank(f.port) {
		addr = ":" + f.port
	}

	// 服务配置项未设定时，默认请求体最大500M
	if f.server == nil {
		f.server = &fasthttp.Server{
			MaxRequestBodySize: 500 * 1024 * 1024,
		}
	}

	// 使用内置路由器
	f.server.Handler = f.router.Handler

	// 启动服务
	return f.server.ListenAndServe(addr)
}
