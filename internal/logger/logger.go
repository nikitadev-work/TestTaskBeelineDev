package logger

import (
	"log"
	"os"
)

func InitLogger(addr string) {
	file, err := os.OpenFile(addr, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}
