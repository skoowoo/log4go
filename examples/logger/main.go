package main

import (
	"log4go"
)

func SetLog() (logger *log4go.Logger) {
	w := log4go.NewFileWriter()
	w.SetPathPattern("/tmp/logs/error%Y%M%D%H.log")

	logger = log4go.NewLogger()
	logger.Register(w)
	logger.SetLevel(log4go.ERROR)
	return
}

func main() {
	logger := SetLog()
	defer logger.Close()

	var name = "skoo"
	logger.Debug("log4go by %s", name)
	logger.Info("log4go by %s", name)
	logger.Warn("log4go by %s", name)
	logger.Error("log4go by %s", name)
	logger.Fatal("log4go by %s", name)
}
