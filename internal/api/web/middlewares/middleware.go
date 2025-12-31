package middlewares

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

func RequestIDMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID := r.Header.Get("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}

		ctx := context.WithValue(r.Context(), "XRequestID", requestID)
		ctx = context.WithValue(ctx, "url", r.URL)
		ctx = context.WithValue(ctx, "method", r.Method)

		w.Header().Set("X-Request-ID", requestID)
		

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
