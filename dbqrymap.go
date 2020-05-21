package lego

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

/**/

/* 查询where字段参数类定义 */
type QMapWhereField struct {
	FieldName   string
	Opt         string
	Param       interface{}
	CombineWith string // AND OR
	OmitEmpty   bool   // 是否忽略空或0值
}

type QryConditions struct {
	CombineWith string // AND OR
	FieldList   []QMapWhereField
}

func NewQryConditions() QryConditions {
	cnd := QryConditions{}
	cnd.FieldList = make([]QMapWhereField, 0)
	cnd.CombineWith = "AND"
	return cnd
}

func (obj *QryConditions) AddFeild(field QMapWhereField) {
	obj.FieldList = append(obj.FieldList, field)
}

/*********************************************************************************************/
/* 构建sql语句 */
func QueryListByMap(qmap map[string]interface{}) ([]map[string]string, int64, error) {
	sql, csql, params, err := BuildSqlByQueryMap(qmap)
	if err != nil {
		//logger.Error(err)
		return make([]map[string]string, 0), 0, err
	}

	if GetDBEngine().Logger().IsShowSQL() {
		LogInfo(fmt.Sprintf("\n%s\t%v", sql, params))
	}

	total := int64(0)
	//list, err := GetDBEngine().SQL(sql).QueryInterface(params...)
	list, err := GetDBEngine().SQL(sql).QueryString(params...)
	if err != nil {
		LogError(err)
		return list, total, errors.New(err.Error() + sql)
	}

	total, err = GetDBEngine().SQL(csql, params...).Count()
	if err != nil {
		LogError(err)
		return list, total, errors.New(err.Error() + sql)
	}

	if list == nil {
		return make([]map[string]string, 0), total, nil
	} else {
		return list, total, nil
	}
}

/*内部功能函数-组装mysql 语句*/
func BuildSqlByQueryMap(qmap map[string]interface{}) (string, string, []interface{}, error) {
	// 检查table字段的设置
	params := make([]interface{}, 0)
	tablestr, err := getTableFromMap(qmap)
	if nil != err {
		return "", "", params, err
	}
	//fmt.Println(tablestr)

	// 检查field字段的配置
	fieldstr, err := getFieldsFromMap(qmap)
	if nil != err {
		return "", "", params, err
	}
	//fmt.Println(fieldstr)

	// 检查Join 字段
	joinstr, err := getJoinsFromMap(qmap)
	if nil != err {
		return "", "", params, err
	}
	//fmt.Println(joinstr)

	// 获取where字段
	wherestr, params, err := getWhereFromQmap(qmap)
	//fmt.Println(wherestr)
	//fmt.Println(params)

	orderandpage := getOrderAndLimit(qmap)
	//fmt.Println(orderandpage)

	qsql := fmt.Sprintf("SELECT\n%s\nFROM\n%s\n%s\n%s\n%s\n", fieldstr, tablestr, joinstr, wherestr, orderandpage)
	csql := fmt.Sprintf("SELECT\nCOUNT(*)\nFROM\n%s\n%s\n%s\n", tablestr, joinstr, wherestr)
	//fmt.Println(finalsql)
	return qsql, csql, params, nil
}

func getTableFromMap(qmap map[string]interface{}) (string, error) {
	table, ok := qmap["table"]
	if !ok {
		return "", errors.New("表名未设置1")
	}
	tableMap, ok := table.(map[string]string)
	if !ok {
		return "", errors.New("表名类型错误")
	}
	tname, ok := tableMap["table"]
	if !ok {
		return "", errors.New("表名未设置2")
	}
	talias, ok := tableMap["alias"]
	if ok {
		return fmt.Sprintf(" %s as %s ", tname, talias), nil
	} else {
		return fmt.Sprintf(" %s ", tname), nil
	}
}

func getFieldsFromMap(qmap map[string]interface{}) (string, error) {
	fields, ok := qmap["fields"]
	if !ok {
		return "", errors.New("查询字段未设置")
	}

	fieldsArray, ok := fields.([]string)
	if !ok {
		return "", errors.New("查询字段非数组")
	}
	return fmt.Sprintf(" %s ", strings.Join(fieldsArray, ",")), nil
}

