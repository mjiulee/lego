package utils

import (
	"os"
	"runtime"
)

var ostype = runtime.GOOS

func GetPwd() string {
	var projectPath string
	projectPath, _ = os.Getwd()
	return projectPath
}

func GetOsType() string {
	return ostype
}

func GetPathSeparter() string {
	if ostype == "windows" {
		return "\\"
	} else {
		return "/"
	}
}
