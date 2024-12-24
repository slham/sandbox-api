package middlewares

import (
	"net/http"

	"github.com/slham/sandbox-api/auth"
	"github.com/slham/sandbox-api/request"
)

func Establish(store *auth.StandardSessionStore) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			store.EstablishSession(w, r)
			if request.GetStop(r.Context()) {
				return
			}
			f(w, r)
		}
	}
}

func Verify(store *auth.StandardSessionStore) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			store.VerifySession(w, r)
			if request.GetStop(r.Context()) {
				return
			}
			f(w, r)
		}
	}
}

func Terminate(store *auth.StandardSessionStore) Middleware {
	return func(f http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			store.TerminateSession(w, r)
			if request.GetStop(r.Context()) {
				return
			}
			f(w, r)
		}
	}
}
