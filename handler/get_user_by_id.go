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
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

type getUserRequest struct {
	ID string
}

func handleGetUserError(w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrNotFound) {
		slog.Warn("error getting user by id", "err", err)
		request.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	slog.Error("error getting user by id", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
}

func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	slog.Debug("get user by id request")
	ctx := r.Context()
	vars := mux.Vars(r)
	userID := vars["user_id"]

	req := getUserRequest{ID: userID}
	user, err := c.getUserByID(ctx, req)
	if err != nil {
		handleGetUserError(w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, user)
	return
}

func (c *UserController) getUserByID(ctx context.Context, req getUserRequest) (model.User, error) {
	user, err := dao.GetUserByID(ctx, req.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return user, NewApiError(404, ApiErrNotFound)
		}
		return user, fmt.Errorf("failed to get user by id. %w", err)
	}
	return user, nil
}
