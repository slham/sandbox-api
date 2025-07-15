package dao

import (
	"context"
	"fmt"

	"github.com/slham/sandbox-api/model"
)

func InsertSnapshot(ctx context.Context, snapshot model.Snapshot) (model.Snapshot, error) {
	_, err := getDB().ExecContext(ctx,
		`INSERT INTO sandbox.snapshot(
			id,
			calendar_id,
			done,
			workout
		)
		VALUES(
			$1,
			$2,
			$3,
			$4
		)`,
		snapshot.ID,
		snapshot.CalendarID,
		snapshot.Done,
		snapshot.Workout,
	)
	if err != nil {
		return snapshot, fmt.Errorf("failed to insert snapshot. %w", err)
	}
	return snapshot, nil
}

type SnapshotQuery struct {
	ID     string
	CalendarID string
	Query
}

func GetSnapshotByCalendarID(ctx context.Context, calendarID string) (model.Snapshot, error) {
	q := SnapshotQuery{CalendarID: calendarID}
	w, err := GetSnapshot(ctx, q)
	if err != nil {
		return model.Snapshot{}, fmt.Errorf("failed to get snapshot by calendar id. %w", err)
	}
	return w, nil
}

func GetSnapshotByID(ctx context.Context, snapshotID string) (model.Snapshot, error) {
	q := SnapshotQuery{ID: snapshotID}
	w, err := GetSnapshot(ctx, q)
	if err != nil {
		return model.Snapshot{}, fmt.Errorf("failed to get snapshot by id. %w", err)
	}
	return w, nil
}

func GetSnapshot(ctx context.Context, q SnapshotQuery) (model.Snapshot, error) {
	snapshots, err := GetSnapshots(ctx, q)
	if err != nil {
		return model.Snapshot{}, fmt.Errorf("failed to get snapshots. %w", err)
	}

	if len(snapshots) != 1 {
		return model.Snapshot{}, ErrSnapshotNotFound
	}

	return snapshots[0], nil
}

func GetSnapshots(ctx context.Context, q SnapshotQuery) ([]model.Snapshot, error) {
	stmt := `
		SELECT
			id,
			calendar_id,
			done,
			workout,
			created,
			updated
		FROM
			sandbox.snapshot
		WHERE`

	if q.ID != "" {
		stmt = fmt.Sprintf("%s %s='%s'", stmt, "id", q.ID)
	}
	if q.CalendarID != "" {
		stmt = checkWhereClause(stmt)
		stmt = fmt.Sprintf("%s %s='%s'", stmt, "calendar_id", q.CalendarID)
	}

	stmt = addDefaultQuery(stmt, q.Query)

	snapshots := []model.Snapshot{}
	rows, err := getDB().QueryContext(ctx, stmt)
	if err != nil {
		return snapshots, fmt.Errorf("failed to query snapshots. %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var s model.Snapshot
		if err := rows.Scan(&s.ID, &s.CalendarID, &s.Done, &s.Workout, &s.Created, &s.Updated); err != nil {
			return snapshots, fmt.Errorf("failed to scan. %w", err)
		}

		snapshots = append(snapshots, s)
	}

	return snapshots, nil
}

func UpdateSnapshot(ctx context.Context, snapshot model.Snapshot) error {
	_, err := getDB().ExecContext(ctx,
		`UPDATE sandbox.snapshot
		SET done = $1, workout = $2
		WHERE id = $3`,
		snapshot.Done,
		snapshot.Workout,
		snapshot.ID,
	)
	if err != nil {
		return fmt.Errorf("failed to update snapshot. %w", err)
	}

	return nil
}

func DeleteSnapshot(ctx context.Context, snapshotID string) error {
	_, err := getDB().ExecContext(ctx,
		`DELETE FROM sandbox.snapshot
		WHERE id = $1`,
		snapshotID)
	if err != nil {
		return fmt.Errorf("failed to delete snapshot. %w", err)
	}

	return nil
}
