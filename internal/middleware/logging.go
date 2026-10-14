package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type statusRecoder struct{
	http.ResponseWriter
	status int
}

func (r *statusRecoder) WriteHeader(code int) {
	r.status = code
	r.ResponseWriter.WriteHeader(code)
}

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		rec := &statusRecoder{ResponseWriter: w, status: 200}

		next.ServeHTTP(rec, r) /*WriteHeader зашит в rec и если вызывается WriteHeader, 
								то вместо стандартного используется наш кастомный с сохранением статус-кода*/
		slog.Info("request",
			"method", r.Method,
			"path", r.URL.Path,
			"status", rec.status,
			"duration", time.Since(start), 
		)
	})
}