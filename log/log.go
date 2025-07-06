package log

import "context"

var (
	DefLogger Interface = &emptyLog{}
)

// Interface is structured log interface
type Interface interface {
	CtxDebugf(ctx context.Context, msg string, args ...interface{})
	CtxInfof(ctx context.Context, msg string, args ...interface{})
	CtxWarnf(ctx context.Context, msg string, args ...interface{})
	CtxErrorf(ctx context.Context, msg string, args ...interface{})
	CtxFatalf(ctx context.Context, msg string, args ...interface{})

	Debugf(msg string, args ...interface{})
	Infof(msg string, args ...interface{})
	Warnf(msg string, args ...interface{})
	Errorf(msg string, args ...interface{})
	Fatalf(msg string, args ...interface{})
}

type emptyLog struct {
}

func (e *emptyLog) Debugf(msg string, args ...interface{}) {
}

func (e *emptyLog) Infof(msg string, args ...interface{}) {
}

func (e *emptyLog) Warnf(msg string, args ...interface{}) {
}

func (e *emptyLog) Errorf(msg string, args ...interface{}) {
}

func (e *emptyLog) Fatalf(msg string, args ...interface{}) {
}

func (e *emptyLog) CtxDebugf(ctx context.Context, msg string, args ...interface{}) {
}

func (e *emptyLog) CtxInfof(ctx context.Context, msg string, args ...interface{}) {
}

func (e *emptyLog) CtxWarnf(ctx context.Context, msg string, args ...interface{}) {
}

func (e *emptyLog) CtxErrorf(ctx context.Context, msg string, args ...interface{}) {
}

func (e *emptyLog) CtxFatalf(ctx context.Context, msg string, args ...interface{}) {
}

func Debugf(msg string, args ...interface{}) {
	DefLogger.Debugf(msg, args...)
}
func Infof(msg string, args ...interface{}) {
	DefLogger.Infof(msg, args...)
}
func Warnf(msg string, args ...interface{}) {
	DefLogger.Warnf(msg, args...)
}
func Errorf(msg string, args ...interface{}) {
	DefLogger.Errorf(msg, args...)
}
func Fatalf(msg string, args ...interface{}) {
	DefLogger.Fatalf(msg, args...)
}
func CtxDebugf(ctx context.Context, msg string, args ...interface{}) {
	DefLogger.CtxDebugf(ctx, msg, args...)
}
func CtxInfof(ctx context.Context, msg string, args ...interface{}) {
	DefLogger.CtxInfof(ctx, msg, args...)
}
func CtxWarnf(ctx context.Context, msg string, args ...interface{}) {
	DefLogger.CtxWarnf(ctx, msg, args...)
}
func CtxErrorf(ctx context.Context, msg string, args ...interface{}) {
	DefLogger.CtxErrorf(ctx, msg, args...)
}
func CtxFatalf(ctx context.Context, msg string, args ...interface{}) {
	DefLogger.CtxFatalf(ctx, msg, args...)
}
