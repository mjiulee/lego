package lego

import (
	"sync"

	"github.com/go-ini/ini"
	"github.com/mjiulee/lego/utils"
)

//var _envinifile *ini.File
var _appinifile *ini.File
var _inifileName string
var once sync.Once

func LoadConfigFrom(fname string) {
	once.Do(func() {
		_inifileName = fname
		doLoad()
	})
}

/*
 * env
 */
func GetIniByKey(section string, key string) string {
	if _appinifile == nil {
		return ""
	}
	value := _appinifile.Section(section).Key(key)
	if value == nil {
		return ""
	}
	return value.String()
}

/*
 * 有时候改了配置，想不重启的情况，调用这个
 */
func GetIniByKeyWithReload(section string, key string, reloadini bool) string {
	if reloadini {
		doLoad()
	}

	if _appinifile == nil {
		return ""
	}
	value := _appinifile.Section(section).Key(key)
	if value == nil {
		return ""
	}
	return value.String()
}

func doLoad() {
	prjpath := utils.GetPwd()
	inifilepaty := prjpath + utils.GetPathSeparter() + _inifileName

	cfg, err := ini.Load([]byte(""), inifilepaty)
	if err != nil {
		panic("read config from ini fail..." + err.Error())
	}
	_appinifile = cfg
}
