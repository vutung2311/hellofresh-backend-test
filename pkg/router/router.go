package router

import "net/http"

type Router struct {
	middlewareList []func(next http.Handler) http.Handler

	mux *http.ServeMux
}

func New(middlewareList []func(next http.Handler) http.Handler) *Router {
	return &Router{
		middlewareList: middlewareList,
		mux:            http.NewServeMux(),
	}
}

func (router *Router) wrapWithMiddlewareList(handler http.Handler) http.Handler {
	finalHandler := handler
	for _, middleware := range router.middlewareList {
		finalHandler = middleware(finalHandler)
	}
	return finalHandler
}

func (router *Router) AddRoute(pattern string, handler http.Handler) {
	router.mux.Handle(pattern, router.wrapWithMiddlewareList(handler))
}

func (router *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	router.mux.ServeHTTP(w, r)
}
