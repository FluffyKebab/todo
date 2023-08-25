package console

import (
	"log"

	logging "github.com/FluffyKebab/todo/app/log"
)

type Logger struct {
	LogErrors   bool
	LogWarnings bool
	LogInfo     bool
}

var _ logging.Logger = Logger{}

func (l Logger) Error(err error) {
	if l.LogErrors {
		log.Println("error: ", err.Error())
	}
}

func (l Logger) Warning(warn string) {
	if l.LogWarnings {
		log.Println("warning: ", warn)
	}
}

func (l Logger) Info(info string) {
	if l.LogInfo {
		log.Println("info", info)
	}
}
