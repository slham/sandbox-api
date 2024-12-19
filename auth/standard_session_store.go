package auth

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
)

var (
	cookieName = "sandbox-cookie"
)

type StandardSessionStore struct {
	cookieStore *sessions.CookieStore
}

func NewStandardSessionStore() *StandardSessionStore {
	key := []byte(os.Getenv("SANDBOX_STANDARD_SESSION_KEY"))
	store := sessions.NewCookieStore(key)

	return &StandardSessionStore{
		cookieStore: store,
	}
}

func (store *StandardSessionStore) EstablishSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.cookieStore.Get(r, cookieName)
	if err != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	session.Values["authenticated"] = true
	session.Save(r, w)
}

func (store *StandardSessionStore) VerifySession(w http.ResponseWriter, r *http.Request) {
	session, err := store.cookieStore.Get(r, cookieName)
	if err != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	slog.Info("The cake is a lie!")
}

func (store *StandardSessionStore) TerminateSession(w http.ResponseWriter, r *http.Request) {
	session, err := store.cookieStore.Get(r, cookieName)
	if err != nil {
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	session.Values["authenticated"] = false
	session.Save(r, w)
}
