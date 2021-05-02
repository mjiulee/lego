package lego

import (
	"path"
	"path/filepath"
	"strings"

	"github.com/mjiulee/lego/logger"
	"github.com/valyala/fasthttp"
)

/** 服务器类封装 */
type Server struct {
	server *fasthttp.Server
	port   string
}

/** 启动服务
 * params
 * @port -- 端口号
 */
func (sv *Server) Start(port, appname string) error {
	// 启动http服务
	sv.server = &fasthttp.Server{
		Handler:            GetRequestHandler,
		Name:               appname,
		MaxRequestBodySize: 1024 * 300 * 1024, //byte
	}
	return sv.server.ListenAndServe(":" + port)
}

/** 关闭服务
 * params
 * @port -- 端口号
 */
func (sv *Server) Shutdown() error {
	return sv.server.Shutdown()
}

/** 设置静态文件访问路径及文件根目录
* params
	@prefix -- 静态资源访问路由前缀
	@root   -- 存放文件的系统跟目录
*/
func (sv *Server) Static(prefix, root string) {
	if root == "" {
		root = "."
	}
	sv.routeStatic(prefix, root)
}

/** 使用通配符，设置目录下的文件访问路由
* params
	@prefix -- 静态资源访问路由前缀
	@root   -- 存放文件的系统跟目录
*/
func (sv *Server) routeStatic(prefix, root string) {
	defer func() {
		if err := recover(); err != nil {
			logger.LogError(err)
			logger.LogPanicTrace(8)
		}
	}()

	//1. 发送文件处理函数handler
	h := fasthttp.RequestHandler(func(ctx *fasthttp.RequestCtx) {
		req := string(ctx.Path())
		// fmt.Println(req)

		suffix := path.Ext(req) //获取文件名带后缀
		// fmt.Println(suffix)
		if len(suffix) <= 0 {
			ctx.SetStatusCode(404)
			return
		}
		if false == sv.acceptFileType(suffix) {
			ctx.SetStatusCode(404)
			return
		}

		req = req[len(prefix):]
		name := filepath.Join(root, path.Clean("/"+req)) // "/"+ for security
		if suffix == ".xlsx" {
			ctx.Request.Header.Add("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
		}
		ctx.SendFile(name)
		// logger.Info("静态资源" + name)
	})

	//2. 设置路由通配符
	if prefix == "/" {
		GetRouter().GET("/*", h)
	} else {
		GetRouter().GET(prefix+"/*.", h)
	}
}

/**支持的文件后缀
 * params
 */
func (sv *Server) acceptFileType(suffix string) (accept bool) {
	exts := []string{".wav", ".zip", ".html", ".gif", ".css", ".js", ".txt", ".jpeg", ".jpg", ".bmp", ".png", ".mp3", ".mp4", ".pcm", ".silk", ".xlsx", ".woff2", ".map", ".woff", "ttf", ".pem", ".apk", ".pdf", ".json"}

	canaccept := false
	for _, ext := range exts {
		if strings.EqualFold(ext, suffix) {
			canaccept = true
			break
		}
	}
	return canaccept
}
