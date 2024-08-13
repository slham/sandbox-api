package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
	"github.com/slham/sandbox-api/valid"
)

type updateUserRequest struct {
	UserID   string
	Username string `json:"username"`
	Email    string `json:"email"`
}

func (c *UserController) UpdateUser(w http.ResponseWriter, r *http.Request) {
	slog.Debug("update user request")
	ctx := r.Context()
	req := updateUserRequest{}
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("error decoding update user request", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	req.UserID = userID

	user, err := c.updateUser(ctx, req)
	if err != nil {
		if errors.Is(err, ApiErrBadRequest) {
			slog.Warn("error updating user", "err", err)
			request.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, ApiErrConflict) {
			slog.Warn("error updating user", "err", err)
			request.RespondWithError(w, http.StatusConflict, err.Error())
			return
		}

		slog.Error("error updating user", "err", err)
		request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	request.RespondWithJSON(w, http.StatusOK, user)
	return
}

func (c *UserController) updateUser(ctx context.Context, req updateUserRequest) (model.User, error) {
	user, err := c.getUserByID(ctx, getUserRequest{ID: req.UserID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, NewApiError(404, ApiErrNotFound)
		}
		return user, fmt.Errorf("failed to update user. %w", err)
	}

	if err := validateUpdateUserRequest(ctx, req); err != nil {
		return user, fmt.Errorf("failed to validate update user request. %w", err)
	}

	if req.Username != "" {
		user.Username = req.Username
	}

	if req.Email != "" {
		user.Email = req.Email
	}

	user.Updated = time.Now()
	user.Password = ""

	err = dao.UpdateUser(ctx, user)
	if err != nil {
		if errors.Is(err, dao.ErrConflictUsername) {
			return user, NewApiError(409, ApiErrConflict).Append("username already exists")
		}
		if errors.Is(err, dao.ErrConflictEmail) {
			return user, NewApiError(409, ApiErrConflict).Append("email already exists")
		}
		return user, fmt.Errorf("failed to update user. %w", err)
	}

	return user, nil
}

func validateUpdateUserRequest(ctx context.Context, req updateUserRequest) error {
	apiErr := NewApiError(400, ApiErrBadRequest)

	if len(req.Username) < 4 {
		apiErr = apiErr.Append("username must be at leat four characters long")
	}

	if err := valid.IsEmail(req.Email); err != nil {
		apiErr = apiErr.Append("invalid email")
	}

	if apiErr.HasError() {
		return apiErr
	}

	return nil
}
