package daemon

import "log"

type Logger interface {
	Trace(args ...interface{})
	Tracef(format string, args ...interface{})

	Debug(args ...interface{})
	Debugf(format string, args ...interface{})

	Info(args ...interface{})
	Infof(format string, args ...interface{})

	Warn(args ...interface{})
	Warnf(format string, args ...interface{})

	Error(args ...interface{})
	Errorf(format string, args ...interface{})

	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

type DefaultLog struct{}

func (l DefaultLog) Tracef(format string, args ...interface{}) {
	log.Printf("[trace] "+format, args...)
}

func (l DefaultLog) Trace(args ...interface{}) {
	log.Print(append([]interface{}{"[trace] "}, args...)...)
}

func (l DefaultLog) Debugf(format string, args ...interface{}) {
	log.Printf("[debug] "+format, args...)
}

func (l DefaultLog) Debug(args ...interface{}) {
	log.Print(append([]interface{}{"[debug] "}, args...)...)
}

func (l DefaultLog) Infof(format string, args ...interface{}) {
	log.Printf("[info] "+format, args...)
}

func (l DefaultLog) Info(args ...interface{}) {
	log.Print(append([]interface{}{"[info] "}, args...)...)
}

func (l DefaultLog) Warnf(format string, args ...interface{}) {
	log.Printf("[warn] "+format, args...)
}

func (l DefaultLog) Warn(args ...interface{}) {
	log.Print(append([]interface{}{"[warn] "}, args...)...)
}

func (l DefaultLog) Errorf(format string, args ...interface{}) {
	log.Printf("[error] "+format, args...)
}

func (l DefaultLog) Error(args ...interface{}) {
	log.Print(append([]interface{}{"[error] "}, args...)...)
}

func (l DefaultLog) Fatalf(format string, args ...interface{}) {
	log.Printf("[fatal] "+format, args...)
}

func (l DefaultLog) Fatal(args ...interface{}) {
	log.Print(append([]interface{}{"[fatal] "}, args...)...)
}
