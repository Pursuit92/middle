package middle

import (
	"fmt"
	"net/http"
)

// RedirHTTPS provides a middleware to redirect the client to https on a
// specified port.
func RedirHTTPS(port int) WareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.TLS == nil {
				host, _ := splitAddr(r.Host)
				url := fmt.Sprintf("https://%s:%d%s", host, port, r.RequestURI)
				http.Redirect(w, r, url, http.StatusMovedPermanently)
				return
			}
			h.ServeHTTP(w, r)
		})
	}
}
