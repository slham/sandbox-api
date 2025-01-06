package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/request"
)

type deleteUserRequest struct {
	UserID string
}

func handleDeleteUserError(ctx context.Context, w http.ResponseWriter, err error) {
	slog.ErrorContext(ctx, "error deleting user", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
	return
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "delete user request")
	vars := mux.Vars(r)
	userID := vars["user_id"]

	req := deleteUserRequest{UserID: userID}

	err := c.deleteUser(ctx, req)
	if err != nil {
		handleDeleteWorkoutError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusNoContent, nil)
	return
}

func (c *UserController) deleteUser(ctx context.Context, req deleteUserRequest) error {
	_, err := c.getUserByID(ctx, getUserRequest{ID: req.UserID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewApiError(404, ApiErrNotFound)
		}
		return fmt.Errorf("failed to delete user. %w", err)
	}

	err = dao.DeleteUser(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("failed to delete user. %w", err)
	}

	return nil
}
