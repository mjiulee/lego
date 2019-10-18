package lego

import (
	"bytes"
	"fmt"
	"github.com/issue9/logs"
	"github.com/mjiulee/lego/utils"
	"os"
	"runtime"
)

func SetupLoggerBy(fname string) {
	prjpath := utils.GetPwd()
	cfgpaty := prjpath + utils.GetPathSeparter() + fname

	err := logs.InitFromXMLFile(cfgpaty)
	if err != nil {
		os.Stderr.WriteString(err.Error())
		os.Exit(1)
	}
}

func Flush() {
	defer logs.Flush()
}

// const (
// 	_VER string = "1.0.2"
// )

/* Warn - 输出信息 */
func LogInfo(info string) {
	file, line := whereLogPrint()
	logs.Debugf("%-2v|%s|%20s %4d \n", "->", info, file, line)
}

/* Warn - 输出信息 */
func LogError(info interface{}) {
	file, line := whereLogPrint()
	logs.Error(fmt.Sprintf("%-2v|%s|%20s %4d \n", "->", info, file, line))
}

/* Warn - 输出信息 */
func LogPrintln(info interface{}) {
	file, line := whereLogPrint()
	logs.Debugf("%-2v|%s|%20s %4d \n", "->", info, file, line)
}

/* whereLogPrint
 * 获取日志发生的文件以及所在行数
 */
func whereLogPrint() (name string, the int) {
	_, file, line, _ := runtime.Caller(2)
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	return file, line
}

func LogPanicTrace(size int) {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine ")
	line := []byte("\n")
	stack := make([]byte, size<<10) //4KB
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, "\n")
	logs.Error(fmt.Sprintf("%s\n", string(stack)))
}
