package httputil

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func Logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqLog := fmt.Sprintf("%s: %s %s %s", strings.Split(r.RemoteAddr, ":")[0], r.Method, r.URL.String(), r.Proto)
		resp := NewSnifferWriter(w)
		h.ServeHTTP(resp, r)
		log.Printf("%s %d", reqLog, resp.Status)
	})
}
