package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/samber/lo"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

type updateWorkoutRequest struct {
	UserID    string
	WorkoutID string
	Name      string          `json:"name"`
	Exercises model.Exercises `json:"exercises"`
}

func handleUpdateWorkoutError(w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrBadRequest) {
		slog.Warn("error creating workout", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}

	if errors.Is(err, ApiErrConflict) {
		slog.Warn("error creating workout", "err", err)
		request.RespondWithError(w, http.StatusConflict, err.Error())
		return
	}

	slog.Error("error creating workout", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
	return
}

func (c *WorkoutController) UpdateWorkout(w http.ResponseWriter, r *http.Request) {
	slog.Debug("update workout request")
	ctx := r.Context()
	req := updateWorkoutRequest{}
	vars := mux.Vars(r)
	userID := vars["user_id"]
	workoutID := vars["workout_id"]
	req.UserID = userID
	req.WorkoutID = workoutID

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.Warn("error decoding create workout request", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	if req.UserID != userID {
		slog.Warn("user not allowed to modify this workout")
		request.RespondWithError(w, http.StatusUnauthorized, "UNAUTHORIZED")
		return
	}

	workout, err := c.updateWorkout(ctx, req)
	if err != nil {
		handleUpdateWorkoutError(w, err)
		return
	}

	request.RespondWithJSON(w, http.StatusOK, workout)
	return
}

func (c *WorkoutController) updateWorkout(ctx context.Context, req updateWorkoutRequest) (model.Workout, error) {
	workout := model.Workout{}
	if _, err := dao.GetUserByID(ctx, req.UserID); err != nil {
		slog.Warn("failed to find user", "err", err)
		return workout, NewApiError(404, ApiErrNotFound).Append("user does not exist")
	}

	workout, err := dao.GetWorkoutByID(ctx, req.UserID, req.WorkoutID)
	if err != nil {
		slog.Warn("failed to find workout", "err", err)
		return workout, NewApiError(404, ApiErrNotFound).Append("workout does not exist")
	}

	if err := validateUpdateWorkoutRequest(ctx, req); err != nil {
		return workout, fmt.Errorf("failed to validate update workout request. %w", err)
	}

	workout.Name = req.Name
	workout.Exercises = req.Exercises

	if err := dao.UpdateWorkout(ctx, workout); err != nil {
		if errors.Is(err, dao.ErrConflictWorkoutName) {
			return workout, NewApiError(409, ApiErrConflict).Append("workout name already exists")
		}
		return workout, fmt.Errorf("failed to update workout. %w", err)
	}

	return workout, nil
}

func validateUpdateWorkoutRequest(ctx context.Context, req updateWorkoutRequest) error {
	apiErr := NewApiError(400, ApiErrBadRequest)

	if req.Name == "" {
		apiErr = apiErr.Append("workout must have a name")
	}

	for _, exercise := range req.Exercises {
		if exercise.Name == "" {
			apiErr = apiErr.Append("exercise must have a name")
		}

		for _, muscle := range exercise.Muscles {
			if muscle.Name == "" {
				apiErr = apiErr.Append("muscle must have a name")
			}
			if !lo.Contains(model.MuscleGroups, model.MuscleGroup(muscle.MuscleGroup)) {
				apiErr = apiErr.Append(fmt.Sprintf("invalid muscle group. valid options: %v", model.MuscleGroups))
			}
		}
	}

	if apiErr.HasError() {
		return apiErr
	}

	return nil
}
