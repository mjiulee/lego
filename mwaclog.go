package lego

import (
	"github.com/valyala/fasthttp"
)

/*
 * 请求日志处理中间件
 ————————————————————————
*/

func middlewareAccessLog(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		next(ctx)
		if ctx.IsPost() {
			//logger.Info("POST" + string(ctx.Path()))
		} else {
			//logger.Info("*GET" + string(ctx.Path()))
		}
	})
}
