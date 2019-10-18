package lego

import (
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"sync"
	"errors"
)

var _router *fasthttprouter.Router
var _routerOnce sync.Once

// 如果是websocket的情况下，走这个处理
var _websockMaper map[string]fasthttp.RequestHandler

/* 获取httprouter指针
 */
func GetRouter() *fasthttprouter.Router {
	_routerOnce.Do(func() {
		_router = fasthttprouter.New()
		_router.PanicHandler = func(ctx *fasthttp.RequestCtx, i interface{}) {
			LogError(i)
			LogPanicTrace(8)
		}
		// websocket handeler
		_websockMaper = make(map[string]fasthttp.RequestHandler)
	})
	return _router
}

/* */
func GetRequestHandler(ctx *fasthttp.RequestCtx) {
	path := string(ctx.Path())
	if h, ok := _websockMaper[path]; ok {
		h(ctx)
	}else{
		_router.Handler(ctx)
	}
}

/* 注册路由，GET的方式
* parmas
	@path -- 访问路径
	@h 	  -- 处理handler
	@checktoken -- 是否校验token
*/
func GET(path string, h fasthttp.RequestHandler, sessionType int, checksession bool) {
	LogPrintln("Route Register: GET:\t" + path)
	if checksession {
		if sessionType == LEGO_SESSION_TYPE_ADMIN {
			h = middlewareAdminCheckSession(h)
		}else if sessionType == LEGO_SESSION_TYPE_WEB {
			h = middlewareWebCheckSession(h)
		}else if sessionType == LEGO_SESSION_TYPE_WEB {
			h = middlewareWapCheckSession(h)
		}else{
			panic(errors.New("PLEASE SET THE SESSION TYPE"))
		}

	}
	h = middlewareIPBlock(h)
	h = middlewareAccessLog(h)
	GetRouter().GET(path, h)
}

/* 注册路由，POST的方式
* parmas
	@path -- 访问路径
	@h 	  -- 处理handler
	@checktoken -- 是否校验token
*/
func POST(path string, h fasthttp.RequestHandler, sessionType int,checksession bool) {
	LogPrintln("Route Register: POST:\t" + path)
	if checksession {
		if sessionType == LEGO_SESSION_TYPE_ADMIN {
			h = middlewareAdminCheckSession(h)
		}else if sessionType == LEGO_SESSION_TYPE_WEB {
			h = middlewareWebCheckSession(h)
		}else if sessionType == LEGO_SESSION_TYPE_WEB {
			h = middlewareWapCheckSession(h)
		}else{
			panic(errors.New("PLEASE SET THE SESSION TYPE"))
		}
	}
	h = middlewareIPBlock(h)
	h = middlewareCROS(h)
	h = middlewareAccessLog(h)
	GetRouter().POST(path, h)
}

/* 注册路由，GET的方式
* parmas
	@path -- 访问路径
	@h 	  -- 处理handler
	@checktoken -- 是否校验token
	@checktoken -- 是否校验资源权限
	@checkip 	-- 是否检查调用方ip
*/
func APIGET(path string, h fasthttp.RequestHandler, checktoken bool) {
	LogPrintln("API Register: GET:\t" + path)
	if checktoken {
		h = middlewareCheckAuthToken(h)
	}

	h = middlewareCheckApiSign(h)
	h = middlewareCROS(h)
	h = middlewareAccessLog(h)

	GetRouter().GET(path, h)
}

/* 注册路由，Post的方式
* parmas
	@path -- 访问路径
	@h 	  -- 处理handler
	@checktoken -- 是否校验token
*/
func APIPOST(path string, h fasthttp.RequestHandler, checktoken bool) {
	LogPrintln("API Register: POST:\t" + path)
	if checktoken {
		h = middlewareCheckAuthToken(h)
	}

	h = middlewareCheckApiSign(h)
	h = middlewareCROS(h)
	h = middlewareAccessLog(h)

	GetRouter().POST(path, h)
	GetRouter().OPTIONS(path, h)
}

/*  一些默认不需要签名的接口 */
func APIPOSTWITHOUTSIGN(path string, h fasthttp.RequestHandler) {
	LogPrintln("API Register: POST:\t" + path)
	h = middlewareIPBlock(h)
	h = middlewareAccessLog(h)
	GetRouter().POST(path, h)
}

/*  一些默认不需要签名的接口 */
func APIGETWITHOUTSIGN(path string, h fasthttp.RequestHandler) {
	LogPrintln("API Register: GET:\t" + path)
	h = middlewareIPBlock(h)
	h = middlewareAccessLog(h)
	GetRouter().GET(path, h)
}

/* websocket，注册
* parmas
	@path -- 访问路径
	@h 	  -- 处理handler
	@checktoken -- 是否校验token
*/
func WEBSOCKET(path string, h fasthttp.RequestHandler) {
	LogPrintln("API Register: WEBSOCKET :\t" + path)
	_websockMaper[path] = h
}