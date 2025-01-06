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

type getWorkoutRequest struct {
	UserID    string
	WorkoutID string
}

func handleGetWorkoutError(ctx context.Context, w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrNotFound) {
		slog.WarnContext(ctx, "error getting workout by id", "err", err)
		request.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	slog.ErrorContext(ctx, "error getting workout by id", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
}

func (c *WorkoutController) GetWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "get workout by id request")
	vars := mux.Vars(r)
	userID := vars["user_id"]
	workoutID := vars["workout_id"]

	req := getWorkoutRequest{UserID: userID, WorkoutID: workoutID}
	workout, err := c.getWorkoutByID(ctx, req)
	if err != nil {
		handleGetWorkoutError(ctx, w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, workout)
	return
}

func (c *WorkoutController) getWorkoutByID(ctx context.Context, req getWorkoutRequest) (model.Workout, error) {
	workout, err := dao.GetWorkoutByID(ctx, req.UserID, req.WorkoutID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return workout, NewApiError(404, ApiErrNotFound)
		}
		return workout, fmt.Errorf("failed to get workout by id. %w", err)
	}
	return workout, nil
}
