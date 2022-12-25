package main

import (
	"fmt"
	"net/http"
)

func indexHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "GET /")
	})
}

func main() {
	myCustomRouter := NewRouter()

	myCustomRouter.Methods(http.MethodGet).Handler(`/`, indexHandler())

	http.ListenAndServe(":8000", myCustomRouter)
}
