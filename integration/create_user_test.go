//go:build integration
// +build integration

package integration

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	matcher "github.com/panta/go-json-matcher"
	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name    string
		req     string
		method  string
		resp    string
		code    int
		comment string
	}{
		{
			name:   "create fail validations",
			method: "POST",
			req:    `{"username": "bad", "password": "bad", "email": "bad"}`,
			resp:   `{"errors": "failed to validate create user request. username must be at leat four characters long. password must be at least 8 characters long and contain at least one number, one special character, one upper case character, and one lower case character. invalid email"}`,
			code:   http.StatusBadRequest,
		},
		{
			name:   "create happy path 1",
			method: "POST",
			req:    `{"username": "test_user_1%s", "password": "thisIsAG00dPassword!", "email": "a@b.c"}`,
			resp:   `{"id":"#regex ^user_[a-zA-Z0-9]{27}$","username": "test_user_1_test_create_user", "email": "a@b.c","created":"#datetime","updated":"#datetime"}`,
			code:   http.StatusCreated,
		},
		{
			name:   "create happy path 2",
			method: "POST",
			req:    `{"username": "test_user_2%s", "password": "thisIsAG00dPassword!", "email": "c@d.e"}`,
			resp:   `{"id":"#regex ^user_[a-zA-Z0-9]{27}$","username": "test_user_2_test_create_user", "email": "c@d.e","created":"#datetime","updated":"#datetime"}`,
			code:   http.StatusCreated,
		},
		{
			name:   "create happy path 3",
			method: "POST",
			req:    `{"username": "test_user_3%s", "password": "thisIsAG00dPassword!", "email": "f@g.h"}`,
			resp:   `{"id":"#regex ^user_[a-zA-Z0-9]{27}$","username": "test_user_3_test_create_user", "email": "f@g.h","created":"#datetime","updated":"#datetime"}`,
			code:   http.StatusCreated,
		},
		{
			name:   "create fail username conflict",
			method: "POST",
			req:    `{"username": "test_user_2%s", "password": "thisIsAG00dPassword!", "email": "good@gmail.com"}`,
			resp:   `{"errors": "username already exists"}`,
			code:   http.StatusConflict,
		},
		{
			name:   "create fail email conflict",
			method: "POST",
			req:    `{"username": "test_user_4%s", "password": "thisIsAG00dPassword!", "email": "c@d.e"}`,
			resp:   `{"errors": "email already exists"}`,
			code:   http.StatusConflict,
		},
	}

	suffix := "_test_create_user"
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := "http://localhost:8080/users"
			tc.req = fmt.Sprintf(tc.req, suffix)
			resp, err := sendJSONHttpRequest(tc.method, url, tc.req, nil)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			assert.Equal(t, resp.StatusCode, tc.code)

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			respString := string(bodyBytes)

			matches, err := matcher.JSONStringMatches(respString, tc.resp)
			if !matches || err != nil {
				t.Log("err", err)
				t.Log("matches", matches)
				t.Log("expected", tc.resp)
				t.Log("got", respString)
				t.Fail()
			}
		})
	}

	if err := cleanUpTestUsers(suffix); err != nil {
		t.Log(err)
		t.Fail()
	}
}
