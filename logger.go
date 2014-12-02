package httputil

import (
	"fmt"
	"log"
	"net/http"
)

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLog := fmt.Sprintf("%s: %s %s %s", r.RemoteAddr, r.Method, r.URL.String(), r.Proto)
		resp := NewSnifferWriter(w)
		h.ServeHTTP(resp, r)
		if resp.Status == 0 {
			resp.Status = 200
		}
		log.Printf("%s %d", reqLog, resp.Status)
	})
}
