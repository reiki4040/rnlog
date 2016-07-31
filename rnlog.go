package rnlog

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"
)

const (
	LEVEL_TRACE = iota
	LEVEL_DEBUG
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR

	LEVEL_STR_NOTICE = "NOTICE"
	LEVEL_STR_FATAL  = "FATAL"

	FATAL_LOG = `{"level":"FATAL","msg":"failed marshal to json"}`
)

var (
	levelMap = map[int]string{
		LEVEL_TRACE: "TRACE",
		LEVEL_DEBUG: "DEBUG",
		LEVEL_INFO:  "INFO",
		LEVEL_WARN:  "WARN",
		LEVEL_ERROR: "ERROR",
	}

	EMPTY_ITEMS = make(map[string]interface{}, 0)

	ERR_INVALID_LOG_LEVEL = errors.New("invalid log level")
)

type Log struct {
	Time  string                 `json:"time"`
	Level string                 `json:"level"`
	Msg   string                 `json:"msg"`
	Items map[string]interface{} `json:"items,omitempty"`
}

func NewRNLoggerWithLogger(logger func(args ...interface{})) *RNLogger {
	return &RNLogger{
		logger: logger,
	}
}

type RNLogger struct {
	level  int
	logger func(args ...interface{})
}

func (l *RNLogger) ChangeLevel(level int) error {
	err := ValidateLogLevel(level)
	if err != nil {
		return err
	} else {
		l.level = level
	}

	return nil
}

func (l *RNLogger) Trace(msg string) {
	l.log(LEVEL_TRACE, msg, nil)
}
func (l *RNLogger) Debug(msg string) {
	l.log(LEVEL_DEBUG, msg, nil)
}
func (l *RNLogger) Info(msg string) {
	l.log(LEVEL_INFO, msg, nil)
}
func (l *RNLogger) Warn(msg string) {
	l.log(LEVEL_WARN, msg, nil)
}
func (l *RNLogger) Error(msg string) {
	l.log(LEVEL_ERROR, msg, nil)
}

func (l *RNLogger) Tracef(format string, args ...interface{}) {
	l.log(LEVEL_TRACE, fmt.Sprintf(format, args...), nil)
}
func (l *RNLogger) Debugf(format string, args ...interface{}) {
	l.log(LEVEL_DEBUG, fmt.Sprintf(format, args...), nil)
}
func (l *RNLogger) Infof(format string, args ...interface{}) {
	l.log(LEVEL_INFO, fmt.Sprintf(format, args...), nil)
}
func (l *RNLogger) Warnf(format string, args ...interface{}) {
	l.log(LEVEL_WARN, fmt.Sprintf(format, args...), nil)
}
func (l *RNLogger) Errorf(format string, args ...interface{}) {
	l.log(LEVEL_ERROR, fmt.Sprintf(format, args...), nil)
}

func (l *RNLogger) TraceItem(items map[string]interface{}, msg string) {
	l.log(LEVEL_TRACE, msg, items)
}
func (l *RNLogger) DebugItem(items map[string]interface{}, msg string) {
	l.log(LEVEL_DEBUG, msg, items)
}
func (l *RNLogger) InfoItem(items map[string]interface{}, msg string) {
	l.log(LEVEL_INFO, msg, items)
}
func (l *RNLogger) WarnItem(items map[string]interface{}, msg string) {
	l.log(LEVEL_WARN, msg, items)
}
func (l *RNLogger) ErrorItem(items map[string]interface{}, msg string) {
	l.log(LEVEL_ERROR, msg, items)
}

func (l *RNLogger) TracefItem(items map[string]interface{}, format string, args ...interface{}) {
	l.log(LEVEL_TRACE, fmt.Sprintf(format, args...), items)
}
func (l *RNLogger) DebugfItem(items map[string]interface{}, format string, args ...interface{}) {
	l.log(LEVEL_DEBUG, fmt.Sprintf(format, args...), items)
}
func (l *RNLogger) InfofItem(items map[string]interface{}, format string, args ...interface{}) {
	l.log(LEVEL_INFO, fmt.Sprintf(format, args...), items)
}
func (l *RNLogger) WarnfItem(items map[string]interface{}, format string, args ...interface{}) {
	l.log(LEVEL_WARN, fmt.Sprintf(format, args...), items)
}
func (l *RNLogger) ErrorfItem(items map[string]interface{}, format string, args ...interface{}) {
	l.log(LEVEL_ERROR, fmt.Sprintf(format, args...), items)
}

