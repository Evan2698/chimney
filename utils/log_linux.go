package utils

import (
	"log"
)

func init() {
	LOG = &mylog{}
	log.SetOutput(setlogglobalnull())
}
