package utils

import (
	"fmt"
	"sort"
)

// 计算签名
func SignParam(signkey string, params map[string]string) string {
	keys := make([]string, 0)
	for k, _ := range params {
		keys = append(keys, k)
	}

	sort.Strings(keys)
	ptext := "" + signkey
	for i := 0; i < len(keys); i++ {
		if keys[i] != "sign" {
			ptext += keys[i] + params[keys[i]]
		}
	}
	ptext += signkey
	fmt.Println("ptext:=" + ptext)
	sersign := Md5(ptext)
	return sersign
}
