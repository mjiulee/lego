package lego

import (
	"github.com/valyala/fasthttp"
	"strings"
)

/*
 * IP白名单处理中间件
 ————————————————————————
 * 需要在ini文件中配置ip白名单，白名单内的请求方处理，否则block掉
*/
func middlewareIPBlock(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {

		whiteIPlist := GetIniByKey("SECURITY", "WHITEIPLIST")

		list := strings.Split(whiteIPlist, ",")
		if len(whiteIPlist) <= 0 || len(list) <= 0 {
			next(ctx)
		} else {
			found := false
			// 从nginx转发过来的http请求，从头里面获取IP信息
			rqip := string(ctx.Request.Header.Peek("Remote_addr"))
			if rqip == "" {
				rqip = ctx.RemoteIP().String()
			}
			for i := 0; i < len(list); i++ {
				if rqip == list[i] {
					found = true
					break
				}
			}

			if found {
				next(ctx)
			} else {
				ctxExt := RequestCtxExtent{ctx}
				data := make(map[string]interface{})
				data["code"] = 1
				data["msg"] = "非白名单内ip"
				ctxExt.JSON(200, data)
			}
		}
	})
}
