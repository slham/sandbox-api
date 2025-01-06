package handler

import (
	"github.com/gorilla/sessions"
	"github.com/slham/sandbox-api/auth"
)

type AuthController struct {
	cookieStore *sessions.CookieStore
}

func NewAuthController(store *auth.StandardSessionStore) AuthController {
	return AuthController{
		cookieStore: store.GetCookieStore(),
	}
}

func (c *AuthController) GetCookieStore() *sessions.CookieStore {
	return c.cookieStore
}
