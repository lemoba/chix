package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

type ctxKey string

const RequestIDKey ctxKey = "chix_request_id"

var RequestIDHeader = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reqID := r.Header.Get(RequestIDHeader)
		if reqID == "" {
			reqID = uuid.New().String()
		}

		ctx := r.Context()

		ctx = context.WithValue(ctx, RequestIDKey, reqID)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
