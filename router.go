package chix

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type HandlerFunc func(c *Context)

func Wrap(h HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		c := &Context{W: w, R: r}
		h(c)
	}
}

func NewRouter() *chi.Mux {
	return chi.NewRouter()
}
