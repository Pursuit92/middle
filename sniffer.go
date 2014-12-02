package httputil

import (
	"bytes"
	"net/http"
)

type SnifferWriter struct {
	w      http.ResponseWriter
	Output bytes.Buffer
	Status int
	Size   int
}

func NewSnifferWriter(w http.ResponseWriter) *SnifferWriter {
	return &SnifferWriter{w: w}
}

func (l *SnifferWriter) Header() http.Header {
	return l.w.Header()
}

func (l *SnifferWriter) WriteHeader(i int) {
	l.Status = i
	l.w.WriteHeader(i)
}

func (l *SnifferWriter) Write(b []byte) (int, error) {
	n, err := l.w.Write(b)
	l.Output.Write(b[:n])
	l.Size += n
	return n, err
}
