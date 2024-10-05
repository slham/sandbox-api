package handler

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

type getWorkoutsQuery struct {
	APIQuery
}

type getWorkoutsRequest struct {
	userID string
	query  getWorkoutsQuery
}

func getWorkoutsQueryParams(q url.Values) (getWorkoutsQuery, error) {
	gwq := getWorkoutsQuery{}
	apiQuery, err := getStandardQueryParams(q)
	if err != nil {
		return gwq, fmt.Errorf("failed to gather query params. %w", err)
	}
	gwq.APIQuery = apiQuery
	return gwq, nil
}

func handleGetWorkoutsError(w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrBadRequest) {
		slog.Warn("error getting workouts", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	slog.Error("error getting workouts", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
	return
}

func (c *WorkoutController) GetWorkouts(w http.ResponseWriter, r *http.Request) {
	slog.Debug("get workouts request")
	ctx := r.Context()
	vars := mux.Vars(r)
	userID := vars["user_id"]
	req := getWorkoutsRequest{userID: userID}
	query := r.URL.Query()
	q, err := getWorkoutsQueryParams(query)
	if err != nil {
		handleGetWorkoutsError(w, err)
		return
	}

	req.query = q

	workouts, err := c.getWorkouts(ctx, req)
	if err != nil {
		handleGetWorkoutsError(w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, workouts)
	return
}

func (c *WorkoutController) getWorkouts(ctx context.Context, req getWorkoutsRequest) ([]model.Workout, error) {
	workouts := []model.Workout{}
	q := dao.WorkoutQuery{
		UserID: req.userID,
		Query: dao.Query{
			SortCol: req.query.APIQuery.SortCol,
			Sort:    req.query.APIQuery.Sort,
			Limit:   req.query.APIQuery.Limit,
			Offset:  req.query.APIQuery.Offset,
		},
	}
	workouts, err := dao.GetWorkouts(ctx, q)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return workouts, nil
		}
		return workouts, fmt.Errorf("failed to get workouts. %w", err)
	}
	return workouts, nil
}
