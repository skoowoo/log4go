package log4go

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"
)

var (
	LEVEL_FLAGS = [...]string{"[DEBUG]", "[ INFO]", "[ WARN]", "[ERROR]", "[FATAL]"}
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

const tunnel_size_default = 1024

type Record struct {
	time  string
	code  string
	info  string
	level int
}

func (r *Record) String() string {
	return fmt.Sprintf("%s %s <%s> %s\n", r.time, LEVEL_FLAGS[r.level], r.code, r.info)
}

type Writer interface {
	Init() error
	Write(*Record) error
}

type Rotater interface {
	Rotate() error
	SetPathPattern(string) error
}

type Flusher interface {
	Flush() error
}

type Logger struct {
	writers     []Writer
	tunnel      chan *Record
	level       int
	lastTime    int64
	lastTimeStr string
}

func NewLogger() *Logger {
	if logger_default != nil && takeup == false {
		takeup = true
		return logger_default
	}

	l := new(Logger)
	l.writers = make([]Writer, 0, 2)
	l.tunnel = make(chan *Record, tunnel_size_default)
	l.level = DEBUG

	go boostrapLogWriter(l)

	return l
}

func (l *Logger) Register(w Writer) {
	if err := w.Init(); err != nil {
		panic(err)
	}
	l.writers = append(l.writers, w)
}

func (l *Logger) SetLevel(lvl int) {
	l.level = lvl
}

func (l *Logger) Debug(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(DEBUG, fmt, args...)
}

func (l *Logger) Warn(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(WARNING, fmt, args...)
}

func (l *Logger) Info(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(INFO, fmt, args...)
}

func (l *Logger) Error(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(ERROR, fmt, args...)
}

func (l *Logger) Fatal(fmt string, args ...interface{}) {
	l.deliverRecordToWriter(FATAL, fmt, args...)
}

func (l *Logger) Close() {
	for {
		select {
		case r := <-l.tunnel:
			for _, w := range l.writers {
				if err := w.Write(r); err != nil {
					log.Println(err)
				}
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Println(err)
					}
				}
			}
		default:
			return
		}
	}
}

func (l *Logger) deliverRecordToWriter(level int, format string, args ...interface{}) {
	var inf, code string

	if level < l.level {
		return
	}

	if format != "" {
		inf = fmt.Sprintf(format, args...)
	} else {
		inf = fmt.Sprint(args...)
	}

	// source code, function and line num
	pc, _, line, ok := runtime.Caller(2)
	if ok {
		code = runtime.FuncForPC(pc).Name() + ":" + strconv.Itoa(line)
	}

	// format time
	now := time.Now()
	if now.Unix() != l.lastTime {
		l.lastTime = now.Unix()
		l.lastTimeStr = now.Format("2006/01/02 15:04:05")
	}

	r := &Record{
		info:  inf,
		code:  code,
		time:  l.lastTimeStr,
		level: level,
	}

	l.tunnel <- r
}

func boostrapLogWriter(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}

	r := <-logger.tunnel
	for _, w := range logger.writers {
		if err := w.Write(r); err != nil {
			log.Println(err)
		}
	}

	for {
		select {
		case r := <-logger.tunnel:
			for _, w := range logger.writers {
				if err := w.Write(r); err != nil {
					log.Println(err)
				}
			}

		case <-time.After(time.Millisecond * 500):
			for _, w := range logger.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Println(err)
					}
				}
			}

		case <-time.After(time.Second * 10):
			for _, w := range logger.writers {
				if r, ok := w.(Rotater); ok {
					if err := r.Rotate(); err != nil {
						log.Println(err)
					}
				}
			}
		}
	}
}

// default
var (
	logger_default *Logger
	takeup         = false
)

func SetLevel(lvl int) {
	logger_default.level = lvl
}

func Debug(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(DEBUG, fmt, args...)
}

func Warn(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(WARNING, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(INFO, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(ERROR, fmt, args...)
}

func Fatal(fmt string, args ...interface{}) {
	logger_default.deliverRecordToWriter(FATAL, fmt, args...)
}

func Register(w Writer) {
	logger_default.Register(w)
}

func Close() {
	logger_default.Close()
}

func init() {
	logger_default = NewLogger()
}
