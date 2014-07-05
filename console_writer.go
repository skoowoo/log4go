package log4go

import (
	"fmt"
	"os"
)

type ConsoleWriter struct {
}

func NewConsoleWriter() *ConsoleWriter {
	return &ConsoleWriter{}
}

func (w *ConsoleWriter) Write(r *Record) error {
	fmt.Fprint(os.Stdout, r.String())
	return nil
}

func (w *ConsoleWriter) Init() error {
	return nil
}
