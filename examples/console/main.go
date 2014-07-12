package main

import (
	log "log4go"
)

func SetLog() {
	w := log.NewConsoleWriter()
	w.SetColor(true)

	log.Register(w)
	log.SetLevel(log.DEBUG)
	log.SetLayout("2006-01-02 15:04:05")
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
