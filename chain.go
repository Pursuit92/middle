package middle

import "net/http"

// ChainMiddleware combines a list of middlewares and returns a
// single middleware. The middlewares are combined with the first
// in the list as the innermost middleware.
// So ChainMiddleware(m1,m2,m3).WrapHandler(h) results in m3(m2(m1(h)))
// This may be changed in a future release
func ChainMiddleware(warez ...Ware) Ware {
	return WareFunc(func(h http.Handler) http.Handler {
		for _, v := range warez {
			h = v.WrapHandler(h)
		}
		return h
	})
}

// ChainMetaMiddleware combines meta-middleware in the same manner
// as chained middleware.
func ChainMetaMiddleware(warez ...Meta) Meta {
	return MetaFunc(func(m Ware) Ware {
		for _, v := range warez {
			m = v.WrapWare(m)
		}
		return m
	})
}
