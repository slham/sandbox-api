package integration

import (
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/go-playground/assert.v1"
)

func TestDeleteUser(t *testing.T) {
	testCases := []struct {
		name string
		code int
	}{
		{
			name: "delete happy path",
			code: http.StatusNoContent,
		},
	}

	suffix := "_test_update_user"
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
			resp, err := sendJSONHttpRequest("DELETE", url, "", testCookie)
			if err != nil {
				t.Log("err", err)
				t.Fail()
			}

			assert.Equal(t, resp.StatusCode, tc.code)
		})
	}
}
