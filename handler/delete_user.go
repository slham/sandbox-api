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

func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	slog.Debug("delete user request")
	ctx := r.Context()
	vars := mux.Vars(r)
	userID := vars["user_id"]

	err := c.deleteUser(ctx, userID)
	if err != nil {
		slog.Error("error deleting user", "err", err)
		request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	request.RespondWithJSON(w, http.StatusNoContent, nil)
	return
}

func (c *UserController) deleteUser(ctx context.Context, userID string) error {
	_, err := c.getUserByID(ctx, getUserRequest{ID: userID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewApiError(404, ApiErrNotFound)
		}
		return fmt.Errorf("failed to delete user. %w", err)
	}

	err = dao.DeleteUser(ctx, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user. %w", err)
	}

	return nil
}
