package log

import (
	"log"
	"os"
	"testing"
)

type Foo struct{}
type Bar struct {
	one string
	two int
}

func TestNewLogger(t *testing.T) {
	log := NewLogger(os.Stdout, "", log.Lshortfile|log.LstdFlags)
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

func TestLoggerFileName(t *testing.T) {
	log := NewLoggerByFileName("./test.log")
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
