package middle

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

// Logger is an (actually useful) example of a simple middleware. It writes
// Apache-style logs to a log.Logger of your choice.
type Logger struct {
	l *log.Logger
}

// WrapHandler impelments the Ware interface for Logger
func (l *Logger) WrapHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Drop the port from the request, don't really care about it
		splitAddr := strings.Split(r.RemoteAddr, ":")
		noPort := strings.Join(splitAddr[:len(splitAddr)-1], ":")
		reqLog := fmt.Sprintf("%s: %s %s %s", noPort, r.Method, r.URL.String(), r.Proto)

		// Need to keep track of the status, so use a sniffer.
		// The actual response data doesn't matter, so just pass itself
		// as the Writer
		resp := NewRespSniffer(w, w)
		h.ServeHTTP(resp, r)
		if l.l == nil {
			log.Printf("%s %d", reqLog, resp.Status)
		} else {
			l.l.Printf("%s %d", reqLog, resp.Status)
		}
	})
}

// NewLogger creates and returns a new Logger middleware.
// If l is nil, the returned Logger will write to the default logger
func NewLogger(l *log.Logger) *Logger {
	return &Logger{l: l}
}
