package dao

import (
	"context"
	"fmt"

	"github.com/slham/sandbox-api/model"
)

func InsertUser(ctx context.Context, u model.User) (model.User, error) {
	err := getDB().QueryRow(
		`INSERT INTO sandbox.user(
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
			$5)
		RETURNING id`,
		u.Username,
		u.Password,
		u.Email,
		u.Created,
		u.Updated).Scan(&u.ID)

	return u, fmt.Errorf("failed to insert user. %w", err)
}
