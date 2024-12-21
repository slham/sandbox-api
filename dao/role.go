package dao

import (
	"context"
	"fmt"

	"github.com/slham/sandbox-api/model"
)

func InsertRole(ctx context.Context, role model.Role) (model.Role, error) {
	stmt := `
		INSERT INTO sandbox.role
			(name)
		VALUES
			($1)
		RETURNING id, created, updated`
	err := getDB().QueryRow(stmt, role.Name).Scan(&role.ID, &role.Created, &role.Updated)
	if err != nil {
		return role, fmt.Errorf("failed to insert role. %w", err)
	}

	return role, nil
}

type RoleQuery struct {
	ID     int
	Name   string
	UserID string
	Query
}

func GetRoleByID(ctx context.Context, id int) (model.Role, error) {
	q := RoleQuery{ID: id}
	role, err := GetRole(ctx, q)
	if err != nil {
		return model.Role{}, fmt.Errorf("failed to get role by id. %w", err)
	}

	return role, nil
}

func GetRoleByName(ctx context.Context, name string) (model.Role, error) {
	q := RoleQuery{Name: name}
	role, err := GetRole(ctx, q)
	if err != nil {
		return model.Role{}, fmt.Errorf("failed to get role by name. %w", err)
	}

	return role, nil
}

func GetUserRoles(ctx context.Context, userID string) ([]model.Role, error) {
	stmt := `
		SELECT 
			r.id, r.name, r.created, r.updated
		FROM
			sandbox.role r
		INNER JOIN
			sandbox.user_role ur
		ON
			ur.role_id = r.id
		INNER JOIN
			sandbox.user u
		ON
			u.id = ur.user_id
		WHERE
			u.id = $1`

	roles := []model.Role{}
	rows, err := getDB().QueryContext(ctx, stmt)
	if err != nil {
		return roles, fmt.Errorf("failed to query roles. %w", err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		return roles, fmt.Errorf("failed to query roles. rows. %w", err)
	}

	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Created, &role.Updated); err != nil {
			return roles, fmt.Errorf("failed to scan.  %w", err)
		}

		roles = append(roles, role)
	}

	return roles, nil

}

func GetRole(ctx context.Context, q RoleQuery) (model.Role, error) {
	roles, err := GetRoles(ctx, q)
	if err != nil {
		return model.Role{}, fmt.Errorf("failed to get role. %w", err)
	}

	return roles[0], nil
}

func GetRoles(ctx context.Context, q RoleQuery) ([]model.Role, error) {
	stmt := `
		SELECT
			id, name, created, updated
		FROM
			sandbox.role`

	if q.ID != 0 || q.Name != "" || q.UserID != "" {
		stmt = fmt.Sprintf("%s %s", stmt, "WHERE")
		if q.ID != 0 {
			stmt = fmt.Sprintf("%s %s='%d'", stmt, "id", q.ID)
		}
		if q.Name != "" {
			stmt = fmt.Sprintf("%s %s='%s'", stmt, "name", q.Name)
		}
	}

	stmt = addDefaultQuery(stmt, q.Query)

	roles := []model.Role{}
	rows, err := getDB().QueryContext(ctx, stmt)
	if err != nil {
		return roles, fmt.Errorf("failed to query roles. %w", err)
	}

	defer rows.Close()

	if rows.Err() != nil {
		return roles, fmt.Errorf("failed to query roles. rows. %w", err)
	}

	for rows.Next() {
		var role model.Role
		if err := rows.Scan(&role.ID, &role.Name, &role.Created, &role.Updated); err != nil {
			return roles, fmt.Errorf("failed to scan.  %w", err)
		}

		roles = append(roles, role)
	}

	return roles, nil
}

func DeleteRole(ctx context.Context, id string) error {
	_, err := getDB().ExecContext(ctx,
		`DELETE FROM sandbox.role 
		WHERE id = $1`,
		id)
	if err != nil {
		return fmt.Errorf("failed to delete role. %w", err)
	}

	return nil
}
