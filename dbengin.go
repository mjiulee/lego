package lego

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"os"
	"sync"
	"xorm.io/core"
)

// 多数据源的情况
var _engineMap map[string]*xorm.EngineGroup
var _beanlist []interface{}
var _defaultEnginName string // 默认引擎，用来兼容

var __cconce sync.Once

func init() {
	__cconce.Do(func() {
		_engineMap = make(map[string]*xorm.EngineGroup, 0)
	})
}

/* 获取数据库引擎对象 */
func GetDBEngine() *xorm.EngineGroup {
	_engine, ok := _engineMap[_defaultEnginName]
	if ok {
		return _engine
	}
	return nil
}

/* 按名字，获取数据库引擎对象 */
func GetDBEngineByName(sourceName string) *xorm.EngineGroup {
	_engine, ok := _engineMap[sourceName]
	if ok {
		return _engine
	}
	return nil
}

func SetUpSourceDatabase(sourceName string, dbUrls []string, prefix, isdefault bool) {
	// 使用xorm-plus
	egroup, err := xorm.NewEngineGroup("mysql", dbUrls, xorm.LeastConnPolicy())
	if err != nil {
		//log.Panic(err)
		LogError(err)
		os.Exit(-1)
	}

	if prefix {
		tbMapper := core.NewPrefixMapper(core.SnakeMapper{}, "tb_")
		egroup.SetTableMapper(tbMapper)
	}

	if err = egroup.Ping(); err != nil {
		LogError(err)
		os.Exit(-1)
	}

	egroup.ShowSQL(false)
	_engineMap[sourceName] = egroup
	if isdefault {
		_defaultEnginName = sourceName
	}
}

///* 把要同步的表结构，传到初始化列表里面 */
func AddBeanToSynList(bean interface{}) {
	_beanlist = append(_beanlist, bean)
}

/* 同步数据库 */
func DoSynBeans(sourceName string) {
	// 初始化数据库引擎
	err := GetDBEngine().Sync2(_beanlist...)
	if nil != err {
		LogError("数据库同步失败，请检查model配置" + err.Error())
		os.Exit(-1)
	}
}

/* 获取结构列表 */
func GetSyncModelBeanlist() []interface{} {
	return _beanlist
}
