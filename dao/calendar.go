package dao

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"
	"github.com/slham/sandbox-api/model"
)

func InsertCalendar(ctx context.Context, calendar model.Calendar) (model.Calendar, error) {
	_, err := getDB().ExecContext(ctx,
		`INSERT INTO sandbox.calendar(
    id,
    name,
    user_id
  )
    VALUES(
    $1,
    $2,
    $3
  )`,
		calendar.ID,
		calendar.Name,
		calendar.UserID,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "u_calendar_name") {
					return calendar, ErrConflictCalendarName
				}
				return calendar, fmt.Errorf("failed to insert calendar. conflict. %w", err)
			}
		}
		return calendar, fmt.Errorf("failed to insert calendar. %w", err)
	}
	return calendar, nil
}

type CalendarQuery struct {
	ID     string
	UserID string
	Query
}

func GetCalendarByUserID(ctx context.Context, userID string) (model.Calendar, error) {
	q := CalendarQuery{UserID: userID}
	w, err := GetCalendar(ctx, q)
	if err != nil {
		return model.Calendar{}, fmt.Errorf("failed to get calendar by user id. %w", err)
	}
	return w, nil
}

func GetCalendarByID(ctx context.Context, userID string, calendarID string) (model.Calendar, error) {
	q := CalendarQuery{ID: calendarID, UserID: userID}
	w, err := GetCalendar(ctx, q)
	if err != nil {
		return model.Calendar{}, fmt.Errorf("failed to get calendar by id. %w", err)
	}
	return w, nil
}

func GetCalendar(ctx context.Context, q CalendarQuery) (model.Calendar, error) {
	calendars, err := GetCalendars(ctx, q)
	if err != nil {
		return model.Calendar{}, fmt.Errorf("failed to get calendars. %w", err)
	}

	if len(calendars) != 1 {
		return model.Calendar{}, ErrCalendarNotFound
	}

	return calendars[0], nil
}

func GetCalendars(ctx context.Context, q CalendarQuery) ([]model.Calendar, error) {
	stmt := `
		SELECT
			id,
			name,
			user_id,
			created,
			updated
		FROM
			sandbox.calendar
		WHERE`

	if q.ID != "" {
		stmt = fmt.Sprintf("%s %s='%s'", stmt, "id", q.ID)
	}
	if q.UserID != "" {
		stmt = checkWhereClause(stmt)
		stmt = fmt.Sprintf("%s %s='%s'", stmt, "user_id", q.UserID)
	}

	stmt = addDefaultQuery(stmt, q.Query)

	calendars := []model.Calendar{}
	rows, err := getDB().QueryContext(ctx, stmt)
	if err != nil {
		return calendars, fmt.Errorf("failed to query calendas. %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var c model.Calendar
		if err := rows.Scan(&c.ID, &c.Name, &c.UserID, &c.Created, &c.Updated); err != nil {
			return calendars, fmt.Errorf("failed to scan. %w", err)
		}

		calendars = append(calendars, c)
	}

	return calendars, nil
}

func UpdateCalendar(ctx context.Context, calendar model.Calendar) error {
	_, err := getDB().ExecContext(ctx,
		`UPDATE sandbox.calendar
		SET name = $1
		WHERE id = $2`,
		calendar.Name,
		calendar.ID,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "u_user_name") {
					return ErrConflictCalendarName
				}
				return fmt.Errorf("failed to update calendar. conflict. %w", err)
			}
		}
		return fmt.Errorf("failed to update calendar. %w", err)
	}

	return nil
}

func DeleteCalendar(ctx context.Context, userID string, calendarID string) error {
	_, err := getDB().ExecContext(ctx,
		`DELETE FROM sandbox.calendar
		WHERE user_id = $1 AND id = $2`,
		userID, calendarID)
	if err != nil {
		return fmt.Errorf("failed to delete calendar. %w", err)
	}

	return nil
}
