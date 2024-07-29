package dao

import (
	"context"
	"fmt"

	"github.com/slham/sandbox-api/model"
)

func InsertWorkout(ctx context.Context, workout model.workout) (model.Workout, error) {
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
		return workout, fmt.Errorf("failed to insert workout. %w", err)
	}
	return workout, nil
}
