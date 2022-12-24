package main

import "net/http"

type Router struct {
	tree *node
}

type route struct {
	methods []string
	path    string
	handler http.Handler
}

func NewRouter() *Router {
	return &Router{
		tree: newNode(),
	}
}
