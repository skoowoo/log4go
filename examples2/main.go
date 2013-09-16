package main

import (
    log "log4go"
)

func main() {
    conf := &log.ConfigWriter{}
    conf.Enable = true
    conf.Level = "debug"
    conf.LogPath = "/tmp/test2.log"
    conf.Name = "file"
    conf.Rotate = true

    log.NewLogger(conf)
    defer log.Close()

    var str = "this is a test"
    log.Debug("test: %s", str)
    log.Info("test: %s", str)
    log.Warn("test: %s", str)
    log.Error("test: %s", str)
    log.Critical("test: %s", str)
}
