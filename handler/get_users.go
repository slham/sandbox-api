package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

type getUsersQuery struct {
	ID       string
	Username string
	Email    string
	APIQuery
}

type getUsersRequest struct {
	query getUsersQuery
}

func getUserQueryParams(ctx context.Context, q url.Values) (getUsersQuery, error) {
	guq := getUsersQuery{}
	if qID := q.Get("id"); qID != "" {
		guq.ID = qID
	}
	if qUsername := q.Get("username"); qUsername != "" {
		guq.Username = qUsername
	}
	if qEmail := q.Get("email"); qEmail != "" {
		guq.Email = qEmail
	}
	apiQuery, err := getStandardQueryParams(ctx, q)
	if err != nil {
		return guq, fmt.Errorf("failed to gather query params. %w", err)
	}
	guq.APIQuery = apiQuery
	return guq, nil
}

func handleGetUsersError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrBadRequest) {
		slog.WarnContext(ctx, "error getting users", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	slog.ErrorContext(ctx, "error getting users", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
	return
}

func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "get users request")
	query := r.URL.Query()

	req := getUsersRequest{}
	q, err := getUserQueryParams(ctx, query)
	if err != nil {
		handleGetUsersError(ctx, w, err)
		return
	}

	req.query = q

	users, err := c.getUsers(ctx, req)
	if err != nil {
		handleGetUsersError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, users)
	return
}

func (c *UserController) getUsers(ctx context.Context, req getUsersRequest) ([]model.User, error) {
	q := dao.UserQuery{
		ID:       req.query.ID,
		Username: req.query.Username,
		Email:    req.query.Email,
		Query: dao.Query{
			SortCol: req.query.APIQuery.SortCol,
			Sort:    req.query.APIQuery.Sort,
			Limit:   req.query.APIQuery.Limit,
			Offset:  req.query.APIQuery.Offset,
		},
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
