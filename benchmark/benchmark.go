package main

import (
	log "log4go"
	"runtime"
	"time"
)

const N = 1000000

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	log.LoadConfigFile("log4go.json")
	defer log.Close()

	var s = "abcdefg"

	start := time.Now()

	for i := 0; i < N; i++ {
		log.Debug("start benchmar, i: %d, this is a benchmark test. %s", i, s)
	}

	delay := time.Now().Sub(start)

	println(delay)
}
