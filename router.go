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

// NewRouterWith wraps an existing chi.Router
func NewRouterWith(r chi.Router) *Router {
	return &Router{r}
}

// wrapHandler binds a Context-aware handler to current Chix instance
func (c *Chix) wrapHandler(h func(Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r, c)
		if cc, ok := any(ctx).(*context); ok {
			cc.handler = h
		}
		h(ctx)
	}
}
