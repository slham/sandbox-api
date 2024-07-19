package handler

import (
	"database/sql"

	"github.com/slham/sandbox-api/dao"
)

type UserController struct {
	DB dao.Dao
}

func NewUserController() UserController {
	return UserController{
		DB: dao.GetDao(),
	}
}

func (c *UserController) GetDB() *sql.DB {
	return c.DB.DB
}
