package handlers

type Logger interface {
	Infof(format string, args ...interface{})
	Info(...interface{})
}
