package middle

import (
	"fmt"
	"log"
	"net/http"
)

// NewLogger creates and returns a new Logger middleware.
// If l is nil, the returned Logger will write to the default logger
func NewLogger(l *log.Logger) WareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Drop the port from the request, don't really care about it
			noPort, _ := splitAddr(r.RemoteAddr)
			reqLog := fmt.Sprintf("%s: %s %s %s", noPort, r.Method, r.URL.String(), r.Proto)

			// Need to keep track of the status, so use a sniffer.
			// The actual response data doesn't matter, so just pass itself
			// as the Writer
			resp := NewRespSniffer(w, w)
			h.ServeHTTP(resp, r)
			if l == nil {
				log.Printf("%s %d", reqLog, resp.Status)
			} else {
				l.Printf("%s %d", reqLog, resp.Status)
			}
		})
	}
}
