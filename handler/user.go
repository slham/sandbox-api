package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
	"github.com/slham/sandbox-api/valid"
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
		slog.Warn("error decoding create user request", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, []string{"malformed request body"})
		return
	}

	user, err := c.createUser(ctx, req)
	if err != nil {
		if errors.Is(err, ApiError) {
			apiErr, _ := err.(ApiError)
			slog.Warn("error creating user", "err", apiErr)
			request.RespondWithError(w, apiErr.StatusCode, apiErr.Errs)
			return
		}

		slog.Error("error creating user", "err", err)
		request.RespondWithError(w, http.StatusInternalServerError, []string{"internal server error"})
		return
	}

	request.RespondWithJSON(w, http.StatusCreated, user)
	return
}

func (c *UserController) createUser(ctx context.Context, req createUserRequest) (model.User, error) {
	user := model.User{}

	if err := validateCreateUserRequest(ctx, req); err != nil {
		return user, fmt.Errorf("failed to validate create user request. %w", err)
	}

	user.Password = ""
	return user, nil
}

func validateCreateUserRequest(ctx context.Context, req createUserRequest) error {
	apiErr := ApiError{StatusCode: 400}

	if len(req.Username) < 4 {
		apiErr = apiErr.Append("username must be at leat four characters long")
	}

	if ok, _ := valid.IsMediumPassword(req.Password); !ok {
		apiErr = apiErr.Append("password must be at least 8 characters long and contain at least one number, one special character, one upper case character, and one lower case character")
	}

	if ok, _ := valid.IsEmail(req.Email); !ok {
		apiErr = apiErr.Append("invalid email")
	}

	if apiError.HasError() {
		return apiError
	}

	return nil
}
