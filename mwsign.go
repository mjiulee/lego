package lego

import (
	"sort"

	"github.com/mjiulee/lego/logger"
	"github.com/mjiulee/lego/utils"
	"github.com/valyala/fasthttp"
)

const ()

/*
 * 接口签名 中间件
 ————————————————————————
*/

func middlewareCheckApiSign(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		var args *fasthttp.Args
		if ctx.IsGet() {
			args = ctx.QueryArgs()
		} else if ctx.IsPost() {
			args = ctx.PostArgs()
		}

		keys := make([]string, 0)
		args.VisitAll(func(key, val []byte) {
			keys = append(keys, string(key))
		})

		if doCheckSign(keys, args) {
			next(ctx)
		} else {
			ctxExt := RequestCtxExtent{ctx}
			data := make(map[string]interface{})
			data["code"] = 1
			data["msg"] = "sign验证失败"
			ctxExt.JSON(200, data)
			return
		}
	})
}

func SignKey() string {
	key := GetIniByKey(K_GGF_CONFIG_CODE_SECTION, K_GGF_CONFIG_CODE_REQ_SIGNKEY)
	return key
}

func doCheckSign(keys []string, args *fasthttp.Args) bool {
	if len(keys) <= 0 {
		return true
	}

	sign := string(args.Peek("sign"))
	if len(sign) <= 0 {
		logger.LogError("sign参数确失")
		return false
	}

	signkey := SignKey()
	sort.Strings(keys)
	ptext := "" + signkey
	for i := 0; i < len(keys); i++ {
		if keys[i] != "sign" {
			ptext += keys[i] + string(args.Peek(keys[i]))
		}
	}
	ptext += signkey

	logger.LogError("ser-ptext=" + ptext)
	sersign := utils.Md5(ptext)
	if sersign != sign {
		logger.LogError("ser-sign=" + sersign)
		logger.LogError("clt-sign=" + sign)
		return false
	}
	return true
}
