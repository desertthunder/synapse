package main

import (
	"context"
	"log/slog"
	"sync"
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
)

type ANSIColor string

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
	coloredText := fg.AddForegroundColor(s)
	return bg.AddBackgroundColor(coloredText)
}

type HOpts struct {
	Level slog.Level
}
type CustomHandler struct {
	Options HOpts
	Mu      *sync.Mutex
}

const (
	LevelDebug    slog.Level = -4
	LevelInfo     slog.Level = 0
	LevelWarn     slog.Level = 4
	LevelError    slog.Level = 8
	LevelCritical slog.Level = 12
)

func (h CustomHandler) Handle(ctx context.Context, r slog.Record) error {
	return nil
}

func (h CustomHandler) Enabled(ctx context.Context, l slog.Level) bool {
	return h.Options.Level <= l
}

func (h CustomHandler) WithAttrs(attrs []slog.Attr) CustomHandler {
	return h
}

func (h CustomHandler) WithGroup(name string) CustomHandler {
	return h
}
