package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/segmentio/ksuid"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

func (c *WorkoutController) CreateWorkout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	workout := model.Workout{}
	vars := mux.Vars(r)
	userID := vars["user_id"]

	if err := json.NewDecoder(r.Body).Decode(&workout); err != nil {
		slog.Warn("error decoding create workout request", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	workout.UserID = userID
	workout, err := c.createWorkout(ctx, workout)
	if err != nil {
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

	request.RespondWithJSON(w, http.StatusCreated, workout)
	return
}

func (c *WorkoutController) createWorkout(ctx context.Context, workout model.Workout) (model.Workout, error) {
	if _, err := dao.GetUserByID(ctx, workout.UserID); err != nil {
		return workout, NewApiError(404, ApiErrNotFound).Append("user does not exist")
	}

	if err := validateCreateWorkoutRequest(ctx, workout); err != nil {
		return workout, fmt.Errorf("failed to validate create workout request. %w", err)
	}

	workout.ID = newWorkoutID()
	workout.Created = time.Now()
	workout.Updated = time.Now()

	return workout, nil
}

func validateCreateWorkoutRequest(ctx context.Context, workout model.Workout) error {
	apiErr := NewApiError(400, ApiErrBadRequest)

	if workout.Name == "" {
		apiErr = apiErr.Append("workout must have a name")
	}

	for _, exercise := range workout.Exercises {
		if exercise.Name == "" {
			apiErr = apiErr.Append("exercise must have a name")
		}

		for _, muscle := range exercise.Muscles {
			if muscle.Name == "" {
				apiErr = apiErr.Append("muscle must have a name")
			}
		}
	}

	if apiErr.HasError() {
		return apiErr
	}

	return nil
}

func newWorkoutID() string {
	return fmt.Sprintf("work_%s", ksuid.New().String())
}
