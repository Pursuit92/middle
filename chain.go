package httputil

import "net/http"

type Middleware func(http.Handler) http.Handler

func ChainMiddleware(handler http.Handler, warez ...Middleware) http.Handler {
	var ret http.Handler
	for _, v := range warez {
		ret = v(ret)
	}

	return ret
}
