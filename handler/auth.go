package handler

import (
	"database/sql"

	"github.com/slham/sandbox-api/dao"
)

type AuthController struct {
	DB dao.Dao
}

func NewAuthController() AuthController {
	return AuthController{
		DB: dao.GetDao(),
	}
}

func (c *AuthController) GetDB() *sql.DB {
	return c.DB.DB
}
