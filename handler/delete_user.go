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

func handleDeleteUserError(w http.ResponseWriter, err error) {
	slog.Error("error deleting user", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
	return
}

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	slog.Debug("delete user request")
	ctx := r.Context()
	vars := mux.Vars(r)
	userID := vars["user_id"]

	req := deleteUserRequest{UserID: userID}

	err := c.deleteUser(ctx, req)
	if err != nil {
		handleDeleteWorkoutError(w, err)
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
