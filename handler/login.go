package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/slham/sandbox-api/crypt"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func handleLoginError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrBadRequest) {
		slog.WarnContext(ctx, "error login", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if errors.Is(err, ApiErrForbidden) {
		slog.ErrorContext(ctx, "unauthenticated login attempt", "err", err)
		request.RespondWithError(w, http.StatusForbidden, err.Error())
		return
	} else if errors.Is(err, ApiErrNotFound) {
		slog.ErrorContext(ctx, "user not found", "err", err)
		request.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}

	slog.ErrorContext(ctx, "error login", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
	return
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "login request")
	loginRequest := LoginRequest{}

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		slog.WarnContext(ctx, "error decoding login request", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	user, err := handleLogin(ctx, loginRequest)
	if err != nil {
		handleLoginError(ctx, w, err)
		return
	}

	c.hydrateSession(w, r, user)

	request.RespondWithJSON(w, http.StatusOK, user)
}

func handleLogin(ctx context.Context, req LoginRequest) (model.User, error) {
	if err := validateLoginRequest(ctx, req); err != nil {
		return model.User{}, fmt.Errorf("failed to validate login request. %w", err)
	}

	user, err := dao.GetUserByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, dao.ErrUserNotFound) {
			return user, NewApiError(404, ApiErrNotFound)
		}
		return user, fmt.Errorf("failed to get user. %w", err)
	}

	plainText, err := crypt.Decrypt(user.Password)
	if err != nil {
		return user, fmt.Errorf("failed to check password. %w", err)
	}

	if plainText != req.Password {
		slog.WarnContext(ctx, "failed login attempt for user", "user_id", user.ID)
		return user, NewApiError(403, ApiErrForbidden)
	}

	user.Password = ""
	return user, nil
}

func validateLoginRequest(ctx context.Context, req LoginRequest) error {
	apiErr := NewApiError(400, ApiErrBadRequest)

	if req.Username == "" {
		apiErr = apiErr.Append("username must be present")
	}

	if req.Password == "" {
		apiErr = apiErr.Append("password must be present")
	}

	if apiErr.HasError() {
		return apiErr
	}

	return nil
}

func (c *AuthController) hydrateSession(w http.ResponseWriter, r *http.Request, user model.User) error {
	ctx := r.Context()
	session, err := c.cookieStore.Get(r, "sandbox-cookie")
	if err != nil {
		slog.ErrorContext(ctx, "failed to establish session")
		return fmt.Errorf("failes to establish session. %w", err)
	}

	rc := request.GetRequestContext(ctx)
	roles := make([]string, len(user.Roles))
	for i := range user.Roles {
		role := user.Roles[i]
		roles[i] = role.Name
	}

	rc.UserID = user.ID
	rc.Roles = roles
	ctx = request.WithRequestContext(ctx, rc)
	r = r.WithContext(ctx)

	slog.DebugContext(ctx, "HYDRATING SESSION", "user_id", user.ID, "roles", roles)
	session.Values["authenticated"] = true
	session.Values["user_id"] = user.ID
	session.Values["roles"] = roles
	session.Save(r, w)

	return nil
}
