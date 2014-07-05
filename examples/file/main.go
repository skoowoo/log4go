package main

import (
	log "log4go"
	"time"
)

func SetLog() {
	w := log.NewFileWriter()
	w.SetPathPattern("/tmp/logs/error%Y%M%D%H%m.log")

	log.Register(w)
	log.SetLevel(log.ERROR)
}

func main() {
	SetLog()
	defer log.Close()

	var name = "skoo"

	for {
		log.Debug("log4go by %s", name)
		log.Info("log4go by %s", name)
		log.Warn("log4go by %s", name)
		log.Error("log4go by %s", name)
		log.Fatal("log4go by %s", name)

		time.Sleep(time.Second * 1)
	}
}
