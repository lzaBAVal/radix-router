package main

import "net/http"

var (
	tmpRoute = &route{}
)

type route struct {
	methods []string
	path    string
	handler http.Handler
}

type Router struct {
	tree *RadixTree
}

func NewRouter() *Router {
	return &Router{
		tree: newRadixTree(),
	}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	method := req.Method
	path := req.URL.Path
	handler, err := r.tree.Search(path, method)
	if err != nil {
		w.WriteHeader(handleError(err))
		return
	}
	handler.ServeHTTP(w, req)

}

func handleError(err error) int {
	var status int
	switch err {
	case ErrMethodNotAllowed:
		status = http.StatusMethodNotAllowed
	case ErrNotFound:
		status = http.StatusNotFound
	}
	return status
}

func (r *Router) Methods(methods ...string) *Router {
	tmpRoute.methods = append(tmpRoute.methods, methods...)
	return r
}

// Handler sets a handler.
func (r *Router) Handler(path string, handler http.Handler) {
	tmpRoute.handler = handler
	tmpRoute.path = path
	r.Handle()
}

// Handle handles a route.
func (r *Router) Handle() {
	r.tree.Insert(tmpRoute.methods, tmpRoute.path, tmpRoute.handler)
	tmpRoute = &route{}
}
