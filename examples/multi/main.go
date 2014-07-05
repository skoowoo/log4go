package main

import (
	log "log4go"
)

func SetLog() {
	w1 := log.NewFileWriter()
	w1.SetPathPattern("/tmp/logs/error%Y%M%D%H.log")

	w2 := log.NewConsoleWriter()

	log.Register(w1)
	log.Register(w2)
	log.SetLevel(log.ERROR)
}

func main() {
	SetLog()
	defer log.Close()

	var name = "skoo"
	log.Debug("log4go by %s", name)
	log.Info("log4go by %s", name)
	log.Warn("log4go by %s", name)
	log.Error("log4go by %s", name)
	log.Fatal("log4go by %s", name)
}
