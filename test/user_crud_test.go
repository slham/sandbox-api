package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func TestUser(t *testing.T) {
	testCases := []struct {
		name   string
		req    string
		method string
		url    string
		resp   map[string]string
		code   int
	}{
		{
			name:   "create fail validations",
			method: "POST",
			url:    "/users",
			req:    `{"username": "bad", "password": "bad", "email": "bad"}`,
			resp:   map[string]string{"errors": "failed to validate create user request. username must be at leat four characters long. password must be at least 8 characters long and contain at least one number, one special character, one upper case character, and one lower case character. invalid email"},
			code:   http.StatusBadRequest,
		},
		{
			name:   "create happy path 1",
			method: "POST",
			url:    "/users",
			req:    `{"username": "test_user_1", "password": "thisIsAG00dPassword!", "email": "a@b.c"}`,
			resp:   map[string]string{"username": "test_user_1", "email": "a@b.c"},
			code:   http.StatusCreated,
		},
		{
			name:   "create happy path 2",
			method: "POST",
			url:    "/users",
			req:    `{"username": "test_user_2", "password": "thisIsAG00dPassword!", "email": "c@d.e"}`,
			resp:   map[string]string{"username": "test_user_2", "email": "c@d.e"},
			code:   http.StatusCreated,
		},
		{
			name:   "create fail username conflict",
			method: "POST",
			url:    "/users",
			req:    `{"username": "test_user_2", "password": "thisIsAG00dPassword!", "email": "good@gmail.com"}`,
			resp:   map[string]string{"errors": "username already exists"},
			code:   http.StatusConflict,
		},
		{
			name:   "create fail email conflict",
			method: "POST",
			url:    "/users",
			req:    `{"username": "test_user_3", "password": "thisIsAG00dPassword!", "email": "c@d.e"}`,
			resp:   map[string]string{"errors": "email already exists"},
			code:   http.StatusConflict,
		},
		{
			name:   "update fail validations",
			method: "PATCH",
			url:    "/users/%s",
			req:    `{"username": "bad", "email": "bad"}`,
			resp:   map[string]string{"errors": "failed to validate update user request. username must be at leat four characters long. invalid email"},
			code:   http.StatusBadRequest,
		},
		{
			name:   "update fail username conflict",
			method: "PATCH",
			url:    "/users/%s",
			req:    `{"username": "test_user_2", "email": "good@gmail.com"}`,
			resp:   map[string]string{"errors": "username already exists"},
			code:   http.StatusConflict,
		},
		{
			name:   "update fail email conflict",
			method: "PATCH",
			url:    "/users/%s",
			req:    `{"username": "test_user_3", "email": "c@d.e"}`,
			resp:   map[string]string{"errors": "email already exists"},
			code:   http.StatusConflict,
		},
		{
			name:   "update happy path",
			method: "PATCH",
			url:    "/users/%s",
			req:    `{"username": "test_user_3", "password": "thisIsAG00dPassword!", "email": "f@g.h"}`,
			resp:   map[string]string{"username": "test_user", "email": "a@b.c"},
			code:   http.StatusOK,
		},
		{
			name:   "delete happy path",
			method: "DELETE",
			url:    "/users/%s",
			resp:   map[string]string{},
			code:   http.StatusNoContent,
		},
	}

	userIDs := []string{}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println("name", tc.name, "userIDs", userIDs)
			bodyReader := bytes.NewReader([]byte(tc.req))
			if tc.method != "POST" {
				tc.url = fmt.Sprintf(tc.url, userIDs[0])
			}
			url := fmt.Sprintf("http://localhost:8080%s", tc.url)
			req, err := http.NewRequest(tc.method, url, bodyReader)
			if err != nil {
				t.Log(err)
				t.Fail()
			}
			req.Header.Set("Content-Type", "application/json")
			client := http.Client{
				Timeout: 10 * time.Second,
			}
			resp, err := client.Do(req)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			assert.Equal(t, resp.StatusCode, tc.code)

			if resp.StatusCode == http.StatusNoContent {
				continue
			}

			bodyBytes, err := io.ReadAll(resp.Body)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			respMap := map[string]string{}
			err = json.Unmarshal(bodyBytes, &respMap)
			if err != nil {
				t.Log(err)
				t.Fail()
			}

			for k, v := range tc.resp {
				assert.Equal(t, respMap[k], v)
			}

			if tc.method == "POST" && resp.StatusCode == http.StatusCreated {
				userIDs = append(userIDs, respMap["id"])
			}
		})
	}
}
