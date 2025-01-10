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
	ErrUserNotFound     = errors.New("user does not exist")
	ErrWorkoutNotFound  = errors.New("workout does not exist")
	ErrRoleNotFound     = errors.New("role does not exist")
)

func InsertUser(ctx context.Context, user model.User) (model.User, error) {
	_, err := getDB().ExecContext(ctx,
		`INSERT INTO sandbox.user(
			id,
			username,
			password,
			email
		)
		VALUES(
			$1,
			$2,
			$3,
			$4
		)`,
		user.ID,
		user.Username,
		user.Password,
		user.Email,
	)
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				if strings.Contains(pgErr.Message, "username") {
					return user, ErrConflictUsername
				} else if strings.Contains(pgErr.Message, "email") {
					return user, ErrConflictEmail
				}
				return user, fmt.Errorf("failed to insert user. conflict. %w", err)
			}
		}
		return user, fmt.Errorf("failed to insert user. %w", err)
	}

	if err := insertUserRoles(ctx, user); err != nil {
		return user, fmt.Errorf("faild to insert user roles. %w", err)
	}

	return user, nil
}

func insertUserRoles(ctx context.Context, user model.User) error {
	for i := range user.Roles {
		role := user.Roles[i]
		if err := insertUserRole(ctx, user.ID, role.ID); err != nil {
			return fmt.Errorf("failed to insert user role. %w", err)
		}
	}
	return nil
}

func insertUserRole(ctx context.Context, userID string, roleID int) error {
	stmt := `
		INSERT INTO sandbox.user_role
			(user_id, role_id)
		VALUES
			($1, $2)`
	_, err := getDB().ExecContext(ctx, stmt, userID, roleID)
	if err != nil {
		return fmt.Errorf("failed to insert user role. %w", err)
	}

	return nil
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

	if len(users) != 1 {
		return model.User{}, ErrUserNotFound
	}

	return users[0], nil
}

func GetUsers(ctx context.Context, q UserQuery) ([]model.User, error) {
	stmt := `
		SELECT
			id, username, password, email, created, updated
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

	stmt = addDefaultQuery(stmt, q.Query)

	users := []model.User{}
	rows, err := getDB().QueryContext(ctx, stmt)
	if err != nil {
		return users, fmt.Errorf("failed to query users. %w", err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		return users, fmt.Errorf("failed to query users. rows. %w", err)
	}

	for rows.Next() {
		var user model.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password, &user.Email, &user.Created, &user.Updated); err != nil {
			return users, fmt.Errorf("failed to scan.  %w", err)
		}

		if q.HidePassword {
			user.Password = ""
		}
		roles, err := GetUserRoles(ctx, user.ID)
		if err != nil {
			return users, fmt.Errorf("failed to get user (%s) roles. %w", user.ID, err)
		}
		user.Roles = roles
		users = append(users, user)
	}

	return users, nil
}

func UpdateUser(ctx context.Context, user model.User) error {
	_, err := getDB().ExecContext(ctx,
		`UPDATE sandbox.user 
		SET username = $1, email = $2
		WHERE id = $3`,
		user.Username,
		user.Email,
		user.ID,
	)
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
