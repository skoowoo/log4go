package log4go

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"path"
)

type WConf interface {
	Name() string
	Init(*ConfigWriter) error
}

type ConfigWriter struct {
	Name    string
	Enable  bool
	Level   string
	LogPath string
	Rotate  bool
}

type ConfigFile struct {
	Writers []ConfigWriter
}

var (
	configStruct ConfigFile
	writers      map[string]WConf = make(map[string]WConf, 2)
	writerLevels map[string]int   = make(map[string]int, 5)
)

func init() {
	writerLevels["debug"] = DEBUG
	writerLevels["info"] = INFO
	writerLevels["warning"] = WARNING
	writerLevels["error"] = ERROR
	writerLevels["critical"] = CRITICAL

}

func addWriter(w WConf) {
	name := w.Name()
	if name == "" {
		panic("writer must have a name")
	}

	if _, ok := writers[name]; ok {
		panic(fmt.Errorf("\"%s\" writer exist", name))
	}
	writers[name] = w
}

func convLevel(l string) int {
	return writerLevels[l]
}

func LoadConfigFile(path string) {
	tmp, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(tmp, &configStruct.Writers); err != nil {
		panic(err)
	}

	if err := verifyConfig(&configStruct); err != nil {
		panic(err)
	}

	if err := bootstrapLogger(&configStruct); err != nil {
		panic(err)
	}
}

func verifyConfig(conf *ConfigFile) error {
	for _, wc := range conf.Writers {
		if _, ok := writers[wc.Name]; !ok {
			return fmt.Errorf("\"%s\" writer don't exist", wc.Name)
		}

		if _, ok := writerLevels[wc.Level]; !ok {
			return fmt.Errorf("\"%s\" writer, level \"%s\" invalid", wc.Name, wc.Level)
		}

		wc.LogPath = path.Clean(wc.LogPath)
	}
	return nil
}

func bootstrapLogger(conf *ConfigFile) error {
	logger := NewLoggerDefault()
	if logger == nil {
		return errors.New("new logger failed")
	}

	for _, wc := range conf.Writers {
		if !wc.Enable {
			continue
		}
		w := writers[wc.Name]
		if err := w.Init(&wc); err != nil {
			return err
		}

		// register into logger
		logger.RegisterWriter(w.Name(), w.(Writer))
	}

	return nil
}
