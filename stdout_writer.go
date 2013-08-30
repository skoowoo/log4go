package log4go

import (
	"fmt"
	"os"
)

type StdoutW struct {
	name   string
	level  int
	rotate bool
}

func (w *StdoutW) Write(r *Record) error {
	fmt.Fprint(os.Stdout, r.String())
	return nil
}

func (w *StdoutW) RotateOrNot() bool {
	return w.rotate
}

func (w *StdoutW) Name() string {
	return w.name
}

func (w *StdoutW) Level() int {
	return w.level
}

func (w *StdoutW) Init(c *ConfigWriter) error {
	w.level = convLevel(c.Level)
	w.rotate = false
	return nil
}

func init() {
	addWriter(&StdoutW{name: "stdout"})
}