func getJoinsFromMap(qmap map[string]interface{}) (string, error) {
	join, ok := qmap["join"]
	if ok {
		joinMapList, ok := join.([]map[string]string)
		if ok {
			if len(joinMapList) <= 0 {
				return "", nil
			}
			joinarray := make([]string, 0)
			for _, am := range joinMapList {
				table, ok := am["table"]
				if !ok {
					return "", errors.New("Join table 错误 ")
				}
				on, ok := am["on"]
				if !ok {
					return "", errors.New("Join on 错误 ")
				}

				joinwith, ok := am["joinwith"]
				if !ok {
					joinwith = "LEFT JOIN"
				}

				joinSql := fmt.Sprintf("%s %s ON %s", joinwith, table, on)
				joinarray = append(joinarray, joinSql)
			}
			return strings.Join(joinarray, "\n"), nil
		}
	}
	return "", nil
}

func getWhereFromQmap(qmap map[string]interface{}) (string, []interface{}, error) {
	wherearray := make([]string, 0)
	searcharray := make([]string, 0)
	params := make([]interface{}, 0)

	where, ok := qmap["where"]
	if ok {
		for {
			fieldlist, ok1 := where.([]QMapWhereField)
			condition, ok2 := where.(QryConditions)
			conditionList, ok3 := where.([]QryConditions)
			if !ok1 && !ok2 && !ok3 {
				break
			}
			if ok1 {
				tsql, tparam, err := GetSqlByWhereFeildList(fieldlist)
				if err != nil {
					return "", params, err
				}
				if len(tsql) > 0 {
					wherearray = append(wherearray, tsql)
					params = append(params, tparam...)
				}
			}
			if ok2 {
				tsql, tparam, err := GetSqlByCondition(&condition)
				if err != nil {
					return "", params, err
				}
				if len(tsql) > 0 {
					wherearray = append(wherearray, tsql)
					params = append(params, tparam...)
				}
			}

			if ok3 {
				tsql, tparam, err := GetSqlByConditionGroup(conditionList)
				if err != nil {
					return "", params, err
				}
				if len(tsql) > 0 {
					wherearray = append(wherearray, tsql)
					params = append(params, tparam...)
				}
			}
			break
		}
	}

	// 数据权限过滤
	datafilter, ok := qmap["dataright"]
	if ok {
		for {
			fieldlist, ok1 := datafilter.([]QMapWhereField)
			condition, ok2 := datafilter.(QryConditions)
			conditionList, ok3 := datafilter.([]QryConditions)
			if !ok1 && !ok2 && !ok3 {
				break
			}
			if ok1 {
				tsql, tparam, err := GetSqlByWhereFeildList(fieldlist)
				if err != nil {
					return "", params, err
				}
				if len(tsql) > 0 {
					wherearray = append(wherearray, tsql)
					params = append(params, tparam...)
				}
			}
			if ok2 {
				tsql, tparam, err := GetSqlByCondition(&condition)
				if err != nil {
					return "", params, err
				}
				if len(tsql) > 0 {
					wherearray = append(wherearray, tsql)
					params = append(params, tparam...)
				}
			}

			if ok3 {
				tsql, tparam, err := GetSqlByConditionGroup(conditionList)
				if err != nil {
					return "", params, err
				}
				if len(tsql) > 0 {
					wherearray = append(wherearray, tsql)
					params = append(params, tparam...)
				}
			}
			break
		}
	}

	searchptr, ok := qmap["search"]
	if ok {
		searchmap, ok := searchptr.(map[string]string)
		if ok {
			for k, v := range searchmap {
				if len(v) > 0 {
					tv := v
					if strings.Index(v, "%") < 0 {
						tv = "%" + v + "%"
					}
					tw := fmt.Sprintf("( %s LIKE '%s' )", k, tv)
					searcharray = append(searcharray, tw)
					//if strings.Index(v, "%") > 0 {
					//	params = append(params, v)
					//} else {
					//	params = append(params, "%"+v+"%")
					//}
				}
			}
		}
	}

	fstr := ""
	if len(wherearray) > 0 && len(searcharray) > 0 {
		wstr := strings.Join(wherearray, " AND ") //
		sstr := strings.Join(searcharray, " OR ") //
		fstr = fmt.Sprintf("WHERE %s AND ( %s )", wstr, sstr)
	} else if len(wherearray) > 0 {
		fstr = fmt.Sprintf("WHERE %s ", strings.Join(wherearray, " AND "))
	} else if len(searcharray) > 0 {
		fstr = fmt.Sprintf("WHERE %s ", strings.Join(searcharray, " OR "))
	}
	return fstr, params, nil
}

