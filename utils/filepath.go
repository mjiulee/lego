package utils

import (
	"os"
	"path/filepath"
	"runtime"
)

var ostype = runtime.GOOS

func GetPwd() string {
	var projectPath string
	projectPath, _ = os.Getwd()
	if ostype == "windows" {
		projectPath = filepath.ToSlash(projectPath)
	}
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
