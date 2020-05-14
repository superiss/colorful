package colorful

import (
	"fmt"
	"io"
	"os"
)

const (
	ldate = 1 << iota
	lstd  = ldate
)

// Logger struct
type Logger struct {
	out io.Writer
}

// New logger
func New() *Logger {
	return &Logger{
		out: os.Stderr,
	}
}

func itoa(buf *[]byte, i int, wid int) {
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}

func (l *Logger) outputColorful(colorDef, s string) error {

	var buff []byte
	buff = append(buff, s...)
	if colorDef == "red" {
		buff = []byte(Sprint("@r%v", string(buff)))
	} else if colorDef == "yellow" {
		buff = []byte(Sprint("@y%v", string(buff)))
	} else if colorDef == "green" {
		buff = []byte(Sprint("@g%v", string(buff)))
	}
	_, err := l.out.Write(buff)
	return err
}

// Info func
func (l *Logger) Info(v ...interface{}) {
	l.outputColorful("", fmt.Sprint(v...))
}

// Warn func
func (l *Logger) Warn(v ...interface{}) {
	l.outputColorful("yellow", fmt.Sprint(v...))
}

// Error func
func (l *Logger) Error(v ...interface{}) {
	l.outputColorful("red", fmt.Sprint(v...))
}

// Success func
func (l *Logger) Success(v ...interface{}) {
	l.outputColorful("green", fmt.Sprint(v...))
}

// Fatal func
func (l *Logger) Fatal(v ...interface{}) {
	l.outputColorful("red", fmt.Sprint(v...))
	os.Exit(1)
}