func getOrderAndLimit(qmap map[string]interface{}) string {
	orderstrs := make([]string, 0)
	page := -1
	psize := -1
	orders, ok := qmap["orders"]
	if ok {
		ordersList, ok := orders.([]string)
		if ok {
			for _, v := range ordersList {
				orderstrs = append(orderstrs, v)
			}
		}
	}

	pageptr, ok := qmap["page"]
	if ok {
		pagemap, ok := pageptr.(map[string]int)
		if ok {
			apage, ok1 := pagemap["page"]
			apsize, ok2 := pagemap["psize"]
			if ok1 && ok2 {
				page = apage
				psize = apsize
			}
		}
	}

	sql := ""
	if len(orderstrs) > 0 {
		sql += fmt.Sprintf("ORDER BY %s ", strings.Join(orderstrs, ","))
	}

	if page >= 0 && psize >= 0 {
		sql += fmt.Sprintf("LIMIT %d, %d ", (page-1)*psize, psize)
	}

	return sql
}

/*
 *  从field列表中获取，该条件组条件的sql及参数
 */
func GetSqlByWhereFeildList(wheremap []QMapWhereField) (string, []interface{}, error) {
	wherearray := make([]string, 0)
	joinArray := make([]string, 0)
	params := make([]interface{}, 0)
	for _, v := range wheremap {
		if v.Opt == "IN" || v.Opt == "NOT IN" {
			switch reflect.TypeOf(v.Param).Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(v.Param)
				if s.Len() <= 0 && v.OmitEmpty == true {
					continue
				}
				//logger.Error(fmt.Sprintf("IN length = %d", s.Len()))
				qList := make([]string, 0)
				for i := 0; i < s.Len(); i++ {
					switch s.Index(i).Kind() {
					case reflect.String:
						qList = append(qList, fmt.Sprintf("'%v'", s.Index(i).Interface()))
					default:
						qList = append(qList, fmt.Sprintf("%v", s.Index(i).Interface()))
					}
					//params = append(params, s.Index(i).Interface())
				}
				tw := fmt.Sprintf("%s %s (%s)", v.FieldName, v.Opt, strings.Join(qList, ","))
				wherearray = append(wherearray, tw)
				joinArray = append(joinArray, v.CombineWith)
			}
		} else if v.Opt == "LIKE" || v.Opt == "like" {
			switch reflect.TypeOf(v.Param).Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(v.Param)
				if s.Len() <= 0 && v.OmitEmpty == true {
					continue
				}
				torlist := make([]string, 0)
				for i := 0; i < s.Len(); i++ {
					if len(s.Index(i).Interface().(string)) <= 0 {
						continue
					}
					tw := fmt.Sprintf(" %s LIKE '%s' ", v.FieldName, s.Index(i).String())
					torlist = append(torlist, tw)
					//params = append(params, s.Index(i).Interface())
				}
				if len(torlist) > 0 {
					wherearray = append(wherearray, fmt.Sprintf("(%s)", strings.Join(torlist, " OR ")))
					joinArray = append(joinArray, v.CombineWith)
				}
			case reflect.String:
				s := reflect.ValueOf(v.Param)
				if len(s.String()) <= 0 {
					continue
				}
				tw := fmt.Sprintf(" %s LIKE '%s' ", v.FieldName, v.Param)
				wherearray = append(wherearray, tw)
				//params = append(params, v.Param)
				joinArray = append(joinArray, v.CombineWith)
			}

		} else {
			switch reflect.TypeOf(v.Param).Kind() {
			case reflect.Slice:
				s := reflect.ValueOf(v.Param)
				if s.Len() <= 0 {
					continue
				}
				torlist := make([]string, 0)
				for i := 0; i < s.Len(); i++ {
					tw := fmt.Sprintf("%s %s %v", v.FieldName, v.Opt, s.Index(i).Interface())
					torlist = append(torlist, tw)
					//params = append(params, s.Index(i).Interface())
				}
				wherearray = append(wherearray, fmt.Sprintf("(%s)", strings.Join(torlist, " OR ")))
				joinArray = append(joinArray, v.CombineWith)
			case reflect.String:
				s := reflect.ValueOf(v.Param)
				if len(s.String()) <= 0 && v.OmitEmpty {
					continue
				}
				//tw := fmt.Sprintf("%s %s '%s'", v.FieldName, v.Opt)
				tw := fmt.Sprintf("%s %s '%s'", v.FieldName, v.Opt, v.Param.(string))
				wherearray = append(wherearray, tw)
				//params = append(params, v.Param)
				joinArray = append(joinArray, v.CombineWith)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				s := reflect.ValueOf(v.Param)
				if s.Int() == 0 && v.OmitEmpty {
					continue
				}
				//tw := fmt.Sprintf("%s %s ?", v.FieldName, v.Opt)
				tw := fmt.Sprintf("%s %s %v", v.FieldName, v.Opt, v.Param)
				wherearray = append(wherearray, tw)
				//params = append(params, v.Param)
				joinArray = append(joinArray, v.CombineWith)
			case reflect.Float32, reflect.Float64:
				s := reflect.ValueOf(v.Param)
				if s.Float() == 0.0 && v.OmitEmpty {
					continue
				}
				tw := fmt.Sprintf("%s %s %v", v.FieldName, v.Opt, v.Param)
				wherearray = append(wherearray, tw)
				//params = append(params, v.Param)
				joinArray = append(joinArray, v.CombineWith)
			default:
				//s := reflect.ValueOf(v.Param)
				//if s.IsNil() && v.OmitEmpty {
				//	continue
				//}
				tw := fmt.Sprintf("%s %s %v", v.FieldName, v.Opt, v.Param)
				wherearray = append(wherearray, tw)
				//params = append(params, v.Param)
				joinArray = append(joinArray, v.CombineWith)
			}
		}
	}

	if len(wherearray) != len(joinArray) {
		return "", params, errors.New("参数和条件不匹配")
	}

	sql := ""
	for i := 0; i < len(wherearray) && i < len(joinArray); i++ {
		if i == 0 {
			sql += wherearray[i]
		} else {
			sql += fmt.Sprintf(" %s %s ", joinArray[i], wherearray[i])
		}
	}
	return sql, params, nil
}

/*
 *  从condition对象中拼接sql
 */
func GetSqlByCondition(cond *QryConditions) (string, []interface{}, error) {
	return GetSqlByWhereFeildList(cond.FieldList)
}

/*
 *  从condition对象组中拼接sql
 */
func GetSqlByConditionGroup(condlist []QryConditions) (string, []interface{}, error) {
	sql := ""
	params := make([]interface{}, 0)

	for i, cnd := range condlist {
		tsql, tp, err := GetSqlByCondition(&cnd)
		if err != nil {
			return "", params, err
		}

		if i == 0 {
			sql += fmt.Sprintf("(%s) ", tsql)
		} else {
			sql += fmt.Sprintf(" %s (%s) ", cnd.CombineWith, tsql)
		}
		params = append(params, tp...)
	}
	if len(condlist) > 0 {
		sql = fmt.Sprintf("(%s)", sql)
	}

	return sql, params, nil
}
