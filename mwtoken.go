package lego

import (
	//"fmt"
	"errors"
	"github.com/mjiulee/lego/utils"
	"github.com/dgrijalva/jwt-go"
	"github.com/valyala/fasthttp"
	"strings"
)

const (
	TOKEN_EXPIRE_CODE = 10001
)

/* JWT-中间件
 * CheckAuthToken
 * 检查上传的token是否有效，用来获取保存用户id的校验接口
 */
func middlewareCheckAuthToken(next fasthttp.RequestHandler) fasthttp.RequestHandler {
	return fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		data := make(map[string]interface{})
		ctxExt := RequestCtxExtent{ctx}
		code := 0
		msg := ""

		for {
			passToken := strings.TrimSpace(string(ctx.FormValue("token")))
			timestamp := strings.TrimSpace(string(ctx.FormValue("timestamp")))
			if len(passToken) <= 0 || len(timestamp) <= 0 {
				code = 1
				msg = "token/timestamp参数不能为空"
				break
			}

			token, err := jwt.Parse(passToken, func(token *jwt.Token) (interface{}, error) {
				return []byte(TokenSecretKey()), nil
			})

			if err != nil || token == nil {
				code = 1
				if err.Error() == "Token is expired" {
					code = TOKEN_EXPIRE_CODE
				}
				msg = "token验证失败"
				break
			}
			if _, ok := token.Claims.(jwt.MapClaims); !ok || !token.Valid {
				code = 1
				msg = "token验证失败"
				break
			}

			// 若需要在redis中存放token，同时检查redis中的是否存在
			ifredis := GetIniByKey("SESSION", "IFREDIS_SESSION")
			if ifredis != "true" {
				break
			}

			tokenredikey := TokenPrefix() + passToken
			_, err = RedisGetKey(tokenredikey)
			if err != nil {
				code = 1
				msg = "token验证失败-redis"
				break
			}
			break
		}
		if code == 0 {
			next(ctx)
		} else {
			data["code"] = code
			data["msg"] = msg
			ctxExt.JSON(200, data)
		}
	})
}

/* JWT-中间件
 * GetUserIdFromToken
 * 从jwt字符串中解析获得userid
 */
func GetUserIdFromToken(passToken string) (uid int64, err error) {
	token, err := jwt.Parse(passToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(TokenSecretKey()), nil
	})

	if err != nil {
		return -1, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		uid , ok := claims["uid"]
		if ok {
			uidint ,_:= utils.StringToInt64(uid.(string))
			return uidint , nil
		}else{
			return -1, errors.New("请先登录")
		}

	} else {
		return -1, err
	}
}

func TokenSecretKey() string {
	key := GetIniByKey("CODE", "PRJ_TOKEN_SECRETKEY")
	return key
}

func TokenPrefix() string {
	prefix := GetIniByKey("CODE", "PRJ_TOKEN_PREFIX")
	return prefix
}
