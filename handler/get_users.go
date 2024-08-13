package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"

	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	slog.Debug("get users request")
	ctx := r.Context()
	query := r.URL.Query()

	users, err := c.getUsers(ctx, query)
	if err != nil {
		if errors.Is(err, ApiErrBadRequest) {
			slog.Warn("error getting users", "err", err)
			request.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		slog.Error("error getting users", "err", err)
		request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	request.RespondWithJSON(w, http.StatusOK, users)
	return
}

func (c *UserController) getUsers(ctx context.Context, query url.Values) ([]model.User, error) {
	users := []model.User{}
	q := dao.UserQuery{}

	if qID := query.Get("id"); qID != "" {
		q.ID = qID
	}
	if qUsername := query.Get("username"); qUsername != "" {
		q.Username = qUsername
	}
	if qEmail := query.Get("email"); qEmail != "" {
		q.Email = qEmail
	}
	if qSortCol := query.Get("sort_column"); qSortCol != "" {
		q.SortCol = qSortCol
	}
	if qSort := query.Get("sort"); qSort != "" {
		q.Sort = qSort
	}
	if qLimit := query.Get("limit"); qLimit != "" {
		limit, err := strconv.Atoi(qLimit)
		if err != nil {
			slog.Warn("invalid limit", "limit", qLimit)
			return users, NewApiError(400, ApiErrBadRequest).Append("invalid limit")
		}
		q.Limit = limit
	}
	if qOffset := query.Get("offset"); qOffset != "" {
		offset, err := strconv.Atoi(qOffset)
		if err != nil {
			slog.Warn("invalid offset", "offset", qOffset)
			return users, NewApiError(400, ApiErrBadRequest).Append("invalid offset")
		}
		q.Offset = offset
	}

	users, err := dao.GetUsers(ctx, q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return users, nil
		}
		return users, fmt.Errorf("failed to get users. %w", err)
	}

	return users, nil
}
