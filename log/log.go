package log

var (
	DefLogger Interface = emptyLog{}
)

// Interface is structured log interface
type Interface interface {
	Debug(msg string, fields ...interface{})
	Debugf(msg string, args ...interface{})
	Info(msg string, fields ...interface{})
	Infof(msg string, args ...interface{})
	Warn(msg string, fields ...interface{})
	Warnf(msg string, args ...interface{})
	Error(msg string, fields ...interface{})
	Errorf(msg string, args ...interface{})
	Fatal(msg string, fields ...interface{})
	Fatalf(msg string, args ...interface{})
}

type emptyLog struct {
}

func (e emptyLog) Debug(msg string, fields ...interface{}) {
}

func (e emptyLog) Debugf(msg string, args ...interface{}) {
}

func (e emptyLog) Info(msg string, fields ...interface{}) {
}

func (e emptyLog) Infof(msg string, args ...interface{}) {
}

func (e emptyLog) Warn(msg string, fields ...interface{}) {
}

func (e emptyLog) Warnf(msg string, args ...interface{}) {
}

func (e emptyLog) Error(msg string, fields ...interface{}) {
}

func (e emptyLog) Errorf(msg string, args ...interface{}) {
}

func (e emptyLog) Fatal(msg string, fields ...interface{}) {
}

func (e emptyLog) Fatalf(msg string, args ...interface{}) {
}

func Debug(msg string, fields ...interface{}) {
	DefLogger.Debug(msg, fields...)
}
func Debugf(msg string, args ...interface{}) {
	DefLogger.Debugf(msg, args...)
}
func Info(msg string, fields ...interface{}) {
	DefLogger.Info(msg, fields...)

}
func Infof(msg string, args ...interface{}) {
	DefLogger.Infof(msg, args...)

}
func Warn(msg string, fields ...interface{}) {
	DefLogger.Warn(msg, fields...)
}
func Warnf(msg string, args ...interface{}) {
	DefLogger.Warnf(msg, args...)
}
func Error(msg string, fields ...interface{}) {
	DefLogger.Error(msg, fields...)
}
func Errorf(msg string, args ...interface{}) {
	DefLogger.Debugf(msg, args...)
}
func Fatal(msg string, fields ...interface{}) {
	DefLogger.Fatal(msg, fields...)
}
func Fatalf(msg string, args ...interface{}) {
	DefLogger.Fatalf(msg, args...)
}
