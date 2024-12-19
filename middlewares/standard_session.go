package middlewares

import (
	"net/http"

	"github.com/slham/sandbox-api/auth"
)

func Establish(store *auth.StandardSessionStore) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			store.EstablishSession(w, r)
			f(w, r)
		}
	}
}

func Verify(store *auth.StandardSessionStore) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			store.VerifySession(w, r)
			f(w, r)
		}
	}
}

func Terminate(store *auth.StandardSessionStore) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			store.TerminateSession(w, r)
			f(w, r)
		}
	}
}
