package integration

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	matcher "github.com/panta/go-json-matcher"
	"gopkg.in/go-playground/assert.v1"
)

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name string
		req  string
		resp string
		code int
	}{
		{
			name: "update fail validations",
			req:  `{"username": "bad", "email": "bad"}`,
			resp: `{"errors": "failed to validate update user request. username must be at leat four characters long. invalid email"}`,
			code: http.StatusBadRequest,
		},
		{
			name: "update fail username and email conflict",
			req:  `{"username": "%s", "email": "%s"}`,
			resp: `{"errors": "email already exists"}`,
			code: http.StatusConflict,
		},
		{
			name: "update happy path",
			req:  `{"username": "%s", "password": "thisIsAG00dPassword!", "email": "%s"}`,
			resp: `{"id":"#regex ^user_[a-zA-Z0-9]{27}$","username": "%s", "email": "%s","created":"#datetime","updated":"#datetime"}`,
			code: http.StatusOK,
		},
	}

	suffix := "_test_update_user"
	email1 := randomEmail()
	username1 := randomUsername(suffix)
	userID, err := createTestUser(username1, email1)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	email2 := randomEmail()
	username2 := randomUsername(suffix)
	_, err = createTestUser(username2, email2)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	email3 := randomEmail()
	username3 := randomUsername(suffix)

	testCookie, err := loginTestUser(username1)
	if err != nil {
		t.Log("err", err)
		t.Fail()
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:8080/users/%s", userID)
			if strings.HasSuffix(tc.name, "conflict") {
				tc.req = fmt.Sprintf(tc.req, username2, email2)
			} else if strings.HasSuffix(tc.name, "happy path") {
				tc.req = fmt.Sprintf(tc.req, username3, email3)
			}
			resp, err := sendJSONHttpRequest("PATCH", url, tc.req, testCookie)
			if err != nil {
				t.Log("err", err)
				t.Log("req", tc.req)
				t.Fail()
			}

			assert.Equal(t, resp.StatusCode, tc.code)

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			respString := string(bodyBytes)

			if strings.HasSuffix(tc.name, "happy path") {
				tc.resp = fmt.Sprintf(tc.resp, username3, email3)
			}
			matches, err := matcher.JSONStringMatches(respString, tc.resp)
			if !matches || err != nil {
				t.Log("err", err)
				t.Log("got", respString)
				t.Log("wanted", tc.resp)
				t.Fail()
			}
		})
	}
	cleanUpTestUsers(suffix)
}
