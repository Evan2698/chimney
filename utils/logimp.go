package utils

import (
	"fmt"
	"io"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// Ilog ..
type Ilog interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}

// LOG ...
var LOG Ilog

type mylog struct {
}

func (*mylog) Print(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	log.Print(file, line, v)
}

func (*mylog) Printf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	var tmp strings.Builder
	tmp.WriteString(file)
	tmp.WriteString(", ")
	tmp.WriteString(strconv.Itoa(line))
	tmp.WriteString(": ")
	tmp.WriteString(format)
	log.Printf(tmp.String(), v)
}

func (*mylog) Println(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	log.Println(file, line, v)
}

func setlogglobal() io.Writer {
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	var logpath = "log_" + timestamp + ".txt"
	var file io.Writer
	var err1 error
	file, err1 = os.Create(logpath)
	if err1 != nil {
		fmt.Print("can not create log file", err1)
		file = &highspeeddevice{}
	}
	return io.MultiWriter(os.Stdout, file)
}

type highspeeddevice struct {
}

func (high *highspeeddevice) Write(p []byte) (n int, err error) {
	return len(p), nil
}

func setlogglobalnull() io.Writer {
	return &highspeeddevice{}
}
