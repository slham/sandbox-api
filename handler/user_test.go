package handler

import (
	"net/http"
	"testing"

	"github.com/go-resty/resty/v2"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name string
		req  string
		resp string
		code int
	}{
		{
			name: "happy path",
			req:  `{"username": "user_one", "password": "thisIsAG00dPassword!", "email": "a@b.c"}`,
			resp: `{"username": "user_one", "email": "a@b.c"}`,
			code: http.StatusCreated,
		},
		{
			name: "fail validations",
			req:  `{"username": "bad", "password": "bad", "email": "bad"}`,
			resp: `{"errors":"failed to validate create user request. username must be at leat four characters long. password must be at least 8 characters long and contain at least one number, one special character, one upper case character, and one lower case character. invalid email"}`,
			code: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		result := map[string]any{}
		client := resty.New()
		resp, err := client.R().
			SetHeader("Content-Type", "application/json").
			SetBody(tc.req).
			SetResult(&result).
			Post("http://localhost:8080/users")
		if err != nil {
			t.Log(err)
			t.Fail()
		}

		t.Logf("resp: %+v\n", resp)
		t.Logf("result: %+v\n", result)
	}
}
