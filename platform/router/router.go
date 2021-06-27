package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Handler func(w http.ResponseWriter, r *http.Request) error

type Router struct {
	mux *chi.Mux
}

func NewRouter() *Router {
	return &Router{mux: chi.NewRouter()}
}

func (r *Router) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	r.mux.ServeHTTP(writer, request)
}

func (r *Router) Handle(method, path string, customHandler Handler) {

	handler := func(w http.ResponseWriter, r *http.Request) {
		if err := customHandler(w, r); err != nil {

		}
	}

	r.mux.MethodFunc(method, path, handler)
}
