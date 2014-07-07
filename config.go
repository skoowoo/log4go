package log4go

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type ConfFileWriter struct {
	LogPath string `json:"LogPath"`
	On      bool   `json:"On"`
}

type ConfConsoleWriter struct {
	On    bool `json:"On"`
	Color bool `json:"Color"`
}

type LogConfig struct {
	Level string            `json:"LogLevel"`
	FW    ConfFileWriter    `json:"FileWriter"`
	CW    ConfConsoleWriter `json:"ConsoleWriter"`
}

func SetupLogWithConf(file string) (err error) {
	var lc LogConfig

	cnt, err := ioutil.ReadFile(file)

	if err = json.Unmarshal(cnt, &lc); err != nil {
		return
	}

	if lc.FW.On {
		w := NewFileWriter()
		w.SetPathPattern(lc.FW.LogPath)
		Register(w)
	}

	if lc.CW.On {
		w := NewConsoleWriter()
		w.SetColor(lc.CW.Color)
		Register(w)
	}

	switch lc.Level {
	case "debug":
		SetLevel(DEBUG)

	case "info":
		SetLevel(INFO)

	case "warning":
		SetLevel(WARNING)

	case "error":
		SetLevel(ERROR)

	case "fatal":
		SetLevel(FATAL)

	default:
		err = errors.New("Invalid log level")
	}
	return
}
