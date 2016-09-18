package log

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
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

const (
	DEFAULT_MAX_DAYS = 6 //default remain 6 days log file
)

type Config struct {
	Level int

	DispFuncCall bool
	highlighting bool

	// Rotate daily for log file
	DailyRotate bool
	MaxDays     int64
}

type logger struct {
	log *log.Logger

	Config

	// private info for log file
	fileName string
	fd       *os.File

	// like "project.log.20160913"
	// "project" is fileNameOnly
	// ".log" is suffix
	// "20160913" is dailySuffix
	fileNameOnly string
	suffix       string
	dailySuffix  string

	lock sync.Mutex
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
	h, _ := formatTimeHeader(time.Now())
	if l.DispFuncCall {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		funcLineInfo = "[" + filename + ":" + strconv.FormatInt(int64(line), 10) + "] "
	}
	if l.fileName == "" && l.highlighting {
		logColor := highlightTypeByLevel(level)
		s = "\033" + logColor + "m" + string(h) + funcLineInfo + "[" + LOG_LEVEL_MAP[level] + "] " + fmt.Sprint(args...) + "\033[0m"
	} else {
		// rotate only for log file
		l.dailyRotate()
		s = string(h) + funcLineInfo + "[" + LOG_LEVEL_MAP[level] + "] " + fmt.Sprint(args...)
	}
	l.log.Output(2, s)
}

func (l *logger) outputf(level int, f string, args ...interface{}) {
	if l.Level|l.Level > level {
		return
	}

	var s, funcLineInfo string
	h, _ := formatTimeHeader(time.Now())
	if l.DispFuncCall {
		_, file, line, ok := runtime.Caller(2)
		if !ok {
			file = "???"
			line = 0
		}
		_, filename := path.Split(file)
		funcLineInfo = "[" + filename + ":" + strconv.FormatInt(int64(line), 10) + "] "
	}
	if l.fileName == "" && l.highlighting {
		logColor := highlightTypeByLevel(level)
		s = "\033" + logColor + "m" + string(h) + funcLineInfo + "[" + LOG_LEVEL_MAP[level] + "] " + fmt.Sprintf(f, args...) + "\033[0m"
	} else {
		// rotate only for log file
		l.dailyRotate()
		s = string(h) + funcLineInfo + "[" + LOG_LEVEL_MAP[level] + "] " + fmt.Sprintf(f, args...)
	}
	l.log.Output(2, s)
}

func (l *logger) dailyRotate() error {
	if l.DailyRotate == false {
		return nil
	}

	l.lock.Lock()
	defer l.lock.Unlock()

	suffix := genDayTime(time.Now())
	if suffix != l.dailySuffix {
		err := l.doRotate(suffix)
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *logger) doRotate(suffix string) error {
	l.fd.Close()

	lastFileName := l.fileName + "." + l.dailySuffix
	err := os.Rename(l.fileName, lastFileName)
	if err != nil {
		return err
	}

	err = l.SetOutputByFile(l.fileName)
	if err != nil {
		return err
	}
	l.dailySuffix = suffix

	l.deleteOldLog()
	return nil
}

func (l *logger) deleteOldLog() {
	dir := filepath.Dir(l.fileName)
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) (returnErr error) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Fprintf(os.Stderr, "Unable to delete old log '%s', error: %v\n", path, r)
			}
		}()

		if info == nil {
			return
		}

		if !info.IsDir() && info.ModTime().Add(24*time.Hour*time.Duration(l.MaxDays)).Before(time.Now()) {
			if strings.HasPrefix(filepath.Base(path), filepath.Base(l.fileNameOnly)) &&
				strings.HasSuffix(filepath.Base(path), l.suffix) {
				os.Remove(path)
			}
		}
		return
	})
	return
}

func (l *logger) SetOutputByFile(fileName string) error {
	//return &logger{log.New(w, prefix, log.LstdFlags)}
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	l.log = log.New(f, "", 0)
	l.fd = f
	l.suffix = filepath.Ext(l.fileName)
	l.fileNameOnly = strings.TrimSuffix(l.fileName, l.suffix)
	if l.suffix == "" {
		l.suffix = ".log"
	}
	return nil
}

// NewLoggerByPath creates a new logger by fileName
func NewLoggerByFileName(fileName string, config Config) *logger {
	log := &logger{}
	log.SetOutputByFile(fileName)
	log.Config = config
	if log.Level == 0 {
		log.Level = LOG_DEBUG
	}

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
		log: log.New(w, "", 0),
		//log: log.New(w, "", 0),
	}
	log.Config = config
	if log.Level == 0 {
		log.Level = LOG_INFO
	}

	if log.DailyRotate && log.MaxDays == 0 {
		log.MaxDays = DEFAULT_MAX_DAYS
	}

	return log
}

var Logger *logger

func init() {
	//if you don't want to use log, you can call NewLoggerDiscard()
	// Logger = NewLoggerDiscard()

	//output log info to stdout
	Logger = NewLogger(os.Stdout, Config{})
	Logger.SetLogLevel(LOG_INFO)
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
