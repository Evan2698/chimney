//  +build !DRCLO

package utils

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"
)

type MyLog struct {
	LogImp *log.Logger
}

func (l *MyLog) Print(v ...interface{}) {
	//_, file, line, _ := runtime.Caller(1)
	//l.LogImp.Print(file, line, v)
}

func (l *MyLog) Printf(format string, v ...interface{}) {
	//_, file, line, _ := runtime.Caller(1)
	//l.LogImp.Print(file, line)
	//l.LogImp.Printf(format, v)
}

func (l *MyLog) Println(v ...interface{}) {
	//_, file, line, _ := runtime.Caller(1)
	//l.LogImp.Println(file, line, v)
}

var (
	Logger *MyLog
)

func NewLog(l *log.Logger) *MyLog {
	return &MyLog{l}
}

func init() {

	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	var logpath = "log_" + timestamp + ".txt"
	var file, err1 = os.Create(logpath)
	if err1 != nil {
		fmt.Print("can not create log file")
		panic(err1)
	}

	thislog := log.New(file, "", log.LstdFlags|log.Lshortfile)

	Logger = NewLog(thislog)
	Logger.Println("LogFile : " + logpath)
}
