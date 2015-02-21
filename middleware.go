package middle

import "net/http"

// A Ware (middleware) is used to wrap an http.Handler and return a new one.
// See Logger for a simple middleware example.
type Ware interface {
	WrapHandler(http.Handler) http.Handler
}

// WrapFunc applies a middleware to an http.HandlerFunc.
func WrapFunc(w Ware, f http.HandlerFunc) http.Handler {
	return w.WrapHandler(http.HandlerFunc(f))
}

// A WareFunc is a middleware defined using just a function. This
// wrapper type allows it to be used as a Ware interface type.
type WareFunc func(http.Handler) http.Handler

// WrapHandler implements the Ware interface for WareFunc.
// It simply applies f to h.
func (f WareFunc) WrapHandler(h http.Handler) http.Handler {
	return f(h)
}

// WareFuncFunc is for when you insist on defining your middleware as
// func(func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request)
// for whatever reason. It's got a funny name because you made poor choices.
type WareFuncFunc func(http.HandlerFunc) http.HandlerFunc

// WrapHandler implements Ware for WareFuncFunc. So now it's a WareFuncFunc Ware.
// Are you pleased with yourself?
func (f WareFuncFunc) WrapHandler(h http.Handler) http.Handler {
	return http.HandlerFunc(f(h.ServeHTTP))
}

// A Meta (meta-middleware) wraps a middleware the same way a middleware wraps
// an http.Handler. Currently only implemented by Limiter, but the interface is
// here nevertheless.
type Meta interface {
	WrapWare(Ware) Ware
}

// A MetaFunc is analagous to a WareFunc. It's used to convert a function taking
// and returning a Ware into a Meta.
type MetaFunc func(Ware) Ware

// WrapWare implements the Meta interface for MetaFunc.
func (f MetaFunc) WrapWare(m Ware) Ware {
	return f(m)
}
