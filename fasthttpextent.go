package lego

import (
	"encoding/json"

	"github.com/mjiulee/lego/logger"
	"github.com/mjiulee/lego/utils"
	"github.com/valyala/fasthttp"
)

type RequestCtxExtent struct {
	*fasthttp.RequestCtx
}

/* 获取参数
* @parsm
	code --- 错误码
	data --- 数据
*/
func (self *RequestCtxExtent) Peek(key string) []byte {
	if self.IsPost() {
		return self.PostArgs().Peek(key)
	} else {
		return self.QueryArgs().Peek(key)
	}
}

/* 输出json
* @parsm
	code --- 错误码
	data --- 数据
*/
func (self *RequestCtxExtent) JSON(code int, dataMap map[string]interface{}) {
	str, err := json.Marshal(dataMap)
	if err != nil {
		self.Response.SetStatusCode(505)
		self.Write([]byte("系统错误"))
	} else {
		self.Response.SetStatusCode(code)
		self.Write(str)
	}
}

/* 输出xml
* @parsm
	code --- 错误码
	data --- 数据
*/
func (self *RequestCtxExtent) XML(dataMap map[string]string) {
	str := utils.Map2Xml(dataMap)
	self.Response.SetStatusCode(200)
	self.Write([]byte(str))
}

/* ******************************************** */
/* HTML
* @parsm
	name --- 文件路径
*/
func (self *RequestCtxExtent) PureHTML(fpath string) {
	self.Response.Header.Set("Content-Type", "text/html; charset=utf-8")
	// TODO: when debug uncache the html file
	// err := self.Response.SendFile("modules/" + name)
	err := self.Response.SendFile(fpath)
	if err != nil {
		logger.LogError(err)
	}
}

// /* ******************************************** */
// /* pongo2 模板支持，定义rander类，以做界面模板显示处理 */

// /* HTML
// * @parsm
// 	name --- 前端模板文件的名称
// 	data --- 数据
// */
// func (self *RequestCtxExtent) HTML(name string, dataMap map[string]interface{}) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			logger.LogError(err)
// 		}
// 	}()

// 	render := getRender()
// 	render.Template = render.gettmpl(name)

// 	self.SetContentType(render.Options.ContentType)
// 	err := render.html(dataMap, self.Response.BodyWriter())
// 	if err != nil {
// 		logger.LogError(err)
// 	}
// }

// /* HTML
// * @parsm
// 	name --- Wap端渲染
// 	data --- 数据
// */
// func (self *RequestCtxExtent) WAPHTML(name string, dataMap map[string]interface{}) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			logger.LogError(err)
// 		}
// 	}()

// 	render := getWapRender()
// 	render.Template = render.gettmpl(name)

// 	self.SetContentType(render.Options.ContentType)
// 	err := render.html(dataMap, self.Response.BodyWriter())
// 	if err != nil {
// 		logger.LogError(err)
// 	}
// }

// /* HTML
// * @parsm
// 	name --- Web端渲染
// 	data --- 数据
// */
// func (self *RequestCtxExtent) WEBHTML(name string, dataMap map[string]interface{}) {
// 	defer func() {
// 		if err := recover(); err != nil {
// 			logger.LogError(err)
// 		}
// 	}()

// 	render := getWebRender()
// 	render.Template = render.gettmpl(name)

// 	self.SetContentType(render.Options.ContentType)
// 	err := render.html(dataMap, self.Response.BodyWriter())
// 	if err != nil {
// 		logger.LogError(err)
// 	}
// }

// /* html模板rander的一些显示配置 */
// type Options struct {
// 	TemplateDir string
// 	ContentType string
// }

// type Render struct {
// 	Options  *Options
// 	Template *pongo2.Template
// 	Context  pongo2.Context
// }

// /* 通过文件名称和数据，load取pongo2的模板文件并输出给浏览器
// * @parsm
// 	name --- 前端模板文件的名称
// */
// func (render *Render) gettmpl(name string) *pongo2.Template {
// 	filename := path.Join(render.Options.TemplateDir, name)
// 	return pongo2.Must(pongo2.FromFile(filename))
// }

// /* 输出给浏览器
// * @parsm
// 	data --- 输出给模板的数据
// */
// func (render *Render) html(data interface{}, w io.Writer) error {
// 	render.Context = pongo2.Context{
// 		"data":            data,
// 		"GoFloatWrap":     GoFloatWrap,
// 		"ExtentIntString": ExtentIntString,
// 	}
// 	err := render.Template.ExecuteWriter(render.Context, w)
// 	return err
// }

// /* 全局静态rander变量，调用请通过getRender
//  * @parsm
//  */
// var _rander *Render
// var _waprander *Render

// func getRender() *Render {
// 	if nil != _rander {
// 		return _rander
// 	}
// 	prjpath := utils.GetPwd()
// 	return &Render{Options: &Options{
// 		TemplateDir: prjpath + utils.GetPathSeparter() + "modules",
// 		ContentType: "text/html; charset=utf-8",
// 	}}
// }

// func getWapRender() *Render {
// 	if nil != _waprander {
// 		return _waprander
// 	}
// 	prjpath := utils.GetPwd()
// 	return &Render{Options: &Options{
// 		TemplateDir: prjpath + utils.GetPathSeparter() + "modules",
// 		ContentType: "text/html; charset=utf-8",
// 	}}
// }

// func getWebRender() *Render {
// 	if nil != _waprander {
// 		return _waprander
// 	}
// 	prjpath := utils.GetPwd()
// 	return &Render{Options: &Options{
// 		TemplateDir: prjpath + utils.GetPathSeparter() + "frontend",
// 		ContentType: "text/html; charset=utf-8",
// 	}}
// }

// func GoFloatWrap(val float64, jd int) string {
// 	fstr := "%0." + fmt.Sprintf("%d", jd) + "f"
// 	logger.LogInfo(fstr)
// 	return fmt.Sprintf(fstr, val)
// }

// func ExtentIntString(val int64) string {
// 	return fmt.Sprintf("%0d", val)
// }
