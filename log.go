package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	LOG_DEBUG    = 1
	LOG_INFO     = 2
	LOG_WARNINIG = 3
	LOG_ERROR    = 4
	LOG_FATAL    = 5
)

var LOG_LEVEL_MAP = map[int]string{
	LOG_DEBUG:    "DEBUG",
	LOG_INFO:     "INFO",
	LOG_WARNINIG: "WARN",
	LOG_FATAL:    "FATAL",
}

var Logger *logger

func init() {
	Logger = NewLogger(os.Stdout, "", log.LstdFlags)
	Logger.SetLogLevel(LOG_DEBUG)
}

// NewLoggerByPath creates a new logger by fileName
func NewLoggerByFileName(fileName string) *logger {
	//return &logger{log.New(w, prefix, log.LstdFlags)}
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("NewLoggerByFileName error: %+v", err)
		return nil
	}
	return NewLogger(f, "opentsdb", log.LstdFlags)
}

// NewLoggerDiscard create a new logger.
// But do without anything
func NewLoggerDiscard() *logger {
	return NewLogger(ioutil.Discard, "", log.Lshortfile|log.LstdFlags)
}

// NewLogger creates a new logger.
// it write to io.Writer
// the prefix appears at the beginning of each generaged log line
func NewLogger(w io.Writer, prefix string, flag int) *logger {
	return &logger{
		log: log.New(w, prefix, flag),
	}
}

type logger struct {
	log   *log.Logger
	level int
}

func (l *logger) SetLogLevel(logLevel int) {
	l.level = logLevel
}

func (l *logger) Debug(args ...interface{}) {
	l.output(LOG_DEBUG, args...)
}

func (l *logger) Debugf(f string, args ...interface{}) {
	l.outputf(LOG_DEBUG, f, args...)
}

func (l *logger) Warn(args ...interface{}) {
	l.output(LOG_WARNINIG, args...)
}

func (l *logger) Warnf(f string, args ...interface{}) {
	l.outputf(LOG_WARNINIG, f, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.output(LOG_ERROR, args...)
}

func (l *logger) Errorf(f string, args ...interface{}) {
	l.outputf(LOG_ERROR, f, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.output(LOG_FATAL, args)
}

func (l *logger) Fatalf(f string, args ...interface{}) {
	l.outputf(LOG_FATAL, f, "")
}

func (l *logger) output(level int, args ...interface{}) {
	if l.level|l.level > level {
		return
	}
	s := "[" + LOG_LEVEL_MAP[level] + "] " + fmt.Sprint(args...)
	l.log.Output(2, s)
}

func (l *logger) outputf(level int, f string, args ...interface{}) {
	if l.level|l.level > level {
		return
	}
	s := "[" + LOG_LEVEL_MAP[level] + "] " + fmt.Sprintf(f, args...)
	l.log.Output(2, s)
}
