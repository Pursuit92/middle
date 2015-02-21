package middle

// WrapFunc turns a func(http.Handler) http.Handler middleware into its function equivalent.
// func WrapFunc(orig Ware) func(func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
// 	return func(f func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
// 		return orig(http.HandlerFunc(f)).ServeHTTP
// 	}
// }
