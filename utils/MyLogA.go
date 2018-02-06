//  +build DRCLO


package utils

import (
	"log"
)

type MyLog struct {
}

func (*MyLog) Print(v ...interface{}) {
	log.Print(v...)
}

func (*MyLog) Printf(v ...interface{}) {
	log.Print(v...)
}

func (*MyLog) Println(v ...interface{}) {
	log.Println(v...)
}

var (
	//Logger *log.Logger
	Logger *MyLog
)

func init() {
	Logger = &MyLog{}
	Logger.Println("android log")
}
