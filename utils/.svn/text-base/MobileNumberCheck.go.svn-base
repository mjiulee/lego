package utils

import (
	"regexp"
)

const (
	regular = "^(13[0-9]|14[579]|15[0-3,5-9]|17[0135678]|18[0-9])\\d{8}$"
	idcard  = "(^[1-9]\\d{5}(18|19|([23]\\d))\\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\\d{3}[0-9Xx]$)|(^[1-9]\\d{5}\\d{2}((0[1-9])|(10|11|12))(([0-2][1-9])|10|20|30|31)\\d{2}[0-9Xx]$)"
)

func MobilePhoneValidate(mobileNum string) bool {
	reg := regexp.MustCompile(regular)
	return reg.MatchString(mobileNum)
}

func IdCardNumberValidate(carNum string) bool {
	reg := regexp.MustCompile(idcard)
	return reg.MatchString(carNum)
}

// func main() {
//     if validate("13888888888") {
//         println("是手机号")
//         return
//     }

//     println("不是手机号")
// }
