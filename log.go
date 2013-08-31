package log4go

import (
	"fmt"
	"log"
	"runtime"
	"strconv"
	"time"
)

var LEVEL_FLAGS = [...]string{"DEBG", "INFO", "WARN", "EROR", "CRIT"}

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	CRITICAL
)

const tunnel_size_default = 1024

type Record struct {
	time  string
	code  string
	info  string
	level int
}

func (r *Record) String() string {
	return fmt.Sprintf("%s [%s] %s %s\n", r.time, LEVEL_FLAGS[r.level], r.code, r.info)
}

type Writer interface {
	Write(*Record) error
	RotateOrNot() bool
	Name() string
	Level() int
}

type Rotater interface {
	Rotate(string)
}

type Flusher interface {
	Flush() error
}

var logger_default *Logger

type Logger struct {
	writers         map[string]Writer
	tunnel          chan *Record
	rotate          chan string
	exit            chan bool
	lastLogTimeSecs int64
	lastLogTimeStr  string
	currentDay      int
}

// deliver all records to every writer
func logRecordToWriters(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}

	for {
		select {
		case r := <-logger.tunnel:
			for _, w := range logger.writers {
				if w.Level() > r.level {
					continue
				}
				if err := w.Write(r); err != nil {
					log.Println(err)
				}
			}
		case suffix := <-logger.rotate:
			for _, w := range logger.writers {
				if w.RotateOrNot() == false {
					continue
				}
				if r, ok := w.(Rotater); ok {
					r.Rotate(suffix)
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
		case <-logger.exit:
			return
		}
	}
}

func NewLoggerDefault() *Logger {
	logger_default = new(Logger)
	logger_default.writers = make(map[string]Writer, 1)
	logger_default.rotate = make(chan string, 1)
	logger_default.exit = make(chan bool) // blocking channel
	logger_default.tunnel = make(chan *Record, tunnel_size_default)
	logger_default.currentDay = time.Now().Day()

	go logRecordToWriters(logger_default)

	return logger_default
}

// register a writer instance into logger
func (l *Logger) RegisterWriter(name string, w Writer) {
	if _, exist := l.writers[name]; exist {
		panic(fmt.Errorf("\"%s\" duplicate writer", name))
	}

	l.writers[name] = w
}

func (l *Logger) formatRecordToTunnel(level int, format string, args ...interface{}) {
	var inf, code string

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
	if now.Unix() != l.lastLogTimeSecs {
		l.lastLogTimeSecs = now.Unix()
		l.lastLogTimeStr = now.Format("2006/01/02 15:04:05")
	}

	// rotate
	if now.Day() != l.currentDay {
		l.currentDay = now.Day()
		l.rotate <- now.Format("2006-01-02")
	}

	r := &Record{
		info:  inf,
		code:  code,
		time:  l.lastLogTimeStr,
		level: level,
	}

	l.tunnel <- r
}

// flush, send exit signal to writer goroutine, then handle all buffered records. 
func (l *Logger) flush() {
	l.exit <- true
	for {
		select {
		case r := <-l.tunnel:
			for _, w := range l.writers {
				if w.Level() > r.level {
					continue
				}
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

// before program exit, you should call the Close() to write buffered records.
func Close() {
	logger_default.flush()
}

func Debug(fmt string, args ...interface{}) {
	logger_default.formatRecordToTunnel(DEBUG, fmt, args...)
}

func Info(fmt string, args ...interface{}) {
	logger_default.formatRecordToTunnel(INFO, fmt, args...)
}

func Warn(fmt string, args ...interface{}) {
	logger_default.formatRecordToTunnel(WARNING, fmt, args...)
}

func Error(fmt string, args ...interface{}) {
	logger_default.formatRecordToTunnel(ERROR, fmt, args...)
}

func Critical(fmt string, args ...interface{}) {
	logger_default.formatRecordToTunnel(CRITICAL, fmt, args...)
}

func Log(level int, args ...interface{}) {
	logger_default.formatRecordToTunnel(level, "", args...)
}

func LogStruct() {

}

func LogSlice() {

}

func LogMap() {

}

func init() {

}
