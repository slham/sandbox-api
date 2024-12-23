package auth

import (
	"log/slog"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/slham/sandbox-api/request"
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
	ctx := r.Context()
	session, err := store.cookieStore.Get(r, cookieName)
	if err != nil {
		slog.ErrorContext(ctx, "failed to establish session")
		request.RespondWithError(w, http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	session.Values["authenticated"] = true
	session.Save(r, w)
}

func (store *StandardSessionStore) VerifySession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, err := store.cookieStore.Get(r, cookieName)
	if err != nil {
		slog.ErrorContext(ctx, "failed to verify session")
		request.RespondWithError(w, http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		slog.WarnContext(ctx, "INTRUDER!")
		request.RespondWithError(w, http.StatusForbidden, "Forbidden")
		return
	}

	slog.InfoContext(ctx, "The cake is a lie!")
}

func (store *StandardSessionStore) TerminateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, err := store.cookieStore.Get(r, cookieName)
	if err != nil {
		slog.ErrorContext(ctx, "failed to terminate session")
		request.RespondWithError(w, http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	session.Values["authenticated"] = false
	session.Save(r, w)
}

func (store *StandardSessionStore) HydrateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	rc := request.GetRequestContext(ctx)
	if rc == nil {
		slog.ErrorContext(ctx, "cannot track user")
		request.RespondWithError(w, http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	session, err := store.cookieStore.Get(r, cookieName)
	if err != nil {
		slog.ErrorContext(ctx, "failed to hydrate session")
		request.RespondWithError(w, http.StatusUnauthorized, "Invalid Credentials")
		return
	}

	session.Values["authenticated"] = true
	session.Values["user_id"] = rc.UserID
	session.Values["roles"] = rc.Roles
	session.Save(r, w)
}
