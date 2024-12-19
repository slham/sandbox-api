//go:build integration
// +build integration

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	matcher "github.com/panta/go-json-matcher"
	"gopkg.in/go-playground/assert.v1"
)

func TestUser(t *testing.T) {
	testCases := []struct {
		name    string
		req     string
		method  string
		url     string
		resp    string
		code    int
		comment string
	}{
		{
			name:   "create fail validations",
			method: "POST",
			url:    "/users",
			req:    `{"username": "bad", "password": "bad", "email": "bad"}`,
			resp:   `{"errors": "failed to validate create user request. username must be at leat four characters long. password must be at least 8 characters long and contain at least one number, one special character, one upper case character, and one lower case character. invalid email"}`,
			code:   http.StatusBadRequest,
		},
		{
			name:   "create happy path 1",
			method: "POST",
			url:    "/users",
			req:    `{"username": "test_user_1", "password": "thisIsAG00dPassword!", "email": "a@b.c"}`,
			resp:   `{"id":"#regex ^user_[a-zA-Z0-9]{27}$","username": "test_user_1", "email": "a@b.c","created":"#datetime","updated":"#datetime"}`,
			code:   http.StatusCreated,
		},
		{
			name:   "create happy path 2",
			method: "POST",
			url:    "/users",
			req:    `{"username": "test_user_2", "password": "thisIsAG00dPassword!", "email": "c@d.e"}`,
			resp:   `{"id":"#regex ^user_[a-zA-Z0-9]{27}$","username": "test_user_2", "email": "c@d.e","created":"#datetime","updated":"#datetime"}`,
			code:   http.StatusCreated,
		},
		{
			name:   "create happy path 3",
			method: "POST",
			url:    "/users",
			req:    `{"username": "test_user_3", "password": "thisIsAG00dPassword!", "email": "f@g.h"}`,
			resp:   `{"id":"#regex ^user_[a-zA-Z0-9]{27}$","username": "test_user_3", "email": "f@g.h","created":"#datetime","updated":"#datetime"}`,
			code:   http.StatusCreated,
		},
		{
			name:   "create fail username conflict",
			method: "POST",
			url:    "/users",
			req:    `{"username": "test_user_2", "password": "thisIsAG00dPassword!", "email": "good@gmail.com"}`,
			resp:   `{"errors": "username already exists"}`,
			code:   http.StatusConflict,
		},
		{
			name:   "create fail email conflict",
			method: "POST",
			url:    "/users",
			req:    `{"username": "test_user_4", "password": "thisIsAG00dPassword!", "email": "c@d.e"}`,
			resp:   `{"errors": "email already exists"}`,
			code:   http.StatusConflict,
		},
		{
			name:   "update fail validations",
			method: "PATCH",
			url:    "/users/%s",
			req:    `{"username": "bad", "email": "bad"}`,
			resp:   `{"errors": "failed to validate update user request. username must be at leat four characters long. invalid email"}`,
			code:   http.StatusBadRequest,
		},
		{
			name:   "update fail username conflict",
			method: "PATCH",
			url:    "/users/%s",
			req:    `{"username": "test_user_2", "email": "good@gmail.com"}`,
			resp:   `{"errors": "username already exists"}`,
			code:   http.StatusConflict,
		},
		{
			name:   "update fail email conflict",
			method: "PATCH",
			url:    "/users/%s",
			req:    `{"username": "test_user_4", "email": "f@g.h"}`,
			resp:   `{"errors": "email already exists"}`,
			code:   http.StatusConflict,
		},
		{
			name:   "update happy path",
			method: "PATCH",
			url:    "/users/%s",
			req:    `{"username": "test_user_4", "password": "thisIsAG00dPassword!", "email": "i@j.k"}`,
			resp:   `{"id":"#regex ^user_[a-zA-Z0-9]{27}$","username": "test_user_4", "email": "i@j.k","created":"#datetime","updated":"#datetime"}`,
			code:   http.StatusOK,
		},
		{
			name:   "delete happy path",
			method: "DELETE",
			url:    "/users/%s",
			code:   http.StatusNoContent,
		},
		{
			name:    "get all happy path",
			method:  "GET",
			url:     "/users",
			code:    http.StatusOK,
			comment: "clean up test data",
		},
	}

	var testCookie http.Cookie
	userIDs := []string{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			bodyReader := bytes.NewReader([]byte(tc.req))
			if tc.method == "PATCH" || tc.method == "DELETE" {
				tc.url = fmt.Sprintf(tc.url, userIDs[0])
			}
			url := fmt.Sprintf("http://localhost:8080%s", tc.url)
			req, err := http.NewRequest(tc.method, url, bodyReader)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			req.Header.Set("Content-Type", "application/json")
			req.AddCookie(&testCookie)
			client := http.Client{
				Timeout: 10 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			assert.Equal(t, resp.StatusCode, tc.code)

			cookies := resp.Cookies()
			for _, cookie := range cookies {
				if cookie.Name == "sandbox-cookie" {
					testCookie = *cookie
					break
				}
			}

			if resp.StatusCode == http.StatusNoContent {
				t.Skipf("skipping test success")
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			respString := string(bodyBytes)

			if tc.method == "GET" {
				usersList := []map[string]string{}
				err = json.Unmarshal(bodyBytes, &usersList)
				if err != nil {
					t.Log(err)
					t.Fail()
				}

				for _, user := range usersList {
					u := fmt.Sprintf("http://localhost:8080/users/%s", user["id"])
					rq, err := http.NewRequest("DELETE", u, nil)
					if err != nil {
						t.Log(err)
						t.Fail()
					}
					rq.Header.Set("Content-Type", "application/json")
					req.AddCookie(&testCookie)
					resp, err := client.Do(rq)
					if err != nil {
						t.Log(err)
						t.Fail()
					}

					assert.Equal(t, resp.StatusCode, http.StatusNoContent)
				}

				t.Skipf("skipping test success")
			}

			matches, err := matcher.JSONStringMatches(respString, tc.resp)
			if !matches || err != nil {
				t.Log(err)
				t.Fail()
			}

			if tc.method == "POST" && resp.StatusCode == http.StatusCreated {
				respMap := map[string]any{}
				err = json.Unmarshal(bodyBytes, &respMap)
				if err != nil {
					t.Log(err)
					t.Fail()
				}
				if id, ok := respMap["id"].(string); ok {
					userIDs = append(userIDs, id)
				}
			}
		})
	}
}
