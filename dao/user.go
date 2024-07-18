package dao

import (
	"context"

	"github.com/slham/sandbox-api/model"
)

func InsertUser(ctx context.Context, u model.User) model.User {
	getDB().QueryRowContext(ctx,
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
	return u
}
