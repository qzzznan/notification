package log

import (
	"io"
	"sync"
)

type LoggerWriter struct {
	io.Writer
	sync.Mutex
}

func (l *LoggerWriter) SetWriter(w io.Writer) {
	l.Lock()
	defer l.Unlock()
	l.Writer = w
}
