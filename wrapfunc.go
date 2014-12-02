package httputil

import "net/http"

func WrapFunc(orig func(http.Handler) http.Handler) func(func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
		return orig(http.HandlerFunc(f)).ServeHTTP
	}
}
