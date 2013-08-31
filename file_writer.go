package log4go

import (
	"bufio"
	"errors"
	"log"
	"os"
	"path"
)

type FileW struct {
	name          string
	level         int
	rotate        bool
	dir           string
	fileName      string
	file          *os.File
	fileBufWriter *bufio.Writer
	lastSuffix    string
}

func (w *FileW) Write(r *Record) error {
	if w.fileBufWriter == nil {
		return errors.New("no opened file")
	}
	if _, err := w.fileBufWriter.WriteString(r.String()); err != nil {
		return err
	}
	return nil
}

func (w *FileW) RotateOrNot() bool {
	return w.rotate
}

func (w *FileW) Name() string {
	return w.name
}

func (w *FileW) Level() int {
	return w.level
}

func (w *FileW) Init(c *ConfigWriter) error {
	w.level = convLevel(c.Level)
	w.rotate = c.Rotate
	w.dir = path.Dir(c.LogPath)
	w.fileName = path.Base(c.LogPath)

	if err := os.MkdirAll(w.dir, 0755); err != nil {
		if !os.IsExist(err) {
			return err
		}
	}

	fileName := w.dir + "/" + w.fileName
	if file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		return err
	} else {
		w.file = file
	}

	w.fileBufWriter = bufio.NewWriterSize(w.file, 8192)
	if w.fileBufWriter == nil {
		return errors.New("new fileBufWriter failed.")
	}

	return nil
}

func (w *FileW) Rotate(suffix string) {
	if w.lastSuffix == suffix {
		return
	}

	if err := w.file.Close(); err != nil {
		log.Println(err)
	}
	w.file = nil
	w.fileBufWriter = nil

	fileName := w.dir + "/" + w.fileName
	newName := w.dir + "/" + w.fileName + "." + suffix
	if err := os.Rename(fileName, newName); err != nil {
		log.Println(err)
	}

	w.lastSuffix = suffix

	if file, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644); err != nil {
		log.Println(err)
		return
	} else {
		w.file = file
	}

	w.fileBufWriter = bufio.NewWriterSize(w.file, 8192)
	if w.fileBufWriter == nil {
		log.Println("new fileBufWriter failed.")
	}
}

func (w *FileW) Flush() error {
	if w.fileBufWriter != nil {
		return w.fileBufWriter.Flush()
	}
	return nil
}

func init() {
	addWriter(&FileW{name: "file"})
}
