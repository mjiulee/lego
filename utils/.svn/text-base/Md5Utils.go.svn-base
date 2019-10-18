package utils

import (
	"crypto/md5"
	"fmt"
	"io"
)

/*	Md5
 *  根据传入的参数计算MD5的值
 */
func Md5(str string) string {
	w := md5.New()
	io.WriteString(w, str)                   //将str写入到w中
	md5str2 := fmt.Sprintf("%x", w.Sum(nil)) //w.Sum(nil)将w的hash转成[]byte格式
	return md5str2
}
