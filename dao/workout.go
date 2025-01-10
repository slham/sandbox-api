package dao

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/slham/sandbox-api/model"
)

var ErrConflictWorkoutName = errors.New("workout name already exists")

func InsertWorkout(ctx context.Context, workout model.Workout) (model.Workout, error) {
	_, err := getDB().ExecContext(ctx,
		`INSERT INTO sandbox.workout(
			id,
			name,
			user_id,
			exercises
		)
		VALUES(
			$1,
			$2,
			$3,
			$4
		)`,
		workout.ID,
		workout.Name,
		workout.UserID,
		workout.Exercises,
	)
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

type WorkoutQuery struct {
	ID     string
	UserID string
	Query
}

func GetWorkoutByUserID(ctx context.Context, userID string) (model.Workout, error) {
	q := WorkoutQuery{UserID: userID}
	w, err := GetWorkout(ctx, q)
	if err != nil {
		return model.Workout{}, fmt.Errorf("failed to get workout by user id. %w", err)
	}
	return w, nil
}

func GetWorkoutByID(ctx context.Context, userID string, workoutID string) (model.Workout, error) {
	q := WorkoutQuery{ID: workoutID, UserID: userID}
	w, err := GetWorkout(ctx, q)
	if err != nil {
		return model.Workout{}, fmt.Errorf("failed to get workout by id. %w", err)
	}
	return w, nil
}

func GetWorkout(ctx context.Context, q WorkoutQuery) (model.Workout, error) {
	workouts, err := GetWorkouts(ctx, q)
	if err != nil {
		return model.Workout{}, fmt.Errorf("failed to get workouts. %w", err)
	}

	if len(workouts) != 1 {
		return model.Workout{}, ErrWorkoutNotFound
	}

	return workouts[0], nil
}

func GetWorkouts(ctx context.Context, q WorkoutQuery) ([]model.Workout, error) {
	stmt := `
		SELECT
			id,
			name,
			user_id,
			exercises,
			created,
			updated
		FROM
			sandbox.workout
		WHERE`

	if q.ID != "" {
		stmt = fmt.Sprintf("%s %s='%s'", stmt, "id", q.ID)
	}
	if q.UserID != "" {
		stmt = checkWhereClause(stmt)
		stmt = fmt.Sprintf("%s %s='%s'", stmt, "user_id", q.UserID)
	}

	stmt = addDefaultQuery(stmt, q.Query)

	workouts := []model.Workout{}
	rows, err := getDB().QueryContext(ctx, stmt)
	if err != nil {
		return workouts, fmt.Errorf("failed to query users. %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var w model.Workout
		if err := rows.Scan(&w.ID, &w.Name, &w.UserID, &w.Exercises, &w.Created, &w.Updated); err != nil {
			return workouts, fmt.Errorf("failed to scan. %w", err)
		}

		workouts = append(workouts, w)
	}

	return workouts, nil
}

func UpdateWorkout(ctx context.Context, workout model.Workout) error {
	_, err := getDB().ExecContext(ctx,
		`UPDATE sandbox.workout
		SET name = $1, exercises = $2
		WHERE id = $3`,
		workout.Name,
		workout.Exercises,
		workout.ID,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "u_user_name") {
					return ErrConflictWorkoutName
				}
				return fmt.Errorf("failed to update workout. conflict. %w", err)
			}
		}
		return fmt.Errorf("failed to update workout. %w", err)
	}

	return nil
}

func DeleteWorkout(ctx context.Context, userID string, workoutID string) error {
	_, err := getDB().ExecContext(ctx,
		`DELETE FROM sandbox.workout
		WHERE user_id = $1 AND id = $2`,
		userID, workoutID)
	if err != nil {
		return fmt.Errorf("failed to delete workout. %w", err)
	}

	return nil
}
