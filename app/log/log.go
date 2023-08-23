package log

type Logger interface {
	Info(string)
	Warning(string)
	Error(error)
}
