package utils

import (
	"log"
	"os"
)

func init() {
	LOG = &mylog{}
	log.SetOutput(os.Stdout)
}
