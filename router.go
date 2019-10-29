package lego

import (
	"fmt"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
	"encoding/json"
	"log"
	"reflect"
	"runtime"
	"strings"
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

func POST_4(path string, fn interface{}, sessionType int,checksession bool) {
	LogPrintln("Route Register: POST_4:\t" + path)


	//useName := true
	f, ok := fn.(reflect.Value)
	if !ok {
		f = reflect.ValueOf(fn)
	}
	if f.Kind() != reflect.Func {
		return //"", errors.New("function must be func or bound method")
	}

	fname := runtime.FuncForPC(reflect.Indirect(f).Pointer()).Name()
	if fname != "" {
		i := strings.LastIndex(fname, ".")
		if i >= 0 {
			fname = fname[i+1:]
		}
	}
	/*if useName {
		fname = name
	}*/
	if fname == "" {
		errorStr := "rpcx.registerFunction: no func name for type " + f.Type().String()
		log.Println(errorStr)
		return //fname, errors.New(errorStr)
	}

	t := f.Type()
	if t.NumIn() != 2 {
		return //fname, fmt.Errorf("rpcx.registerFunction: has wrong number of ins: %s", f.Type().String())
	}
	if t.NumOut() != 0 {
		return //fname, fmt.Errorf("rpcx.registerFunction: has wrong number of outs: %s", f.Type().String())
	}

	// First arg must be context.Context
	ctxType := t.In(0)
	ctxType = ctxType
	/*if !ctxType.Implements(typeOfContext) {
		return //fname, fmt.Errorf("function %s must use context as  the first parameter", f.Type().String())
	}*/

	argType := t.In(1)
	argType = argType
	var argv reflect.Value
	if argType.Kind() == reflect.Ptr { // reply must be ptr
		argv = reflect.New(argType.Elem())
	} else {
		argv = reflect.New(argType)
	}
	fmt.Printf("argType: %v", argv.Interface())
	/*if !isExportedOrBuiltinType(argType) {
		return //fname, fmt.Errorf("function %s parameter type not exported: %v", f.Type().String(), argType)
	}*/

	//replyType := t.In(2)
	//replyType = replyType
	/*if replyType.Kind() != reflect.Ptr {
		return //fname, fmt.Errorf("function %s reply type not a pointer: %s", f.Type().String(), replyType)
	}*/
	/*if !isExportedOrBuiltinType(replyType) {
		return //fname, fmt.Errorf("function %s reply type not exported: %v", f.Type().String(), replyType)
	}*/

	// The return type of the method must be error.
	/*if returnType := t.Out(0); returnType != typeOfError {
		return //fname, fmt.Errorf("function %s returns %s, not error", f.Type().String(), returnType.String())
	}*/
    type vcHandle func(*fasthttp.RequestCtx, interface{})


	h := func(ctx *fasthttp.RequestCtx) {

		requestByte := ctx.PostBody()
		v := argv.Interface()
		err := json.Unmarshal(requestByte, &v)
		err = err
		//放置解析函数
		fmt.Printf("POST_4 f:%v \n", argType)

		//f1 := fn.(vcHandle)
		f.Call([]reflect.Value{reflect.ValueOf(ctx), reflect.ValueOf(v)})
		//f1(ctx, argType)
		//h(ctx)
	}
	POST(path, h, sessionType, checksession)
}
