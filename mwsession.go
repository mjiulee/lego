package lego

import (
	//"encoding/json"
	"github.com/mjiulee/go-sessions"
	"github.com/valyala/fasthttp"
)

const (
	LEGO_SESSION_TYPE_ADMIN = 1
	LEGO_SESSION_TYPE_WEB   = 2
	LEGO_SESSION_TYPE_WAP   = 3
)

/*
 * SESSION 会话管理中间件
 ————————————————————————
 * session保存地方分2种，
	1、redies
	2.gosession中
*/

func middlewareAdminCheckSession(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		sess := sessions.StartFasthttp(ctx) // init the session
		// 如果是单应用情况下
		sessValues := sess.GetAll() // get all values from this session
		isRedirect := false
		for {
			userid := sessValues["user_id"]
			if userid != nil {
				next(ctx)
			} else {
				isRedirect = true
				break
			}
			break
		}

		if isRedirect {
			//domain := GetIniByKey("HTTP", "DOMAIN")
			ctx.Redirect("/admin/login", 302)
		}
	})
}

func middlewareWebCheckSession(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		sess := sessions.StartFasthttp(ctx) // init the session
		// 如果是单应用情况下
		sessValues := sess.GetAll() // get all values from this session
		isRedirect := false
		for {
			userid := sessValues["user_id"]
			if userid != nil {
				next(ctx)
			} else {
				isRedirect = true
				break
			}
			break
		}

		if isRedirect {
			//domain := GetIniByKey("HTTP", "DOMAIN")
			ctx.Redirect("/user/login", 302)
		}
	})
}

func middlewareWapCheckSession(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		sess := sessions.StartFasthttp(ctx) // init the session
		// 如果是单应用情况下
		sessValues := sess.GetAll() // get all values from this session
		isRedirect := false
		for {
			userid := sessValues["user_id"]
			if userid != nil {
				next(ctx)
			} else {
				isRedirect = true
				break
			}
			break
		}

		if isRedirect {
			//domain := GetIniByKey("HTTP", "DOMAIN")
			ctx.Redirect("/wap/login", 302)
		}
	})
}

// 通过content获取user_id
func SessionGetKeyIntVal(key string, ctx *fasthttp.RequestCtx) int64 {
	sess := sessions.StartFasthttp(ctx) // init the session
	sessValues := sess.GetAll()         // get all values from this session
	// fmt.Println("session: %s\n", sessValues)

	keyval := sessValues[key]
	if keyval != nil {
		return keyval.(int64)
	} else {
		return -1
	}
}

// 通过content获取user_id
func SessionGetKeyStringVal(key string, ctx *fasthttp.RequestCtx) string {
	sess := sessions.StartFasthttp(ctx) // init the session
	sessValues := sess.GetAll()         // get all values from this session
	// fmt.Println("session: %s\n", sessValues)

	keyval := sessValues[key]
	if keyval != nil {
		return string(keyval.(string))
	} else {
		return ""
	}
}
