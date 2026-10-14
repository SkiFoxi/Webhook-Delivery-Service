package middleware

import (
	"net/http"

)

type Middleware func(http.Handler) http.Handler

func Chain (h http.Handler, middlewareы ...Middleware) http.Handler {
	for i := len(middlewareы) - 1; i >= 0; i-- {
		h = middlewareы[i](h)
	}
	return h
}