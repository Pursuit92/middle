## Simple Go Middleware

This is a simple implementation of middleware for Go's ```http.Handler```.
Its API mimics that of the ```http.Handler``` system and allows middlewares
to be implemented as either functions or types.

The basic middleware interface (```Ware```) has a single
```WrapHandler``` method that takes an http.Handler and returns a new
one.

```go
type Ware interface {
  WrapHandler(http.Handler) http.Handler
}
```

Much like http.Handler, it also provides a type that can be used to
turn middleware functions into usable middleware interface types.

```go
func identity(h http.Handler) http.Handler {
  return h
}

var identityWare Ware = WareFunc(identity)
```

Middlewares implemented in this manner can easily be chained using the
```ChainMiddleware``` function. This expects middlewares to be ordered
innermost-first.

```go
var m1, m2, m3 Ware

combined := ChainMiddleware(m1, m2, m3)
// combined.WrapHandler(h) is the equivalent of
// m3.WrapHandler(m2.WrapHandler(m1.WrapHandler(h)))
```

### Meta-middleware

Meta-middleware are essentially middleware for middleware. Currently,
only one is implemented.

The ```Limiter``` meta-middleware is used to filter the paths that a
middleware gets applied to.

```go
var mw Ware
var h http.Handler
// Create a Limiter that won't run the middleware on paths starting
// with "/assets"
filterAssets := Exclude(`^/assets`)

newHandler := filterAssets.WrapWare(mw).WrapHandler(h)
```

### Provided Middleware

Currently, only a simple logging middleware is provided. It will log
the address, path, protocol, and response status of a request.

```go
var h http.Handler
// nil makes the middleware use the standard logger.
l := NewLogger(nil)

l.WrapHandler(h)
```

### Godoc

View the full godocs at [godoc.org](https://godoc.org/github.com/Pursuit92/middle)
