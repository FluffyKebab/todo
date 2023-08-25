package testlog

import "github.com/FluffyKebab/todo/app/log"

type Logger struct {
	ErrorFunc func(error)
}

var _ log.Logger = Logger{}

func (l Logger) Info(string)    {}
func (l Logger) Warning(string) {}
func (l Logger) Error(err error) {
	l.ErrorFunc(err)
}
