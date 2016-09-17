package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"runtime"
	"strconv"
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
	//if you don't want to use log, you can call NewLoggerDiscard()
	// Logger = NewLoggerDiscard()

	//output log info to stdout
	Logger = NewLogger(os.Stdout, Config{})
	Logger.SetLogLevel(LOG_INFO)
}

type Config struct {
	Level int

	DispFuncLineInfo bool
	highlighting     bool

	// Rotate daily for log file
	Daily   bool
	MaxDays int64
}

type logger struct {
	log *log.Logger

	Config

	// for log file
	FileName string
}

func (l *logger) SetLogLevel(logLevel int) {
	l.Level = logLevel
}

func (l *logger) Debug(args ...interface{}) {
	l.output(LOG_DEBUG, args...)
}

func (l *logger) Debugf(f string, args ...interface{}) {
	l.outputf(LOG_DEBUG, f, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.output(LOG_INFO, args...)
}

func (l *logger) Infof(f string, args ...interface{}) {
	l.outputf(LOG_INFO, f, args...)
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
	if l.Level|l.Level > level {
		return
	}

	var s, funcLineInfo string
	if l.DispFuncLineInfo {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		funcLineInfo = "[" + filename + ":" + strconv.FormatInt(int64(line), 10) + "] "
	}
	if l.FileName == "" && l.highlighting {
		logColor := highlightTypeByLevel(level)
		s += "\033" + logColor + "m" + funcLineInfo + "[" + LOG_LEVEL_MAP[level] + "] " + fmt.Sprint(args...) + "\033[0m"
	} else {
		s += "[" + LOG_LEVEL_MAP[level] + "] " + fmt.Sprint(args...)
	}
	l.log.Output(2, s)
}

func (l *logger) outputf(level int, f string, args ...interface{}) {
	if l.Level|l.Level > level {
		return
	}

	var s, funcLineInfo string
	if l.DispFuncLineInfo {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		funcLineInfo = "[" + filename + ":" + strconv.FormatInt(int64(line), 10) + "] "
	}
	if l.FileName == "" && l.highlighting {
		logColor := highlightTypeByLevel(level)
		s = "\033" + logColor + "m" + funcLineInfo + "[" + LOG_LEVEL_MAP[level] + "] " + fmt.Sprintf(f, args...) + "\033[0m"
	} else {
		s += "[" + LOG_LEVEL_MAP[level] + "] " + fmt.Sprintf(f, args...)
	}
	l.log.Output(2, s)
}

func highlightTypeByLevel(t int) string {
	switch t {
	case LOG_DEBUG, LOG_INFO:
		return "[0;36"
	case LOG_WARNINIG:
		return "[0;33"
	case LOG_FATAL, LOG_ERROR:
		return "[0;31"
	}
	return "[0;37"
}

// NewLoggerByPath creates a new logger by fileName
func NewLoggerByFileName(fileName string, config Config) *logger {
	//return &logger{log.New(w, prefix, log.LstdFlags)}
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("NewLoggerByFileName error: %+v", err)
		return nil
	}
	//return NewLogger(f, config)
	log := NewLogger(f, config)
	log.FileName = fileName
	return log
}

// NewLoggerDiscard create a new logger.
// But do without anything
func NewLoggerDiscard() *logger {
	return NewLogger(ioutil.Discard, Config{})
}

// NewLogger creates a new logger.
// it write to io.Writer
// the prefix appears at the beginning of each generaged log line
func NewLogger(w io.Writer, config Config) *logger {
	log := &logger{
		log: log.New(w, "", log.LstdFlags),
	}
	if config.Level == 0 {
		config.Level = LOG_INFO
	}
	log.Config = config
	return log
}
