package handler

import (
	"database/sql"

	"github.com/slham/sandbox-api/dao"
)

type WorkoutController struct {
	DB dao.Dao
}

func NewWorkoutController() WorkoutController {
	return WorkoutController{
		DB: dao.GetDao(),
	}
}

func (c *WorkoutController) GetDB() *sql.DB {
	return c.DB.DB
}
