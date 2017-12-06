package utils
import (
	"strconv"
	"time"
	"log"
	"os"
	"fmt"
)

var (
	Logger *log.Logger
)

func init(){
	
	t := time.Now()
	timestamp := strconv.FormatInt(t.UTC().UnixNano(), 10)
	var logpath = "log_" + timestamp + ".txt"
	var file, err1 = os.Create(os.DevNull)
	if err1 != nil {
		fmt.Print("can not create log file")
		panic(err1)
	}

	Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Logger.Println("LogFile : " + logpath)
}