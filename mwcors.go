package lego

import (
	"github.com/valyala/fasthttp"
)

/*
 * 跨域处理中间件
 ————————————————————————
*/

var (
	corsAllowHeaders     = "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With"
	corsAllowMethods     = "HEAD, GET, POST, PUT, DELETE, OPTIONS"
	corsAllowOrigin      = "*"
	corsAllowCredentials = "true"
)

func middlewareCROS(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.Set("Access-Control-Allow-Credentials", corsAllowCredentials)
		ctx.Response.Header.Set("Access-Control-Allow-Headers", corsAllowHeaders)
		ctx.Response.Header.Set("Access-Control-Allow-Methods", corsAllowMethods)
		ctx.Response.Header.Set("Access-Control-Allow-Origin", corsAllowOrigin)
		next(ctx)
	})
}
