package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"math"
	"os"
	"runtime"
	"sync"
	"time"
)

const (
	Black      ANSIColor = "\033[0;30m"
	Red        ANSIColor = "\033[0;31m"
	Green      ANSIColor = "\033[0;32m"
	Yellow     ANSIColor = "\033[0;33m"
	Blue       ANSIColor = "\033[0;34m"
	Purple     ANSIColor = "\033[0;35m"
	Cyan       ANSIColor = "\033[0;36m"
	White      ANSIColor = "\033[0;37m"
	BoldBlack  ANSIColor = "\033[1;30m"
	BoldRed    ANSIColor = "\033[1;31m"
	BoldGreen  ANSIColor = "\033[1;32m"
	BoldYellow ANSIColor = "\033[1;33m"
	BoldBlue   ANSIColor = "\033[1;34m"
	BoldPurple ANSIColor = "\033[1;35m"
	BoldCyan   ANSIColor = "\033[1;36m"
	BoldWhite  ANSIColor = "\033[1;37m"
	BgBlack    ANSIColor = "\033[40m"
	BgRed      ANSIColor = "\033[41m"
	BgGreen    ANSIColor = "\033[42m"
	BgYellow   ANSIColor = "\033[43m"
	BgBlue     ANSIColor = "\033[44m"
	BgPurple   ANSIColor = "\033[45m"
	BgCyan     ANSIColor = "\033[46m"
	BgWhite    ANSIColor = "\033[47m"
	Reset      ANSIColor = "\033[0m"
	Bold       ANSIColor = "\033[1m"
	Underline  ANSIColor = "\033[4m"
	Inverse    ANSIColor = "\033[7m"

	DebugLevel    LogLevel = -4
	InfoLevel     LogLevel = 0
	WarnLevel     LogLevel = 4
	ErrorLevel    LogLevel = 8
	CriticalLevel LogLevel = 12
	FatalLevel    LogLevel = 16
	NoLevel       LogLevel = math.MaxInt32
)

type LogLevel int32

type ANSIColor string

type HandlerOpts struct {
	Level      slog.Level
	TimeLayout string
	Prefix     string
}

// Logger is the Custom Structured Logging Handler implementation for synapse
//
// It uses a pointer to a [sync.Mutex] to ensure that no other goroutines access the
// data passed through an io.Writer
type Logger struct {
	Options HandlerOpts
	Writer  io.Writer
	Ctx     context.Context
	Mu      *sync.Mutex
	Attrs   []slog.Attr
}

func (c ANSIColor) String() string {
	return string(c)
}

func (c ANSIColor) AddForegroundColor(s string) string {
	return c.String() + s + Reset.String()
}

func (c ANSIColor) AddBackgroundColor(s string) string {
	return c.String() + s + Reset.String()
}

func Colorize(s string, fg ANSIColor, bg ANSIColor) string {
	return fg.String() + bg.String() + s + Reset.String()
}

func (l LogLevel) String() string {
	switch l {
	case ErrorLevel:
		return "ERROR"
	case CriticalLevel:
		return "CRITICAL"
	case DebugLevel:
		return "DEBUG"
	case WarnLevel:
		return "WARNING"
	case FatalLevel:
		return "FATAL"
	}

	return "INFO"
}

func (l LogLevel) TagColor() (fg ANSIColor, bg ANSIColor) {
	switch l {
	case ErrorLevel:
	case CriticalLevel:
		return White, BgRed
	case DebugLevel:
		return BoldWhite, BgBlue
	case WarnLevel:
		return Black, BgYellow
	}
	return White, BgCyan
}

func (l LogLevel) Tag() string {
	tagFg, tagBg := l.TagColor()
	return Colorize(" "+l.String()[:4]+" ", tagFg, tagBg)
}

// func Handle takes a [slog.Record] and appends data to a buffer that is
// written to an output stream via an [io.Writer]
//
// Implements [slog.Handler].Handle
func (l Logger) Handle(ctx context.Context, r slog.Record) error {
	if !l.Enabled(ctx, r.Level) {
		return nil
	}

	data := bytes.NewBuffer([]byte(LogLevel(r.Level).Tag() + " "))

	// Format Time
	data.WriteString(r.Time.Format(l.Options.TimeLayout) + " ")

	if len(l.Options.Prefix) > 0 {
		data.WriteString(l.Options.Prefix + " ")
	}
	// Format Message
	data.WriteString(r.Message + " ")

	// Format Attrs
	r.Attrs(func(a slog.Attr) bool {
		data.WriteString(fmt.Sprintf("%v: %v ", a.Key, a.Value))
		return true
	})

	data.WriteRune('\n')
	l.Mu.Lock()

	defer l.Mu.Unlock()
	if c, err := data.WriteTo(l.Writer); err != nil {
		return fmt.Errorf("write failed with code %v %v", c, err.Error())
	}

	return nil
}

// Implements [slog.Handler].Enabled
func (l Logger) Enabled(ctx context.Context, lvl slog.Level) bool {
	return lvl >= l.Options.Level
}

// Implements [slog.Handler].WithAttrs
func (l Logger) WithAttrs(attrs []slog.Attr) slog.Handler {
	l.Attrs = append(l.Attrs, attrs...)
	return l
}

// func WithGroup appends the prefix to the buffer
//
// Implements [slog.Handler].WithGroup
func (l Logger) WithGroup(name string) slog.Handler {
	l.Options.Prefix = fmt.Sprintf("%v.%v", l.Options.Prefix, name)
	return l
}

func (l Logger) Log(lvl LogLevel, msg interface{}, args ...interface{}) error {
	r := slog.NewRecord(time.Now(), slog.Level(lvl), fmt.Sprint(msg), 0)
	pc, _, _, ok := runtime.Caller(0)
	if ok {
		r.PC = pc
	}

	if err := l.Handle(l.Ctx, r); err != nil {
		return err
	} else {
		return nil
	}
}

func (l *Logger) Debug(msg interface{}, keyvals ...interface{}) {
	l.Log(DebugLevel, msg, keyvals...)
}

func (l *Logger) Info(msg interface{}, keyvals ...interface{}) {
	l.Log(InfoLevel, msg, keyvals...)
}

func (l *Logger) Warn(msg interface{}, keyvals ...interface{}) {
	l.Log(WarnLevel, msg, keyvals...)
}

func (l *Logger) Error(msg interface{}, keyvals ...interface{}) {
	l.Log(ErrorLevel, msg, keyvals...)
}

func (l *Logger) Fatal(msg interface{}, keyvals ...interface{}) {
	l.Log(FatalLevel, msg, keyvals...)
	os.Exit(1)
}

func (l *Logger) Print(msg interface{}, keyvals ...interface{}) {
	l.Log(NoLevel, msg, keyvals...)
}

func (l *Logger) SetLevel(lvl LogLevel) {
	l.Options.Level = slog.Level(lvl)
}

func NewLogger(w io.Writer, l LogLevel, p string) *Logger {
	return &Logger{
		Options: HandlerOpts{
			Level:      slog.Level(l),
			TimeLayout: time.RFC822Z,
			Prefix:     p,
		},
		Writer: w,
		Ctx:    context.Background(),
		Mu:     &sync.Mutex{},
	}
}

func DefaultLogger() *Logger {
	return NewLogger(os.Stdout, InfoLevel, "")
}