func (l *RNLogger) Fatal(msg string) {
	l.Logging(LEVEL_STR_FATAL, msg, nil)
	os.Exit(1)
}
func (l *RNLogger) Fatalf(format string, args ...interface{}) {
	l.Logging(LEVEL_STR_FATAL, fmt.Sprintf(format, args...), nil)
	os.Exit(1)
}
func (l *RNLogger) FatalItem(items map[string]interface{}, msg string) {
	l.Logging(LEVEL_STR_FATAL, msg, items)
	os.Exit(1)
}
func (l *RNLogger) FatalfItem(items map[string]interface{}, format string, args ...interface{}) {
	l.Logging(LEVEL_STR_FATAL, fmt.Sprintf(format, args...), items)
	os.Exit(1)
}

func (l *RNLogger) Notice(msg string) {
	l.Logging(LEVEL_STR_NOTICE, msg, nil)
}
func (l *RNLogger) Noticef(format string, args ...interface{}) {
	l.Logging(LEVEL_STR_NOTICE, fmt.Sprintf(format, args...), nil)
}
func (l *RNLogger) NoticeItem(items map[string]interface{}, msg string) {
	l.Logging(LEVEL_STR_NOTICE, msg, items)
}
func (l *RNLogger) NoticefItem(items map[string]interface{}, format string, args ...interface{}) {
	l.Logging(LEVEL_STR_NOTICE, fmt.Sprintf(format, args...), items)
}

func (l *RNLogger) log(level int, msg string, items map[string]interface{}) {
	if l.level > level {
		return
	}

	levelStr := levelMap[level]

	l.Logging(levelStr, msg, items)
}

func (l *RNLogger) Logging(level string, msg string, items map[string]interface{}) {
	if items == nil {
		items = EMPTY_ITEMS
	}

	log := Log{
		Time:  time.Now().Format(time.RFC3339),
		Level: level,
		Msg:   msg,
		Items: items,
	}

	json, err := json.Marshal(log)
	if err != nil {
		l.logger(FATAL_LOG)
		return
	}

	l.logger(string(json))
}

func ValidateLogLevel(level int) error {
	if LEVEL_TRACE > level && LEVEL_ERROR < level {
		return ERR_INVALID_LOG_LEVEL
	} else {
		return nil
	}
}

// Default logger Stdout with fmt.Println()
var std = &RNLogger{logger: func(args ...interface{}) { fmt.Println(args...) }}

func ChangeLevel(level int) error {
	return std.ChangeLevel(level)
}

func Trace(msg string) {
	std.Trace(msg)
}
func Debug(msg string) {
	std.Debug(msg)
}
func Info(msg string) {
	std.Info(msg)
}
func Warn(msg string) {
	std.Warn(msg)
}
func Error(msg string) {
	std.Error(msg)
}

func Tracef(format string, args ...interface{}) {
	std.Tracef(format, args...)
}
func Debugf(format string, args ...interface{}) {
	std.Debugf(format, args...)
}
func Infof(format string, args ...interface{}) {
	std.Infof(format, args...)
}
func Warnf(format string, args ...interface{}) {
	std.Warnf(format, args...)
}
func Errorf(format string, args ...interface{}) {
	std.Errorf(format, args...)
}

func TraceItem(items map[string]interface{}, msg string) {
	std.TraceItem(items, msg)
}
func DebugItem(items map[string]interface{}, msg string) {
	std.DebugItem(items, msg)
}
func InfoItem(items map[string]interface{}, msg string) {
	std.InfoItem(items, msg)
}
func WarnItem(items map[string]interface{}, msg string) {
	std.WarnItem(items, msg)
}
func ErrorItem(items map[string]interface{}, msg string) {
	std.ErrorItem(items, msg)
}

func TracefItem(items map[string]interface{}, format string, args ...interface{}) {
	std.TracefItem(items, format, args...)
}
func DebugfItem(items map[string]interface{}, format string, args ...interface{}) {
	std.DebugfItem(items, format, args...)
}
func InfofItem(items map[string]interface{}, format string, args ...interface{}) {
	std.InfofItem(items, format, args...)
}
func WarnfItem(items map[string]interface{}, format string, args ...interface{}) {
	std.WarnfItem(items, format, args...)
}
func ErrorfItem(items map[string]interface{}, format string, args ...interface{}) {
	std.ErrorfItem(items, format, args...)
}

func Fatal(msg string) {
	std.Fatal(msg)
}
func Fatalf(format string, args ...interface{}) {
	std.Fatalf(format, args...)
}
func FatalItem(items map[string]interface{}, msg string) {
	std.FatalItem(items, msg)
}
func FatalfItem(items map[string]interface{}, format string, args ...interface{}) {
	std.FatalfItem(items, format, args...)
}

func Notice(msg string) {
	std.Notice(msg)
}
func Noticef(format string, args ...interface{}) {
	std.Noticef(format, args...)
}
func NoticeItem(items map[string]interface{}, msg string) {
	std.NoticeItem(items, msg)
}
func NoticefItem(items map[string]interface{}, format string, args ...interface{}) {
	std.NoticefItem(items, format, args...)
}

func Logging(level string, msg string, items map[string]interface{}) {
	std.Logging(level, msg, items)
}
