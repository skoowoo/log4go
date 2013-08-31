package main

import (
    log "log4go"
)

func main() {
    log.LoadConfigFile("log4go.json")
    defer log.Close()

    var str = "this is a test"
    log.Debug("test: %s", str)
    log.Info("test: %s", str)
    log.Warn("test: %s", str)
    log.Error("test: %s", str)
    log.Critical("test: %s", str)
}
