package middle

import "net/http"

// ChainMiddleware combines a list of middlewares and returns a
// single middleware. The middlewares are combined with the first
// in the list as the outermost middleware.
// So ChainMiddleware(m1,m2,m3).WrapHandler(h) results in m1(m2(m3(h)))
func ChainMiddleware(warez ...Ware) WareFunc {
	return WareFunc(func(h http.Handler) http.Handler {
		for i := len(warez) - 1; i >= 0; i-- {
			h = warez[i].WrapHandler(h)
		}
		return h
	})
}

// ChainMetaMiddleware combines meta-middleware in the same manner
// as chained middleware.
func ChainMetaMiddleware(warez ...Meta) MetaFunc {
	return MetaFunc(func(m Ware) Ware {
		for i := len(warez) - 1; i >= 0; i-- {
			m = warez[i].WrapWare(m)
		}
		return m
	})
}
