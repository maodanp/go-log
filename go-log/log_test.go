package log

import (
	"os"
	"testing"
)

type Foo struct{}
type Bar struct {
	one string
	two int
}

func TestNewLogger(t *testing.T) {
	log := NewLogger(os.Stdout, Config{highlighting: true, DispFuncCall: true})
	log.SetLogLevel(LOG_DEBUG)
	f := &Foo{}
	b := &Bar{"asfd", 3}
	for _, test := range []interface{}{
		1234,
		"asdf",
		f,
		b,
	} {
		log.Info("hello", test)
		log.Debugf("something, %v", test)
		log.Warn(test)
		log.Warnf("something, %v", test)
	}
}

func TestLoggerFileName(t *testing.T) {
	log := NewLoggerByFileName("/tmp/test.log", Config{})
	log.SetLogLevel(LOG_DEBUG)
	f := &Foo{}
	b := &Bar{"asfd", 3}
	for _, test := range []interface{}{
		1234,
		"asdf",
		f,
		b,
	} {
		log.Debug("hello", test)
		log.Debugf("something, %v", test)
		log.Warn(test)
		log.Warnf("something, %v", test)
	}
}
