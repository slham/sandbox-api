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
	ErrConflictUsername = errors.New("username already exists")
	ErrConflictEmail    = errors.New("email already exists")
)

func InsertUser(ctx context.Context, u model.User) (model.User, error) {
	row := getDB().QueryRowContext(ctx,
		`INSERT INTO sandbox.user(
			id,
			username,
			password,
			email,
			created,
			updated
		)
		VALUES(
			$1,
			$2,
			$3,
			$4,
			$5,
			$6
		)`,
		u.ID,
		u.Username,
		u.Password,
		u.Email,
		u.Created,
		u.Updated)
	if err := row.Err(); err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "username") {
					return u, ErrConflictUsername
				} else if strings.Contains(pgErr.Message, "email") {
					return u, ErrConflictEmail
				}
				return u, fmt.Errorf("failed to insert user. conflict. %w", err)
			}
		}
		return u, fmt.Errorf("failed to insert user. %w", err)
	}

	return u, nil
}

type UserQuery struct {
	ID           string
	Username     string
	Email        string
	Sort         string
	SortCol      string
	Limit        int
	Offset       int
	HidePassword bool
}

func checkWhereClause(stmt string) string {
	if !strings.HasSuffix(stmt, "WHERE") {
		return fmt.Sprintf("%s AND", stmt)
	}
	return stmt
}

func GetUserByEmail(ctx context.Context, email string) (model.User, error) {
	q := UserQuery{Email: email}
	u, err := GetUser(ctx, q)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user by email. %w", err)
	}
	return u, nil
}

func GetUserByUsername(ctx context.Context, username string) (model.User, error) {
	q := UserQuery{Username: username}
	u, err := GetUser(ctx, q)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user by username. %w", err)
	}
	return u, nil
}

func GetUserByID(ctx context.Context, id string) (model.User, error) {
	q := UserQuery{ID: id, HidePassword: true}
	u, err := GetUser(ctx, q)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user by id. %w", err)
	}
	return u, nil
}

func GetUser(ctx context.Context, q UserQuery) (model.User, error) {
	users, err := GetUsers(ctx, q)
	if err != nil {
		return model.User{}, fmt.Errorf("failed to get user. %w", err)
	}

	return users[0], nil
}

func GetUsers(ctx context.Context, q UserQuery) ([]model.User, error) {
	stmt := `
		SELECT
			id,
			username,
			password,
			email,
			created,
			updated
		FROM
			sandbox.user
		WHERE`

	if q.ID != "" {
		stmt = fmt.Sprintf("%s %s='%s'", stmt, "id", q.ID)
	}
	if q.Username != "" {
		stmt = checkWhereClause(stmt)
		stmt = fmt.Sprintf("%s %s='%s'", stmt, "username", q.Username)
	}
	if q.Email != "" {
		stmt = checkWhereClause(stmt)
		stmt = fmt.Sprintf("%s %s='%s'", stmt, "email", q.Email)
	}
	if q.SortCol != "" {
		stmt = fmt.Sprintf("%s ORDER BY %s", stmt, q.SortCol)
	} else {
		stmt = fmt.Sprintf("%s ORDER BY id", stmt)
	}
	if q.Sort != "" {
		stmt = fmt.Sprintf("%s %s", stmt, q.Sort)
	} else {
		stmt = fmt.Sprintf("%s ASC", stmt)
	}
	if q.Limit > 0 {
		stmt = fmt.Sprintf("%s LIMIT %d", stmt, q.Limit)
	} else {
		stmt = fmt.Sprintf("%s LIMIT 100", stmt)
	}
	if q.Offset > 0 {
		stmt = fmt.Sprintf("%s OFFSET %d", stmt, q.Offset)
	} else {
		stmt = fmt.Sprintf("%s OFFSET 0", stmt)
	}

	users := []model.User{}
	rows, err := getDB().QueryContext(ctx, stmt)
	if err != nil {
		return users, fmt.Errorf("failed to query users. %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		var u model.User
		if err := rows.Scan(&u.ID, &u.Username, &u.Password, &u.Email, &u.Created, &u.Updated); err != nil {
			return users, fmt.Errorf("failed to scan.  %w", err)
		}

		if q.HidePassword {
			u.Password = ""
		}
		users = append(users, u)
	}

	return users, nil
}
