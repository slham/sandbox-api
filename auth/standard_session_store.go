package auth

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"slices"

	"github.com/gorilla/mux"
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

func (store *StandardSessionStore) GetCookieStore() *sessions.CookieStore {
	return store.cookieStore
}

func (store *StandardSessionStore) EstablishSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, err := store.cookieStore.Get(r, cookieName)
	if err != nil {
		slog.ErrorContext(ctx, "failed to establish session")
		r = stop(r, ctx)
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
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
		r = stop(r, ctx)
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	if auth, ok := session.Values["authenticated"].(bool); !ok || !auth {
		slog.ErrorContext(ctx, "INTRUDER!", "authenticated", false)
		r = stop(r, ctx)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	roles, ok := session.Values["user_id"].([]string)
	if !ok || roles == nil {
		slog.ErrorContext(ctx, "INTRUDER!", "session_roles", roles)
		r = stop(r, ctx)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	sessionUserID, ok := session.Values["user_id"].(string)
	if !ok || sessionUserID == "" {
		slog.ErrorContext(ctx, "INTRUDER!", "session_user_id", sessionUserID)
		r = stop(r, ctx)
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	userID := vars["user_id"]

	if userID != "" && userID != sessionUserID || !isAdmin(roles) {
		slog.ErrorContext(ctx, "INTRUDER!", "session_user_id", sessionUserID, "client_user_id", userID)
		r = stop(r, ctx)
		http.Error(w, "FUCK OFF!", http.StatusForbidden)
		return
	} else {
		slog.InfoContext(ctx, "SIR, YES, SIR!")
	}

	slog.InfoContext(ctx, "The cake is a lie!")
}

func (store *StandardSessionStore) TerminateSession(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	session, err := store.cookieStore.Get(r, cookieName)
	if err != nil {
		slog.ErrorContext(ctx, "failed to terminate session")
		r = stop(r, ctx)
		http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
		return
	}

	session.Values["authenticated"] = false
	session.Save(r, w)
}

func isAdmin(roles []string) bool {
	return slices.Contains(roles, "ADMIN")
}

func stop(r *http.Request, ctx context.Context) *http.Request {
	ctx = request.SetStop(ctx)
	r = r.WithContext(ctx)

	return r
}
