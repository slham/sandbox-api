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

type deleteWorkoutRequest struct {
	UserID    string
	WorkoutID string
}

func handleDeleteWorkoutError(w http.ResponseWriter, err error) {
	slog.Error("error deleting workout", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
	return
}

func (c *WorkoutController) DeleteWorkout(w http.ResponseWriter, r *http.Request) {
	slog.Debug("delete workout request")
	ctx := r.Context()
	vars := mux.Vars(r)
	userID := vars["user_id"]
	workoutID := vars["workout_id"]
	req := deleteWorkoutRequest{
		UserID:    userID,
		WorkoutID: workoutID,
	}

	err := c.deleteWorkout(ctx, req)
	if err != nil {
		handleDeleteWorkoutError(w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusNoContent, nil)
	return
}

func (c *WorkoutController) deleteWorkout(ctx context.Context, req deleteWorkoutRequest) error {
	_, err := c.getWorkoutByID(ctx, getWorkoutRequest{UserID: req.UserID, WorkoutID: req.WorkoutID})
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return NewApiError(404, ApiErrNotFound)
		}
		return fmt.Errorf("failed to delete workout. %w", err)
	}

	err = dao.DeleteWorkout(ctx, req.UserID, req.WorkoutID)
	if err != nil {
		return fmt.Errorf("failed to delete workout. %w", err)
	}

	return nil
}
