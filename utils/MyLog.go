package utils

import (
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

type MyLog struct {
	LogImp *log.Logger
}

func (l *MyLog) Print(v ...interface{}) {
	l.LogImp.Print(v)
}

func (l *MyLog) Printf(format string, v ...interface{}) {
	l.LogImp.Printf(format, v)
}

func (l *MyLog) Println(v ...interface{}) {
	l.LogImp.Println(v)
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
	/*var file, err1 = os.Create(logpath)
	if err1 != nil {
		fmt.Print("can not create log file")
		panic(err1)
	}*/

	thislog := log.New(ioutil.Discard, "", log.LstdFlags|log.Lshortfile)
	//thislog = log.New(file, "", log.LstdFlags|log.Lshortfile)

	Logger = NewLog(thislog)
	Logger.Println("LogFile : " + logpath)
}
