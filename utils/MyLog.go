//  +build !DRCLO

package utils

import (
	"io/ioutil"
	"log"
	"runtime"
	"strconv"
	"time"
)

// MyLog struct
//
type MyLog struct {
	LogImp *log.Logger
}

// Print function
//
func (l *MyLog) Print(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	l.LogImp.Print(file, line, v)
}

// Printf function
//
func (l *MyLog) Printf(format string, v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	l.LogImp.Print(file, line)
	l.LogImp.Printf(format, v)
}

// Println function
//
func (l *MyLog) Println(v ...interface{}) {
	_, file, line, _ := runtime.Caller(1)
	l.LogImp.Println(file, line, v)
}

var (

	// Logger MyLog
	Logger *MyLog
)

// NewLog function
//
func NewLog(l *log.Logger) *MyLog {
	return &MyLog{l}
}

func init() {

	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	var logpath = "log_" + timestamp + ".txt"
	/*var file, err1 = os.Create(logpath)
	if err1 != nil {
		fmt.Print("can not create log file")
		panic(err1)
	}*/

	thislog := log.New(ioutil.Discard, "", log.LstdFlags|log.Lshortfile)

	//thislog := log.New(file, "", log.LstdFlags|log.Lshortfile)
	Logger = NewLog(thislog)
	Logger.Println("LogFile : " + logpath)
}
