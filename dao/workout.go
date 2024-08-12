package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/slham/sandbox-api/model"
)

var (
	ErrConflictWorkoutName = errors.New("workout name already exists")
)

func InsertWorkout(ctx context.Context, workout model.Workout) (model.Workout, error) {
	_, err := getDB().ExecContext(ctx,
		`INSERT INTO sandbox.workout(
			id,
			name,
			user_id,
			created,
			updated,
			exercises
		)
		VALUES(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)`,
		workout.ID,
		workout.Name,
		workout.UserID,
		workout.Created,
		workout.Updated,
		workout.Exercises)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "u_user_name") {
					return workout, ErrConflictWorkoutName
				}
				return workout, fmt.Errorf("failed to insert workout. conflict. %w", err)
			}
		}
		return workout, fmt.Errorf("failed to insert workout. %w", err)
	}
	return workout, nil
}
