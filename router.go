package chix

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type Router struct {
	chi.Router
}

func NewRouter() *Router {
	return &Router{chi.NewRouter()}
}

func (r *Router) Connect(pattern string, handlerFunc func(*Context)) {
	r.Router.Connect(pattern, wrapHandler(handlerFunc))
}

func (r *Router) Delete(pattern string, handlerFunc func(*Context)) {
	r.Router.Delete(pattern, wrapHandler(handlerFunc))
}

func (r *Router) Get(pattern string, handlerFunc func(*Context)) {
	r.Router.Get(pattern, wrapHandler(handlerFunc))
}

func (r *Router) Head(pattern string, handlerFunc func(*Context)) {
	r.Router.Head(pattern, wrapHandler(handlerFunc))
}

func (r *Router) Options(pattern string, handlerFunc func(*Context)) {
	r.Router.Options(pattern, wrapHandler(handlerFunc))
}

func (r *Router) Patch(pattern string, handlerFunc func(*Context)) {
	r.Router.Patch(pattern, wrapHandler(handlerFunc))
}

func (r *Router) Post(pattern string, handlerFunc func(*Context)) {
	r.Router.Post(pattern, wrapHandler(handlerFunc))
}

func (r *Router) Put(pattern string, handlerFunc func(*Context)) {
	r.Router.Put(pattern, wrapHandler(handlerFunc))
}

func (r *Router) Trace(pattern string, handlerFunc func(*Context)) {
	r.Router.Trace(pattern, wrapHandler(handlerFunc))
}

func wrapHandler(h func(*Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		h(ctx)
	}
}
