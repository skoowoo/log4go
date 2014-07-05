package main

import (
	log "log4go"
)

func SetLog() {
	w := log.NewFileWriter()
	w.SetPathPattern("/tmp/logs/error%Y%M%D%H.log")

	log.Register(w)
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
