package integration

import (
	"fmt"
	"io"
	"net/http"
	"testing"

	matcher "github.com/panta/go-json-matcher"
	"gopkg.in/go-playground/assert.v1"
)

func TestGetUser(t *testing.T) {
	testCases := []struct {
		name    string
		req     string
		method  string
		resp    string
		code    int
		comment string
	}{
		{
			name:   "read happy path",
			method: "GET",
			resp:   `{"id":"#regex ^user_[a-zA-Z0-9]{27}$","username": "%s", "email": "%s","created":"#datetime","updated":"#datetime"}`,
			code:   http.StatusOK,
		},
	}

	suffix := "_test_read_user"
	email := randomEmail()
	username := randomUsername(suffix)
	userID, err := createTestUser(username, email)
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	testCookie, err := loginTestUser(username)
	if err != nil {
		t.Log("err", err)
		t.Fail()
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			url := fmt.Sprintf("http://localhost:8080/users/%s", userID)
			resp, err := sendJSONHttpRequest(tc.method, url, tc.req, testCookie)
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

			tc.resp = fmt.Sprintf(tc.resp, username, email)
			matches, err := matcher.JSONStringMatches(respString, tc.resp)
			if !matches || err != nil {
				t.Log(err)
				t.Fail()
			}
		})
	}
}
