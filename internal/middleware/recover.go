package middleware

import (
	"log/slog"
	"net/http"
)

func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func()  {
			if p := recover(); p != nil {
				slog.Error("panic recover",
					"err", p,
					"path", r.URL.Path,
					"method", r.Method,
				)
				http.Error(w, "internal error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
} 