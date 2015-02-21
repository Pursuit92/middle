package middle

import (
	"compress/gzip"
	"net/http"
	"strings"
)

// Gzip compresses the response body using gzip. It only does this if gzip
// is included in the Accept-Encoding header of the request
var Gzip = WareFunc(gzipFunc)

func gzipFunc(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		encoding := r.Header.Get("Accept-Encoding")

		if strings.Contains(encoding, "gzip") {
			gz := gzip.NewWriter(w)
			defer gz.Close()
			resp := NewRespSniffer(w, gz)
			resp.Header().Set("Content-Encoding", "gzip")
			h.ServeHTTP(resp, r)
		} else {
			h.ServeHTTP(w, r)
		}
	})
}
