package utils

import (
	"io/ioutil"
	"log"
	"strconv"
	"time"
)

var (
	Logger *log.Logger
)

func init() {

	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	var logpath = "log_" + timestamp + ".txt"
	/*var file, err1 = os.Create(logpath)
	if err1 != nil {
		fmt.Print("can not create log file")
		panic(err1)
	}*/

	Logger = log.New(ioutil.Discard, "", log.LstdFlags|log.Lshortfile)
	//Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Logger.Println("LogFile : " + logpath)
}
