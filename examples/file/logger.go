package logger

import (
	log4go "github.com/skoo87/log4go"
)

func GetLogger(filename string, level int) (logger *log4go.Logger) {
	w := log4go.NewFileWriter()
	//	filename = filename + "_%Y%M%D-%H"
	//	w.SetPathPattern(filename)
	pattern := "-%Y%M%D%H"
	w.SetFileName(filename)
	w.SetPathPattern(pattern)

	logger = log4go.NewLogger()

	if w == nil || logger == nil {
		panic("Init logger failed!")
	}

	logger.Register(w)
	logger.SetLevel(level)
	return
}
