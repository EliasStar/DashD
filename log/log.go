package log

import (
	"fmt"
	"log"
)

func Info(tag string, msg ...any) {
	fmt.Print("[" + tag + "] ")
	fmt.Println(msg...)
}

func Error(tag string, msg ...any) {
	log.Print("[" + tag + "] ")
	log.Println(msg...)
}

func PanicIf(tag string, err error) {
	if err != nil {
		log.Panic("["+tag+"]", err)
	}
}
