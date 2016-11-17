package log4go

import (
	"fmt"
	"testing"

	"github.com/crius79/log4go"
)

func TestPrefix(t *testing.T) {
	InitLogger()
	defer CloseLogger()

	pre := log4go.PREFIX_CODE | log4go.PREFIX_LEVEL | log4go.PREFIX_TIME
	fmt.Printf("%v,%v,%v,%v \n", log4go.PREFIX_CODE, log4go.PREFIX_LEVEL, log4go.PREFIX_TIME, pre)

	fmt.Printf("%v,%v,%v \n", (pre&1) == 1, (pre&1<<0x1) == 1<<1, (pre&2<<0x1) == 2<<1)

	Logcode.Info("this ia access log")
	Logtime.Info("this ia create log")
	LogDebug.Info("this is debug log")
}

var (
	Logcode  *log4go.Logger
	Logtime  *log4go.Logger
	LogDebug *log4go.Logger
)

func InitLogger() error {
	Logcode = log4go.NewLogger()
	w1 := log4go.NewFileWriter()
	err1 := w1.SetPathPattern("logcode%Y%M%D.log")
	if err1 != nil {
		fmt.Printf("InitLogger err1= %s \n", err1.Error())
		return err1
	}
	Logcode.Register(w1)
	Logcode.SetLevel(log4go.INFO)
	Logcode.SetLayout("2006-01-02 15:04:05")
	Logcode.SetPrefix(log4go.PREFIX_CODE)

	Logtime = log4go.NewLogger()
	w2 := log4go.NewFileWriter()
	err2 := w2.SetPathPattern("logtime%Y%M%D.log")
	if err2 != nil {
		fmt.Printf("InitLogger err2= %s \n", err1.Error())
		return err1
	}
	Logtime.Register(w2)
	Logtime.SetLevel(log4go.INFO)
	Logtime.SetLayout("2006-01-02 15:04:05")
	Logtime.SetPrefix(log4go.PREFIX_TIME)

	LogDebug = log4go.NewLogger()
	w3 := log4go.NewConsoleWriter()
	LogDebug.Register(w3)
	LogDebug.SetLevel(log4go.DEBUG)
	LogDebug.SetLayout("2006-01-02 15:04:05")
	LogDebug.SetPrefix(log4go.PREFIX_TIME | log4go.PREFIX_LEVEL | log4go.PREFIX_CODE)

	return nil
}

func CloseLogger() {
	fmt.Println("CloseLogger")
	Logcode.Close()
	Logtime.Close()
	LogDebug.Close()
}
