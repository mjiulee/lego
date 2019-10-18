package lego

import (
	"fmt"
	"sync"
)

/** 为了不同业务能以插件方式进行引入，需要定义一些接口，让注册到vctrl */

type GGModule interface {
	Name() string
	Version() string
	Info() string
}

var _ggmodulsList []GGModule
var _ggmodulsOnce sync.Once

func init() {
	_ggmodulsOnce.Do(func() {
		_ggmodulsList = make([]GGModule, 0)
	})
}

func ModuleRegist(m GGModule) {
	_ggmodulsList = append(_ggmodulsList, m)
}

/* 显示加载了那些模块 */
func ShowSysModules() {
	fmt.Println("module init")
	fmt.Println("********************************************************")
	for i := 0; i < len(_ggmodulsList); i++ {
		fmt.Printf("模块:%s\n版本：%s\n备注：%s\r\n", _ggmodulsList[i].Name(), _ggmodulsList[i].Version(), _ggmodulsList[i].Info())
	}
	fmt.Println("********************************************************")
}
