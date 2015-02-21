package middle

import (
	"net/http"
	"regexp"
)

// Limiter is a meta-middleware that limits when middleware
// gets applied based on the URL path requested.
// This is useful for filtering static content paths from logging middleware,
// requiring login on every path but the login page, etc.
type limiter struct {
	patterns []*regexp.Regexp
	black    bool
}

// WrapWare implements the Meta interface for Limiter.
func (l limiter) WrapWare(m Ware) Ware {
	return WareFunc(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			matched, _ := matchAny(l.patterns, r.URL.Path)
			if (!l.black && matched) || (l.black && !matched) {
				m.WrapHandler(h).ServeHTTP(w, r)
			} else {
				h.ServeHTTP(w, r)
			}
		})
	})
}

// convert a slice of strings to a slice of regexp matchers.
// panics if any fail to compile.
func makeMatchers(patterns []string) []*regexp.Regexp {
	matchers := make([]*regexp.Regexp, len(patterns))
	for i, v := range patterns {
		matchers[i] = regexp.MustCompile(v)
	}
	return matchers
}

// tests s against a slice of regexps. Returns true on first match,
// false otherwise
func matchAny(ms []*regexp.Regexp, s string) (bool, int) {
	for n, v := range ms {
		if v.MatchString(s) {
			return true, n
		}
	}
	return false, -1
}

// Exclude returns a Limiter that will only run middleware on paths that don't
// match any of the listed regular expression patterns.
// Panics if one of the expressions fails to compile.
func Exclude(patterns ...string) MetaFunc {
	matchers := makeMatchers(patterns)
	return limiter{patterns: matchers, black: true}.WrapWare
}

// Allow returns a Limiter that will only run middleware on paths that
// match one of the listed regular expression patterns.
// Panics if one of the expressions fails to compile.
func Allow(patterns ...string) MetaFunc {
	matchers := makeMatchers(patterns)
	return limiter{patterns: matchers, black: false}.WrapWare
}
