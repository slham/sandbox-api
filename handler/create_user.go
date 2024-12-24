package handler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/segmentio/ksuid"
	"github.com/slham/sandbox-api/crypt"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
	"github.com/slham/sandbox-api/valid"
)

type createUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	slog.DebugContext(ctx, "create user request")
	req := createUserRequest{}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		slog.WarnContext(ctx, "error decoding create user request", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, "malformed request body")
		return
	}

	newContext, user, err := c.createUser(ctx, req)
	r = r.WithContext(newContext)
	if err != nil {
		if errors.Is(err, ApiErrBadRequest) {
			slog.WarnContext(ctx, "error creating user", "err", err)
			request.RespondWithError(w, http.StatusBadRequest, err.Error())
			return
		}

		if errors.Is(err, ApiErrConflict) {
			slog.WarnContext(ctx, "error creating user", "err", err)
			request.RespondWithError(w, http.StatusConflict, err.Error())
			return
		}

		slog.ErrorContext(ctx, "error creating user", "err", err)
		request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
		return
	}

	request.RespondWithJSON(w, http.StatusCreated, user)
}

func (c *UserController) createUser(ctx context.Context, req createUserRequest) (context.Context, model.User, error) {
	user := model.User{}

	if err := validateCreateUserRequest(ctx, req); err != nil {
		return ctx, user, fmt.Errorf("failed to validate create user request. %w", err)
	}

	var err error
	user.Password, err = crypt.Encrypt(req.Password)
	if err != nil {
		return ctx, user, fmt.Errorf("failed to encrypt password. %w", err)
	}

	role, err := dao.GetRole(ctx, dao.RoleQuery{Name: "CIVILIAN"})
	if err != nil {
		return ctx, user, fmt.Errorf("failed to get default user role. %w", err)
	}

	user.ID = newUserID()
	user.Username = req.Username
	user.Email = req.Email
	user.Roles = []model.Role{role}

	user, err = dao.InsertUser(ctx, user)
	if err != nil {
		if errors.Is(err, dao.ErrConflictUsername) {
			return ctx, user, NewApiError(409, ApiErrConflict).Append("username already exists")
		}
		if errors.Is(err, dao.ErrConflictEmail) {
			return ctx, user, NewApiError(409, ApiErrConflict).Append("email already exists")
		}
		return ctx, user, fmt.Errorf("failed to insert user. %w", err)
	}

	user.Password = ""
	user.Roles = []model.Role{}

	return ctx, user, nil
}

func validateCreateUserRequest(ctx context.Context, req createUserRequest) error {
	apiErr := NewApiError(400, ApiErrBadRequest)

	if len(req.Username) < 4 {
		apiErr = apiErr.Append("username must be at leat four characters long")
	}

	if ok := valid.IsMediumPassword(req.Password); !ok {
		passwordErrMsg := "password must be at least 8 characters long and contain at least one number, one special character, one upper case character, and one lower case character"
		apiErr = apiErr.Append(passwordErrMsg)
	}

	if err := valid.IsEmail(req.Email); err != nil {
		apiErr = apiErr.Append("invalid email")
	}

	if apiErr.HasError() {
		return apiErr
	}

	return nil
}

func newUserID() string {
	return fmt.Sprintf("user_%s", ksuid.New().String())
}
