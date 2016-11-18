package log4go

import (
	"fmt"
	"log"
	"path"
	"runtime"
	"strconv"
	"sync"
	"time"
)

var (
	LEVEL_FLAGS = [...]string{"DEBUG", " INFO", " WARN", "ERROR", "FATAL"}
	recordPool  *sync.Pool
)

const (
	DEBUG = iota
	INFO
	WARNING
	ERROR
	FATAL
)

const (
	//   1    1     1
	// TIME LEVEL CODE INFO
	PREFIX_NONE  = iota
	PREFIX_CODE  = 0x1
	PREFIX_LEVEL = 1 << 0x1
	PREFIX_TIME  = 2 << 0x1
)

const tunnel_size_default = 1024

type Record struct {
	time   string
	code   string
	info   string
	level  int
	prefix int
}

func (r *Record) String() string {
	//fmt.Printf("r.prfix=%v\n ", r.prefix)
	var verbose string = r.info + "\n"
	if (r.prefix & 0x1) == 0x1 {
		verbose = r.code + " " + verbose
	}

	if (r.prefix & (1 << 0x1)) == (1 << 0x1) {
		verbose = "[" + LEVEL_FLAGS[r.level] + "] " + verbose
	}

	if (r.prefix & (2 << 0x1)) == (2 << 0x1) {
		verbose = "[" + r.time + "] " + verbose
	}
	return verbose

	//return fmt.Sprintf("%s [%s] <%s> %s\n", r.time, LEVEL_FLAGS[r.level], r.code, r.info)
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
	prefix      int
	lastTime    int64
	lastTimeStr string
	c           chan bool
	layout      string
}

func NewLogger() *Logger {
	if logger_default != nil && takeup == false {
		takeup = true
		return logger_default
	}

	l := new(Logger)
	l.writers = make([]Writer, 0, 2)
	l.tunnel = make(chan *Record, tunnel_size_default)
	l.c = make(chan bool, 1)
	l.level = DEBUG
	l.layout = "2006/01/02 15:04:05"
	l.prefix = PREFIX_TIME | PREFIX_LEVEL | PREFIX_CODE

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

func (l *Logger) SetLayout(layout string) {
	l.layout = layout
}

func (l *Logger) SetPrefix(prefix int) {
	l.prefix = prefix
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
	close(l.tunnel)
	<-l.c

	for _, w := range l.writers {
		if f, ok := w.(Flusher); ok {
			if err := f.Flush(); err != nil {
				log.Println(err)
			}
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

	// source code, file and line num
	_, file, line, ok := runtime.Caller(2)
	if ok {
		code = path.Base(file) + ":" + strconv.Itoa(line)
	}

	// format time
	now := time.Now()
	if now.Unix() != l.lastTime {
		l.lastTime = now.Unix()
		l.lastTimeStr = now.Format(l.layout)
	}

	r := recordPool.Get().(*Record)
	r.info = inf
	r.code = code
	r.time = l.lastTimeStr
	r.level = level
	r.prefix = l.prefix

	l.tunnel <- r
}

func boostrapLogWriter(logger *Logger) {
	if logger == nil {
		panic("logger is nil")
	}

	var (
		r  *Record
		ok bool
	)

	if r, ok = <-logger.tunnel; !ok {
		logger.c <- true
		return
	}

	for _, w := range logger.writers {
		if err := w.Write(r); err != nil {
			log.Println(err)
		}
	}

	flushTimer := time.NewTimer(time.Millisecond * 500)
	rotateTimer := time.NewTimer(time.Second * 10)

	for {
		select {
		case r, ok = <-logger.tunnel:
			if !ok {
				logger.c <- true
				return
			}

			for _, w := range logger.writers {
				if err := w.Write(r); err != nil {
					log.Println(err)
				}
			}

			recordPool.Put(r)

		case <-flushTimer.C:
			for _, w := range logger.writers {
				if f, ok := w.(Flusher); ok {
					if err := f.Flush(); err != nil {
						log.Println(err)
					}
				}
			}
			flushTimer.Reset(time.Millisecond * 1000)

		case <-rotateTimer.C:
			for _, w := range logger.writers {
				if r, ok := w.(Rotater); ok {
					if err := r.Rotate(); err != nil {
						log.Println(err)
					}
				}
			}
			rotateTimer.Reset(time.Second * 10)
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

func SetLayout(layout string) {
	logger_default.layout = layout
}

func SetPrefix(prefix int) {
	logger_default.prefix = prefix
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
	recordPool = &sync.Pool{New: func() interface{} {
		return &Record{}
	}}
}
