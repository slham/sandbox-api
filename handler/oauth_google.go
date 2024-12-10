package handler

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log/slog"
	"math/big"
	"net/http"
	"os"
	"time"

	"github.com/segmentio/ksuid"
	"github.com/slham/sandbox-api/crypt"
	"github.com/slham/sandbox-api/dao"
	"github.com/slham/sandbox-api/model"
	"github.com/slham/sandbox-api/request"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

var (
	// Scopes: OAuth 2.0 scopes provide a way to limit the amount of access that is granted to an access token.
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8000/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_OAUTH_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}

	oauthFlowMap = map[string]func(context.Context, GoogleOAuthUserInfo) (model.User, error){
		"login":    handleOauthGoogleLogin,
		"register": handleOauthGoogleRegister,
	}
)

const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

type GoogleOAuthUserInfo struct {
	ID            string `json:"id"`
	Email         string `json:"email"`
	VerifiedEmail bool   `json:"verified_email"`
	Name          string `json:"name"`
	GivenName     string `json:"give_name"`
	FamilyName    string `json:"family_name"`
	Picture       string `json:"picture"`
}

func handleOauthGoogleError(w http.ResponseWriter, err error) {
	if errors.Is(err, ApiErrBadRequest) {
		slog.Warn("error oauth google", "err", err)
		request.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	} else if errors.Is(err, ApiErrConflict) {
		slog.Error("user already exists", "err", err)
		request.RespondWithError(w, http.StatusConflict, err.Error())
		return
	}

	slog.Error("error oauth google", "err", err)
	request.RespondWithError(w, http.StatusInternalServerError, "internal server error")
	return
}

func (c *AuthController) OauthGoogleLogin(w http.ResponseWriter, r *http.Request) {

	// Create oauthState cookie
	oauthState := generateStateOauthCookie(w)

	/*
		AuthCodeURL receive state that is a token to protect the user from CSRF attacks. You must always provide a non-empty string and
		validate that it matches the the state query parameter on your redirect callback.
	*/
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func (c *AuthController) OauthGoogleCallback(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	query := r.URL.Query()
	oauthFlow := query.Get("oauth-flow")

	oauthState, err := r.Cookie("oauthstate")
	if err != nil {
		slog.Error("failed to read oauthstate cookie", "err", err)
		handleOauthGoogleError(w, err)
		return
	}

	state := r.FormValue("state")
	code := r.FormValue("code")

	if state != oauthState.Value {
		slog.Error("invalid oauth google state", "have", state, "wanted", oauthState.Value)
		handleOauthGoogleError(w, err)
		return
	}

	data, err := getUserDataFromGoogle(ctx, code)
	if err != nil {
		slog.Error("failed to get user data from google", "err", err)
		handleOauthGoogleError(w, err)
		return
	}

	userInfo := GoogleOAuthUserInfo{}
	err := json.Unmarshal(data, &userInfo)
	if err != nil {
		slog.Error("failed to unmarshal user data", "err", err)
		handleOauthGoogleError(w, err)
		return
	}

	user, err := oauthFlowMap[flow](ctx, userInfo)
	if err != nil {
		slog.Error("failed to check user", "err", err)
		handleOauthGoogleError(w, err)
		return
	}
	//TODO: set up user session
	fmt.Fprintf(w, "UserInfo: %s\n", data)
}

func handleOauthGoogleRegister(ctx context.Context, userInfo GoogleOAuthUserInfo) (model.User, error) {
	user, err := dao.GetUserByEmail(ctx, userInfo.Email)
	if errors.Is(err, sql.ErrNoRows) {
		user, err = makeUser(ctx, userInfo)
		if err != nil {
			return user, fmt.Errorf("failed to create new user. %w", err)
		}
		return user, nil
	} else if err != nil {
		return user, fmt.Errorf("failed to get user. %w", err)
	} else {
		return user, NewApiError(409, ApiErrConflict)
	}
}

func handleOauthGoogleLogin(ctx context.Context, userInfo GoogleOAuthUserInfo) (model.User, error) {
	user, err := dao.GetUserByEmail(ctx, userInfo.Email)
	if err != nil {
		return user, fmt.Errorf("failed to get user. %w", err)
	}

	return user, nil
}

func generateStateOauthCookie(w http.ResponseWriter) string {
	var expiration = time.Now().Add(20 * time.Minute)

	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	cookie := http.Cookie{Name: "oauthstate", Value: state, Expires: expiration}
	http.SetCookie(w, &cookie)

	return state
}

func getUserDataFromGoogle(ctx context.Context, code string) ([]byte, error) {
	// Use code to get token and get user info from Google.

	token, err := googleOauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("code exchange wrong: %s", err.Error())
	}
	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, fmt.Errorf("failed getting user info: %s", err.Error())
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed read response: %s", err.Error())
	}
	return contents, nil
}

func generatePassword(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"

	var password []byte
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		password = append(password, charset[randomIndex.Int64()])
	}
	return string(password), nil
}

func newPassword() (string, error) {
	password, passwordErr := generatePassword(8)
	if passwordErr != nil {
		return password, fmt.Errorf("failed to generate new user password. %w", passwordErr)
	}
	password, encryptErr := crypt.Encrypt(password)
	if encryptErr != nil {
		return password, fmt.Errorf("failed to encrypt password. %w", encryptErr)
	}

	return password, nil
}

func makeUser(ctx context.Context, userInfo GoogleOAuthUserInfo) (model.User, error) {
	password, passwordErr := newPassword()
	if passwordErr != nil {
		return user, fmt.Errorf("failed to generate new user password. %w", passwordErr)
	}
	now := time.Now()
	newUser := model.User{
		ID:       fmt.Sprintf("user_%s", ksuid.New().String()),
		Username: userInfo.Name,
		Password: password,
		Email:    userInfo.Email,
		Created:  now,
		Updated:  now,
	}
	user, err = dao.InsertUser(ctx, user)
	if err != nil {
		if errors.Is(err, dao.ErrConflictUsername) {
			return user, NewApiError(409, ApiErrConflict).Append("username already exists")
		}
		if errors.Is(err, dao.ErrConflictEmail) {
			return user, NewApiError(409, ApiErrConflict).Append("email already exists")
		}
		return user, fmt.Errorf("failed to insert user. %w", err)
	}

	user.Password = ""
	return user, nil
}
