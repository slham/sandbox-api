package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
)

type createUserRequest struct {
	UserName string `json:"name"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

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

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	req := createUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return
	}

	user, err := c.createUser(ctx, req)
	if err != nil {
		return
	}

	request.RespondWithJSON(w, http.StatusCreated, user)
	return
}

func (c *UserController) createUser(ctx context.Context, req createUserRequest) (model.User, error) {
	user := model.User{}
	return user, nil
}
