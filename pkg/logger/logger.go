package logger

import (
	"github.com/hashicorp/go-hclog"
	"github.com/readystock/golog"
	"io"
	"log"
)

type Logger interface {
	// Args are alternating key, val pairs
	// keys must be strings
	// vals can be any type, but display is implementation specific
	// Emit a message and key/value pairs at the TRACE level
	Trace(msg string, args ...interface{})

	// Emit a message and key/value pairs at the DEBUG level
	Debug(msg string, args ...interface{})

	// Emit a message and key/value pairs at the INFO level
	Info(msg string, args ...interface{})

	// Emit a message and key/value pairs at the WARN level
	Warn(msg string, args ...interface{})

	// Emit a message and key/value pairs at the ERROR level
	Error(msg string, args ...interface{})

	// Indicate if TRACE logs would be emitted. This and the other Is* guards
	// are used to elide expensive logging code based on the current level.
	IsTrace() bool

	// Indicate if DEBUG logs would be emitted. This and the other Is* guards
	IsDebug() bool

	// Indicate if INFO logs would be emitted. This and the other Is* guards
	IsInfo() bool

	// Indicate if WARN logs would be emitted. This and the other Is* guards
	IsWarn() bool

	// Indicate if ERROR logs would be emitted. This and the other Is* guards
	IsError() bool

	// Creates a sublogger that will always have the given key/value pairs
	With(args ...interface{}) hclog.Logger

	// Create a logger that will prepend the name string on the front of all messages.
	// If the logger already has a name, the new value will be appended to the current
	// name. That way, a major subsystem can use this to decorate all it's own logs
	// without losing context.
	Named(name string) hclog.Logger

	// Create a logger that will prepend the name string on the front of all messages.
	// This sets the name of the logger to the value directly, unlike Named which honor
	// the current name as well.
	ResetNamed(name string) hclog.Logger

	// Updates the level. This should affect all sub-loggers as well. If an
	// implementation cannot update the level on the fly, it should no-op.
	SetLevel(level hclog.Level)

	// Return a value that conforms to the stdlib log.Logger interface
	StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger

	// Return a value that conforms to io.Writer, which can be passed into log.SetOutput()
	StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer
}

func NewLogger() Logger {
	l := golog.NewWithDepth(6)
	l.SetLevel("trace")
	return &lggr{
		gologger: l,
	}
}

type lggr struct {
	gologger *golog.Logger
}

func (l *lggr) Trace(msg string, args ...interface{}) {
	l.gologger.Tracef(msg, args...)
}

func (l *lggr) Debug(msg string, args ...interface{}) {
	l.gologger.Debugf(msg, args...)
}

func (l *lggr) Info(msg string, args ...interface{}) {
	l.gologger.Infof(msg, args...)
}

func (l *lggr) Warn(msg string, args ...interface{}) {
	l.gologger.Warnf(msg, args...)
}

func (l *lggr) Error(msg string, args ...interface{}) {
	l.gologger.Errorf(msg, args...)
}

func (l *lggr) IsTrace() bool {
	return l.gologger.Level >= golog.TraceLevel
}

func (l *lggr) IsDebug() bool {
	return l.gologger.Level >= golog.DebugLevel
}

func (l *lggr) IsInfo() bool {
	return l.gologger.Level >= golog.InfoLevel
}

func (l *lggr) IsWarn() bool {
	return l.gologger.Level >= golog.WarnLevel
}

func (l *lggr) IsError() bool {
	return l.gologger.Level >= golog.ErrorLevel
}

func (l *lggr) With(args ...interface{}) hclog.Logger {
	return NewLogger()
}

func (l *lggr) Named(name string) hclog.Logger {
	return NewLogger()
}

func (l *lggr) ResetNamed(name string) hclog.Logger {
	return NewLogger()
}

func (l *lggr) SetLevel(level hclog.Level) {
	lvl := func() golog.Level {
		switch level {
		case hclog.NoLevel:
			return golog.DisableLevel
		case hclog.Trace:
			return golog.TraceLevel
		case hclog.Debug:
			return golog.DebugLevel
		case hclog.Info:
			return golog.InfoLevel
		case hclog.Warn:
			return golog.WarnLevel
		case hclog.Error:
			return golog.ErrorLevel
		default:
			return golog.TraceLevel
		}
	}()
	l.gologger.SetLevel(golog.Levels[lvl].Name)
}

func (l *lggr) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	return nil
}

func (l *lggr) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	panic("implement me")
}
