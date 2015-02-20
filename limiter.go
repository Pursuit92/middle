package httputil

import (
	"net/http"
	"regexp"
)

func makeMatchers(patterns []string) []*regexp.Regexp {
	matchers := make([]*regexp.Regexp, len(patterns))
	for i, v := range patterns {
		matchers[i] = regexp.MustCompile(v)
	}
	return matchers
}

func matchAny(ms []*regexp.Regexp, s string) (bool, int) {
	for n, v := range ms {
		if v.MatchString(s) {
			return true, n
		}
	}
	return false, -1
}

// Run middleware only on non-matching paths
func ExcludeHandler(handler func(http.Handler) http.Handler, patterns ...string) func(http.Handler) http.Handler {
	matchers := makeMatchers(patterns)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matched, _ := matchAny(matchers, r.URL.Path)
			if !matched {
				handler(h).ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}

		})
	}
}

// Run middleware only on matching paths
func AllowHandler(handler func(http.Handler) http.Handler, patterns ...string) func(http.Handler) http.Handler {
	matchers := makeMatchers(patterns)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matched, _ := matchAny(matchers, r.URL.Path)
			if matched {
				handler(h).ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}
		})
	}
}
