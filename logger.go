package daemon

import "log"

type Logger interface {
	Trace(format string)
	Tracef(format string, args ...interface{})

	Debug(format string)
	Debugf(format string, args ...interface{})

	Info(format string)
	Infof(format string, args ...interface{})

	Warn(format string)
	Warnf(format string, args ...interface{})

	Error(format string)
	Errorf(format string, args ...interface{})

	Fatal(format string)
	Fatalf(format string, args ...interface{})
}

type DefaultLog struct{}

func (l DefaultLog) Tracef(format string, args ...interface{}) {
	log.Printf("[trace] "+format, args...)
}

func (l DefaultLog) Trace(format string) {
	l.Tracef(format)
}

func (l DefaultLog) Debugf(format string, args ...interface{}) {
	log.Printf("[debug] "+format, args...)
}

func (l DefaultLog) Debug(format string) {
	l.Debugf(format)
}

func (l DefaultLog) Infof(format string, args ...interface{}) {
	log.Printf("[info] "+format, args...)
}

func (l DefaultLog) Info(format string) {
	l.Infof(format)
}

func (l DefaultLog) Warnf(format string, args ...interface{}) {
	log.Printf("[warn] "+format, args...)
}

func (l DefaultLog) Warn(format string) {
	l.Warnf(format)
}

func (l DefaultLog) Errorf(format string, args ...interface{}) {
	log.Printf("[error] "+format, args...)
}

func (l DefaultLog) Error(format string) {
	l.Errorf(format)
}

func (l DefaultLog) Fatalf(format string, args ...interface{}) {
	log.Printf("[fatal] "+format, args...)
}

func (l DefaultLog) Fatal(format string) {
	l.Fatalf(format)
}
