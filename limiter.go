package httputil

import (
	"net/http"
	"regexp"
)

func ExcludeHandler(handler func(http.Handler) http.Handler, patterns ...string) func(http.Handler) http.Handler {
	matchers := make([]*regexp.Regexp, len(patterns))
	for i, v := range patterns {
		matchers[i] = regexp.MustCompile(v)
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matched := false
			for _, v := range matchers {
				if v.MatchString(r.URL.Path) {
					matched = true
				}
			}
			if !matched {
				handler(h).ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}

		})
	}

}

func AllowHandler(handler func(http.Handler) http.Handler, patterns ...string) func(http.Handler) http.Handler {
	matchers := make([]*regexp.Regexp, len(patterns))
	for i, v := range patterns {
		matchers[i] = regexp.MustCompile(v)
	}
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matched := false
			for _, v := range matchers {
				if v.MatchString(r.URL.Path) {
					matched = true
				}
			}
			if matched {
				handler(h).ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}

		})
	}

}
