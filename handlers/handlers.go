package handlers

import (
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/newrelic/go-agent/v3/newrelic"
)

func NewRelic(app *newrelic.Application) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		_, handler := newrelic.WrapHandle(app, "", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			txn := newrelic.FromContext(r.Context())
			txn.SetName(r.Method + " " + r.URL.Path)

			next.ServeHTTP(w, r)
		}))
		return handler
	}
}

func ChiNewRelic() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log.Println("chi handler")

			next.ServeHTTP(w, r)

			txn := newrelic.FromContext(r.Context())
			if txn == nil {
				log.Println("nil newrelic.FromContext")
				return
			}

			ctx := chi.RouteContext(r.Context())
			if ctx == nil {
				log.Println("nil chi.RouteContext")
				return
			}

			// update the pattern
			pattern := ctx.RoutePattern()
			txn.SetName(r.Method + " " + pattern)
			log.Println("SetName:", r.Method+" "+pattern)
		})
	}
}
