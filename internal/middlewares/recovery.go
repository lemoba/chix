package middlewares

import (
	"net/http"

	"github.com/go-chi/render"

	"github.com/lemoba/chix/pkg/errors"
)

func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				switch e := rec.(type) {
				case *errors.Error:
					render.Status(r, http.StatusOK)
					render.JSON(w, r, e)
				default:
					render.Status(r, http.StatusInternalServerError)
					render.JSON(w, r, errors.ErrInternalServerError)
				}
			}
		}()
		next.ServeHTTP(w, r)
	})
}
