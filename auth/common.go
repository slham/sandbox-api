package auth

import "net/http"

type SessionStore interface {
	EstablishSession(w http.ResponseWriter, r *http.Request)
	VerifySession(w http.ResponseWriter, r *http.Request)
	TerminateSession(w http.ResponseWriter, r *http.Request)
}
