package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"testing"
	"time"

	"gopkg.in/go-playground/assert.v1"
)

func TestCreateUser(t *testing.T) {
	testCases := []struct {
		name string
		req  string
		resp map[string]string
		code int
	}{
		{
			name: "happy path",
			req:  `{"username": "test_user", "password": "thisIsAG00dPassword!", "email": "a@b.c"}`,
			resp: map[string]string{"username": "test_user", "email": "a@b.c"},
			code: http.StatusCreated,
		},
		{
			name: "fail validations",
			req:  `{"username": "bad", "password": "bad", "email": "bad"}`,
			resp: map[string]string{"errors": "failed to validate create user request. username must be at leat four characters long. password must be at least 8 characters long and contain at least one number, one special character, one upper case character, and one lower case character. invalid email"},
			code: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		bodyReader := bytes.NewReader([]byte(tc.req))
		url := "http://localhost:8080/users"
		req, err := http.NewRequest(http.MethodPost, url, bodyReader)
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

	}
}
