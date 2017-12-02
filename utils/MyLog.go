package utils
import (
	"log"
	"os"
	"fmt"
)

var (
	Logger *log.Logger
)

func init(){
	var logpath = "log.txt"
	var file, err1 = os.Create(logpath)
	if err1 != nil {
		fmt.Print("can not create log file")
		panic(err1)
	}

	Logger = log.New(file, "", log.LstdFlags|log.Lshortfile)
	Logger.Println("LogFile : " + logpath)
}