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
	_, err := getDB().ExecContext(ctx,
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
	if err != nil {
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
	HidePassword bool
	Query
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
			sandbox.user`

	if q.ID != "" || q.Username != "" || q.Email != "" {
		stmt = fmt.Sprintf("%s %s", stmt, "WHERE")
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

func UpdateUser(ctx context.Context, user model.User) error {
	_, err := getDB().ExecContext(ctx,
		`UPDATE sandbox.user 
		SET username = $1, email = $2, updated = $3
		WHERE id = $4`,
		user.Username,
		user.Email,
		user.Updated,
		user.ID)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "username") {
					return ErrConflictUsername
				} else if strings.Contains(pgErr.Message, "email") {
					return ErrConflictEmail
				}
				return fmt.Errorf("failed to update user. conflict. %w", err)
			}
		}
		return fmt.Errorf("failed to update user. %w", err)
	}

	return nil
}

func DeleteUser(ctx context.Context, id string) error {
	_, err := getDB().ExecContext(ctx,
		`DELETE FROM sandbox.user 
		WHERE id = $1`,
		id)
	if err != nil {
		return fmt.Errorf("failed to delete user. %w", err)
	}

	return nil
}
